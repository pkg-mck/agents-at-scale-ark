import logging
from typing import Dict, Any
from .types import EvaluationRequest, EvaluationResponse
from .llm_client import LLMClient

logger = logging.getLogger(__name__)


class LLMEvaluator:
    def __init__(self, session=None):
        self.llm_client = LLMClient(session=session)
    
    async def evaluate(self, request: EvaluationRequest) -> EvaluationResponse:
        """
        Evaluate query performance using LLM-as-a-Judge approach
        """
        try:
            logger.info(f"Starting evaluation for query {request.queryId}")
            
            # Prepare evaluation prompt
            evaluation_prompt = self._build_evaluation_prompt(request)
            
            # Get LLM evaluation
            evaluation_result = await self.llm_client.evaluate(
                prompt=evaluation_prompt,
                model=request.model
            )
            
            # Parse evaluation result
            score, passed, metadata = self._parse_evaluation_result(evaluation_result)
            
            logger.info(f"Evaluation completed for query {request.queryId}: score={score}, passed={passed}")
            
            return EvaluationResponse(
                score=score,
                passed=passed,
                metadata=metadata
            )
            
        except Exception as e:
            logger.error(f"Evaluation failed for query {request.queryId}: {str(e)}")
            return EvaluationResponse(
                error=str(e),
                passed=False
            )
    
    def _build_evaluation_prompt(self, request: EvaluationRequest) -> str:
        """
        Build evaluation prompt using LLM-as-a-Judge pattern
        """
        responses_text = "\n".join([
            f"Response from {resp.target.type} '{resp.target.name}':\n{resp.content}"
            for resp in request.responses
        ])
        
        prompt = f"""You are an AI evaluator tasked with assessing the quality of responses to user queries.

USER QUERY:
{request.input}

RESPONSES TO EVALUATE:
{responses_text}

Please evaluate the responses based on the following criteria:
1. Relevance: How well do the responses address the user's query?
2. Accuracy: Are the responses factually correct and reliable?
3. Completeness: Do the responses provide comprehensive information?
4. Clarity: Are the responses clear and easy to understand?
5. Usefulness: How helpful are the responses to the user?

Provide your evaluation in the following format:
SCORE: [0-100]
PASSED: [true/false] (true if score >= 70)
REASONING: [Brief explanation of your evaluation]
CRITERIA_SCORES: relevance=[0-100], accuracy=[0-100], completeness=[0-100], clarity=[0-100], usefulness=[0-100]

Be objective and thorough in your assessment."""

        return prompt
    
    def _parse_evaluation_result(self, result: str) -> tuple[str, bool, Dict[str, str]]:
        """
        Parse LLM evaluation result into structured format
        """
        lines = result.strip().split('\n')
        score = "0"
        passed = False
        metadata = {}
        
        for line in lines:
            line = line.strip()
            if line.startswith('SCORE:'):
                score = line.split(':', 1)[1].strip()
                try:
                    score_int = int(score)
                    passed = score_int >= 70
                except ValueError:
                    score = "0"
                    passed = False
            elif line.startswith('PASSED:'):
                passed_str = line.split(':', 1)[1].strip().lower()
                passed = passed_str == 'true'
            elif line.startswith('REASONING:'):
                metadata['reasoning'] = line.split(':', 1)[1].strip()
            elif line.startswith('CRITERIA_SCORES:'):
                criteria_str = line.split(':', 1)[1].strip()
                metadata['criteria_scores'] = criteria_str
        
        return score, passed, metadata