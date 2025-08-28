#!/usr/bin/env python3
"""Generate OpenAPI schema without running the server."""
import json
import sys
from pathlib import Path

# Add the src directory to Python path
sys.path.insert(0, str(Path(__file__).parent / "src"))

from ark_api.main import app

# Generate OpenAPI schema
openapi_schema = app.openapi()

# Write to file
with open("openapi.json", "w") as f:
    json.dump(openapi_schema, f, indent=2)

print("OpenAPI schema written to openapi.json")