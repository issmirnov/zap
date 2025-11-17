#!/bin/bash
# Zap Kubernetes Quickstart Script
# This script helps you quickly deploy Zap to Kubernetes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "\n${GREEN}===${NC} $1 ${GREEN}===${NC}\n"
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"

    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl not found. Please install kubectl first."
        exit 1
    fi
    print_info "✓ kubectl found"

    # Check kubectl access
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster. Please configure kubectl."
        exit 1
    fi
    print_info "✓ kubectl configured and cluster accessible"

    # Check for helm (optional)
    if command -v helm &> /dev/null; then
        HELM_AVAILABLE=true
        print_info "✓ helm found (optional)"
    else
        HELM_AVAILABLE=false
        print_warn "helm not found (optional - will use plain manifests)"
    fi
}

# Deploy using Helm
deploy_helm() {
    print_header "Deploying Zap with Helm"

    local release_name="${1:-zap}"
    local namespace="${2:-default}"

    if [ "$namespace" != "default" ]; then
        kubectl create namespace "$namespace" --dry-run=client -o yaml | kubectl apply -f -
    fi

    print_info "Installing Zap using Helm..."
    helm install "$release_name" ./helm/zap \
        --namespace "$namespace" \
        --wait

    print_info "✓ Zap deployed successfully"
}

# Deploy using plain manifests
deploy_manifests() {
    print_header "Deploying Zap with Kubernetes Manifests"

    print_info "Creating namespace..."
    kubectl apply -f kubernetes/namespace.yaml

    print_info "Creating ConfigMap..."
    kubectl apply -f kubernetes/configmap.yaml

    print_info "Creating Deployment..."
    kubectl apply -f kubernetes/deployment.yaml

    print_info "Creating Service..."
    kubectl apply -f kubernetes/service.yaml

    print_info "✓ Zap deployed successfully"
}

# Wait for LoadBalancer IP
wait_for_lb() {
    print_header "Waiting for LoadBalancer IP"

    local namespace="${1:-zap}"
    local service="${2:-zap}"
    local max_wait=300  # 5 minutes
    local elapsed=0

    print_info "Waiting for LoadBalancer IP assignment (timeout: ${max_wait}s)..."

    while [ $elapsed -lt $max_wait ]; do
        LB_IP=$(kubectl get svc "$service" -n "$namespace" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")

        if [ -n "$LB_IP" ]; then
            print_info "✓ LoadBalancer IP assigned: $LB_IP"
            return 0
        fi

        echo -n "."
        sleep 5
        elapsed=$((elapsed + 5))
    done

    echo ""
    print_warn "LoadBalancer IP not assigned after ${max_wait}s"
    print_warn "This is normal for:"
    print_warn "  - Self-hosted clusters without MetalLB/Cilium configured"
    print_warn "  - Cloud clusters during initial provisioning"
    return 1
}

# Show status
show_status() {
    local namespace="${1:-zap}"

    print_header "Zap Status"

    print_info "Pods:"
    kubectl get pods -n "$namespace" -l app=zap

    echo ""
    print_info "Service:"
    kubectl get svc -n "$namespace"

    LB_IP=$(kubectl get svc zap -n "$namespace" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")

    if [ -n "$LB_IP" ]; then
        echo ""
        print_info "Zap is available at: $LB_IP"
        echo ""
        print_info "Test health:"
        echo "  curl http://$LB_IP/healthz"
        echo ""
        print_info "View configuration:"
        echo "  curl http://$LB_IP/varz"
        echo ""
        print_info "Test shortcut (using Host header):"
        echo "  curl -I -H 'Host: g' http://$LB_IP/z"
        echo ""
        print_info "Next steps:"
        echo "  1. Configure DNS to point shortcuts to $LB_IP"
        echo "  2. See deploy/README.md for DNS configuration options"
    else
        echo ""
        print_warn "LoadBalancer IP not yet assigned"
        print_info "Check status with: kubectl get svc -n $namespace"
    fi
}

# Show next steps
show_next_steps() {
    local namespace="${1:-zap}"

    print_header "Next Steps"

    echo "1. Configure DNS:"
    echo "   - Option A: Wildcard DNS (e.g., *.zap.local -> LoadBalancer IP)"
    echo "   - Option B: Individual records (e.g., g.lan -> LoadBalancer IP)"
    echo "   - See deploy/README.md for detailed DNS setup"
    echo ""
    echo "2. Update shortcuts:"
    echo "   kubectl edit configmap zap-config -n $namespace"
    echo ""
    echo "3. View logs:"
    echo "   kubectl logs -n $namespace -l app=zap -f"
    echo ""
    echo "4. Restart deployment (if needed):"
    echo "   kubectl rollout restart deployment/zap -n $namespace"
    echo ""
    echo "For more information, see:"
    echo "  - deploy/README.md - Comprehensive deployment guide"
    echo "  - deploy/helm/zap/README.md - Helm chart documentation"
    echo "  - https://github.com/issmirnov/zap"
}

# Main menu
main() {
    echo "╔════════════════════════════════════════╗"
    echo "║   Zap Kubernetes Quickstart Script    ║"
    echo "╚════════════════════════════════════════╝"
    echo ""

    check_prerequisites

    echo ""
    echo "Select deployment method:"
    echo "  1) Helm (recommended)"
    echo "  2) Plain Kubernetes manifests"
    echo "  3) Show status only"
    echo "  4) Exit"
    echo ""
    read -p "Enter choice [1-4]: " choice

    case $choice in
        1)
            if [ "$HELM_AVAILABLE" = false ]; then
                print_error "Helm is not available. Please install Helm or choose option 2."
                exit 1
            fi

            read -p "Release name [zap]: " release_name
            release_name=${release_name:-zap}

            read -p "Namespace [default]: " namespace
            namespace=${namespace:-default}

            deploy_helm "$release_name" "$namespace"
            wait_for_lb "$namespace" "$release_name"
            show_status "$namespace"
            show_next_steps "$namespace"
            ;;
        2)
            deploy_manifests
            wait_for_lb "zap" "zap"
            show_status "zap"
            show_next_steps "zap"
            ;;
        3)
            read -p "Namespace [zap]: " namespace
            namespace=${namespace:-zap}
            show_status "$namespace"
            ;;
        4)
            print_info "Exiting..."
            exit 0
            ;;
        *)
            print_error "Invalid choice"
            exit 1
            ;;
    esac
}

# Run main menu
main
