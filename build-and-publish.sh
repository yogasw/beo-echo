#!/bin/bash

# Beo Echo - Simple Docker Build and Publish Script
# =================================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Print colored output
print_color() {
    echo -e "${1}${2}${NC}"
}

# Check if GitHub CLI is installed and authenticated
check_github_auth() {
    print_color $YELLOW "🔍 Checking GitHub authentication..."
    
    if ! command -v gh &> /dev/null; then
        print_color $RED "❌ GitHub CLI (gh) not installed"
        echo "Install with: brew install gh"
        exit 1
    fi
    
    if ! gh auth status &> /dev/null; then
        print_color $RED "❌ Not authenticated with GitHub"
        echo "Run: gh auth login"
        exit 1
    fi
    
    print_color $YELLOW "🔐 Checking GitHub token permissions..."
    
    # Get detailed auth status
    local auth_output=$(gh auth status 2>&1)
    local scopes_line=$(echo "$auth_output" | grep "Token scopes")
    local account=$(echo "$auth_output" | grep "Logged in to" | awk '{print $4}' 2>/dev/null || echo "unknown")
    
    echo "Account: $account"
    echo "Scopes info: $scopes_line"
    
    # Check required scopes for Container Registry
    local has_write_packages=false
    local has_repo=false
    
    if echo "$scopes_line" | grep -q "write:packages"; then
        has_write_packages=true
    fi
    
    if echo "$scopes_line" | grep -q "repo"; then
        has_repo=true
    fi
    
    # Check if we can access the current repository
    if ! gh repo view &> /dev/null; then
        print_color $RED "❌ Cannot access current repository"
        echo "Make sure you have access to this repository"
        exit 1
    fi
    
    # Verify permissions for Container Registry
    if [[ "$has_write_packages" == "true" ]] || [[ "$has_repo" == "true" ]]; then
        print_color $GREEN "✅ Token has required permissions for Container Registry"
    else
        print_color $RED "❌ Token lacks required permissions for Container Registry"
        echo ""
        echo "Required permissions:"
        echo "  - write:packages (to push container images)"
        echo "  - repo (includes write:packages and more)"
        echo ""
        echo "Current token info: $scopes_line"
        echo ""
        echo "To fix this, run one of:"
        echo "  gh auth refresh -s write:packages"
        echo "  gh auth refresh -s repo"
        echo "  gh auth login (and select appropriate scopes)"
        echo ""
        exit 1
    fi
    
    # Test Docker registry login capability
    print_color $YELLOW "🐳 Testing Docker registry access..."
    local token=$(gh auth token)
    if echo "$token" | docker login ghcr.io -u "$account" --password-stdin &> /dev/null; then
        print_color $GREEN "✅ Docker registry login successful"
        docker logout ghcr.io &> /dev/null
    else
        print_color $RED "❌ Docker registry login failed"
        echo "This might be due to:"
        echo "  1. Insufficient token permissions"
        echo "  2. Docker daemon not running"
        echo "  3. Network connectivity issues"
        echo ""
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    print_color $GREEN "✅ GitHub authentication and permissions verified"
}

# Check if Docker is running
check_docker() {
    print_color $YELLOW "🐳 Checking Docker status..."
    
    if ! command -v docker &> /dev/null; then
        print_color $RED "❌ Docker not installed"
        echo "Install Docker from: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        print_color $RED "❌ Docker daemon is not running"
        echo "Please start Docker Desktop or Docker daemon"
        echo ""
        echo "On macOS: Open Docker Desktop application"
        echo "On Linux: sudo systemctl start docker"
        echo ""
        exit 1
    fi
    
    # Check Docker version
    local docker_version=$(docker --version | cut -d ' ' -f3 | tr -d ',')
    echo "Docker version: $docker_version"
    
    # Check if buildx is available for multi-platform builds
    if docker buildx version &> /dev/null; then
        echo "Docker Buildx: Available"
    else
        echo "Docker Buildx: Not available (single-platform builds only)"
    fi
    
    print_color $GREEN "✅ Docker is running and ready"
}

# Get repository info
get_repo_info() {
    OWNER=$(gh repo view --json owner --jq '.owner.login' 2>/dev/null || echo "unknown")
    REPO_NAME=$(gh repo view --json name --jq '.name' 2>/dev/null || basename $(pwd))
    
    if [[ -f "VERSION" ]]; then
        VERSION=$(cat VERSION | tr -d '\n' | tr -d ' ')
    else
        VERSION="latest"
    fi
    
    # Get git branch and commit for nightly tag
    local branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    local timestamp=$(date +%s)
    
    # Generate nightly tags: nightly-branch-commit_id-unixtime and nightly-latest
    NIGHTLY_TAG="nightly-${branch}-${commit}-${timestamp}"
    NIGHTLY_LATEST_TAG="nightly-latest"
    
    IMAGE_NAME="ghcr.io/${OWNER}/${REPO_NAME}"
    
    print_color $BLUE "📋 Build Info:"
    echo "   Repository: ${OWNER}/${REPO_NAME}"
    echo "   Image: ${IMAGE_NAME}"
    echo "   Version: ${VERSION}"
    echo "   Branch: ${branch}"
    echo "   Commit: ${commit}"
    echo "   Nightly: ${NIGHTLY_TAG}"
    echo "   Nightly Latest: ${NIGHTLY_LATEST_TAG}"
}

# Build Docker image
build_image() {
    print_color $YELLOW "🔨 Building Docker image..."
    
    docker build \
        --tag "${IMAGE_NAME}:${NIGHTLY_TAG}" \
        --tag "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" \
        --label "org.opencontainers.image.source=https://github.com/${OWNER}/${REPO_NAME}" \
        .
    
    print_color $GREEN "✅ Docker image built successfully"
    echo "   Tagged: ${IMAGE_NAME}:${NIGHTLY_TAG}"
    echo "   Tagged: ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
}

# Login to GitHub Container Registry
login_registry() {
    print_color $YELLOW "🔐 Logging into GitHub Container Registry..."
    
    gh auth token | docker login ghcr.io -u ${OWNER} --password-stdin
    
    print_color $GREEN "✅ Logged in to ghcr.io"
}

# Push image to registry
push_image() {
    print_color $YELLOW "📤 Pushing to GitHub Container Registry..."
    
    docker push "${IMAGE_NAME}:${NIGHTLY_TAG}"
    docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
    
    print_color $GREEN "✅ Images pushed successfully"
    echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
    echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
}

# Check repository permissions
check_repo_permissions() {
    print_color $YELLOW "📋 Checking repository permissions..."
    
    # Check if we can read repository info
    if ! gh repo view &> /dev/null; then
        print_color $RED "❌ Cannot access repository"
        echo "Make sure you have access to this repository"
        exit 1
    fi
    
    # Get repository visibility
    local repo_visibility=$(gh repo view --json visibility --jq '.visibility' 2>/dev/null || echo "unknown")
    echo "Repository visibility: $repo_visibility"
    
    # Check if we can create packages for this repository
    local repo_full_name=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
    echo "Repository: $repo_full_name"
    
    # For private repos, ensure we have appropriate permissions
    if [[ "$repo_visibility" == "PRIVATE" ]]; then
        print_color $YELLOW "⚠️  Private repository detected"
        echo "Ensure your token has access to this private repository"
    fi
    
    print_color $GREEN "✅ Repository permissions verified"
}

# Show build summary
show_build_summary() {
    print_color $CYAN "📋 Build Summary"
    echo "================"
    echo "Repository: $OWNER/$REPO_NAME"
    echo "Version: $VERSION"
    echo "Nightly: $NIGHTLY_TAG"
    echo "Registry: ghcr.io"
    echo "Images to be built and pushed:"
    echo "  - ghcr.io/$OWNER/$REPO_NAME:$NIGHTLY_TAG"
    echo "  - ghcr.io/$OWNER/$REPO_NAME:$NIGHTLY_LATEST_TAG"
    echo ""
    echo "This will:"
    echo "  1. Build Docker image locally"
    echo "  2. Tag with $NIGHTLY_TAG and $NIGHTLY_LATEST_TAG"
    echo "  3. Push to GitHub Container Registry"
    echo ""
    read -p "Continue? (Y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        print_color $YELLOW "Build cancelled by user"
        exit 0
    fi
}

# Main execution
main() {
    print_color $BLUE "🚀 Beo Echo - Docker Build & Publish"
    echo "=================================="
    
    check_github_auth
    check_docker
    check_repo_permissions
    get_repo_info
    show_build_summary
    build_image
    login_registry
    push_image
    
    print_color $GREEN "🎉 Build and publish completed successfully!"
}

# Run main function
main
