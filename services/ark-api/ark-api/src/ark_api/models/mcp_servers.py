from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class MCPServerResponse(BaseModel):
    name: str
    namespace: str
    description: Optional[str] = None
    labels: Optional[Dict[str, str]] = None
    address: Optional[str] = None
    annotations: Optional[Dict[str, str]] = None
    transport: Optional[str] = None
    ready: Optional[bool] = None
    discovering: Optional[bool] = None
    status_message: Optional[str] = None
    tool_count: Optional[int] = None


class MCPServerListResponse(BaseModel):
    items: List[MCPServerResponse]
    total: int


class MCPServerDetailResponse(BaseModel):
    name: str
    namespace: str
    description: Optional[str] = None
    labels: Optional[Dict[str, str]] = None
    annotations: Optional[Dict[str, str]] = None
    spec: Optional[Dict[str, Any]] = None
    status: Optional[Dict[str, Any]] = None


class MCPTransport(BaseModel):
    type: str
    image: str
    env: Optional[Dict[str, str]] = None
    args: Optional[List[str]] = None
    command: Optional[List[str]] = None


class AddressModel(BaseModel):
    value: str

class ValueSource(BaseModel):
    """ValueSource for configuration (supports direct value or valueFrom)."""
    value: Optional[str] = None
    value_from: Optional[Dict[str, Dict[str, str]]] = Field(None, alias="valueFrom")

class Header(BaseModel):
    name:str
    value: ValueSource

class MCPServerSpec(BaseModel):
    transport: str
    description: Optional[str] = None
    tools: Optional[List[str]] = None
    address: AddressModel
    headers: Optional[List[Header]] = None


class MCPServerCreateRequest(BaseModel):
    name: str
    namespace: str
    labels: Optional[Dict[str, str]] = None
    annotations: Optional[Dict[str, str]] = None
    spec: MCPServerSpec


class MCPServerUpdateRequest(BaseModel):
    labels: Optional[Dict[str, str]] = None
    annotations: Optional[Dict[str, str]] = None
    spec: Optional[MCPServerSpec] = None