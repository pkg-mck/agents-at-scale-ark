#!/bin/bash

# Build and push PostgreSQL Memory Service Docker image to local Kubernetes cluster

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
IMAGE_NAME="postgres-memory"
IMAGE_TAG="${1:-latest}"
TARGET_CLUSTER="${2:-auto}"

echo "Building and pushing PostgreSQL Memory Service Docker image..."
echo "Image: ${IMAGE_NAME}:${IMAGE_TAG}"
echo "Target cluster: ${TARGET_CLUSTER}"

"${PROJECT_ROOT}/scripts/build-and-push.sh" \
    --image "${IMAGE_NAME}" \
    --dockerfile "${SCRIPT_DIR}/Dockerfile" \
    --context "${SCRIPT_DIR}" \
    --tag "${IMAGE_TAG}" \
    --cluster "${TARGET_CLUSTER}"

echo ""
echo "Next steps:"
echo "  Deploy with Helm:"
echo "    helm install postgres-memory ./chart"
echo ""
echo "  Test the deployment:"
echo "    kubectl get pods -l app.kubernetes.io/name=postgres-memory"