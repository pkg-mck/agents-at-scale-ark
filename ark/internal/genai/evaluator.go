/* Copyright 2025. McKinsey & Company */

package genai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	arkv1alpha1 "mckinsey.com/ark/api/v1alpha1"
	"mckinsey.com/ark/internal/common"
)

type EvaluationRequest struct {
	QueryID   string                 `json:"queryId"`
	Input     string                 `json:"input"`
	Responses []arkv1alpha1.Response `json:"responses"`
	Query     arkv1alpha1.Query      `json:"query"`
	Model     any                    `json:"model,omitempty"`
}

type EvaluationResponse struct {
	Score    string            `json:"score,omitempty"`
	Passed   bool              `json:"passed,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Error    string            `json:"error,omitempty"`
}

func CallSingleEvaluator(ctx context.Context, k8sClient client.Client, query arkv1alpha1.Query, evaluatorRef arkv1alpha1.EvaluatorRef, recorder EventEmitter) (*arkv1alpha1.EvaluationResult, error) {
	tracker := NewOperationTracker(recorder, ctx, "Evaluation", query.Name, map[string]string{
		"namespace": query.Namespace,
		"evaluator": evaluatorRef.Name,
	})

	evaluator, err := loadEvaluator(ctx, k8sClient, evaluatorRef, query.Namespace)
	if err != nil {
		tracker.Fail(err)
		return nil, err
	}

	address, err := resolveEvaluatorAddress(ctx, k8sClient, evaluator)
	if err != nil {
		tracker.Fail(err)
		return nil, err
	}

	model, err := loadEvaluatorModel(ctx, k8sClient, evaluator)
	if err != nil {
		tracker.Fail(err)
		return nil, err
	}

	request := buildEvaluationRequest(query, convertModelConfig(model))
	response, err := callEvaluatorHTTP(ctx, address, request)
	if err != nil {
		tracker.Fail(err)
		return nil, err
	}

	result := &arkv1alpha1.EvaluationResult{
		Score:    response.Score,
		Passed:   response.Passed,
		Metadata: response.Metadata,
	}

	tracker.Complete(fmt.Sprintf("score: %s, passed: %t", response.Score, response.Passed))
	return result, nil
}

func CallEvaluators(ctx context.Context, k8sClient client.Client, query arkv1alpha1.Query, evaluatorRefs []arkv1alpha1.EvaluatorRef, recorder EventEmitter) ([]arkv1alpha1.EvaluationResult, error) {
	if len(evaluatorRefs) == 0 {
		return nil, nil
	}

	results := make([]arkv1alpha1.EvaluationResult, len(evaluatorRefs))
	var wg sync.WaitGroup

	for i, evaluatorRef := range evaluatorRefs {
		wg.Add(1)
		go func(idx int, evalRef arkv1alpha1.EvaluatorRef) {
			defer wg.Done()
			results[idx] = callEvaluatorWithErrorHandling(ctx, k8sClient, query, evalRef, recorder)
		}(i, evaluatorRef)
	}

	wg.Wait()
	return results, nil
}

func loadEvaluator(ctx context.Context, k8sClient client.Client, evaluatorRef arkv1alpha1.EvaluatorRef, defaultNamespace string) (*arkv1alpha1.Evaluator, error) {
	namespace := evaluatorRef.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	var evaluator arkv1alpha1.Evaluator
	key := types.NamespacedName{Name: evaluatorRef.Name, Namespace: namespace}

	if err := k8sClient.Get(ctx, key, &evaluator); err != nil {
		return nil, fmt.Errorf("failed to get evaluator %s: %w", evaluatorRef.Name, err)
	}

	return &evaluator, nil
}

func resolveEvaluatorAddress(ctx context.Context, k8sClient client.Client, evaluator *arkv1alpha1.Evaluator) (string, error) {
	resolver := common.NewValueSourceResolver(k8sClient)
	address, err := resolver.ResolveValueSource(ctx, evaluator.Spec.Address, evaluator.Namespace)
	if err != nil {
		return "", fmt.Errorf("failed to resolve evaluator address: %w", err)
	}
	return address, nil
}

func loadEvaluatorModel(ctx context.Context, k8sClient client.Client, evaluator *arkv1alpha1.Evaluator) (*Model, error) {
	modelName := "default"
	modelNamespace := evaluator.Namespace

	// Check parameters for model configuration
	for _, param := range evaluator.Spec.Parameters {
		switch param.Name {
		case "model.name":
			if param.Value != "" {
				modelName = param.Value
			}
		case "model.namespace":
			if param.Value != "" {
				modelNamespace = param.Value
			}
		}
	}

	// Create a temporary struct for LoadModel compatibility
	modelRef := &struct {
		Name      string
		Namespace string
	}{
		Name:      modelName,
		Namespace: modelNamespace,
	}

	model, err := LoadModel(ctx, k8sClient, modelRef, evaluator.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to load model %s: %w", modelName, err)
	}

	return model, nil
}

func convertModelConfig(model *Model) map[string]any {
	config := map[string]any{}

	switch model.Type {
	case ModelTypeAzure:
		if azureProvider, ok := model.Provider.(*AzureProvider); ok {
			config["base_url"] = azureProvider.BaseURL
			config["api_key"] = azureProvider.APIKey
			config["api_version"] = azureProvider.APIVersion
		}
	case ModelTypeOpenAI:
		if openaiProvider, ok := model.Provider.(*OpenAIProvider); ok {
			config["base_url"] = openaiProvider.BaseURL
			config["api_key"] = openaiProvider.APIKey
		}
	case ModelTypeBedrock:
		if bedrockProvider, ok := model.Provider.(*BedrockModel); ok {
			config["region"] = bedrockProvider.Region
			config["access_key_id"] = bedrockProvider.AccessKeyID
			config["secret_access_key"] = bedrockProvider.SecretAccessKey
			config["session_token"] = bedrockProvider.SessionToken
			config["model_arn"] = bedrockProvider.ModelArn
			if maxTokens, exists := bedrockProvider.Properties["max_tokens"]; exists {
				config["max_tokens"] = maxTokens
			}
			if temperature, exists := bedrockProvider.Properties["temperature"]; exists {
				config["temperature"] = temperature
			}
		}
	}

	return map[string]any{
		"name":   model.Model,
		"type":   model.Type,
		"config": config,
	}
}

func buildEvaluationRequest(query arkv1alpha1.Query, model map[string]any) EvaluationRequest {
	return EvaluationRequest{
		QueryID:   string(query.UID),
		Input:     query.Spec.Input,
		Responses: query.Status.Responses,
		Query:     query,
		Model:     model,
	}
}

func callEvaluatorHTTP(ctx context.Context, address string, request EvaluationRequest) (*EvaluationResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal evaluation request: %w", err)
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call evaluator: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logf.Log.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("evaluator returned status %d", resp.StatusCode)
	}

	var response EvaluationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode evaluation response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("evaluator returned error: %s", response.Error)
	}

	return &response, nil
}

func callEvaluatorWithErrorHandling(ctx context.Context, k8sClient client.Client, query arkv1alpha1.Query, evaluatorRef arkv1alpha1.EvaluatorRef, recorder EventEmitter) arkv1alpha1.EvaluationResult {
	result, err := CallSingleEvaluator(ctx, k8sClient, query, evaluatorRef, recorder)
	if err != nil {
		return arkv1alpha1.EvaluationResult{
			EvaluatorName: evaluatorRef.Name,
			Score:         "0",
			Passed:        false,
			Metadata:      map[string]string{"error": err.Error()},
		}
	}

	if result != nil {
		result.EvaluatorName = evaluatorRef.Name
		return *result
	}

	return arkv1alpha1.EvaluationResult{
		EvaluatorName: evaluatorRef.Name,
		Score:         "0",
		Passed:        false,
		Metadata:      map[string]string{"error": "no result returned"},
	}
}
