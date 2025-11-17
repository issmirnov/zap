# Helm Chart CI/CD Testing - Summary

## What Was Added

Created `.github/workflows/helm-test.yml` - A comprehensive GitHub Actions workflow that automatically tests the Helm chart using **kind** (Kubernetes in Docker).

## Why This Matters

Before this workflow:
- ❌ Manual testing only
- ❌ Chart errors discovered after merge
- ❌ No automated validation
- ❌ Risky deployments

After this workflow:
- ✅ Automatic testing on every PR
- ✅ Real Kubernetes deployment testing
- ✅ Functional validation (redirects work!)
- ✅ Catches issues before merge

## What Gets Tested

### 3 Jobs Running in Parallel

**Job 1: helm-lint (~30s)**
- Helm chart syntax
- Template rendering
- YAML validation

**Job 2: helm-install-test (~4m)**
- Build Docker image
- Create kind cluster (lightweight K8s)
- Install Helm chart
- Verify pod health
- Test `/healthz` endpoint
- Test `/varz` config endpoint
- Test redirects: `gh` → `github.com`
- Test nested: `g/z` → `github.com/issmirnov/zap`
- Test hot reload (ConfigMap update)

**Job 3: helm-examples-test (~30s)**
- Validate `values-dev.yaml`
- Validate `values-ha.yaml`
- Validate `values-minimal.yaml`

## Example Test Output

```bash
✅ Docker image loaded into kind cluster
✅ kind cluster is ready
✅ Helm chart installed successfully
✅ Deployment is ready
✅ Pod is running
✅ Health check passed: OK
✅ Config endpoint working
✅ Redirect test passed: gh -> https://github.com/
✅ Nested redirect test passed: g/z -> https://github.com/issmirnov/zap
✅ Hot reload working - config updated
```

## Workflow Triggers

Runs automatically when:
- PR is opened/updated
- Changes to `deploy/helm/**`
- Changes to `Dockerfile`
- Push to master or feature branches

Does NOT run for Go code changes (separate ci.yml handles that).

## Key Features

### 1. Real Kubernetes Testing
Uses kind to create actual K8s cluster - not just linting!

### 2. Functional Validation
Doesn't just check if it installs - validates redirects actually work!

### 3. Fast
Completes in ~5 minutes total (jobs run in parallel)

### 4. Comprehensive
Tests health checks, config endpoints, redirects, and hot reload

### 5. Debuggable
Shows full logs, events, and status on failure

## Local Testing

You can run the same tests locally:

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

# Test
kubectl wait --for=condition=available deployment/zap
POD=$(kubectl get pods -l app.kubernetes.io/name=zap -o name)
kubectl port-forward $POD 8927:8927 &
curl http://localhost:8927/healthz

# Cleanup
kind delete cluster --name zap-test
```

## Benefits

1. **Catches Issues Early** - Before they reach production
2. **Confidence** - Know the chart works before merging
3. **Documentation** - Shows how to test the chart
4. **Regression Prevention** - Catches breaking changes
5. **Integration Testing** - Tests real deployment scenario

## Files Added

1. `.github/workflows/helm-test.yml` - Main workflow
2. `.github/workflows/README-HELM-TEST.md` - Detailed documentation

## Integration with Existing CI

```
PR Created
    │
    ├─► ci.yml (Go tests, linting)
    └─► helm-test.yml (K8s deployment tests)
         │
         └─► All pass? → Ready to merge
```

## What This Validates

- ✅ Docker image builds
- ✅ Helm chart syntax
- ✅ Kubernetes deployment
- ✅ Pod starts successfully
- ✅ Health checks work
- ✅ Application functions (redirects)
- ✅ ConfigMap hot reload
- ✅ Example values files work

## Performance

- **Total time:** ~5 minutes
- **Parallel execution:** Yes
- **Cost:** $0 (GitHub Actions free tier sufficient)

## Next Steps

Optional future enhancements:
- Test multiple K8s versions
- Test with different LoadBalancers
- Performance/load testing
- Security scanning
- Upgrade testing

---

**Bottom line:** Every PR now automatically deploys to a real Kubernetes cluster and validates everything works!
