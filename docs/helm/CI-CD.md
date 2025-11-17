# CI/CD Testing for Helm Charts

## Status: ALL TESTS PASSING ✅

**PR:** https://github.com/issmirnov/zap/pull/53
**Branch:** feature/k8s-deployment

### Current Test Results
```
✅ Lint Helm Chart                  - 5s   (pass)
✅ Test Helm Installation in kind   - 2m   (pass)
✅ Test Example Values Files (dev)  - 5s   (pass)
✅ Test Example Values Files (ha)   - 5s   (pass)
✅ Test Example Values Files (minimal) - 5s (pass)
✅ CI (Go tests)                    - 1m   (pass)
✅ lint                             - 27s  (pass)
✅ test                             - 1m7s (pass)
✅ goreleaser-check                 - 10s  (pass)
✅ Docker Build and Push            - (pass)
```

## Overview

Comprehensive GitHub Actions workflow (`.github/workflows/helm-test.yml`) that automatically tests the Helm chart using **kind** (Kubernetes in Docker).

### Why This Matters

**Before:**
- ❌ Manual testing only
- ❌ Chart errors discovered after merge
- ❌ No automated validation
- ❌ Risky deployments

**After:**
- ✅ Automatic testing on every PR
- ✅ Real Kubernetes deployment testing
- ✅ Functional validation (redirects work!)
- ✅ Catches issues before merge

## What Gets Tested

### 3 Parallel Jobs

**Job 1: Lint (~30s)**
- Validates Helm chart syntax
- Tests template rendering
- Verifies YAML validity

**Job 2: Install Test (~2m)**
- Builds Docker image
- Creates kind K8s cluster
- Installs Helm chart
- Tests pod health
- Tests `/healthz` endpoint
- Tests `/varz` config endpoint
- **Tests redirects work**: `gh` → `github.com`
- **Tests nested redirects**: `g/z` → `github.com/issmirnov/zap`
- Tests ConfigMap hot reload

**Job 3: Examples Test (~30s)**
- Validates `values-dev.yaml`
- Validates `values-ha.yaml`
- Validates `values-minimal.yaml`

## Workflow Triggers

```yaml
on:
  push:
    branches: [master]  # Only after merge
  pull_request:
    branches: [master]  # On PRs targeting master
```

**Runs when:**
- PR is opened/updated
- Changes to `deploy/helm/**`
- Changes to `Dockerfile`
- Push to master

**Does NOT run when:**
- Only Go code changes (separate ci.yml handles that)
- Only documentation changes

**Optimized:** Single workflow run per PR (no duplicates)

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

## Issues Found & Fixed

### Issue 1: Multi-Arch Broken
- **Found by:** Bot reviewer
- **Problem:** ARM64 images contained AMD64 binaries
- **Fix:** Use TARGETOS/TARGETARCH build args
- **Commit:** daa38f6

### Issue 2: Dockerfile ENTRYPOINT
- **Found by:** Live deployment testing
- **Problem:** Container failing with "executable file not found: -host"
- **Fix:** Split CMD into ENTRYPOINT+CMD pattern
- **Commit:** f2846a3

### Issue 3: CI Redirect Test Parsing
- **Found by:** CI test failures
- **Problem:** URL parsing broken by colon splitting (`cut -d:`)
- **Fix:** Use `sed 's/^REDIRECT_URL://'` instead
- **Commit:** 4badec3

### Issue 4: Duplicate Workflow Runs
- **Found by:** User observation
- **Problem:** Running twice on PRs (push + pull_request events)
- **Fix:** Properly scoped triggers
- **Commit:** 3597cd6

## Local Testing

Run the same tests locally:

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
kubectl wait --for=condition=available deployment/zap

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
- ✅ Multi-arch images work
- ✅ Helm chart syntax correct
- ✅ Kubernetes deployment works
- ✅ Pod starts healthy
- ✅ Endpoints respond correctly
- ✅ Redirects function properly
- ✅ Hot reload works
- ✅ Example configs valid

## Benefits Delivered

### 1. Confidence
Every PR automatically:
- ✅ Builds Docker image
- ✅ Deploys to real Kubernetes
- ✅ Verifies functionality works
- ✅ Tests redirects actually redirect
- ✅ Validates hot reload

### 2. Speed
- Total test time: ~5 minutes
- Parallel execution where possible
- Fast feedback on PRs

### 3. Coverage
Tests validate end-to-end deployment and functionality

### 4. Cost Efficiency
- GitHub Actions free tier
- Optimized to avoid duplicate runs
- Caches Docker layers

## Performance

- **Total time:** ~5 minutes
- **Parallel execution:** Yes (3 jobs)
- **Cost:** $0 (GitHub Actions free tier sufficient)

## Troubleshooting CI Failures

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

Optional improvements:
- Test multiple K8s versions (1.25, 1.26, 1.27)
- Test with different LoadBalancer implementations
- Performance/load testing
- Security scanning (trivy, kubesec)
- Test upgrades (install v1, upgrade to v2)
- Test rollback scenarios

## Lessons Learned

1. **Test early, test often** - Found multiple issues through CI
2. **Real deployments catch real bugs** - kind testing found ENTRYPOINT issue
3. **Multi-arch is tricky** - Easy to hardcode architecture
4. **URL parsing needs care** - Don't split on delimiters in data
5. **Optimize CI triggers** - Avoid unnecessary duplicate runs

## Result

Production-ready Kubernetes deployment with comprehensive automated testing! Every PR now automatically deploys to a real Kubernetes cluster and validates everything works before merge.
