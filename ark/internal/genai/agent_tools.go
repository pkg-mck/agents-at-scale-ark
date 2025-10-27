package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	arkv1alpha1 "mckinsey.com/ark/api/v1alpha1"
	"mckinsey.com/ark/internal/telemetry"
)

// Add MCP client pool to ToolRegistry
type MCPClientPool struct {
	clients map[string]*MCPClient // key: mcpServerName
}

func NewMCPClientPool() *MCPClientPool {
	return &MCPClientPool{
		clients: make(map[string]*MCPClient),
	}
}

// GetOrCreateClient returns an existing MCP client or creates a new one for the given server
func (p *MCPClientPool) GetOrCreateClient(ctx context.Context, serverName, serverNamespace, serverURL string, headers map[string]string, transport string, timeout time.Duration, mcpSettings map[string]MCPSettings) (*MCPClient, error) {
	key := fmt.Sprintf("%s/%s", serverNamespace, serverName)
	if mcpClient, exists := p.clients[key]; exists {
		return mcpClient, nil
	}

	// Get MCP settings for this server if available
	mcpSetting := mcpSettings[key]

	// Create new client for this MCP server
	mcpClient, err := NewMCPClient(ctx, serverURL, headers, transport, timeout, mcpSetting)
	if err != nil {
		return nil, err
	}

	p.clients[key] = mcpClient
	return mcpClient, nil
}

// Close closes all MCP client connections in the pool
func (p *MCPClientPool) Close() error {
	var lastErr error
	for key, mcpClient := range p.clients {
		if mcpClient != nil && mcpClient.client != nil {
			if err := mcpClient.client.Close(); err != nil {
				lastErr = fmt.Errorf("failed to close MCP client %s: %w", key, err)
			}
		}
		delete(p.clients, key)
	}
	return lastErr
}

func (r *ToolRegistry) registerTools(ctx context.Context, k8sClient client.Client, agent *arkv1alpha1.Agent, telemetryProvider telemetry.Provider) error {
	for _, agentTool := range agent.Spec.Tools {
		if err := r.registerTool(ctx, k8sClient, agentTool, agent.Namespace, telemetryProvider); err != nil {
			return err
		}
	}
	return nil
}

func CreateToolExecutor(ctx context.Context, k8sClient client.Client, tool *arkv1alpha1.Tool, namespace string, mcpPool *MCPClientPool, mcpSettings map[string]MCPSettings, telemetryProvider telemetry.Provider) (ToolExecutor, error) {
	switch tool.Spec.Type {
	case ToolTypeHTTP:
		return createHTTPExecutor(k8sClient, tool, namespace)
	case ToolTypeMCP:
		return createMCPExecutor(ctx, k8sClient, tool, namespace, mcpPool, mcpSettings)
	case ToolTypeAgent:
		return createAgentExecutor(ctx, k8sClient, tool, namespace, telemetryProvider)
	case ToolTypeBuiltin:
		return createBuiltinExecutor(tool)
	default:
		return nil, fmt.Errorf("unsupported tool type %s for tool %s", tool.Spec.Type, tool.Name)
	}
}

func createAgentExecutor(ctx context.Context, k8sClient client.Client, tool *arkv1alpha1.Tool, namespace string, telemetryProvider telemetry.Provider) (ToolExecutor, error) {
	if tool.Spec.Agent.Name == "" {
		return nil, fmt.Errorf("agent spec is required for tool %s", tool.Name)
	}

	agentCRD := &arkv1alpha1.Agent{}
	key := types.NamespacedName{Name: tool.Spec.Agent.Name, Namespace: namespace}
	if err := k8sClient.Get(ctx, key, agentCRD); err != nil {
		return nil, fmt.Errorf("failed to get agent %v: %w", key, err)
	}

	return &AgentToolExecutor{
		AgentName:         tool.Spec.Agent.Name,
		Namespace:         namespace,
		AgentCRD:          agentCRD,
		k8sClient:         k8sClient,
		telemetryProvider: telemetryProvider,
	}, nil
}

func createBuiltinExecutor(tool *arkv1alpha1.Tool) (ToolExecutor, error) {
	switch tool.Name {
	case BuiltinToolNoop:
		return &NoopExecutor{}, nil
	case BuiltinToolTerminate:
		return &TerminateExecutor{}, nil
	default:
		return nil, fmt.Errorf("unsupported builtin tool %s", tool.Name)
	}
}

func createHTTPExecutor(k8sClient client.Client, tool *arkv1alpha1.Tool, namespace string) (ToolExecutor, error) {
	if tool.Spec.HTTP == nil {
		return nil, fmt.Errorf("http spec is required for tool %s", tool.Name)
	}
	return &HTTPExecutor{
		K8sClient:     k8sClient,
		ToolName:      tool.Name,
		ToolNamespace: namespace,
	}, nil
}

func createMCPExecutor(ctx context.Context, k8sClient client.Client, tool *arkv1alpha1.Tool, namespace string, mcpPool *MCPClientPool, mcpSettings map[string]MCPSettings) (ToolExecutor, error) {
	if tool.Spec.MCP == nil {
		return nil, fmt.Errorf("mcp spec is required for tool %s", tool.Name)
	}

	mcpServerNamespace := tool.Spec.MCP.MCPServerRef.Namespace
	if mcpServerNamespace == "" {
		mcpServerNamespace = namespace
	}

	var mcpServerCRD arkv1alpha1.MCPServer
	mcpServerKey := types.NamespacedName{
		Name:      tool.Spec.MCP.MCPServerRef.Name,
		Namespace: mcpServerNamespace,
	}
	if err := k8sClient.Get(ctx, mcpServerKey, &mcpServerCRD); err != nil {
		return nil, fmt.Errorf("failed to get MCP server %v: %w", mcpServerKey, err)
	}

	mcpURL, err := BuildMCPServerURL(ctx, k8sClient, &mcpServerCRD)
	if err != nil {
		return nil, fmt.Errorf("failed to build MCP server URL: %w", err)
	}

	headers := make(map[string]string)
	for _, header := range mcpServerCRD.Spec.Headers {
		value, err := ResolveHeaderValue(ctx, k8sClient, header, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve header %s: %w", header.Name, err)
		}
		headers[header.Name] = value
	}

	// Parse timeout from MCPServer spec (default to 30s if not specified)
	timeout := 30 * time.Second
	if mcpServerCRD.Spec.Timeout != "" {
		parsedTimeout, err := time.ParseDuration(mcpServerCRD.Spec.Timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timeout %s: %w", mcpServerCRD.Spec.Timeout, err)
		}
		timeout = parsedTimeout
	}

	// Use the MCP client pool to get or create the client
	mcpClient, err := mcpPool.GetOrCreateClient(
		ctx,
		tool.Spec.MCP.MCPServerRef.Name,
		mcpServerNamespace,
		mcpURL,
		headers,
		mcpServerCRD.Spec.Transport,
		timeout,
		mcpSettings,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create MCP client for tool %s: %w", tool.Name, err)
	}

	return &MCPExecutor{
		ToolName:  tool.Spec.MCP.ToolName,
		MCPClient: mcpClient,
	}, nil
}

func (r *ToolRegistry) registerTool(ctx context.Context, k8sClient client.Client, agentTool arkv1alpha1.AgentTool, namespace string, telemetryProvider telemetry.Provider) error {
	tool := &arkv1alpha1.Tool{}
	key := client.ObjectKey{Name: agentTool.Name, Namespace: namespace}

	if err := k8sClient.Get(ctx, key, tool); err != nil {
		return fmt.Errorf("failed to get tool %s: %w", agentTool.Name, err)
	}

	toolDef := CreateToolFromCRD(tool)
	executor, err := CreateToolExecutor(ctx, k8sClient, tool, namespace, r.mcpPool, r.mcpSettings, telemetryProvider)
	if err != nil {
		return fmt.Errorf("failed to create executor for tool %s: %w", agentTool.Name, err)
	}

	if agentTool.Partial != nil {
		var err error
		toolDef, err = CreatePartialToolDefinition(toolDef, agentTool.Partial)
		if err != nil {
			return fmt.Errorf("failed to create partial tool definition for tool %s: %w", agentTool.Name, err)
		}
		// Wrap with PartialToolExecutor if partial is specified
		executor = &PartialToolExecutor{
			BaseExecutor: executor,
			Partial:      agentTool.Partial,
		}
	}

	// Apply function filtering if specified
	if len(agentTool.Functions) > 0 {
		executor = &FilteredToolExecutor{
			BaseExecutor: executor,
			Functions:    agentTool.Functions,
		}
	}

	r.RegisterTool(toolDef, executor)
	return nil
}

// AgentToolExecutor executes agent tools by calling other agents via MCP
type AgentToolExecutor struct {
	AgentName         string
	Namespace         string
	AgentCRD          *arkv1alpha1.Agent
	k8sClient         client.Client
	telemetryProvider telemetry.Provider
}

func (a *AgentToolExecutor) Execute(ctx context.Context, call ToolCall, recorder EventEmitter) (ToolResult, error) {
	var arguments map[string]any
	if err := json.Unmarshal([]byte(call.Function.Arguments), &arguments); err != nil {
		log := logf.FromContext(ctx)
		log.Error(err, "Error parsing tool arguments", "ToolCall")
		return ToolResult{
			ID:    call.ID,
			Name:  call.Function.Name,
			Error: "Failed to parse tool arguments",
		}, fmt.Errorf("failed to parse tool arguments: %v", err)
	}

	input, exists := arguments["input"]
	if !exists {
		return ToolResult{
			ID:    call.ID,
			Name:  call.Function.Name,
			Error: "input parameter is required",
		}, fmt.Errorf("input parameter is required for agent tool %s", a.AgentName)
	}

	inputStr, ok := input.(string)
	if !ok {
		return ToolResult{
			ID:    call.ID,
			Name:  call.Function.Name,
			Error: "input parameter must be a string",
		}, fmt.Errorf("input parameter must be a string for agent tool %s", a.AgentName)
	}

	// Log the agent execution
	log := logf.FromContext(ctx)
	log.Info("calling agent directly", "agent", a.AgentName, "namespace", a.Namespace, "input", inputStr)

	// Create the Agent object using the Agent CRD and recorder
	agent, err := MakeAgent(ctx, a.k8sClient, a.AgentCRD, recorder, a.telemetryProvider)
	if err != nil {
		return ToolResult{
			ID:    call.ID,
			Name:  call.Function.Name,
			Error: fmt.Sprintf("failed to create agent %s: %v", a.AgentName, err),
		}, err
	}

	// Prepare user input and history
	userInput := NewSystemMessage(inputStr)
	history := []Message{} // Provide history if applicable

	// Call the agent's Execute function
	// Pass nil for memory and eventStream (agents-as-tools don't use memory or streaming)
	// See ARKQB-137 for discussion on streaming support for agents as tools
	responseMessages, err := agent.Execute(ctx, userInput, history, nil, nil)
	if err != nil {
		log.Info("agent execution error", "agent", a.AgentName, "error", err)
		return ToolResult{
			ID:    call.ID,
			Name:  call.Function.Name,
			Error: fmt.Sprintf("failed to execute agent %s: %v", a.AgentName, err),
		}, err
	}

	lastMessage := responseMessages[len(responseMessages)-1]

	log.Info("agent direct call response", "agent", a.AgentName, "response", lastMessage.OfAssistant.Content.OfString.Value)

	return ToolResult{
		ID:      call.ID,
		Name:    call.Function.Name,
		Content: lastMessage.OfAssistant.Content.OfString.Value,
	}, nil
}
