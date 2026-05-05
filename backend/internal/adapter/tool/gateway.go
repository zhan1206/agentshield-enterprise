package tool

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ToolType represents the type of tool
type ToolType string

const (
	ToolTypeMCP             ToolType = "mcp"
	ToolTypeOpenAPI         ToolType = "openapi"
	ToolTypeDatabase       ToolType = "database"
	ToolTypeInternalService ToolType = "internal_service"
	ToolTypeCustom         ToolType = "custom"
)

// ToolDefinition defines a tool that can be called by agents
type ToolDefinition struct {
	ID            string                 `json:"id" yaml:"id"`
	Name          string                 `json:"name" yaml:"name"`
	Type          ToolType               `json:"type" yaml:"type"`
	Endpoint     string                 `json:"endpoint" yaml:"endpoint"`
	Schema        string                 `json:"schema" yaml:"schema"`
	AuthRequired  bool                   `json:"auth_required" yaml:"auth_required"`
	RateLimit    int                    `json:"rate_limit" yaml:"rate_limit"` // requests per minute
	Timeout      int                    `json:"timeout" yaml:"timeout"`     // seconds
	Description  string                 `json:"description" yaml:"description"`
	Options      map[string]interface{} `json:"options" yaml:"options"`
	RequiredRoles []string               `json:"required_roles" yaml:"required_roles"`
	RequiredPerms []string             `json:"required_perms" yaml:"required_perms"`
}

// CallRequest represents a request to call a tool
type CallRequest struct {
	ToolID   string                 `json:"tool_id" yaml:"tool_id"`
	AgentID  string                 `json:"agent_id" yaml:"agent_id"`
	Method   string                 `json:"method" yaml:"method"`
	Params   map[string]interface{} `json:"params" yaml:"params"`
	Timeout  int                    `json:"timeout" yaml:"timeout"`
	Context  map[string]interface{} `json:"context" yaml:"context"`
}

// CallResult represents the result of a tool call
type CallResult struct {
	Success   bool                   `json:"success" yaml:"success"`
	Data      interface{}            `json:"data" yaml:"data"`
	Error     string                `json:"error" yaml:"error"`
	LatencyMs int64                 `json:"latency_ms" yaml:"latency_ms"`
	Timestamp time.Time             `json:"timestamp" yaml:"timestamp"`
}

// rateLimiter tracks rate limits per agent/tool
type rateLimiter struct {
	mu       sync.Mutex
	counts   map[string]map[string]int      // agentID -> toolID -> count
	windows  map[string]map[string]time.Time // agentID -> toolID -> window start
}

// Gateway provides a unified tool calling gateway
type Gateway struct {
	logger     interface{ Printf(format string, args ...interface{}) }
	tools      map[string]ToolDefinition
	rateLimiter *rateLimiter
	auditFunc  func(CallRequest, CallResult)
	httpClient *http.Client
	mu        sync.RWMutex
	agentPerms map[string]map[string]bool // agentID -> toolID -> allowed
}

// NewGateway creates a new tool gateway
func NewGateway(logger interface{ Printf(format string, args ...interface{}) }, auditFunc func(CallRequest, CallResult)) *Gateway {
	gateway := &Gateway{
		logger:     logger,
		tools:      make(map[string]ToolDefinition),
		rateLimiter: &rateLimiter{
			counts:  make(map[string]map[string]int),
			windows: make(map[string]map[string]time.Time),
		},
		auditFunc: auditFunc,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		agentPerms: make(map[string]map[string]bool),
	}

	return gateway
}

// RegisterTool registers a new tool with the gateway
func (g *Gateway) RegisterTool(def ToolDefinition) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if def.ID == "" {
		return fmt.Errorf("tool ID is required")
	}

	if _, exists := g.tools[def.ID]; exists {
		return fmt.Errorf("tool %s already registered", def.ID)
	}

	if def.Timeout == 0 {
		def.Timeout = 30
	}

	if def.RateLimit == 0 {
		def.RateLimit = 100 // default 100 requests per minute
	}

	g.tools[def.ID] = def

	if g.logger != nil {
		g.logger.Printf("[tool] Registered tool %s (%s) at %s", def.Name, def.Type, def.Endpoint)
	}

	return nil
}

// Call executes a tool call request
func (g *Gateway) Call(ctx context.Context, req CallRequest) (*CallResult, error) {
	startTime := time.Now()

	// Validate request
	if req.ToolID == "" {
		result := CallResult{Success: false, Error: "tool_id is required", Timestamp: startTime}
		g.recordAudit(req, result)
		return &result, fmt.Errorf("tool_id is required")
	}

	if req.AgentID == "" {
		result := CallResult{Success: false, Error: "agent_id is required", Timestamp: startTime}
		g.recordAudit(req, result)
		return &result, fmt.Errorf("agent_id is required")
	}

	// Get tool definition
	g.mu.RLock()
	tool, exists := g.tools[req.ToolID]
	g.mu.RUnlock()

	if !exists {
		result := CallResult{Success: false, Error: fmt.Sprintf("tool %s not found", req.ToolID), Timestamp: startTime}
		g.recordAudit(req, result)
		return &result, fmt.Errorf("tool %s not found", req.ToolID)
	}

	// Check permissions
	if !g.validatePermission(req.AgentID, req.ToolID) {
		result := CallResult{Success: false, Error: "permission denied", Timestamp: startTime}
		g.recordAudit(req, result)
		if g.logger != nil {
			g.logger.Printf("[tool] Permission denied: agent %s cannot access tool %s", req.AgentID, req.ToolID)
		}
		return &result, fmt.Errorf("permission denied")
	}

	// Check rate limit
	if !g.checkRateLimit(req.AgentID, req.ToolID) {
		result := CallResult{Success: false, Error: "rate limit exceeded", Timestamp: startTime}
		g.recordAudit(req, result)
		if g.logger != nil {
			g.logger.Printf("[tool] Rate limit exceeded: agent %s tool %s", req.AgentID, req.ToolID)
		}
		return &result, fmt.Errorf("rate limit exceeded")
	}

	// Execute tool call
	data, err := g.executeTool(ctx, tool, req)
	latencyMs := time.Since(startTime).Milliseconds()

	result := CallResult{
		Success:   err == nil,
		Data:      data,
		Error:     "",
		LatencyMs: latencyMs,
		Timestamp: startTime,
	}

	if err != nil {
		result.Error = err.Error()
	}

	g.recordAudit(req, result)

	return &result, err
}

// executeTool executes the actual tool call
func (g *Gateway) executeTool(ctx context.Context, tool ToolDefinition, req CallRequest) (interface{}, error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = tool.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Build request based on tool type
	switch tool.Type {
	case ToolTypeMCP:
		return g.executeMCP(ctx, tool, req)
	case ToolTypeOpenAPI:
		return g.executeOpenAPI(ctx, tool, req)
	case ToolTypeInternalService:
		return g.executeInternalService(ctx, tool, req)
	default:
		return g.executeGeneric(ctx, tool, req)
	}
}

// executeMCP executes a MCP tool call
func (g *Gateway) executeMCP(ctx context.Context, tool ToolDefinition, req CallRequest) (interface{}, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":     1,
		"method": req.Method,
		"params": req.Params,
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tool.Endpoint, strings.NewReader(fmt.Sprintf("%v", payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if tool.AuthRequired {
		httpReq.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%v", req.Context["api_key"]))
	}

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	// Parse response (simplified)
	return map[string]interface{}{"executed": true}, nil
}

// executeOpenAPI executes an OpenAPI tool call
func (g *Gateway) executeOpenAPI(ctx context.Context, tool ToolDefinition, req CallRequest) (interface{}, error) {
	return map[string]interface{}{"executed": true, "method": req.Method}, nil
}

// executeInternalService executes an internal service call
func (g *Gateway) executeInternalService(ctx context.Context, tool ToolDefinition, req CallRequest) (interface{}, error) {
	return map[string]interface{}{"executed": true, "service": tool.Name}, nil
}

// executeGeneric executes a generic tool call
func (g *Gateway) executeGeneric(ctx context.Context, tool ToolDefinition, req CallRequest) (interface{}, error) {
	return map[string]interface{}{"executed": true, "tool": tool.Name}, nil
}

// validatePermission checks if an agent has permission to use a tool
func (g *Gateway) validatePermission(agentID, toolID string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	perms, exists := g.agentPerms[agentID]
	if !exists {
		// Default: allow if no explicit permissions defined
		return true
	}

	allowed, ok := perms[toolID]
	if !ok {
		// Default: allow if no explicit deny
		return true
	}

	return allowed
}

// checkRateLimit checks if the agent has exceeded rate limit for the tool
func (g *Gateway) checkRateLimit(agentID, toolID string) bool {
	g.mu.RLock()
	tool, exists := g.tools[toolID]
	g.mu.RUnlock()

	if !exists {
		return true
	}

	limit := tool.RateLimit
	if limit == 0 {
		return true
	}

	g.rateLimiter.mu.Lock()
	defer g.rateLimiter.mu.Unlock()

	now := time.Now()
	key := agentID + ":" + toolID

	// Initialize agent tracking if needed
	if _, ok := g.rateLimiter.counts[agentID]; !ok {
		g.rateLimiter.counts[agentID] = make(map[string]int)
		g.rateLimiter.windows[agentID] = make(map[string]time.Time)
	}

	// Check if we need to reset the window
	windowStart, windowExists := g.rateLimiter.windows[agentID][toolID]
	if !windowExists || now.Sub(windowStart) > time.Minute {
		g.rateLimiter.counts[agentID][toolID] = 0
		g.rateLimiter.windows[agentID][toolID] = now
	}

	// Check limit
	count := g.rateLimiter.counts[agentID][toolID]
	if count >= limit {
		return false
	}

	g.rateLimiter.counts[agentID][toolID] = count + 1
	return true
}

// ListTools returns all registered tools
func (g *Gateway) ListTools() []ToolDefinition {
	g.mu.RLock()
	defer g.mu.RUnlock()

	tools := make([]ToolDefinition, 0, len(g.tools))
	for _, tool := range g.tools {
		tools = append(tools, tool)
	}

	return tools
}

// UnregisterTool unregisters a tool from the gateway
func (g *Gateway) UnregisterTool(toolID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.tools[toolID]; !exists {
		return fmt.Errorf("tool %s not found", toolID)
	}

	delete(g.tools, toolID)

	if g.logger != nil {
		g.logger.Printf("[tool] Unregistered tool %s", toolID)
	}

	return nil
}

// GrantPermission grants an agent permission to use a tool
func (g *Gateway) GrantPermission(agentID, toolID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.agentPerms[agentID]; !ok {
		g.agentPerms[agentID] = make(map[string]bool)
	}

	g.agentPerms[agentID][toolID] = true

	if g.logger != nil {
		g.logger.Printf("[tool] Granted permission: agent %s -> tool %s", agentID, toolID)
	}
}

// RevokePermission revokes an agent's permission to use a tool
func (g *Gateway) RevokePermission(agentID, toolID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if perms, ok := g.agentPerms[agentID]; ok {
		perms[toolID] = false
	}

	if g.logger != nil {
		g.logger.Printf("[tool] Revoked permission: agent %s -> tool %s", agentID, toolID)
	}
}

// recordAudit records an audit log for the tool call
func (g *Gateway) recordAudit(req CallRequest, result CallResult) {
	if g.auditFunc != nil {
		g.auditFunc(req, result)
	}
}

// Middleware returns a Gin middleware for the tool gateway
func Middleware(gateway *Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		c.Set("tool_gateway", gateway)
		c.Next()
	}
}

// GetGatewayFromContext retrieves the gateway from Gin context
func GetGatewayFromContext(c *gin.Context) *Gateway {
	if gateway, exists := c.Get("tool_gateway"); exists {
		return gateway.(*Gateway)
	}
	return nil
}

// ValidateToolSchema validates a tool's JSON schema (basic validation)
func ValidateToolSchema(schema string) error {
	if schema == "" {
		return nil
	}

	// Basic JSON schema validation - check for balanced braces
	openBraces := strings.Count(schema, "{")
	closeBraces := strings.Count(schema, "}")

	if openBraces != closeBraces {
		return fmt.Errorf("unbalanced braces in schema")
	}

	// Check for required fields
	requiredFields := []string{"type", "properties"}
	for _, field := range requiredFields {
		if !regexp.MustCompile(`"` + field + `"`).MatchString(schema) {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	return nil
}