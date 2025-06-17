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
    
    # Check system architecture
    local arch=$(uname -m)
    echo "System architecture: $arch"
    
    # Check if buildx is available for multi-platform builds
    if docker buildx version &> /dev/null; then
        echo "Docker Buildx: Available"
        
        # For ARM systems, ensure buildx is properly configured
        if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
            print_color $CYAN "🔧 ARM system detected - configuring buildx for multi-platform builds"
            
            # Enable experimental features for manifest commands
            export DOCKER_CLI_EXPERIMENTAL=enabled
            
            # Create/use a builder instance that supports multi-platform
            if ! docker buildx inspect multiplatform &> /dev/null; then
                print_color $YELLOW "Creating multiplatform builder..."
                docker buildx create --name multiplatform --use --bootstrap
            else
                print_color $YELLOW "Using existing multiplatform builder..."
                docker buildx use multiplatform
            fi
        fi
    else
        echo "Docker Buildx: Not available"
        
        # For ARM systems, buildx is required for multi-platform builds
        if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
            print_color $RED "❌ Docker Buildx is required for ARM systems to build multi-platform images"
            echo "Please update Docker to a version that includes buildx"
            exit 1
        else
            echo "(Single-platform builds only)"
        fi
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

# Build Docker image per platform
build_image() {
    print_color $YELLOW "🔨 Building Docker images per platform..."
    
    # Detect system architecture
    local arch=$(uname -m)
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        print_color $CYAN "🏗️  ARM architecture detected ($arch)"
        print_color $YELLOW "   Building for ARM64 and AMD64 platforms separately..."
        
        # Build ARM64 platform
        print_color $YELLOW "📦 Building ARM64 platform..."
        docker buildx build \
            --platform linux/arm64 \
            --tag "${IMAGE_NAME}:${NIGHTLY_TAG}-arm64" \
            --tag "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-arm64" \
            --label "org.opencontainers.image.source=https://github.com/${OWNER}/${REPO_NAME}" \
            --load \
            .
        
        # Build AMD64 platform
        print_color $YELLOW "📦 Building AMD64 platform..."
        docker buildx build \
            --platform linux/amd64 \
            --tag "${IMAGE_NAME}:${NIGHTLY_TAG}-amd64" \
            --tag "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-amd64" \
            --label "org.opencontainers.image.source=https://github.com/${OWNER}/${REPO_NAME}" \
            --load \
            .
            
        print_color $GREEN "✅ Multi-platform images built successfully"
        echo "   ARM64: ${IMAGE_NAME}:${NIGHTLY_TAG}-arm64"
        echo "   AMD64: ${IMAGE_NAME}:${NIGHTLY_TAG}-amd64"
    else
        print_color $CYAN "🏗️  x86_64 architecture detected ($arch)"
        print_color $YELLOW "   Building for current platform only..."
        
        # Non-ARM system: use regular docker build
        docker build \
            --tag "${IMAGE_NAME}:${NIGHTLY_TAG}" \
            --tag "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" \
            --label "org.opencontainers.image.source=https://github.com/${OWNER}/${REPO_NAME}" \
            .
            
        print_color $GREEN "✅ Docker image built successfully"
        echo "   Tagged: ${IMAGE_NAME}:${NIGHTLY_TAG}"
        echo "   Tagged: ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
    fi
}

# Login to GitHub Container Registry
login_registry() {
    print_color $YELLOW "🔐 Logging into GitHub Container Registry..."
    
    gh auth token | docker login ghcr.io -u ${OWNER} --password-stdin
    
    print_color $GREEN "✅ Logged in to ghcr.io"
}

# Push images per platform
push_image() {
    # Detect system architecture
    local arch=$(uname -m)
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        print_color $YELLOW "📤 Pushing platform-specific images to GitHub Container Registry..."
        
        # Push ARM64 images
        print_color $YELLOW "🚀 Pushing ARM64 images..."
        docker push "${IMAGE_NAME}:${NIGHTLY_TAG}-arm64"
        docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-arm64"
        
        # Push AMD64 images
        print_color $YELLOW "🚀 Pushing AMD64 images..."
        docker push "${IMAGE_NAME}:${NIGHTLY_TAG}-amd64"
        docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-amd64"
        
        print_color $GREEN "✅ Platform-specific images pushed successfully"
        echo "   ARM64: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}-arm64"
        echo "   AMD64: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}-amd64"
    else
        print_color $YELLOW "📤 Pushing to GitHub Container Registry..."
        
        docker push "${IMAGE_NAME}:${NIGHTLY_TAG}"
        docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
        
        print_color $GREEN "✅ Images pushed successfully"
        echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
        echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
    fi
}

# Create multi-platform manifest
create_manifest() {
    # Detect system architecture
    local arch=$(uname -m)
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        print_color $YELLOW "� Creating multi-platform manifests..."
        
        # Create manifest for nightly tag
        print_color $YELLOW "🔗 Creating manifest for ${NIGHTLY_TAG}..."
        docker manifest create "${IMAGE_NAME}:${NIGHTLY_TAG}" \
            "${IMAGE_NAME}:${NIGHTLY_TAG}-arm64" \
            "${IMAGE_NAME}:${NIGHTLY_TAG}-amd64"
        
        # Annotate platform-specific images
        docker manifest annotate "${IMAGE_NAME}:${NIGHTLY_TAG}" \
            "${IMAGE_NAME}:${NIGHTLY_TAG}-arm64" --arch arm64 --os linux
        docker manifest annotate "${IMAGE_NAME}:${NIGHTLY_TAG}" \
            "${IMAGE_NAME}:${NIGHTLY_TAG}-amd64" --arch amd64 --os linux
        
        # Push manifest for nightly tag
        docker manifest push "${IMAGE_NAME}:${NIGHTLY_TAG}"
        
        # Create manifest for nightly-latest tag
        print_color $YELLOW "🔗 Creating manifest for ${NIGHTLY_LATEST_TAG}..."
        docker manifest create "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" \
            "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-arm64" \
            "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-amd64"
        
        # Annotate platform-specific images
        docker manifest annotate "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" \
            "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-arm64" --arch arm64 --os linux
        docker manifest annotate "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" \
            "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-amd64" --arch amd64 --os linux
        
        # Push manifest for nightly-latest tag
        docker manifest push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
        
        print_color $GREEN "✅ Multi-platform manifests created and pushed successfully"
        echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
        echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
        echo ""
        echo "Users can now pull:"
        echo "   docker pull ${IMAGE_NAME}:${NIGHTLY_TAG} (auto-selects platform)"
        echo "   docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}-arm64 (specific ARM64)"
        echo "   docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}-amd64 (specific AMD64)"
    else
        print_color $CYAN "ℹ️  Single-platform build - no manifest creation needed"
        echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
        echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
    fi
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
    
    # Detect system architecture for build info
    local arch=$(uname -m)
    echo "Build architecture: $arch"
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        echo "Build strategy: Multi-platform (ARM64 + AMD64) using buildx"
        echo "Platforms: linux/arm64, linux/amd64"
    else
        echo "Build strategy: Single-platform using standard docker build"
        echo "Platform: Current system architecture"
    fi
    
    echo ""
    echo "Images to be built and pushed:"
    echo "  - ghcr.io/$OWNER/$REPO_NAME:$NIGHTLY_TAG"
    echo "  - ghcr.io/$OWNER/$REPO_NAME:$NIGHTLY_LATEST_TAG"
    echo ""
    echo "This will:"
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        echo "  1. Build platform-specific Docker images (ARM64 and AMD64 separately)"
        echo "  2. Tag with platform suffixes (-arm64, -amd64)"
        echo "  3. Push each platform-specific image to GitHub Container Registry"
        echo "  4. Create multi-platform manifests"
        echo "  5. Push manifests to enable auto-platform selection"
        echo ""
        echo "Final images available:"
        echo "  - ${IMAGE_NAME}:${NIGHTLY_TAG} (multi-platform manifest)"
        echo "  - ${IMAGE_NAME}:${NIGHTLY_TAG}-arm64 (ARM64 specific)"
        echo "  - ${IMAGE_NAME}:${NIGHTLY_TAG}-amd64 (AMD64 specific)"
    else
        echo "  1. Build Docker image locally"
        echo "  2. Tag with $NIGHTLY_TAG and $NIGHTLY_LATEST_TAG"
        echo "  3. Push to GitHub Container Registry"
    fi
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
    create_manifest
    
    print_color $GREEN "🎉 Build and publish completed successfully!"
}

# Run main function
main
