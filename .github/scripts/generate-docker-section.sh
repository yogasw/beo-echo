#!/bin/bash

# Script to generate Docker Images section for release notes
# Usage: ./generate-docker-section.sh <repository> <version>
# Example: ./generate-docker-section.sh yogasw/beo-echo 2.6.0

REPOSITORY="${1}"
VERSION="${2}"

if [ -z "$REPOSITORY" ] || [ -z "$VERSION" ]; then
  echo "Error: Missing required arguments"
  echo "Usage: $0 <repository> <version>"
  echo "Example: $0 yogasw/beo-echo 2.6.0"
  exit 1
fi

# Generate Docker Images section
cat <<EOF

### Docker Images
Use one of these commands to pull the Docker image:

\`\`\`bash
# Pull specific version
docker pull ghcr.io/${REPOSITORY}:${VERSION}

# Pull latest version
docker pull ghcr.io/${REPOSITORY}:latest
\`\`\`
EOF
