# Helm Chart Documentation

Complete documentation for deploying Zap to Kubernetes using Helm.

## Documentation Files

### [QUICKSTART.md](QUICKSTART.md)
Quick reference guide for minimal configuration. Start here if you want to get Zap deployed fast with sensible defaults.

### [KUBERNETES.md](KUBERNETES.md)
Comprehensive Kubernetes deployment guide covering:
- Architecture overview
- Prerequisites and setup
- Network DNS configuration
- Deployment examples (Cilium BGP, MetalLB, HA)
- Configuration management
- Troubleshooting
- Advanced topics

### [CI-CD.md](CI-CD.md)
Complete CI/CD testing documentation:
- Automated Helm chart testing with kind
- 3-job parallel test workflow
- Issues found and fixed
- Local testing instructions
- Integration with existing CI

## Quick Links

- **Helm chart**: `/deploy/helm/zap/`
- **CI/CD workflow**: `/.github/workflows/helm-test.yml`
- **Example values**: `/deploy/helm/zap/examples/`
- **Plain manifests**: `/deploy/kubernetes/`
