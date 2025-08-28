from typing import Dict, List, Any, Optional
from pydantic import BaseModel


class QueryTarget(BaseModel):
    type: str
    name: str


class Response(BaseModel):
    target: QueryTarget
    content: str


class Model(BaseModel):
    name: str
    type: str
    config: Dict[str, Any] = {}


class EvaluationRequest(BaseModel):
    queryId: str
    input: str
    responses: List[Response]
    query: Dict[str, Any]
    model: Model


class EvaluationResponse(BaseModel):
    score: Optional[str] = None
    passed: bool = False
    metadata: Optional[Dict[str, str]] = None
    error: Optional[str] = None