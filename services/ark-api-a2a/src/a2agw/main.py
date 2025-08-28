import logging
import os
import time
from contextlib import asynccontextmanager
from datetime import datetime
from importlib.metadata import version

import uvicorn
from ark_sdk.k8s import init_k8s
from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from opentelemetry import baggage, propagate, trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.httpx import HTTPXClientInstrumentor
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from .manager import DynamicManager
from .registry import get_registry

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

VERSION=version('a2agw')
PORT = int(os.getenv("PORT", "7184"))

manager = DynamicManager()

def setup_telemetry():
    """Initialize OpenTelemetry tracing"""
    otel_endpoint = os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
    if not otel_endpoint:
        logger.info("OTEL_EXPORTER_OTLP_ENDPOINT not set, telemetry disabled")
        return
    
    service_name = os.getenv("OTEL_SERVICE_NAME", "ark-api-a2a")
    
    # Set up resource
    resource = Resource.create({
        "service.name": service_name,
        "service.version": VERSION,
    })
    
    # Set up tracer provider
    tracer_provider = TracerProvider(resource=resource)
    trace.set_tracer_provider(tracer_provider)
    
    # Set up OTLP exporter
    otlp_exporter = OTLPSpanExporter(endpoint=f"{otel_endpoint}/v1/traces")
    span_processor = BatchSpanProcessor(otlp_exporter)
    tracer_provider.add_span_processor(span_processor)
    
    logger.info(f"Telemetry initialized for {service_name} -> {otel_endpoint}")

def extract_session_context(request: Request):
    """Extract OTEL context and session ID from request headers"""
    # Extract OTEL trace context from headers
    ctx = propagate.extract(request.headers)
    
    # Extract session ID from custom header
    session_id = request.headers.get("x-session-id")
    if session_id:
        # Add session to baggage
        ctx = baggage.set_baggage("session.id", session_id, context=ctx)
    
    return ctx, session_id

@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info(f"Starting up ARK A2A GW v{VERSION}")
    
    # Initialize telemetry
    setup_telemetry()
    
    await init_k8s()
    
    # Initialize manager and start periodic sync
    await manager.initialize()
    
    app.mount("/agent", manager.app)
    
    yield
    
    # Shutdown manager on app shutdown
    logger.info("Shutting down ARK A2A GW...")
    await manager.shutdown()

app = FastAPI(
    title="A2AGW - Agent-to-Agent Gateway",
    description="Gateway for agent communication in the ARK ecosystem",
    version=VERSION,
    lifespan=lifespan,
)

# Instrument FastAPI and HTTPx for automatic tracing
FastAPIInstrumentor.instrument_app(app)
HTTPXClientInstrumentor().instrument()

@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    logger.error(f"Validation error for {request.method} {request.url.path}")
    logger.error(f"Validation errors: {exc.errors()}")
    
    try:
        body = await request.body()
        body_str = body.decode() if body else "<empty body>"
        logger.error(f"Request body that failed validation: {body_str}")
    except Exception as e:
        logger.error(f"Could not read body: {e}")
    
    return JSONResponse(
        status_code=422,
        content={"detail": exc.errors(), "body": "Invalid request"}
    )


@app.middleware("http")
async def session_aware_middleware(request: Request, call_next):
    start_time = time.time()
    
    # Extract OTEL context and session ID
    otel_ctx, session_id = extract_session_context(request)
    
    # Add session info to logs
    session_info = f"session={session_id}" if session_id else "no-session"
    logger.info(
        f"Request: {request.method} {request.url.path} - {session_info} - Query: {dict(request.query_params)}"
    )
    
    # Process request with OTEL context
    response = await call_next(request)
    
    process_time = time.time() - start_time
    logger.info(
        f"Response: {request.method} {request.url.path} - {session_info} - Status: {response.status_code} - Time: {process_time:.3f}s"
    )
    
    return response


@app.get("/agents", response_model=list[dict])
async def list_agents():
    agents = await get_registry().list_agents()
    return [
        {
            "name": agent.name,
            "description": agent.description,
            "capabilities": [skill.name for skill in agent.skills],
            "host": "localhost",
            "agent-card": f"/agent/{agent.name}/.well-known/agent.json",
            "created_at": datetime.utcnow().isoformat(),
            "metadata": {"type": "analytical", "version": agent.version},
        }
        for agent in agents
    ]

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "a2agw"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=PORT)
