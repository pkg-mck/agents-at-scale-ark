# MCP Filesystem Server Helm Chart

This directory contains a Helm chart for deploying the MCP (Model Context Protocol) Filesystem Server to Kubernetes.

## Prerequisites

- Docker
- Kubernetes cluster
- Helm 3.x
- kubectl configured to access your cluster

## Build and Deploy

### 1. Build and Load Image

```bash
# Build the image
docker build -t mcp-filesys:latest .

# Load image into cluster - adjust command based on your cluster type
# kind load docker-image mcp-filesys:latest
# minikube image load mcp-filesys:latest
```

### 2. Deploy with Helm

```bash
# Install the chart
helm install mcp-filesys ./chart

# Or upgrade if already installed
helm upgrade mcp-filesys ./chart

# Install with custom values
helm install mcp-filesys ./chart -f custom-values.yaml
```

### 3. Verify Deployment

```bash
# Check pods
kubectl get pods -l app.kubernetes.io/name=mcp-filesys

# Check service
kubectl get svc -l app.kubernetes.io/name=mcp-filesys

# Check PVC
kubectl get pvc
```

## Configuration

Key configuration options in `values.yaml`:

- `persistence.size`: Storage size for the `/data` volume (default: 10Gi)
- `persistence.storageClass`: Storage class for PVC
- `resources`: CPU and memory limits/requests
- `config`: MCP proxy and server configuration

## Services

This deployment creates **two separate services**:

### 1. MCP Server Service
- **Service Name**: `mcp-filesys-server`
- **Port**: 9090
- **Purpose**: Handles MCP (Model Context Protocol) requests for filesystem operations
- **Component**: `mcp-server`

### 2. File Browser Service  
- **Service Name**: `mcp-filesys-filebrowser`
- **Port**: 8080
- **Purpose**: Provides web-based file browser UI for managing files
- **Component**: `filebrowser`

## Accessing the Services

### MCP Server (Port 9090)
```bash
# Port forward MCP server to local machine
kubectl port-forward svc/mcp-filesys-server 9090:9090

# Test MCP endpoint
curl http://localhost:9090/sse
```

### File Browser UI (Port 8080)
```bash
# Port forward file browser to local machine  
kubectl port-forward svc/mcp-filesys-filebrowser 8080:8080

# Open in browser
open http://localhost:8080
```

## Uninstall

```bash
helm uninstall mcp-filesys
```