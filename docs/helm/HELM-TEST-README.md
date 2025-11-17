# Helm Chart Testing Workflow

## Overview

The `helm-test.yml` workflow automatically tests the Zap Helm chart using **kind** (Kubernetes in Docker) on every PR and push. This ensures the chart is production-ready before merging.

## What Gets Tested

### Job 1: `helm-lint`
Basic chart validation:
- ✅ Helm chart syntax validation
- ✅ Template rendering without errors
- ✅ YAML syntax validation

**Runtime:** ~30 seconds

### Job 2: `helm-install-test`
Full deployment test in a real Kubernetes cluster:

1. **Build Docker image** - Builds Zap image from Dockerfile
2. **Create kind cluster** - Spins up lightweight K8s cluster
3. **Load image** - Loads built image into kind
4. **Install Helm chart** - Deploys with test configuration
5. **Wait for ready** - Ensures deployment becomes healthy
6. **Test health endpoint** - Verifies `/healthz` returns OK
7. **Test config endpoint** - Verifies `/varz` returns config
8. **Test redirects** - Verifies actual URL expansion works
9. **Test hot reload** - Updates ConfigMap and checks reload

**Runtime:** ~3-4 minutes

### Job 3: `helm-examples-test`
Validates example values files:
- ✅ `values-dev.yaml` - Development config
- ✅ `values-ha.yaml` - High availability config
- ✅ `values-minimal.yaml` - Minimal config

**Runtime:** ~30 seconds

## Workflow Triggers

```yaml
on:
  push:
    branches: [master, feature/*]
    paths:
      - 'deploy/helm/**'
      - 'Dockerfile'
      - '.github/workflows/helm-test.yml'
  pull_request:
    paths:
      - 'deploy/helm/**'
      - 'Dockerfile'
      - '.github/workflows/helm-test.yml'
```

**Runs when:**
- Helm chart files change
- Dockerfile changes (affects deployment)
- Workflow file itself changes
- PR is opened/updated
- Push to master or feature branches

**Does NOT run when:**
- Only Go code changes (covered by ci.yml)
- Only documentation changes

## What Gets Tested

### Deployment Tests
```bash
✅ Docker image builds successfully
✅ Image loads into kind cluster
✅ Helm chart installs without errors
✅ Pods reach Running state
✅ Health checks pass
```

### Functionality Tests
```bash
✅ /healthz endpoint returns "OK"
✅ /varz endpoint returns valid JSON config
✅ Simple redirect: gh -> https://github.com/
✅ Nested redirect: g/z -> https://github.com/issmirnov/zap
✅ ConfigMap hot reload updates config
```

### Error Detection
```bash
✅ Pod logs checked for errors
✅ Kubernetes events reviewed
✅ Deployment status verified
```

## Example Output

### Successful Run
```
✅ Helm lint passed
✅ Helm template rendered successfully
✅ YAML validation passed
✅ Docker image loaded into kind cluster
✅ kind cluster is ready
✅ Helm chart installed successfully
✅ Deployment is ready
✅ Pod is running
✅ Health check passed: OK
✅ Config endpoint working
✅ Redirect test passed: gh -> https://github.com//
✅ Nested redirect test passed: g/z -> https://github.com/issmirnov/zap
✅ No errors in logs
```

### Failed Run Example
```
❌ Health check failed: connection refused
❌ Redirect test failed
Expected: https://github.com/
Got: 404 Not Found

[Pod logs shown]
[Events shown for debugging]
```

## Local Testing

You can run the same tests locally:

### 1. Install Prerequisites
```bash
# Install kind
brew install kind  # macOS
# OR
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Install kubectl (if not already installed)
brew install kubectl  # macOS

# Install helm
brew install helm  # macOS
```

### 2. Run Tests Manually
```bash
# Create kind cluster
kind create cluster --name zap-test

# Build and load image
docker build -t zap:test .
kind load docker-image zap:test --name zap-test

# Install chart
helm install zap deploy/helm/zap \
  --set image.repository=zap \
  --set image.tag=test \
  --set image.pullPolicy=Never \
  --set service.type=ClusterIP

# Wait for ready
kubectl wait --for=condition=available --timeout=60s deployment/zap

# Test health
POD=$(kubectl get pods -l app.kubernetes.io/name=zap -o jsonpath='{.items[0].metadata.name}')
kubectl port-forward $POD 8927:8927 &
curl http://localhost:8927/healthz  # Should return: OK

# Test redirect
curl -H "Host: gh" http://localhost:8927/  # Should redirect to github.com

# Cleanup
helm uninstall zap
kind delete cluster --name zap-test
```

## CI/CD Integration

This workflow integrates with the existing CI pipeline:

```
┌─────────────────────────────────────────────┐
│           Pull Request Created              │
└──────────────────┬──────────────────────────┘
                   │
      ┌────────────┴────────────┐
      │                         │
      ▼                         ▼
┌──────────────┐         ┌──────────────┐
│   ci.yml     │         │ helm-test.yml│
│              │         │              │
│ • Go tests   │         │ • Lint chart │
│ • Linting    │         │ • Deploy kind│
│ • E2E tests  │         │ • Test deploy│
│ • Coverage   │         │ • Test funcs │
└──────────────┘         └──────────────┘
      │                         │
      └────────────┬────────────┘
                   │
                   ▼
           ┌───────────────┐
           │   All Pass?   │
           └───────┬───────┘
                   │
                   ▼
           ┌───────────────┐
           │ Ready to Merge│
           └───────────────┘
```

## Benefits

### 1. Early Detection
Catches issues before they reach production:
- Docker build failures
- Helm chart syntax errors
- Missing dependencies
- Invalid configurations
- Runtime errors

### 2. Confidence
Validates that:
- Chart installs successfully
- Pods start and become healthy
- Application functions correctly
- Redirects work as expected
- Hot reload works

### 3. Documentation
Serves as executable documentation:
- Shows how to deploy the chart
- Demonstrates testing methodology
- Provides example commands
- Documents expected behavior

### 4. Regression Prevention
Automatically catches:
- Breaking changes to chart
- Docker image issues
- Configuration problems
- API changes

## Performance

Total workflow runtime: **~5 minutes**

Breakdown:
- `helm-lint`: 30s
- `helm-install-test`: 3-4m (includes building Docker image)
- `helm-examples-test`: 30s

All jobs run in parallel where possible.

## Troubleshooting

### Workflow Fails on Image Build
**Cause:** Dockerfile syntax error or missing dependencies
**Fix:** Test Docker build locally: `docker build -t zap:test .`

### Workflow Fails on Helm Install
**Cause:** Chart syntax error or invalid values
**Fix:** Test locally: `helm lint deploy/helm/zap`

### Workflow Fails on Pod Startup
**Cause:** Image issues or missing configuration
**Fix:** Check pod logs in workflow output, test deployment locally

### Tests Pass Locally but Fail in CI
**Cause:** Environment differences (networking, permissions)
**Fix:** Check workflow logs for specific error messages

## Future Enhancements

Potential improvements:
- [ ] Test multiple Kubernetes versions (1.25, 1.26, 1.27)
- [ ] Test with different LoadBalancer implementations
- [ ] Performance testing (load test redirects)
- [ ] Security scanning (trivy, kubesec)
- [ ] Test upgrades (install v1, upgrade to v2)
- [ ] Test rollback scenarios
- [ ] Matrix testing with different values combinations

## Related Workflows

- `.github/workflows/ci.yml` - Go code testing
- `.github/workflows/docker.yml` - Docker image publishing
- `.github/workflows/release.yml` - Release automation

## Resources

- [kind documentation](https://kind.sigs.k8s.io/)
- [Helm testing best practices](https://helm.sh/docs/topics/chart_tests/)
- [GitHub Actions for Kubernetes](https://github.com/marketplace?type=actions&query=kubernetes)
