#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "🔨 Building evaluator-llm Docker image..."

# Build from services directory
cd ..
../scripts/build-and-push.sh -i evaluator-llm -c . -f evaluator-llm/Dockerfile

echo ""
echo "Next steps:"
echo "  Deploy with Helm:"
echo "    helm install evaluator-llm ./chart"
echo ""
echo "  Test the deployment:"
echo "    kubectl get pods -l app.kubernetes.io/name=evaluator-llm"