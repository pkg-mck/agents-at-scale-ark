/* Copyright 2025. McKinsey & Company */

package genai

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	arkv1prealpha1 "mckinsey.com/ark/api/v1prealpha1"
)

// A2AExecutionEngine handles execution for agents with the reserved 'a2a' execution engine
type A2AExecutionEngine struct {
	client client.Client
}

// NewA2AExecutionEngine creates a new A2A execution engine
func NewA2AExecutionEngine(k8sClient client.Client) *A2AExecutionEngine {
	return &A2AExecutionEngine{
		client: k8sClient,
	}
}

// Execute executes a query against an A2A agent
func (e *A2AExecutionEngine) Execute(ctx context.Context, agentName, namespace string, annotations map[string]string, userInput Message) ([]Message, error) {
	log := logf.FromContext(ctx)
	log.Info("executing A2A agent", "agent", agentName)

	// Get the A2A server address from annotations
	a2aAddress, hasAddress := annotations["a2a.server/address"]
	if !hasAddress {
		return nil, fmt.Errorf("A2A agent missing a2a.server/address annotation")
	}

	// Get the A2AServer name from annotations
	a2aServerName, hasServerName := annotations["a2a.server/name"]
	if !hasServerName {
		return nil, fmt.Errorf("A2A agent missing a2a.server/name annotation")
	}

	var a2aServer arkv1prealpha1.A2AServer
	serverKey := client.ObjectKey{Name: a2aServerName, Namespace: namespace}
	if err := e.client.Get(ctx, serverKey, &a2aServer); err != nil {
		return nil, fmt.Errorf("unable to get A2AServer %v: %w", serverKey, err)
	}

	// Extract content from the userInput message
	content := ""
	if userInput.OfUser != nil && userInput.OfUser.Content.OfString.Value != "" {
		content = userInput.OfUser.Content.OfString.Value
	}

	// Execute A2A agent
	response, err := ExecuteA2AAgent(ctx, e.client, a2aAddress, a2aServer.Spec.Headers, namespace, content, agentName)
	if err != nil {
		return nil, fmt.Errorf("A2A agent execution failed: %w", err)
	}

	log.Info("A2A agent execution completed", "agent", agentName, "response_length", len(response))

	// Convert response to genai.Message format
	responseMessage := NewAssistantMessage(response)
	return []Message{responseMessage}, nil
}
