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
    local has_read_packages=false
    local has_delete_packages=false
    local has_repo=false
    
    if echo "$scopes_line" | grep -q "write:packages"; then
        has_write_packages=true
    fi
    
    if echo "$scopes_line" | grep -q "read:packages"; then
        has_read_packages=true
    fi
    
    if echo "$scopes_line" | grep -q "delete:packages"; then
        has_delete_packages=true
    fi
    
    if echo "$scopes_line" | grep -q "repo"; then
        has_repo=true
        # repo scope includes all package permissions
        has_write_packages=true
        has_read_packages=true
        has_delete_packages=true
    fi
    
    # Check if we can access the current repository
    if ! gh repo view &> /dev/null; then
        print_color $RED "❌ Cannot access current repository"
        echo "Make sure you have access to this repository"
        exit 1
    fi
    
    # Verify permissions for Container Registry
    local missing_scopes=()
    
    if [[ "$has_write_packages" != "true" ]]; then
        missing_scopes+=("write:packages")
    fi
    
    if [[ "$has_read_packages" != "true" ]]; then
        missing_scopes+=("read:packages")
    fi
    
    if [[ "$has_delete_packages" != "true" ]]; then
        missing_scopes+=("delete:packages")
    fi
    
    if [ ${#missing_scopes[@]} -eq 0 ]; then
        print_color $GREEN "✅ Token has all required permissions for Container Registry"
        echo "   ✓ write:packages (push container images)"
        echo "   ✓ read:packages (read container images)"
        echo "   ✓ delete:packages (cleanup temporary images)"
    else
        print_color $YELLOW "⚠️  Token is missing some recommended permissions"
        echo ""
        echo "Missing scopes: ${missing_scopes[*]}"
        echo "Current token info: $scopes_line"
        echo ""
        
        if [[ "$has_write_packages" != "true" ]]; then
            print_color $RED "❌ write:packages - Required to push container images"
        else
            print_color $GREEN "✅ write:packages - Can push container images"
        fi
        
        if [[ "$has_read_packages" != "true" ]]; then
            print_color $YELLOW "⚠️  read:packages - Recommended for reading container images"
        else
            print_color $GREEN "✅ read:packages - Can read container images"
        fi
        
        if [[ "$has_delete_packages" != "true" ]]; then
            print_color $YELLOW "⚠️  delete:packages - Recommended for cleanup temporary images"
        else
            print_color $GREEN "✅ delete:packages - Can cleanup temporary images"
        fi
        
        echo ""
        echo "To fix missing permissions, run:"
        echo "  gh auth refresh -s write:packages,read:packages,delete:packages"
        echo "  or"
        echo "  gh auth refresh -s repo (includes all package permissions)"
        echo ""
        
        if [[ "$has_write_packages" != "true" ]]; then
            print_color $RED "❌ Cannot continue without write:packages permission"
            exit 1
        fi
        
        read -p "Continue with limited permissions? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Please refresh your token permissions and try again"
            exit 1
        fi
        
        print_color $YELLOW "⚠️  Continuing with limited permissions - some features may not work"
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

# Push platform-specific images temporarily for manifest creation
push_temp_images() {
    # Detect system architecture
    local arch=$(uname -m)
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        print_color $YELLOW "📤 Pushing temporary platform-specific images for manifest creation..."
        
        # Push ARM64 images (temporary for manifest)
        print_color $YELLOW "🚀 Pushing ARM64 images (temporary)..."
        docker push "${IMAGE_NAME}:${NIGHTLY_TAG}-arm64"
        docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-arm64"
        
        # Push AMD64 images (temporary for manifest)
        print_color $YELLOW "🚀 Pushing AMD64 images (temporary)..."
        docker push "${IMAGE_NAME}:${NIGHTLY_TAG}-amd64"
        docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}-amd64"
        
        print_color $GREEN "✅ Temporary platform-specific images pushed"
        echo "   (These will be cleaned up after manifest creation)"
    else
        print_color $YELLOW "📤 Pushing to GitHub Container Registry..."
        
        docker push "${IMAGE_NAME}:${NIGHTLY_TAG}"
        docker push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
        
        print_color $GREEN "✅ Images pushed successfully"
        echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
        echo "   Available: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
    fi
}

# Delete temporary platform-specific images from registry
cleanup_temp_images() {
    local arch=$(uname -m)
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        print_color $YELLOW "🧹 Cleaning up temporary platform-specific images from registry..."
        
        # Set environment for GitHub CLI
        export GH_PAGER=""
        
        # Define tags to delete (platform-specific images)
        local tags_to_delete=(
            "${NIGHTLY_TAG}-arm64"
            "${NIGHTLY_TAG}-amd64"
            "${NIGHTLY_LATEST_TAG}-arm64"
            "${NIGHTLY_LATEST_TAG}-amd64"
        )
        
        print_color $YELLOW "🔍 Fetching container versions for cleanup..."
        
        # Fetch all container versions using GitHub API
        local versions_json
        if versions_json=$(gh api "/users/$OWNER/packages/container/$REPO_NAME/versions" --paginate \
            -H "Accept: application/vnd.github+json" 2>/dev/null); then
            
            print_color $YELLOW "📋 Found container versions, processing deletions..."
            
            # Loop through tags to delete
            for tag in "${tags_to_delete[@]}"; do
                print_color $YELLOW "� Looking for tag: $tag"
                
                # Find version ID for this specific tag
                local version_id
                version_id=$(echo "$versions_json" | jq -r --arg TAG "$tag" '
                    .[] | select(.metadata.container.tags[]? == $TAG) | .id' | head -n 1)
                
                if [[ -z "$version_id" || "$version_id" == "null" ]]; then
                    echo "   ⚠️  Tag not found or already deleted: $tag"
                else
                    print_color $YELLOW "🗑️  Deleting version ID: $version_id (tag: $tag)"
                    
                    if gh api --method DELETE \
                        -H "Accept: application/vnd.github+json" \
                        "/users/$OWNER/packages/container/$REPO_NAME/versions/$version_id" 2>/dev/null; then
                        print_color $GREEN "   ✅ Deleted: $tag"
                    else
                        echo "   ⚠️  Failed to delete (may not exist): $tag"
                    fi
                fi
            done
            
            print_color $GREEN "✅ Temporary platform-specific images cleanup completed"
            echo "   Only multi-platform manifests remain in registry"
            
        else
            print_color $YELLOW "⚠️  Could not fetch container versions (API access issue)"
            echo "   Platform-specific images may remain in registry"
            echo "   This is not critical - manifests were created successfully"
        fi
    fi
}

# Create multi-platform manifest and push final results
create_manifest() {
    # Detect system architecture
    local arch=$(uname -m)
    
    if [[ "$arch" == "arm64" || "$arch" == "aarch64" ]]; then
        print_color $YELLOW "� Creating multi-platform manifests..."
        
        # Remove existing manifests if they exist (to avoid conflicts)
        print_color $YELLOW "🧹 Cleaning up any existing manifests..."
        docker manifest rm "${IMAGE_NAME}:${NIGHTLY_TAG}" 2>/dev/null || true
        docker manifest rm "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" 2>/dev/null || true
        
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
        print_color $YELLOW "📤 Pushing manifest for ${NIGHTLY_TAG}..."
        if docker manifest push "${IMAGE_NAME}:${NIGHTLY_TAG}"; then
            print_color $GREEN "✅ Manifest pushed successfully: ${NIGHTLY_TAG}"
        else
            print_color $RED "❌ Failed to push manifest: ${NIGHTLY_TAG}"
            exit 1
        fi
        
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
        print_color $YELLOW "📤 Pushing manifest for ${NIGHTLY_LATEST_TAG}..."
        if docker manifest push "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"; then
            print_color $GREEN "✅ Manifest pushed successfully: ${NIGHTLY_LATEST_TAG}"
        else
            print_color $RED "❌ Failed to push manifest: ${NIGHTLY_LATEST_TAG}"
            exit 1
        fi
        
        print_color $GREEN "✅ Multi-platform manifests created and pushed successfully"
        
        # # Wait a moment for manifest propagation
        # print_color $YELLOW "⏳ Waiting for manifest propagation (10 seconds)..."
        # sleep 10
        
        # Verify manifests are accessible
        print_color $YELLOW "🔍 Verifying manifest accessibility..."
        local verification_failed=false
        
        if docker manifest inspect "${IMAGE_NAME}:${NIGHTLY_TAG}" &>/dev/null; then
            print_color $GREEN "✅ Manifest verified: ${NIGHTLY_TAG}"
        else
            print_color $RED "❌ Manifest verification failed: ${NIGHTLY_TAG}"
            verification_failed=true
        fi
        
        if docker manifest inspect "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" &>/dev/null; then
            print_color $GREEN "✅ Manifest verified: ${NIGHTLY_LATEST_TAG}"
        else
            print_color $RED "❌ Manifest verification failed: ${NIGHTLY_LATEST_TAG}"
            verification_failed=true
        fi
        
        if [[ "$verification_failed" == "true" ]]; then
            print_color $YELLOW "⚠️  Manifest verification failed - skipping cleanup to preserve platform images"
            print_color $YELLOW "   Platform-specific images will remain available in registry"
            echo ""
            print_color $CYAN "🎯 Available Images:"
            echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
            echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
            echo "   ARM64 specific: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}-arm64"
            echo "   AMD64 specific: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}-amd64"
            return
        fi
        
        echo ""
        print_color $CYAN "🎯 Final Results - Only these images are available in registry:"
        echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
        echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
        echo ""
        echo "✨ These manifests automatically select the correct platform (ARM64/AMD64)"
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
        echo "  3. Temporarily push platform-specific images to create manifests"
        echo "  4. Create and push multi-platform manifests"
        echo "  5. Clean up temporary platform-specific images"
        echo "  6. Only multi-platform manifests remain in registry"
        echo ""
        echo "Final images available in registry:"
        echo "  - ${IMAGE_NAME}:${NIGHTLY_TAG} (multi-platform manifest)"
        echo "  - ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG} (multi-platform manifest)"
        echo ""
        print_color $CYAN "✨ Platform-specific images (-arm64/-amd64) will NOT be kept in registry"
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

# Calculate and display build duration
show_build_completion() {
    local end_time=$(date +%s)
    local duration=$((end_time - START_TIME))
    local minutes=$((duration / 60))
    local seconds=$((duration % 60))
    
    echo ""
    print_color $GREEN "🎉 Build and publish completed successfully!"
    echo ""
    print_color $CYAN "⏱️  Build Duration:"
    if [ $minutes -gt 0 ]; then
        echo "   Total time: ${minutes} minute(s) and ${seconds} second(s)"
    else
        echo "   Total time: ${seconds} second(s)"
    fi
    echo ""
    print_color $BLUE "📦 Final Results:"
    echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_TAG}"
    echo "   Multi-platform: docker pull ${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}"
}

# Main execution
main() {
    # Record start time
    START_TIME=$(date +%s)
    
    print_color $BLUE "🚀 Beo Echo - Docker Build & Publish"
    echo "=================================="
    echo "Started at: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    
    check_github_auth
    check_docker
    check_repo_permissions
    get_repo_info
    show_build_summary
    build_image
    login_registry
    push_temp_images
    create_manifest
    
    # Only cleanup if manifest verification passed
    # if [[ "$(uname -m)" == "arm64" || "$(uname -m)" == "aarch64" ]]; then
    #     # Check if we should proceed with cleanup
    #     if docker manifest inspect "${IMAGE_NAME}:${NIGHTLY_TAG}" &>/dev/null && \
    #        docker manifest inspect "${IMAGE_NAME}:${NIGHTLY_LATEST_TAG}" &>/dev/null; then
    #         cleanup_temp_images
    #     else
    #         print_color $YELLOW "⚠️  Skipping cleanup due to manifest verification issues"
    #     fi
    # fi
    
    show_build_completion
}

# Run main function
main
