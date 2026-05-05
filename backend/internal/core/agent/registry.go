package agent

import (
	"fmt"
	"sync"
	"time"
)

// SecurityLevel represents the trust level of an agent
type SecurityLevel string

const (
	SecurityTrusted    SecurityLevel = "trusted"
	SecurityStandard   SecurityLevel = "standard"
	SecurityRestricted SecurityLevel = "restricted"
	SecurityUntrusted  SecurityLevel = "untrusted"
)

// Status represents the current state of an agent
type Status string

const (
	StatusActive    Status = "active"
	StatusIdle      Status = "idle"
	StatusSuspended Status = "suspended"
	StatusOffline   Status = "offline"
)

// Agent represents a registered AI agent
type Agent struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Framework     string        `json:"framework"`
	SecurityLevel SecurityLevel `json:"security_level"`
	Status        Status        `json:"status"`
	TeamID        string        `json:"team_id"`
	Labels        []string      `json:"labels,omitempty"`
	RegisteredAt  time.Time     `json:"registered_at"`
	LastHeartbeat time.Time     `json:"last_heartbeat"`
	SandboxID     string        `json:"sandbox_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Registry manages agent registration and lifecycle
type Registry struct {
	mu     sync.RWMutex
	agents map[string]*Agent
}

// NewRegistry creates a new agent registry
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]*Agent),
	}
}

// Register adds a new agent to the registry
func (r *Registry) Register(agent *Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[agent.ID]; exists {
		return fmt.Errorf("agent already registered: %s", agent.ID)
	}

	agent.RegisteredAt = time.Now()
	agent.LastHeartbeat = time.Now()
	if agent.Status == "" {
		agent.Status = StatusActive
	}
	if agent.SecurityLevel == "" {
		agent.SecurityLevel = SecurityStandard
	}
	r.agents[agent.ID] = agent
	return nil
}

// Deregister removes an agent from the registry
func (r *Registry) Deregister(agentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[agentID]; !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	delete(r.agents, agentID)
	return nil
}

// Get retrieves an agent by ID
func (r *Registry) Get(agentID string) (*Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}
	return agent, nil
}

// List returns all registered agents
func (r *Registry) List() []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Agent, 0, len(r.agents))
	for _, a := range r.agents {
		result = append(result, a)
	}
	return result
}

// ListByTeam returns agents belonging to a specific team
func (r *Registry) ListByTeam(teamID string) []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Agent
	for _, a := range r.agents {
		if a.TeamID == teamID {
			result = append(result, a)
		}
	}
	return result
}

// UpdateHeartbeat updates the last heartbeat time for an agent
func (r *Registry) UpdateHeartbeat(agentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	agent.LastHeartbeat = time.Now()
	agent.Status = StatusActive
	return nil
}

// UpdateStatus updates the status of an agent
func (r *Registry) UpdateStatus(agentID string, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	agent.Status = status
	return nil
}

// SetSandboxID associates a sandbox with an agent
func (r *Registry) SetSandboxID(agentID, sandboxID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	agent.SandboxID = sandboxID
	return nil
}
