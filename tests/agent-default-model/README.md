# Agent Default Model Test

Tests that agents automatically use model named "default" when no modelRef is specified.

## What it tests
- Agent creation without explicit modelRef
- Automatic fallback to model named "default"
- Query execution using implicit model reference
- Default model behavior validation

## Running
```bash
chainsaw test
```

Validates that agents can successfully execute queries using the default model when no explicit model is specified.