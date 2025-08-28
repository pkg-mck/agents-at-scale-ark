import logging
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Request
import aiohttp
from .types import EvaluationRequest, EvaluationResponse
from .evaluator import LLMEvaluator

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Global shared session for all evaluator instances
shared_session = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    global shared_session
    shared_session = aiohttp.ClientSession()
    logger.info("LLM Evaluator service started")
    yield
    # Shutdown
    if shared_session:
        await shared_session.close()
    logger.info("LLM Evaluator service stopped")


def create_app() -> FastAPI:
    app = FastAPI(
        title="LLM Evaluator Service",
        description="AI-powered query evaluation service using LLM-as-a-Judge",
        version="0.1.0",
        lifespan=lifespan
    )

    @app.middleware("http")
    async def log_requests(request: Request, call_next):
        if request.url.path == "/health":
            response = await call_next(request)
            return response
        
        logger.info(f"{request.method} {request.url.path}")
        response = await call_next(request)
        return response

    @app.get("/health")
    async def health():
        return {"status": "healthy", "service": "evaluator-llm"}

    @app.get("/ready")
    async def ready():
        return {"status": "ready", "service": "evaluator-llm"}

    @app.post("/evaluate", response_model=EvaluationResponse)
    async def evaluate(request: EvaluationRequest) -> EvaluationResponse:
        """
        Evaluate query responses using LLM-as-a-Judge approach
        """
        try:
            logger.info(f"Received evaluation request for query {request.queryId}")
            
            if not shared_session:
                raise HTTPException(status_code=500, detail="HTTP session not initialized")
            
            # Create a new evaluator instance for each request for thread safety
            evaluator = LLMEvaluator(session=shared_session)
            result = await evaluator.evaluate(request)
            
            logger.info(f"Evaluation completed for query {request.queryId}")
            return result
            
        except Exception as e:
            logger.error(f"Error processing evaluation request: {str(e)}")
            raise HTTPException(status_code=500, detail=str(e))

    return app