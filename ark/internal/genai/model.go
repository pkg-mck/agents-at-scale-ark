package genai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/option"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	arkv1alpha1 "mckinsey.com/ark/api/v1alpha1"
	"mckinsey.com/ark/internal/common"
	"mckinsey.com/ark/internal/telemetry"
)

const defaultModelName = "default"

func ResolveModelSpec(modelSpec any, defaultNamespace string) (string, string, error) {
	if modelSpec == nil {
		return "", "", fmt.Errorf("model spec is nil")
	}
	switch spec := modelSpec.(type) {
	case *arkv1alpha1.AgentModelRef:
		modelName := spec.Name
		namespace := spec.Namespace
		if namespace == "" {
			namespace = defaultNamespace
		}
		return modelName, namespace, nil

	case string:
		modelName := spec
		if modelName == "" {
			modelName = defaultModelName
		}
		return modelName, defaultNamespace, nil

	default:
		return "", "", fmt.Errorf("unsupported model spec type: %T", modelSpec)
	}
}

// LoadModel loads a model by resolving modelSpec and defaultNamespace
func LoadModel(ctx context.Context, k8sClient client.Client, modelSpec interface{}, defaultNamespace string, modelRecorder telemetry.ModelRecorder) (*Model, error) {
	modelName, namespace, err := ResolveModelSpec(modelSpec, defaultNamespace)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model spec: %w", err)
	}
	modelCRD, err := loadModelCRD(ctx, k8sClient, modelName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to load model CRD %s in namespace %s: %w", modelName, namespace, err)
	}

	resolver := common.NewValueSourceResolver(k8sClient)
	model, err := resolver.ResolveValueSource(ctx, modelCRD.Spec.Model, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model: %w", err)
	}

	modelInstance := &Model{
		Model:         model,
		Type:          modelCRD.Spec.Type,
		ModelRecorder: modelRecorder,
	}

	switch modelCRD.Spec.Type {
	case ModelTypeAzure:
		if err := loadAzureConfig(ctx, resolver, modelCRD.Spec.Config.Azure, namespace, modelInstance); err != nil {
			return nil, err
		}
	case ModelTypeOpenAI:
		if err := loadOpenAIConfig(ctx, resolver, modelCRD.Spec.Config.OpenAI, namespace, modelInstance); err != nil {
			return nil, err
		}
	case ModelTypeBedrock:
		if err := loadBedrockConfig(ctx, resolver, modelCRD.Spec.Config.Bedrock, namespace, model, modelInstance); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported model type: %s", modelCRD.Spec.Type)
	}

	return modelInstance, nil
}

func loadModelCRD(ctx context.Context, k8sClient client.Client, name, namespace string) (*arkv1alpha1.Model, error) {
	var modelCRD arkv1alpha1.Model
	key := types.NamespacedName{Name: name, Namespace: namespace}

	if err := k8sClient.Get(ctx, key, &modelCRD); err != nil {
		return nil, fmt.Errorf("failed to get Model %s/%s: %w", namespace, name, err)
	}

	return &modelCRD, nil
}

// resolveModelHeaders resolves custom headers from Model configuration
func resolveModelHeaders(ctx context.Context, k8sClient client.Client, headers []arkv1alpha1.Header, modelName, namespace, providerName string) (map[string]string, error) {
	if len(headers) == 0 {
		return nil, nil
	}

	log := logf.FromContext(ctx)
	resolvedHeaders := make(map[string]string)
	log.Info("resolving custom headers for model", "provider", providerName, "model", modelName, "namespace", namespace, "header_count", len(headers))

	for _, header := range headers {
		value, err := ResolveHeaderValue(ctx, k8sClient, header, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s header %s: %w", providerName, header.Name, err)
		}
		resolvedHeaders[header.Name] = value
		log.Info("resolved custom header for model", "provider", providerName, "model", modelName, "header_name", header.Name)
	}

	return resolvedHeaders, nil
}

// applyHeadersToOptions applies custom headers to OpenAI client options
func applyHeadersToOptions(ctx context.Context, headers map[string]string, options []option.RequestOption, modelName string) []option.RequestOption {
	if len(headers) == 0 {
		return options
	}

	log := logf.FromContext(ctx)
	log.Info("applying custom headers to client", "model", modelName, "header_count", len(headers))
	for name, value := range headers {
		log.V(1).Info("applying custom header", "model", modelName, "header_name", name, "header_value_length", len(value))
		options = append(options, option.WithHeader(name, value))
	}

	return options
}
