# MCP Filesystem Server Test

Comprehensive chainsaw test for the MCP (Model Context Protocol) Filesystem Server functionality with full ARK integration.

## What it tests

### Core MCP Server Functionality
- **MCP Server deployment**: Validates the MCP filesystem server starts and runs correctly
- **SSE endpoint**: Tests Server-Sent Events endpoint for MCP communication
- **Filebrowser UI**: Validates the web-based file management interface on port 8080
- **Service connectivity**: Ensures services are accessible and responding
- **MCPServer resource**: Validates the ARK MCPServer custom resource is created

### ARK Integration Testing
- **Model deployment**: Azure OpenAI model configuration and readiness
- **Agent creation**: File management agent with MCP filesystem tools
- **Tool integration**: MCP tools (read_file, write_file, list_directory, copy_file, delete_file) 
- **Query execution**: End-to-end filesystem operations through ARK agent
- **Response validation**: Confirms agent successfully uses MCP tools

### File Operations Testing
- **Direct operations**: Container-level file operations for baseline validation
- **MCP tool operations**: Filesystem operations through agent and MCP tools
- **Multi-step workflows**: Complex file management tasks (create, read, copy, delete)

## Components tested

### MCP Filesystem Server
- Port 9090 - MCP proxy with SSE transport
- Filesystem operations via `mcp-filesystem-server`
- Persistent storage mounted at `/data`
- Tools: read_file, write_file, list_directory, copy_file, delete_file

### Filebrowser UI
- Port 8080 - Web-based file management interface
- Shared storage with MCP server

### ARK Resources
- **Model**: Azure OpenAI gpt-4.1-mini configuration
- **MCPServer**: Filesystem server registration 
- **Agent**: File management assistant with MCP tools
- **Query**: Multi-step file operation workflow

## Test Workflow

1. **Setup**: Deploy base infrastructure (RBAC, PVC, ConfigMap, Services, Pods)
2. **Validation**: Verify all pods running and model ready
3. **HTTP Testing**: Validate MCP server endpoints with Hurl
4. **ARK Integration**: Deploy and test Agent with MCP filesystem tools
5. **Query Execution**: Run comprehensive file management workflow
6. **Direct Testing**: Validate container-level file operations

## Running

```bash
# Run the complete test suite
chainsaw test

# Or via Makefile
make test-chainsaw
```

**Environment Requirements:**
- `E2E_TEST_AZURE_OPENAI_KEY`: Azure OpenAI API key
- `E2E_TEST_AZURE_OPENAI_BASE_URL`: Azure OpenAI base URL

Validates that the MCP filesystem server provides working file operations through both direct MCP interface and ARK agent integration, ensuring comprehensive filesystem functionality in Kubernetes environments.