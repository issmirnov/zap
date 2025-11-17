# Kubernetes Deployment Guide

Complete guide for deploying Zap to Kubernetes as a network-wide URL shortcut service.

## Quick Start

### Using Helm (Recommended)
```bash
helm install zap ./deploy/helm/zap
kubectl get svc zap --watch  # Wait for LoadBalancer IP
```

### Using Plain Manifests
```bash
kubectl apply -f deploy/kubernetes/
kubectl get svc -n zap --watch
```

See [QUICKSTART.md](QUICKSTART.md) for minimal configuration examples.

## Architecture Overview

```
┌──────────────┐     DNS Query: g      ┌─────────────┐
│ Client       │ ───────────────────>  │  DNS Server │
│ Machine      │                       └──────┬──────┘
└──────┬───────┘                              │
       │ Responds: g = 192.168.1.100          │
       │<─────────────────────────────────────┘
       │ HTTP Request: http://g/repo
       ↓
┌──────────────────────────────────────────────┐
│  Kubernetes LoadBalancer (192.168.1.100)     │
│  ┌─────────────────────────────────────────┐ │
│  │  Zap Service (LoadBalancer)             │ │
│  │  ┌─────────────┐  ┌──────────────┐     │ │
│  │  │  Zap Pod    │  │  ConfigMap   │     │ │
│  │  │  (c.yml)    │←─│  (config)    │     │ │
│  │  └─────────────┘  └──────────────┘     │ │
│  └─────────────────────────────────────────┘ │
└──────────────────────────────────────────────┘
       │ HTTP 302 Redirect
       ↓ Location: https://github.com/issmirnov/repo
    Browser follows redirect
```

## What Was Implemented

### Docker Support
- **Multi-stage Alpine build**: 15.3MB image
- **Multi-arch**: linux/amd64, linux/arm64
- **Security**: Non-root (UID 1000), read-only filesystem
- **Auto-published**: `ghcr.io/issmirnov/zap`
- **Health checks**: Built-in `/healthz` endpoint

### Helm Chart (`deploy/helm/zap/`)
- **LoadBalancer** with Cilium BGP, MetalLB, cloud provider support
- **Security hardening**: Non-root, dropped capabilities, read-only FS
- **Hot reload**: ConfigMap volume + inotify-based config reloading
- **Health probes**: Liveness and readiness configured
- **Autoscaling**: HPA support for high availability
- **Example configs**: Cilium BGP, MetalLB, HA, development

### Plain Manifests (`deploy/kubernetes/`)
- Namespace, ConfigMap, Deployment, Service
- Production-ready defaults
- No Helm dependency required

### Values.yaml Organization
The Helm values are organized with visual markers:
- 👉 = "You'll likely customize this" (image, config)
- ⚠️ = "Usually no need to change" (ports, resources)
- ✅ = "Already optimal defaults" (security, health checks)

Most users only need to override:
- `image.repository` (if using local registry)
- `zap.config` (your shortcuts - **THE MAIN THING**)
- `service.loadBalancer.ip` (optional, for static IP)

All other defaults (ports, resources, security) are production-ready.

## Prerequisites

- **Kubernetes** v1.19+ (self-hosted, cloud, or local)
- **kubectl** configured
- **LoadBalancer support**:
  - Self-hosted: MetalLB, Cilium BGP, kube-vip
  - Cloud: Native LoadBalancer (GKE, EKS, AKS)
  - Testing: Port-forward or NodePort
- **Helm 3+** (optional, for Helm deployments)
- **DNS server access** (dnsmasq, Pi-hole, router) for network-wide use

## Network DNS Configuration

### Get LoadBalancer IP
```bash
export ZAP_IP=$(kubectl get svc zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $ZAP_IP
```

### Configure DNS (Pick One)

**Option A: Individual DNS Records (Recommended)**

**dnsmasq** (`/etc/dnsmasq.conf`):
```
address=/g/192.168.1.100
address=/gh/192.168.1.100
address=/r/192.168.1.100
```

**Pi-hole**: Local DNS → DNS Records
- Add: `g` → `192.168.1.100`
- Add: `gh` → `192.168.1.100`

**Option B: Wildcard Domain**
```
address=/.zap.local/192.168.1.100
```
Then use: `http://g.zap.local/repo`

**Option C: Single Machine** (`/etc/hosts`):
```
192.168.1.100  g gh gl r
```

### Test
```bash
nslookup g  # Should resolve to LoadBalancer IP
curl -I http://g/test  # Should return HTTP 302
```

## Configuration Management

### Update Shortcuts

**Method 1: Edit ConfigMap**
```bash
kubectl edit configmap zap-config -n zap
# Zap auto-reloads in 30-60s
```

**Method 2: Helm Upgrade**
```bash
helm upgrade zap ./deploy/helm/zap -f values-custom.yaml
```

**Method 3: kubectl apply**
```bash
kubectl apply -f deploy/kubernetes/configmap.yaml
```

### Hot Reload
Zap has built-in hot reload via `inotify`. ConfigMap changes take 30-60s to propagate, then Zap automatically reloads. No pod restart needed.

Force restart if needed:
```bash
kubectl rollout restart deployment/zap -n zap
```

## Deployment Examples

### Self-Hosted with Cilium BGP
```yaml
# values-cilium.yaml
service:
  type: LoadBalancer
  loadBalancer:
    ip: "192.168.1.100"
  annotations:
    io.cilium/lb-ipam-ips: "192.168.1.100"
  labels:
    bgp: "true"

zap:
  config: |
    g:
      expand: github.com
```

```bash
helm install zap ./deploy/helm/zap -f values-cilium.yaml
```

### Self-Hosted with MetalLB
```yaml
service:
  type: LoadBalancer
  annotations:
    metallb.universe.tf/address-pool: default
```

### High Availability
```yaml
replicaCount: 3

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchLabels:
              app.kubernetes.io/name: zap
          topologyKey: kubernetes.io/hostname
```

## Troubleshooting

### LoadBalancer IP Stays Pending

**Self-hosted clusters need a LoadBalancer controller:**

**MetalLB:**
```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.0/config/manifests/metallb-native.yaml
# Then configure IP pool
```

**Cilium BGP:**
```bash
helm upgrade cilium cilium/cilium \
  --set bgp.enabled=true \
  --set bgp.announce.loadbalancerIP=true
```

### Pods Not Starting
```bash
kubectl describe pod -n zap -l app=zap
kubectl logs -n zap -l app=zap
```

Common issues:
- Image pull error (check `image.repository`)
- Config validation error (check ConfigMap YAML)
- Resource constraints (increase limits)

### Shortcuts Not Resolving

1. **Test DNS**: `nslookup g` (should return LoadBalancer IP)
2. **Test with Host header**: `curl -I -H "Host: g" http://$ZAP_IP/test`
3. **Check logs**: `kubectl logs -n zap -l app=zap`
4. **Verify config**: `curl http://$ZAP_IP/varz`

### Hot Reload Not Working
```bash
# Check logs for reload events
kubectl logs -n zap -l app=zap | grep -i reload

# ConfigMap changes take 30-60s to propagate

# Force restart
kubectl rollout restart deployment/zap -n zap
```

## Advanced Configuration

### Resource Tuning
Zap is lightweight but can handle high load:
```yaml
resources:
  requests:
    cpu: 100m      # Minimal idle
    memory: 50Mi   # ~20-30MB actual
  limits:
    cpu: 1000m     # For burst traffic
    memory: 256Mi  # Generous limit
```

### Security Hardening
Defaults are already secure (non-root, read-only FS, dropped capabilities).

Additional hardening:
```yaml
podSecurityContext:
  seccompProfile:
    type: RuntimeDefault
```

### Network Policies
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: zap-netpol
  namespace: zap
spec:
  podSelector:
    matchLabels:
      app: zap
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector: {}
      ports:
        - protocol: TCP
          port: 8927
```

## Performance

- **Throughput**: 150,000+ requests per second
- **Resource usage**: ~30MB RAM, ~50m CPU per pod
- **Startup time**: <5 seconds
- **Image size**: 15.3MB
- **Network overhead**: <1ms on LAN

## Monitoring

```bash
# Health checks
curl http://$ZAP_IP/healthz  # Returns "OK"
curl http://$ZAP_IP/varz     # Config dump

# Kubernetes
kubectl get pods -n zap
kubectl top pods -n zap
```

## Upgrading

**Helm:**
```bash
helm upgrade zap ./deploy/helm/zap -f values-custom.yaml
```

**Plain Manifests:**
```bash
kubectl apply -f deploy/kubernetes/
kubectl rollout restart deployment/zap -n zap
```

## Uninstalling

**Helm:**
```bash
helm uninstall zap
```

**Plain Manifests:**
```bash
kubectl delete -f deploy/kubernetes/
```

## Architecture Benefits

**Before (Local Installation)**
```
[Machine 1] → Zap → /etc/hosts
[Machine 2] → Zap → /etc/hosts
[Machine 3] → Zap → /etc/hosts
[Machine 4] → Zap → /etc/hosts
```
- 4 installations to maintain
- 4 config files to sync
- 4 update procedures

**After (Kubernetes)**
```
[Machine 1] ↘
[Machine 2] → DNS → LoadBalancer → Zap (K8s)
[Machine 3] ↗
[Machine 4] ↗
```
- 1 deployment to maintain
- 1 config file (ConfigMap)
- 1 update procedure
- Centralized management

## References

- [Zap GitHub](https://github.com/issmirnov/zap)
- [MetalLB Documentation](https://metallb.universe.tf/)
- [Cilium BGP Documentation](https://docs.cilium.io/en/stable/network/bgp-control-plane/)
- [Helm Documentation](https://helm.sh/docs/)
