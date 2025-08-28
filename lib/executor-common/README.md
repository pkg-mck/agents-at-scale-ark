# Executor Common

Common utilities and types shared across execution engine services (AutoGen, LangChain, CrewAI).

## Contents

- **`types.py`** - Shared type definitions for execution requests/responses
- **`app.py`** - Base FastAPI application setup with common patterns
- **`base.py`** - Abstract base classes for executors

## Usage

This package is used internally by the executor services and provides:

- Consistent API contracts
- Shared FastAPI setup
- Common logging and error handling patterns
- Abstract base classes for executor implementations