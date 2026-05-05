package permission

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Effect defines whether a permission rule allows or denies an action.
type Effect int

const (
	EffectAllow Effect = iota
	EffectDeny
)

func (e Effect) String() string {
	switch e {
	case EffectAllow:
		return "Allow"
	case EffectDeny:
		return "Deny"
	default:
		return fmt.Sprintf("Unknown(%d)", e)
	}
}

// ParseEffect parses a string into an Effect value.
func ParseEffect(s string) (Effect, error) {
	switch strings.ToLower(s) {
	case "allow", "permit", "grant":
		return EffectAllow, nil
	case "deny", "reject", "revoke":
		return EffectDeny, nil
	default:
		return EffectAllow, fmt.Errorf("invalid effect: %s", s)
	}
}

// PermissionLevel represents the hierarchical level of a permission rule.
type PermissionLevel int

const (
	LevelEnvironment PermissionLevel = iota
	LevelResource
	LevelOperation
	LevelData
)

func (l PermissionLevel) String() string {
	switch l {
	case LevelEnvironment:
		return "Environment"
	case LevelResource:
		return "Resource"
	case LevelOperation:
		return "Operation"
	case LevelData:
		return "Data"
	default:
		return fmt.Sprintf("Unknown(%d)", l)
	}
}

// ParsePermissionLevel parses a string into a PermissionLevel.
func ParsePermissionLevel(s string) (PermissionLevel, error) {
	switch strings.ToLower(s) {
	case "environment":
		return LevelEnvironment, nil
	case "resource":
		return LevelResource, nil
	case "operation":
		return LevelOperation, nil
	case "data":
		return LevelData, nil
	default:
		return LevelEnvironment, fmt.Errorf("invalid permission level: %s", s)
	}
}

// --- Core Models ---

// PermissionRule defines a single permission rule within the ABAC+RBAC hybrid model.
type PermissionRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Effect      Effect            `json:"effect"`
	Level       PermissionLevel   `json:"level"`
	AgentRoles  []string          `json:"agentRoles"`
	Actions     []string          `json:"actions"`
	Resources   []string          `json:"resources"`
	Conditions  map[string]string `json:"conditions"`
	Priority    int               `json:"priority"`
	Enabled     bool              `json:"enabled"`
}

// RoleDefinition defines a named role with optional permission inheritance.
type RoleDefinition struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	BasePermissions []string `json:"basePermissions"`
	InheritFrom     string   `json:"inheritFrom"`
}

// AgentContext encapsulates contextual information for permission evaluation.
type AgentContext struct {
	AgentID        string   `json:"agentId"`
	Roles          []string `json:"roles"`
	TaskType       string   `json:"taskType"`
	ExecutionPhase string   `json:"executionPhase"`
	SecurityLevel  string   `json:"securityLevel"`
	TeamID         string   `json:"teamId"`
	Environment    string   `json:"environment"`
}

// NewAgentContext creates a new AgentContext with default values.
func NewAgentContext(agentID string) *AgentContext {
	return &AgentContext{
		AgentID:        agentID,
		Roles:          []string{},
		ExecutionPhase: "execution",
		SecurityLevel:  "Standard",
		Environment:    "production",
	}
}

// HasRole checks if the agent context includes a specific role.
func (c *AgentContext) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Clone creates a deep copy of the AgentContext.
func (c *AgentContext) Clone() *AgentContext {
	rolesCopy := make([]string, len(c.Roles))
	copy(rolesCopy, c.Roles)
	return &AgentContext{
		AgentID:        c.AgentID,
		Roles:          rolesCopy,
		TaskType:       c.TaskType,
		ExecutionPhase: c.ExecutionPhase,
		SecurityLevel:  c.SecurityLevel,
		TeamID:         c.TeamID,
		Environment:    c.Environment,
	}
}

// PermissionDecision represents the outcome of a permission evaluation.
type PermissionDecision struct {
	Allowed     bool   `json:"allowed"`
	MatchedRule string `json:"matchedRule"`
	Reason      string `json:"reason"`
}

// Grant represents a temporary permission grant with expiration.
type Grant struct {
	GrantID      string    `json:"grantId"`
	AgentID      string    `json:"agentId"`
	PermissionID string    `json:"permissionId"`
	GrantedAt    time.Time `json:"grantedAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
	AutoRevoke   bool      `json:"autoRevoke"`
}

// IsExpired checks if the grant has expired.
func (g *Grant) IsExpired() bool {
	return time.Now().After(g.ExpiresAt)
}

// RemainingDuration returns the time until expiration.
func (g *Grant) RemainingDuration() time.Duration {
	return time.Until(g.ExpiresAt)
}

// RuleSet is a thread-safe collection of permission rules.
type RuleSet struct {
	rules map[string]*PermissionRule
	mu    sync.RWMutex
}

// NewRuleSet creates a new empty rule set.
func NewRuleSet() *RuleSet {
	return &RuleSet{rules: make(map[string]*PermissionRule)}
}

// Add inserts a rule into the set.
func (rs *RuleSet) Add(rule *PermissionRule) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.rules[rule.ID] = rule
}

// Remove deletes a rule from the set by ID.
func (rs *RuleSet) Remove(id string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	delete(rs.rules, id)
}

// Get retrieves a rule by ID.
func (rs *RuleSet) Get(id string) (*PermissionRule, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	rule, ok := rs.rules[id]
	return rule, ok
}

// List returns all rules in the set.
func (rs *RuleSet) List() []*PermissionRule {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	result := make([]*PermissionRule, 0, len(rs.rules))
	for _, rule := range rs.rules {
		result = append(result, rule)
	}
	return result
}

// EnabledRules returns only enabled rules, sorted by priority descending.
func (rs *RuleSet) EnabledRules() []*PermissionRule {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	result := make([]*PermissionRule, 0)
	for _, rule := range rs.rules {
		if rule.Enabled {
			result = append(result, rule)
		}
	}
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Priority < result[j].Priority {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

// --- Engine ---

// Engine is the ABAC+RBAC permission evaluation engine
type Engine struct {
	mu    sync.RWMutex
	rules map[string]*PermissionRule
	roles map[string]*RoleDefinition
}

// NewEngine creates a new permission engine
func NewEngine() *Engine {
	return &Engine{
		rules: make(map[string]*PermissionRule),
		roles: make(map[string]*RoleDefinition),
	}
}

// AddRule adds a permission rule
func (e *Engine) AddRule(rule *PermissionRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[rule.ID] = rule
}

// RemoveRule removes a permission rule
func (e *Engine) RemoveRule(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.rules, id)
}

// AddRole adds a role definition
func (e *Engine) AddRole(role *RoleDefinition) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.roles[role.ID] = role
}

// CheckRequest is the input for a permission check
type CheckRequest struct {
	AgentID   string                 `json:"agent_id"`
	AgentRole string                 `json:"agent_role"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// CheckResult is the output of a permission check
type CheckResult struct {
	Allowed      bool     `json:"allowed"`
	Reason       string   `json:"reason,omitempty"`
	MatchedRules []string `json:"matched_rules,omitempty"`
	DeniedBy     string   `json:"denied_by,omitempty"`
}

// Check evaluates whether an action is permitted
func (e *Engine) Check(ctx context.Context, req *CheckRequest) (*CheckResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var applicable []*PermissionRule
	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}
		if e.ruleMatches(rule, req) {
			applicable = append(applicable, rule)
		}
	}

	// Sort by priority (higher first)
	for i := 0; i < len(applicable); i++ {
		for j := i + 1; j < len(applicable); j++ {
			if applicable[j].Priority > applicable[i].Priority {
				applicable[i], applicable[j] = applicable[j], applicable[i]
			}
		}
	}

	if len(applicable) == 0 {
		return &CheckResult{
			Allowed: false,
			Reason:  "no matching rule found, default deny",
		}, nil
	}

	var matchedIDs []string
	for _, rule := range applicable {
		matchedIDs = append(matchedIDs, rule.ID)
		if rule.Effect == EffectDeny {
			return &CheckResult{
				Allowed:      false,
				Reason:       "denied by rule: " + rule.Name,
				MatchedRules: matchedIDs,
				DeniedBy:     rule.ID,
			}, nil
		}
	}

	return &CheckResult{
		Allowed:      true,
		Reason:       "allowed by matching rules",
		MatchedRules: matchedIDs,
	}, nil
}

// ruleMatches checks if a rule applies to the given request
func (e *Engine) ruleMatches(rule *PermissionRule, req *CheckRequest) bool {
	if len(rule.AgentRoles) > 0 {
		roleMatch := false
		for _, r := range rule.AgentRoles {
			if r == req.AgentRole || r == "*" {
				roleMatch = true
				break
			}
		}
		if !roleMatch {
			return false
		}
	}

	if len(rule.Actions) > 0 {
		actionMatch := false
		for _, a := range rule.Actions {
			if a == req.Action || a == "*" {
				actionMatch = true
				break
			}
		}
		if !actionMatch {
			return false
		}
	}

	if len(rule.Resources) > 0 {
		resourceMatch := false
		for _, res := range rule.Resources {
			if res == req.Resource || res == "*" || resourceGlobMatch(res, req.Resource) {
				resourceMatch = true
				break
			}
		}
		if !resourceMatch {
			return false
		}
	}

	return true
}

func resourceGlobMatch(pattern, resource string) bool {
	if pattern == "*" {
		return true
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		if len(resource) >= len(prefix) && resource[:len(prefix)] == prefix {
			return true
		}
	}
	return pattern == resource
}

// ListRules returns all rules
func (e *Engine) ListRules() []*PermissionRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]*PermissionRule, 0, len(e.rules))
	for _, r := range e.rules {
		result = append(result, r)
	}
	return result
}

// ListRoles returns all role definitions
func (e *Engine) ListRoles() []*RoleDefinition {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]*RoleDefinition, 0, len(e.roles))
	for _, r := range e.roles {
		result = append(result, r)
	}
	return result
}
