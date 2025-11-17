# Zap Helm Chart Quick Start

## Minimal Installation

The Helm chart ships with production-ready defaults. For most users, you only need to customize the shortcuts configuration.

### 1. Install with default shortcuts

```bash
helm install zap ./deploy/helm/zap
```

This deploys Zap with example shortcuts (github, reddit, etc.) on a LoadBalancer.

### 2. Install with custom shortcuts only

```bash
helm install zap ./deploy/helm/zap --set-file zap.config=my-shortcuts.yml
```

Or create a minimal values file:

```yaml
# my-values.yaml
zap:
  config: |
    g:
      expand: github.com
    r:
      expand: reddit.com
```

```bash
helm install zap ./deploy/helm/zap -f my-values.yaml
```

### 3. What values typically need customization?

**Almost always customize:**
- `zap.config` - Your shortcuts (THIS IS THE MAIN THING)

**Sometimes customize:**
- `image.repository` - If using local registry
- `image.tag` - Pin to specific version
- `service.loadBalancer.ip` - Set specific IP (self-hosted k8s)
- `service.loadBalancer.sourceRanges` - Restrict access by IP

**Rarely customize:**
- `zap.port` (default: 8927) ✅ Already optimal
- `zap.host` (default: 0.0.0.0) ✅ Already optimal
- `service.port` (default: 80) ✅ Already optimal
- `service.targetPort` (default: 8927) ✅ Already optimal
- `resources.*` ✅ Already production-ready
- `healthCheck.*` ✅ Already configured
- Security contexts ✅ Already hardened

## Examples

See `examples/` directory:
- `values-minimal.yaml` - Absolute minimum override
- `values-dev.yaml` - Development setup
- `values-ha.yaml` - High availability setup
- `values-cilium-bgp.yaml` - Self-hosted with Cilium
- `values-metallb.yaml` - Self-hosted with MetalLB

## What's included by default?

✅ LoadBalancer service on port 80
✅ Zap listening on 0.0.0.0:8927 internally
✅ Health checks (liveness + readiness)
✅ Security contexts (non-root, read-only filesystem)
✅ Resource limits (500m CPU, 128Mi memory)
✅ Hot reload support (built-in via inotify)
✅ Production-ready example shortcuts

## Common Patterns

### Pattern 1: Just change shortcuts (99% of users)

```yaml
zap:
  config: |
    # Your shortcuts here
```

### Pattern 2: Local registry + shortcuts

```yaml
image:
  repository: registry.local/zap
  tag: "v1"

zap:
  config: |
    # Your shortcuts here
```

### Pattern 3: Specific LoadBalancer IP + shortcuts

```yaml
service:
  loadBalancer:
    ip: "192.168.1.100"

zap:
  config: |
    # Your shortcuts here
```

That's it! The defaults handle everything else.
