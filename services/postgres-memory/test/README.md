# PostgreSQL Memory Service Test

Comprehensive chainsaw test for the PostgreSQL Memory Service functionality with full ARK integration.

## What it tests

### Core Memory Service Functionality
- **Memory service deployment**: Validates the memory service starts and runs correctly
- **PostgreSQL integration**: Tests PostgreSQL cluster deployment and connectivity
- **HTTP endpoints**: Validates REST API endpoints for message storage and retrieval
- **Service connectivity**: Ensures services are accessible and responding
- **Memory resource**: Validates the ARK Memory custom resource is created

### ARK Integration Testing
- **Model deployment**: Azure OpenAI model configuration and readiness
- **Agent creation**: Memory-aware agent with conversation persistence
- **Memory integration**: Agent uses Memory resource for conversation storage
- **Query execution**: End-to-end conversation memory through ARK agent
- **Response validation**: Confirms agent successfully uses memory functionality

### Message Operations Testing
- **Single message operations**: Add and retrieve individual messages
- **Bulk message operations**: Add and retrieve multiple messages at once
- **Session isolation**: Verify messages are isolated by session ID
- **Message ordering**: Confirm messages are retrieved in chronological order
- **Complex message types**: Test tool calls and structured message content

## Components tested

### PostgreSQL Memory Service
- Port 8080 - HTTP API for message storage
- PostgreSQL backend with PGO (PostgreSQL Operator)
- Session-based message isolation
- JSON message storage with JSONB support

### Database Operations
- Message persistence across service restarts
- CRUD operations on message storage
- Transaction handling for bulk operations
- Database schema creation and migration

### ARK Resources
- **Model**: Azure OpenAI GPT-4.1-mini configuration
- **Memory**: PostgreSQL-backed conversation memory
- **Agent**: Memory-aware conversational agent
- **Query**: Conversation with memory persistence

## Test Workflow

1. **Setup**: Deploy PostgreSQL cluster, memory service, and ARK resources
2. **Validation**: Verify all pods running, database ready, and model ready
3. **HTTP Testing**: Validate memory service endpoints with Hurl
4. **ARK Integration**: Deploy and test Agent with memory functionality
5. **Query Execution**: Run conversation with memory persistence
6. **Memory Validation**: Verify conversation context is maintained

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
- `POSTGRES_MEMORY_IMAGE`: PostgreSQL Memory service image

Validates that the PostgreSQL Memory service provides working conversation memory through both direct HTTP interface and ARK agent integration, ensuring comprehensive memory functionality in Kubernetes environments.