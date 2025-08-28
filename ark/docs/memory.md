# Memory System

## Overview

Memory provides persistent conversation storage for queries. Uses HTTP protocol with PostgreSQL backend.

## How It Works

### Message Flow
1. Query references Memory resource (or falls back to "default" memory)
2. Controller loads existing messages from memory service
3. User input appended to conversation history
4. Agent/Team executes with full conversation context
5. Response messages saved back to memory service

### Session Isolation
- Each Query gets unique session ID (Query UID)
- Messages stored per session in PostgreSQL
- Sessions prevent conversation mixing between queries

### Fallback Behavior
- No memory reference → uses "default" memory if exists
- No "default" memory → uses NoopMemory (no persistence)
- Memory service unavailable → query fails with error

## Memory CRD

### Spec
```yaml
apiVersion: ark.mckinsey.com/v1alpha1
kind: Memory
metadata:
  name: chat-history
spec:
  address:
    valueFrom:
      serviceRef:
        name: postgres-memory
        port: 8080
```

### Status
Controller resolves service address and updates status with endpoint URL.

## PostgreSQL Memory Service

### Location
`services/postgres-memory/`

### API Endpoints
- `PUT /message/{session_id}` - Store message
- `GET /message/{session_id}` - Retrieve messages
- `GET /health` - Health check

### Database Schema
Table: `ark_messages`
- `session_id` - Query UID for isolation
- `message` - JSONB chat completion message  
- `created_at` - Timestamp for ordering

## Deployment

### Install PostgreSQL Operator
```bash
make install-pgo
```

### Deploy Memory Service
```bash
# With PGO cluster
helm install postgres-memory ./services/postgres-memory/chart --set postgres.enabled=true

# With external database  
helm install postgres-memory ./services/postgres-memory/chart --set postgres.enabled=false
```

### Configuration Options
```yaml
# External PostgreSQL (default)
postgres:
  enabled: false
  host: postgres
  database: memory
  username: postgres

# Crunchy PGO cluster
postgres:
  enabled: true
  storage: "10Gi"
  replicas: 2
```

## Usage

### Query with Memory
```yaml
apiVersion: ark.mckinsey.com/v1alpha1
kind: Query
spec:
  input: "What did we discuss earlier?"
  targets:
    - name: my-agent
  memory:
    name: chat-history
```

### Query without Memory
```yaml
apiVersion: ark.mckinsey.com/v1alpha1
kind: Query
spec:
  input: "Hello"
  targets:
    - name: my-agent
  # No memory specified - uses "default" memory or NoopMemory
```