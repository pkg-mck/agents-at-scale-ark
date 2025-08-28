# Observability and Event Recording

The ark operator provides comprehensive observability through a structured event recording system with configurable verbosity levels.

## Event Recording Architecture

All operations in the system are tracked using two main components:

### OperationTracker
Tracks individual operations with timing and lifecycle management:
- **Start events**: Emitted when operations begin
- **Completion events**: Emitted when operations succeed with duration
- **Error events**: Emitted when operations fail with error details
- **Termination events**: Special completion for graceful termination

### ExecutionRecorder  
Records high-level execution events for teams and agents:
- **Team execution**: Team strategy and lifecycle events
- **Team members**: Individual member execution within teams
- **Team turns**: Round-robin turn tracking
- **Agent execution**: Agent lifecycle and model interactions

## Verbosity Levels

The system uses 4 verbosity levels to control event granularity:

### Level 0 (Always Visible) - Critical Operations
**Always emitted regardless of log configuration**

- **Query Resolution**: Query start, completion, and errors
- **Model Resolution**: Model validation and configuration

**Use Case**: Production monitoring, health checks, SLA tracking

**Example Events**:
```json
{
  "name": "my-query",
  "namespace": "default",
  "targets": "2",
  "component": "query"
}
```

### Level 1 (Standard) - Operational Events  
**Emitted when log verbosity >= 1**

- **Agent Execution**: Agent lifecycle and configuration
- **Team Execution**: Team strategy execution  
- **Tool Calls**: Tool invocation and results
- **Team Members**: Individual team member execution

**Use Case**: Standard operations monitoring, debugging workflows

**Example Events**:
```json
{
  "name": "my-agent", 
  "model": "gpt-4",
  "component": "agent",
  "duration": "2.5s"
}
```

### Level 2 (Detailed) - LLM Interactions
**Emitted when log verbosity >= 2**

- **LLM Calls**: Model API calls and responses
- **Model Interactions**: Request/response cycles

**Use Case**: Debugging model interactions, performance tuning, cost tracking

**Example Events**:
```json
{
  "name": "gpt-4",
  "agent": "my-agent", 
  "model": "gpt-4",
  "component": "llm",
  "duration": "1.2s"
}
```

### Level 3 (Debug) - Response Content
**Emitted when log verbosity >= 3**

- **Response Content**: Full LLM response data
- **Termination Messages**: Detailed termination reasons
- **Detailed Operational Data**: Complete context and metadata

**Use Case**: Full debugging, content analysis, development

**Security Warning**: Contains sensitive data from LLM responses

## Configuration

### Development Environment

```bash
# Level 0 (default) - critical operations only
cd ark && make dev

# Level 1 - add operational events
cd ark && make dev ARGS="--zap-log-level=1"

# Level 2 - add LLM call tracking
cd ark && make dev ARGS="--zap-log-level=2"

# Level 3 - add response content (debug)
cd ark && make dev ARGS="--zap-log-level=3"
```

### Production Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ark-controller
spec:
  template:
    spec:
      containers:
      - name: manager
        env:
        - name: ZAPLOGLEVEL
          value: "1"  # Recommended for production
        args:
        - --zap-log-level=$(ZAPLOGLEVEL)
```

### Environment Variables

| Variable | Values | Description |
|----------|--------|-------------|
| `ZAPLOGLEVEL` | `0-3` | Controls event verbosity level |

## Monitoring and Observability

### Viewing Events

#### All Events
```bash
# View all operator events
kubectl get events --sort-by='.lastTimestamp'

# Filter by event type
kubectl get events --field-selector reason=ResolveStart
```

#### Resource-Specific Events
```bash
# Query execution monitoring
kubectl describe query my-query

# Agent execution details  
kubectl describe agent my-agent

# Model resolution status
kubectl describe model my-model
```

#### Real-Time Monitoring
```bash
# Watch all events live
kubectl get events --watch

# Watch specific resource type
kubectl get events --watch --field-selector involvedObject.kind=Query

# Monitor with filtering
kubectl get events --watch --field-selector reason!=Pulled,reason!=Created
```

### Event Filtering and Analysis

#### Common Event Types

| Event Type | Verbosity | Description |
|------------|-----------|-------------|
| `ResolveStart` | 0 | Query/Model resolution begins |
| `ResolveComplete` | 0 | Query/Model resolution succeeds |
| `ResolveError` | 0 | Query/Model resolution fails |
| `AgentExecutionStart` | 1 | Agent begins execution |
| `AgentExecutionComplete` | 1 | Agent completes execution |
| `LLMCallStart` | 2 | LLM API call begins |
| `LLMCallComplete` | 2 | LLM API call completes |
| `ToolCallStart` | 1 | Tool invocation begins |
| `ToolCallComplete` | 1 | Tool invocation completes |

#### Event Metadata Structure

All events include structured metadata in JSON format:

```json
{
  "name": "resource-name",
  "component": "query|agent|team|tool|llm|model",
  "duration": "1.234s",
  "error": "error message if failed",
  "namespace": "kubernetes-namespace",
  "additionalContext": "varies by component"
}
```

### Troubleshooting Guide

#### Common Issues

**No Events Visible**
- Check verbosity level: `kubectl logs deployment/ark-controller | grep "verbosity"`
- Verify RBAC permissions for event creation
- Ensure resources are being processed: `kubectl get queries,agents,models`

**Missing Detailed Events**  
- Increase verbosity level in deployment configuration
- Check log level: `kubectl logs deployment/ark-controller | head -10`

**Too Many Events**
- Reduce verbosity level to 0 or 1 for production
- Use event filtering: `kubectl get events --field-selector reason!=EventType`

#### Performance Considerations

- **Level 0-1**: Minimal performance impact, suitable for production
- **Level 2**: Moderate impact, adds LLM call tracking overhead  
- **Level 3**: High impact, includes response content serialization

#### Security Considerations

- **Level 3 logs contain sensitive data** from LLM responses
- Use appropriate RBAC to restrict access to events
- Consider log retention policies for sensitive environments
- Monitor event storage consumption in cluster

## Integration with Monitoring Systems

### Prometheus Metrics

Events can be scraped and converted to metrics:

```yaml
# Example ServiceMonitor for event-based metrics
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: ark-events
spec:
  selector:
    matchLabels:
      app: ark-controller
  endpoints:
  - port: metrics
```

### Log Aggregation

Events are structured JSON suitable for log aggregation:

```bash
# Export events for analysis
kubectl get events -o json | jq '.items[] | select(.reason | startswith("Resolve"))'
```

### Alerting

Create alerts based on event patterns:

```yaml
# Example alert for query failures
- alert: QueryResolutionFailure
  expr: increase(kubernetes_events_total{reason="ResolveError"}[5m]) > 0
  labels:
    severity: warning
  annotations:
    summary: "Query resolution failing"
```