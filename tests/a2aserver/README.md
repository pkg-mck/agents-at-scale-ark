# A2AServer Basic Test

Tests basic A2AServer CRD creation and validation. Tests the complete A2A execution engine integration with the LangChain weather agent.

## What it tests

- LangChain weather service deployment via `make install`
- A2AServer resource creation and agent discovery
- Agent creation with `a2a` execution engine and proper annotations
- Templated query execution using parameter substitution
- Query completion through dedicated A2A execution engine

## Running

```bash
chainsaw test
```

The test validates that A2A agents are properly discovered, created with the reserved 'a2a' execution engine, and can successfully execute templated queries through the dedicated A2A execution engine running in-cluster.

Test skipped, it relies on deploying the hosted-langchain-agents in samples/agent-hosting/hosted-langchain-agents (@Dave, check if needed - too specific for e2e tests)