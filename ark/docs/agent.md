# Agent: Triggering Evaluation

To have an agent's responses automatically evaluated by the default evaluation operator, set the `evaluator` label to `judge` on the agent resource.

When a Query targets an agent with this label, an Evaluation will be created and processed automatically.

## Example Agent with Evaluator Label
```yaml
apiVersion: agent.mckinsey/v1
kind: Agent
metadata:
  name: research-agent
  labels:
    evaluator: judge
spec:
  description: "Research and analysis specialist"
  prompt: |
    You are a research analyst. Your role is to gather information, analyze data,
    and provide insights based on research findings. Be thorough and professional.
```

## How It Works
- When a Query targets this agent, the evaluation operator will automatically create and process an Evaluation for it.
- The Evaluation will be named `<query_name>-<agent_name>` (e.g., `sample-query-research-agent`).
- The evaluation logic and status will be managed by the default evaluation operator.

## Notes
- Only agents with `labels.evaluator: judge` will be evaluated by the default operator.
- No manual Evaluation creation is needed.
- See the Evaluation documentation for more details on status fields and results. 