FROM python:3.12.9-slim

# Update package list and upgrade security packages
RUN apt-get update && \
    apt-get upgrade -y libpam0g libsqlite3-0 zlib1g libgnutls30 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

# Copy dependency files
COPY ark-mcp/pyproject.toml ark-mcp/uv.lock ./

# Copy wheel files
COPY ark-mcp/out/ark*.whl /tmp/wheels/

# Copy source code only (not pyproject.toml again)
COPY ark-mcp/src/ ./src/

# Add ark-sdk wheel as a dependency and sync all dependencies
RUN uv add file://$(ls /tmp/wheels/ark*.whl) && \
    uv sync --frozen --no-install-project

# Create non-root user
RUN addgroup --system --gid 1001 arkmcp
RUN adduser --system --uid 1001 arkmcp

# Create cache directory for uv
RUN mkdir -p /app/.cache && \
    chown -R arkmcp:arkmcp /app

# Switch to non-root user
USER arkmcp

# Set HOME environment variable to fix uv cache location
ENV HOME=/app
ENV UV_CACHE_DIR=/app/.cache/uv

# Expose port 2627 (AMCP on dial pad)
EXPOSE 2627

# Run the MCP server
CMD ["uv", "run", "python", "-m", "ark_mcp"]