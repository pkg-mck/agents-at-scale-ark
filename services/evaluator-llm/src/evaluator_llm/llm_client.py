import aiohttp
import json
import logging
from typing import Dict, Any
from .types import Model

logger = logging.getLogger(__name__)


class LLMClient:
    def __init__(self, session=None):
        self.session = session
        self._session_lock = None
        self._owns_session = session is None
    
    async def _get_session(self):
        if self.session is None:
            if self._session_lock is None:
                import asyncio
                self._session_lock = asyncio.Lock()
            
            async with self._session_lock:
                if self.session is None:
                    self.session = aiohttp.ClientSession()
        return self.session
    
    async def close(self):
        if self.session and self._owns_session:
            await self.session.close()
            self.session = None
    
    async def evaluate(self, prompt: str, model: Model) -> str:
        """
        Call the model to evaluate the prompt
        """
        model_type = model.type
        
        if model_type == 'openai':
            return await self._call_openai_compatible(prompt, model)
        elif model_type == 'azure':
            return await self._call_azure_openai(prompt, model)
        else:
            raise ValueError(f"Unsupported model type: {model_type}")
    
    async def _call_openai_compatible(self, prompt: str, model: Model) -> str:
        """
        Call OpenAI-compatible API
        """
        session = await self._get_session()
        
        # Extract model configuration
        url = model.config.get('base_url', '')
        model_name = model.name
        api_key = model.config.get('api_key', '')
        
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {api_key}'
        }
        
        payload = {
            'model': model_name,
            'messages': [
                {'role': 'user', 'content': prompt}
            ],
            'temperature': 0.1,  # Low temperature for consistent evaluation
            'max_tokens': 1000
        }
        
        logger.info(f"Calling OpenAI-compatible API: {url}")
        
        async with session.post(url, headers=headers, json=payload) as response:
            if response.status != 200:
                error_text = await response.text()
                raise Exception(f"API call failed with status {response.status}: {error_text}")
            
            result = await response.json()
            return result['choices'][0]['message']['content']
    
    async def _call_azure_openai(self, prompt: str, model: Model) -> str:
        """
        Call Azure OpenAI API
        """
        session = await self._get_session()
        
        # Extract Azure-specific configuration
        url = model.config.get('base_url', '')
        api_key = model.config.get('api_key', '')
        deployment_name = model.name
        api_version = model.config.get('api_version', '')
        
        # Build Azure OpenAI URL
        full_url = f"{url}/openai/deployments/{deployment_name}/chat/completions?api-version={api_version}"
        
        headers = {
            'Content-Type': 'application/json',
            'api-key': api_key
        }
        
        payload = {
            'messages': [
                {'role': 'user', 'content': prompt}
            ],
            'temperature': 0.1,  # Low temperature for consistent evaluation
            'max_tokens': 1000
        }
        
        logger.info(f"Calling Azure OpenAI API: {full_url}")
        
        async with session.post(full_url, headers=headers, json=payload) as response:
            if response.status != 200:
                error_text = await response.text()
                raise Exception(f"Azure API call failed with status {response.status}: {error_text}")
            
            result = await response.json()
            return result['choices'][0]['message']['content']