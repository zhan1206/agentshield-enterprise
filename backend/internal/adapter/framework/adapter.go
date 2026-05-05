package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// FrameworkType represents the type of agent framework
type FrameworkType string

const (
	FrameworkOpenHands FrameworkType = "openhands"
	FrameworkPlandex  FrameworkType = "plandex"
	FrameworkLangGraph FrameworkType = "langgraph"
	FrameworkCrewAI   FrameworkType = "crewai"
	FrameworkAutoGen  FrameworkType = "autogen"
	FrameworkCustom   FrameworkType = "custom"
)

// AdapterConfig contains configuration for the framework adapter
type AdapterConfig struct {
	Framework  FrameworkType          `json:"framework" yaml:"framework"`
	Endpoint   string                 `json:"endpoint" yaml:"endpoint"`
	APIKey     string                 `json:"api_key" yaml:"api_key"`
	Options    map[string]interface{} `json:"options" yaml:"options"`
	Timeout    int                    `json:"timeout" yaml:"timeout"`
	MaxRetries int                    `json:"max_retries" yaml:"max_retries"`
}

// AgentStatus represents the current status of a registered agent
type AgentStatus struct {
	ID            string    `json:"id" yaml:"id"`
	Name          string    `json:"name" yaml:"name"`
	Status        string    `json:"status" yaml:"status"`
	LastActivity time.Time `json:"last_activity" yaml:"last_activity"`
	Capabilities []string   `json:"capabilities" yaml:"capabilities"`
	ErrorMessage string    `json:"error_message,omitempty" yaml:"error_message,omitempty"`
}

// Adapter provides a unified interface to different agent frameworks
type Adapter struct {
	config     AdapterConfig
	httpClient *http.Client
	connected bool
	agents    map[string]*AgentStatus
	mu        sync.RWMutex
	logger    interface{ Printf(format string, args ...interface{}) }
}

// NewAdapter creates a new framework adapter with the given configuration
func NewAdapter(config AdapterConfig, logger interface{ Printf(format string, args ...interface{}) }) *Adapter {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30
	}

	adapter := &Adapter{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		connected: false,
		agents:    make(map[string]*AgentStatus),
		logger:   logger,
	}

	if config.Options == nil {
		adapter.config.Options = make(map[string]interface{})
	}

	// Apply default options
	if _, ok := config.Options["max_retries"]; !ok {
		adapter.config.Options["max_retries"] = 3
	}

	return adapter
}

// Connect establishes a connection to the agent framework
func (a *Adapter) Connect() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.connected {
		return nil
	}

	// Validate configuration
	if a.config.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	// Validate endpoint is reachable
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.config.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if a.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+a.config.APIKey)
	}
	req.Header.Set("User-Agent", "AgentShield-Enterprise/1.0")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("endpoint returned status code: %d", resp.StatusCode)
	}

	a.connected = true

	if a.logger != nil {
		a.logger.Printf("[framework] Connected to %s at %s", a.config.Framework, a.config.Endpoint)
	}

	return nil
}

// Disconnect closes the connection to the agent framework
func (a *Adapter) Disconnect() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return nil
	}

	a.connected = false
	a.agents = make(map[string]*AgentStatus)

	if a.logger != nil {
		a.logger.Printf("[framework] Disconnected from %s", a.config.Framework)
	}

	return nil
}

// IsConnected returns whether the adapter is currently connected
func (a *Adapter) IsConnected() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.connected
}

// RegisterAgent registers a new agent with the framework
func (a *Adapter) RegisterAgent(agentID, name string, capabilities []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return fmt.Errorf("adapter not connected")
	}

	if agentID == "" {
		return fmt.Errorf("agent ID is required")
	}

	if name == "" {
		name = agentID
	}

	agent := &AgentStatus{
		ID:            agentID,
		Name:          name,
		Status:        "registered",
		LastActivity: time.Now(),
		Capabilities: capabilities,
	}

	a.agents[agentID] = agent

	if a.logger != nil {
		a.logger.Printf("[framework] Registered agent %s with capabilities: %v", agentID, capabilities)
	}

	return nil
}

// SendCommand sends a command to a specific agent
func (a *Adapter) SendCommand(agentID, command string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return "", fmt.Errorf("adapter not connected")
	}

	agent, exists := a.agents[agentID]
	if !exists {
		return "", fmt.Errorf("agent %s not found", agentID)
	}

	// Build request based on framework type
	payload := map[string]interface{}{
		"agent_id": agentID,
		"command":  command,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	endpoint := a.buildEndpoint("/agents/" + agentID + "/execute")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if a.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+a.config.APIKey)
	}
	req.Header.Set("User-Agent", "AgentShield-Enterprise/1.0")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		agent.Status = "error"
		agent.ErrorMessage = err.Error()
		return "", fmt.Errorf("failed to send command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		agent.Status = "error"
		body, _ := json.Marshal(map[string]interface{}{"error": "status code", "code": resp.StatusCode})
		return "", fmt.Errorf("server returned status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	agent.Status = "executing"
	agent.LastActivity = time.Now()

	output, ok := result["output"].(string)
	if !ok {
		output = fmt.Sprintf("%v", result)
	}

	return output, nil
}

// GetAgentStatus retrieves the status of a specific agent
func (a *Adapter) GetAgentStatus(agentID string) (*AgentStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return nil, fmt.Errorf("adapter not connected")
	}

	agent, exists := a.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}

	return agent, nil
}

// ListAgents returns all registered agents
func (a *Adapter) ListAgents() ([]AgentStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return nil, fmt.Errorf("adapter not connected")
	}

	agents := make([]AgentStatus, 0, len(a.agents))
	for _, agent := range a.agents {
		agents = append(agents, *agent)
	}

	return agents, nil
}

// ValidateConnection verifies the connection is still valid
func (a *Adapter) ValidateConnection() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return fmt.Errorf("adapter not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.config.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create validation request: %w", err)
	}

	if a.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+a.config.APIKey)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.connected = false
		return fmt.Errorf("connection validation failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		a.connected = false
		return fmt.Errorf("connection validation failed with status: %d", resp.StatusCode)
	}

	return nil
}

// buildEndpoint constructs the full API endpoint URL
func (a *Adapter) buildEndpoint(path string) string {
	base := strings.TrimRight(a.config.Endpoint, "/")
	path = strings.TrimLeft(path, "/")
	return base + "/" + path
}

// GetFrameworkType returns the framework type
func (a *Adapter) GetFrameworkType() FrameworkType {
	return a.config.Framework
}

// UpdateAgentStatus updates the status of an agent
func (a *Adapter) UpdateAgentStatus(agentID, status, errorMsg string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	agent, exists := a.agents[agentID]
	if !exists {
		return fmt.Errorf("agent %s not found", agentID)
	}

	agent.Status = status
	agent.LastActivity = time.Now()
	if errorMsg != "" {
		agent.ErrorMessage = errorMsg
	}

	return nil
}

// Middleware returns a Gin middleware for framework adapter
func Middleware(adapter *Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		if !adapter.IsConnected() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "framework adapter not connected",
			})
			c.Abort()
			return
		}

		c.Set("framework_adapter", adapter)
		c.Next()
	}
}

// GetAdapterFromContext retrieves the adapter from Gin context
func GetAdapterFromContext(c *gin.Context) *Adapter {
	if adapter, exists := c.Get("framework_adapter"); exists {
		return adapter.(*Adapter)
	}
	return nil
}