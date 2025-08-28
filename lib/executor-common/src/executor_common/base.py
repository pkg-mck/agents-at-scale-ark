"""Base classes for execution engines."""

import logging
from abc import ABC, abstractmethod
from typing import List
from .types import ExecutionEngineRequest, Message

logger = logging.getLogger(__name__)


class BaseExecutor(ABC):
    """Abstract base class for execution engines."""

    def __init__(self, engine_name: str):
        """Initialize the executor with a name."""
        self.engine_name = engine_name
        logger.info(f"{engine_name} executor initialized")

    @abstractmethod
    async def execute_agent(self, request: ExecutionEngineRequest) -> List[Message]:
        """Execute an agent with the given request.
        
        Args:
            request: The execution request containing agent config and user input
            
        Returns:
            List of response messages from the agent execution
            
        Raises:
            Exception: If execution fails
        """
        pass

    def _resolve_prompt(self, agent_config, base_prompt: str = None) -> str:
        """Resolve agent prompt with parameter substitution."""
        prompt = base_prompt or agent_config.prompt or "You are a helpful assistant."
        
        for param in agent_config.parameters:
            placeholder = f"{{{param.name}}}"
            prompt = prompt.replace(placeholder, param.value)

        return prompt