# Plain Kubernetes Manifests for Zap

This directory contains simple Kubernetes manifests for deploying Zap without Helm.

## Quick Start

### 1. Create Namespace

```bash
kubectl apply -f namespace.yaml
```

### 2. Configure Shortcuts

Edit `configmap.yaml` and customize the shortcuts in the `c.yml` section.

### 3. Deploy Zap

```bash
# Apply all manifests
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

Or apply everything at once:

```bash
kubectl apply -f .
```

### 4. Get LoadBalancer IP

```bash
kubectl get svc zap -n zap

# Wait for EXTERNAL-IP to be assigned, then:
export ZAP_IP=$(kubectl get svc zap -n zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "Zap is available at: $ZAP_IP"
```

### 5. Configure DNS

Point your shortcuts to the LoadBalancer IP. See [parent README](../README.md) for DNS configuration options.

## Customization

### LoadBalancer Configuration

Edit `service.yaml` to configure your LoadBalancer:

#### For Cilium BGP

```yaml
metadata:
  annotations:
    io.cilium/lb-ipam-ips: "192.168.1.100"
  labels:
    bgp: "true"
spec:
  loadBalancerIP: 192.168.1.100
```

#### For MetalLB

```yaml
metadata:
  annotations:
    metallb.universe.tf/address-pool: default
    metallb.universe.tf/allow-shared-ip: zap
spec:
  loadBalancerIP: 192.168.1.100  # Optional
```

### Resource Limits

Edit `deployment.yaml` to adjust resource requests/limits:

```yaml
resources:
  requests:
    cpu: 100m
    memory: 50Mi
  limits:
    cpu: 500m
    memory: 128Mi
```

### Replicas

For high availability, increase replicas in `deployment.yaml`:

```yaml
spec:
  replicas: 3
```

## Updating Configuration

### Edit ConfigMap

```bash
kubectl edit configmap zap-config -n zap
```

Changes will be picked up automatically by Zap's hot reload (may take up to 60 seconds for Kubernetes to propagate ConfigMap changes).

### Force Restart

To apply changes immediately:

```bash
kubectl rollout restart deployment/zap -n zap
```

## Monitoring

### Check Status

```bash
kubectl get all -n zap
```

### View Logs

```bash
kubectl logs -n zap -l app=zap -f
```

### Test Health

```bash
export ZAP_IP=$(kubectl get svc zap -n zap -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl http://$ZAP_IP/healthz
curl http://$ZAP_IP/varz
```

## Cleanup

```bash
kubectl delete -f .
```

Or delete the entire namespace:

```bash
kubectl delete namespace zap
```

## Helm Alternative

For more flexibility and easier management, consider using the Helm chart instead:

```bash
helm install zap ../helm/zap
```

See [../README.md](../README.md) for Helm documentation.
