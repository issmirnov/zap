# Kubernetes Deployment for Zap

This directory contains resources for deploying Zap to Kubernetes.

## Contents

- **helm/zap/** - Helm chart for Zap deployment
- This README - Deployment guide and configuration examples

## Quick Start

### Prerequisites

- Kubernetes cluster (self-hosted or cloud)
- `kubectl` configured to access your cluster
- `helm` v3+ installed
- (Optional) MetalLB, Cilium, or cloud LoadBalancer for external access

### Basic Installation

```bash
# Install with default values
helm install zap ./helm/zap

# Or install in a specific namespace
helm install zap ./helm/zap --namespace zap --create-namespace
```

### Custom Configuration

Create a `values-custom.yaml` file:

```yaml
service:
  type: LoadBalancer
  loadBalancer:
    ip: "192.168.1.100"  # Your desired IP
  annotations:
    io.cilium/lb-ipam-ips: "192.168.1.100"

zap:
  config: |
    g:
      expand: github.com
      z:
        expand: issmirnov/zap
    gh:
      expand: github.com
```

Install with custom values:

```bash
helm install zap ./helm/zap -f values-custom.yaml
```

## Configuration Examples

### 1. Self-Hosted Kubernetes with Cilium BGP

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

resources:
  requests:
    cpu: 100m
    memory: 50Mi
  limits:
    cpu: 500m
    memory: 128Mi
```

### 2. Self-Hosted Kubernetes with MetalLB

```yaml
# values-metallb.yaml
service:
  type: LoadBalancer
  loadBalancer:
    ip: "192.168.1.100"  # Optional - MetalLB can auto-assign
  annotations:
    metallb.universe.tf/address-pool: default
    metallb.universe.tf/allow-shared-ip: zap

zap:
  config: |
    # Your shortcuts here
    g:
      expand: github.com
```

### 3. Cloud Kubernetes (GKE, EKS, AKS)

```yaml
# values-cloud.yaml
service:
  type: LoadBalancer
  # Cloud providers auto-assign external IPs
  annotations:
    # AWS example:
    # service.beta.kubernetes.io/aws-load-balancer-type: nlb
    # GCP example:
    # cloud.google.com/load-balancer-type: "Internal"

resources:
  requests:
    cpu: 100m
    memory: 50Mi
  limits:
    cpu: 1000m
    memory: 256Mi
```

### 4. ClusterIP with Ingress

```yaml
# values-ingress.yaml
service:
  type: ClusterIP

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: zap.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: zap-tls
      hosts:
        - zap.example.com
```

### 5. With Stakater Reloader (Auto-restart on ConfigMap changes)

```yaml
# values-reloader.yaml
podAnnotations:
  reloader.stakater.com/auto: "true"

hotReload:
  enabled: true
```

## DNS Configuration

After deploying Zap, you need to configure DNS to point shortcuts to the LoadBalancer IP.

### Get LoadBalancer IP

```bash
kubectl get svc zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

### DNS Options

#### Option 1: Wildcard DNS (Recommended)

Configure your DNS server (dnsmasq, CoreDNS, bind, etc.) with a wildcard:

```
# dnsmasq example
address=/g/192.168.1.100
address=/gh/192.168.1.100
address=/r/192.168.1.100
```

Or use a wildcard domain:

```
*.zap.local -> 192.168.1.100
```

Then access shortcuts as: `http://g.zap.local/repo`

#### Option 2: Individual DNS Records

Create A records for each shortcut:

```
g.lan      A    192.168.1.100
gh.lan     A    192.168.1.100
r.lan      A    192.168.1.100
```

#### Option 3: Local /etc/hosts (Single Machine)

```bash
# Add to /etc/hosts
192.168.1.100  g gh r l
```

### Router/Network-Wide DNS

#### Using dnsmasq

Add to `/etc/dnsmasq.conf`:

```
address=/g/192.168.1.100
address=/gh/192.168.1.100
address=/r/192.168.1.100
```

Restart dnsmasq:

```bash
sudo systemctl restart dnsmasq
```

#### Using Pi-hole

1. Navigate to Local DNS → DNS Records
2. Add A record: `g` -> `192.168.1.100`
3. Repeat for each shortcut

## Updating Configuration

### Method 1: Edit ConfigMap directly

```bash
kubectl edit configmap zap-config
```

**Note**: Changes may take up to 60 seconds to propagate to pods due to Kubernetes' ConfigMap volume update mechanism. Zap's built-in hot reload will pick up changes automatically via inotify.

### Method 2: Update via Helm

Edit your `values-custom.yaml`:

```yaml
zap:
  config: |
    # Updated shortcuts
    g:
      expand: github.com
      z:
        expand: issmirnov/zap
      d:
        expand: issmirnov/dotfiles
```

Apply changes:

```bash
helm upgrade zap ./helm/zap -f values-custom.yaml
```

### Method 3: Force pod restart

If hot reload doesn't work or you want immediate changes:

```bash
kubectl rollout restart deployment/zap
```

## Monitoring and Debugging

### Check pod status

```bash
kubectl get pods -l app.kubernetes.io/name=zap
```

### View logs

```bash
kubectl logs -l app.kubernetes.io/name=zap -f
```

### Test health endpoint

```bash
# Get LoadBalancer IP
export ZAP_IP=$(kubectl get svc zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Health check
curl http://$ZAP_IP/healthz

# View configuration
curl http://$ZAP_IP/varz
```

### Test a shortcut

```bash
# Using Host header to simulate DNS
curl -I -H "Host: g" http://$ZAP_IP/z

# Should return HTTP 302 redirect to github.com/issmirnov/zap
```

### Debug DNS resolution

```bash
# From a client machine
nslookup g
ping g
curl -I http://g/test
```

## Advanced Configuration

### High Availability

For high-traffic environments:

```yaml
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

### Resource Tuning

Based on load testing (Zap can handle 150k+ QPS):

```yaml
resources:
  requests:
    cpu: 200m
    memory: 100Mi
  limits:
    cpu: 2000m
    memory: 512Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

### Security Hardening

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
  seccompProfile:
    type: RuntimeDefault

securityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL

networkPolicy:
  enabled: true
  ingress:
    - from:
      - namespaceSelector: {}
```

## Troubleshooting

### Pods not starting

```bash
kubectl describe pod -l app.kubernetes.io/name=zap
kubectl logs -l app.kubernetes.io/name=zap
```

Common issues:
- **Config validation error**: Check ConfigMap for valid YAML
- **Image pull error**: Ensure GHCR image is accessible
- **Resource constraints**: Increase limits or check node capacity

### LoadBalancer IP pending

```bash
kubectl get svc zap
```

If external IP is `<pending>`:
- **Cloud**: Wait a few minutes for provisioning
- **Self-hosted**: Install MetalLB or configure Cilium BGP
  - MetalLB: https://metallb.universe.tf/
  - Cilium BGP: https://docs.cilium.io/en/stable/network/bgp-control-plane/

### Hot reload not working

1. Check if inotify is supported in container:
   ```bash
   kubectl exec -it deployment/zap -- ls /proc/sys/fs/inotify
   ```

2. Check logs for reload messages:
   ```bash
   kubectl logs -l app.kubernetes.io/name=zap -f
   ```

3. ConfigMap updates may take 30-60 seconds to propagate

4. Consider installing Stakater Reloader for guaranteed restarts:
   ```bash
   helm repo add stakater https://stakater.github.io/stakater-charts
   helm install reloader stakater/reloader
   ```

### Shortcuts not resolving

1. Verify DNS configuration:
   ```bash
   nslookup g
   ```

2. Check LoadBalancer IP is correct:
   ```bash
   kubectl get svc zap -o wide
   ```

3. Test with Host header:
   ```bash
   curl -I -H "Host: g" http://<LOADBALANCER_IP>/test
   ```

4. Check Zap logs for requests:
   ```bash
   kubectl logs -l app.kubernetes.io/name=zap -f
   ```

## Uninstallation

```bash
helm uninstall zap
```

To completely remove all resources:

```bash
helm uninstall zap
kubectl delete namespace zap  # If you used a dedicated namespace
```

## Next Steps

- Configure your network DNS to use Zap
- Add your custom shortcuts to `values.yaml`
- Set up monitoring with Prometheus/Grafana
- Consider setting up automatic backups of your ConfigMap

## References

- [Zap GitHub Repository](https://github.com/issmirnov/zap)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Services](https://kubernetes.io/docs/concepts/services-networking/service/)
- [MetalLB](https://metallb.universe.tf/)
- [Cilium BGP](https://docs.cilium.io/en/stable/network/bgp-control-plane/)
