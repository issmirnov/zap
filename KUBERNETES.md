# Kubernetes Deployment Guide for Zap

This document provides a comprehensive guide for deploying Zap to Kubernetes as a network-wide URL shortcut service.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Deployment Options](#deployment-options)
- [Network DNS Configuration](#network-dns-configuration)
- [Configuration Management](#configuration-management)
- [Advanced Topics](#advanced-topics)
- [Troubleshooting](#troubleshooting)

## Overview

Zap can be deployed to Kubernetes to provide network-wide URL shortcuts. Instead of running Zap on each machine, you:

1. Deploy Zap to Kubernetes with a LoadBalancer service
2. Configure your network DNS to point shortcut domains to the LoadBalancer IP
3. All machines on your network can use shortcuts without local installation

### Architecture

```
┌──────────────┐     DNS Query: g      ┌─────────────┐
│ Client       │ ───────────────────>  │  DNS Server │
│ Machine      │                       │  (dnsmasq)  │
└──────┬───────┘                       └──────┬──────┘
       │                                      │
       │ Responds: g = 192.168.1.100         │
       │<─────────────────────────────────────┘
       │
       │ HTTP Request: http://g/repo
       ↓
┌──────────────────────────────────────────────┐
│  Kubernetes LoadBalancer (192.168.1.100)     │
│  ┌─────────────────────────────────────────┐ │
│  │  Zap Service (type: LoadBalancer)       │ │
│  │  ┌─────────────────┐  ┌──────────────┐ │ │
│  │  │  Zap Pod 1      │  │  Zap Pod 2   │ │ │
│  │  │  (reads c.yml)  │  │  (c.yml)     │ │ │
│  │  └─────────────────┘  └──────────────┘ │ │
│  │           ↑                    ↑         │ │
│  │           └────────┬───────────┘         │ │
│  │                    │                     │ │
│  │           ┌────────▼─────────┐           │ │
│  │           │  ConfigMap       │           │ │
│  │           │  (c.yml config)  │           │ │
│  │           └──────────────────┘           │ │
│  └─────────────────────────────────────────┘ │
└──────────────────────────────────────────────┘
       │
       │ HTTP 302 Redirect
       │ Location: https://github.com/issmirnov/repo
       ↓
    Browser follows redirect
```

## Prerequisites

### Required

- **Kubernetes cluster** (v1.19+)
  - Self-hosted (with MetalLB, Cilium BGP, or similar for LoadBalancer)
  - Cloud-hosted (GKE, EKS, AKS)
  - Local (minikube, kind - for testing)

- **kubectl** configured to access your cluster

- **LoadBalancer support**:
  - Self-hosted: MetalLB, Cilium BGP, or similar
  - Cloud: Native cloud LoadBalancer
  - Testing: Port-forward or NodePort

### Optional

- **Helm 3+** (recommended for easier deployment)
- **DNS server access** (dnsmasq, Pi-hole, router admin, etc.)
- **Stakater Reloader** (for automatic pod restarts on config changes)

## Quick Start

### Using Quickstart Script

The fastest way to get started:

```bash
cd deploy
./quickstart.sh
```

The script will:
1. Check prerequisites
2. Guide you through deployment (Helm or plain manifests)
3. Wait for LoadBalancer IP
4. Show you the next steps

### Manual Quick Start (Helm)

```bash
# Install with default configuration
helm install zap ./deploy/helm/zap

# Wait for LoadBalancer IP
kubectl get svc zap --watch

# Get the IP
export ZAP_IP=$(kubectl get svc zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "Zap LoadBalancer IP: $ZAP_IP"

# Test health
curl http://$ZAP_IP/healthz
```

### Manual Quick Start (Plain Manifests)

```bash
# Apply all manifests
kubectl apply -f deploy/kubernetes/

# Wait for LoadBalancer IP
kubectl get svc zap -n zap --watch

# Get the IP
export ZAP_IP=$(kubectl get svc zap -n zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "Zap LoadBalancer IP: $ZAP_IP"
```

## Deployment Options

### Option 1: Helm Chart (Recommended)

**Best for:**
- Production deployments
- Multiple environments
- Easy configuration management
- Teams wanting standardized deployments

**Features:**
- Parameterized configuration
- Easy upgrades
- Built-in best practices
- Example values files for common scenarios

See [deploy/helm/zap/README.md](deploy/helm/zap/README.md) for detailed Helm documentation.

**Example with Cilium BGP:**

```bash
helm install zap ./deploy/helm/zap \
  --set service.loadBalancer.ip=192.168.1.100 \
  --set service.annotations."io\.cilium/lb-ipam-ips"=192.168.1.100 \
  --set service.labels.bgp=true
```

### Option 2: Plain Kubernetes Manifests

**Best for:**
- Simple deployments
- Learning/testing
- Environments without Helm
- Direct control over resources

**Files:**
- `deploy/kubernetes/namespace.yaml` - Namespace
- `deploy/kubernetes/configmap.yaml` - Configuration
- `deploy/kubernetes/deployment.yaml` - Zap pods
- `deploy/kubernetes/service.yaml` - LoadBalancer service

See [deploy/kubernetes/README.md](deploy/kubernetes/README.md) for details.

## Network DNS Configuration

This is the key to making Zap work network-wide. You need to configure DNS so that shortcut domains resolve to your LoadBalancer IP.

### Step 1: Get LoadBalancer IP

```bash
export ZAP_IP=$(kubectl get svc zap -n zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $ZAP_IP
```

### Step 2: Configure DNS

#### Option A: Individual DNS Records (Recommended)

Add A records for each shortcut:

**dnsmasq (`/etc/dnsmasq.conf`):**

```
address=/g/192.168.1.100
address=/gh/192.168.1.100
address=/gl/192.168.1.100
address=/r/192.168.1.100
```

**Pi-hole:**
1. Go to Local DNS → DNS Records
2. Add: `g` → `192.168.1.100`
3. Add: `gh` → `192.168.1.100`
4. Repeat for each shortcut

**Router (varies by model):**
- Look for "Static DNS" or "DNS Override"
- Add entries for each shortcut

#### Option B: Wildcard Domain

Use a wildcard subdomain:

**dnsmasq:**

```
address=/.zap.local/192.168.1.100
```

Then use shortcuts as: `http://g.zap.local/repo`

#### Option C: Individual Machine /etc/hosts

For testing or single-machine use:

```bash
# Add to /etc/hosts
192.168.1.100  g gh gl r l
```

### Step 3: Test DNS

```bash
# Test DNS resolution
nslookup g
ping g

# Should resolve to your LoadBalancer IP
```

### Step 4: Test Zap

```bash
# Using curl with Host header
curl -I -H "Host: g" http://$ZAP_IP/z

# Should return HTTP 302 redirect to github.com/issmirnov/zap

# Or if DNS is configured
curl -I http://g/z
```

## Configuration Management

### Editing Shortcuts

#### Method 1: Edit ConfigMap Directly

```bash
kubectl edit configmap zap-config -n zap
```

Zap will automatically reload the config (may take up to 60 seconds for Kubernetes to propagate changes).

#### Method 2: Update via Helm

Edit your values file, then:

```bash
helm upgrade zap ./deploy/helm/zap -f values-custom.yaml
```

#### Method 3: kubectl apply

Edit the ConfigMap file, then:

```bash
kubectl apply -f deploy/kubernetes/configmap.yaml
```

### Hot Reload Behavior

Zap has built-in hot reload via `inotify`. When the ConfigMap volume updates:
1. Kubernetes propagates changes (can take 30-60 seconds)
2. Zap's file watcher detects the change
3. Config is reloaded automatically
4. No pod restart needed

**Note:** If hot reload doesn't work, manually restart:

```bash
kubectl rollout restart deployment/zap -n zap
```

### Using Stakater Reloader (Optional)

For guaranteed restarts on ConfigMap changes:

```bash
# Install Reloader
helm repo add stakater https://stakater.github.io/stakater-charts
helm install reloader stakater/reloader

# Enable in Helm values
helm upgrade zap ./deploy/helm/zap \
  --set podAnnotations."reloader\.stakater\.com/auto"="true"
```

## Advanced Topics

### High Availability

For production, run multiple replicas:

```yaml
# Helm values
replicaCount: 3

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchExpressions:
              - key: app.kubernetes.io/name
                operator: In
                values:
                  - zap
          topologyKey: kubernetes.io/hostname
```

### Autoscaling

Enable horizontal pod autoscaler:

```yaml
# Helm values
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

### Resource Tuning

Zap is lightweight but can handle high load:

```yaml
resources:
  requests:
    cpu: 100m      # Minimal for idle
    memory: 50Mi   # ~20-30MB actual usage
  limits:
    cpu: 1000m     # For burst traffic
    memory: 256Mi  # Generous limit
```

### Security Hardening

The default configuration is already secure:
- Non-root user (UID 1000)
- Read-only root filesystem
- Dropped all capabilities
- No privilege escalation

Additional hardening:

```yaml
podSecurityContext:
  seccompProfile:
    type: RuntimeDefault
```

### Network Policies

Restrict traffic to Zap:

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

### Custom Image

Build and use your own image:

```bash
# Build
docker build -t myregistry/zap:custom .

# Push
docker push myregistry/zap:custom

# Deploy
helm install zap ./deploy/helm/zap \
  --set image.repository=myregistry/zap \
  --set image.tag=custom
```

## Troubleshooting

### LoadBalancer IP Stays Pending

**Symptom:** `EXTERNAL-IP` shows `<pending>`

**Causes:**
- Self-hosted cluster without LoadBalancer controller
- Cloud cluster during initial provisioning

**Solutions:**

For self-hosted clusters, install a LoadBalancer controller:

**MetalLB:**
```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.0/config/manifests/metallb-native.yaml
# Then configure IP pool
```

**Cilium BGP:**
```bash
# Enable BGP in Cilium configuration
helm upgrade cilium cilium/cilium \
  --set bgp.enabled=true \
  --set bgp.announce.loadbalancerIP=true
```

### Pods Not Starting

Check pod events:

```bash
kubectl describe pod -n zap -l app=zap
```

Common issues:
- **Image pull error:** Ensure ghcr.io/issmirnov/zap:latest is accessible
- **Config validation error:** Check ConfigMap for valid YAML
- **Resource constraints:** Increase limits or check node capacity

### Hot Reload Not Working

1. Check if file watcher is working:
   ```bash
   kubectl logs -n zap -l app=zap | grep -i reload
   ```

2. ConfigMap changes take time to propagate (up to 60s)

3. Force restart:
   ```bash
   kubectl rollout restart deployment/zap -n zap
   ```

4. Install Stakater Reloader for automatic restarts

### Shortcuts Not Resolving

1. **Test DNS:**
   ```bash
   nslookup g
   ```
   Should return LoadBalancer IP

2. **Test with Host header:**
   ```bash
   curl -I -H "Host: g" http://$ZAP_IP/test
   ```
   Should return HTTP 302 redirect

3. **Check Zap logs:**
   ```bash
   kubectl logs -n zap -l app=zap -f
   ```

4. **Verify config:**
   ```bash
   curl http://$ZAP_IP/varz
   ```

### DNS Not Resolving on Clients

1. Verify DNS server is configured on clients:
   ```bash
   cat /etc/resolv.conf
   ```

2. Check DNS server is receiving queries:
   ```bash
   # On DNS server (e.g., dnsmasq)
   tail -f /var/log/dnsmasq.log
   ```

3. Test from client:
   ```bash
   dig @<dns-server-ip> g
   ```

## Performance

Zap can handle **150,000+ requests per second** with the default configuration. For most networks, even a single replica is sufficient.

Expected resource usage:
- **CPU:** ~50-100m idle, bursts for traffic
- **Memory:** ~20-30MB actual usage
- **Network:** Minimal (redirects only)

## Monitoring

### Basic Health Checks

```bash
# Health endpoint
curl http://$ZAP_IP/healthz

# Configuration dump
curl http://$ZAP_IP/varz

# Kubernetes health
kubectl get pods -n zap
```

### Prometheus Integration (Future)

Zap doesn't currently expose Prometheus metrics, but this could be added. For now, use Kubernetes metrics:

```bash
kubectl top pods -n zap
```

## Upgrading

### Helm

```bash
helm upgrade zap ./deploy/helm/zap -f values-custom.yaml
```

### Plain Manifests

```bash
kubectl apply -f deploy/kubernetes/
kubectl rollout restart deployment/zap -n zap
```

## Uninstalling

### Helm

```bash
helm uninstall zap
```

### Plain Manifests

```bash
kubectl delete -f deploy/kubernetes/
```

## References

- [Zap GitHub](https://github.com/issmirnov/zap)
- [MetalLB Documentation](https://metallb.universe.tf/)
- [Cilium BGP Documentation](https://docs.cilium.io/en/stable/network/bgp-control-plane/)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Services](https://kubernetes.io/docs/concepts/services-networking/service/)

## Contributing

Found an issue or have a suggestion? Please open an issue on GitHub:
https://github.com/issmirnov/zap/issues
