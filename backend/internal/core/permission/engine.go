package permission

import (
	"context"
	"sync"
)

// Effect represents whether a permission allows or denies access
type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)

// Level represents the granularity of a permission
type Level string

const (
	LevelEnvironment Level = "environment"
	LevelResource    Level = "resource"
	LevelOperation   Level = "operation"
	LevelData        Level = "data"
)

// Rule represents a single permission rule
type Rule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Effect      Effect   `json:"effect"`
	Level       Level    `json:"level"`
	AgentRoles  []string `json:"agent_roles"`
	Actions     []string `json:"actions"`
	Resources   []string `json:"resources"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Priority    int      `json:"priority"`
	Enabled     bool     `json:"enabled"`
}

// RoleDefinition represents a role with associated permissions
type RoleDefinition struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	BasePermissions []string `json:"base_permissions"`
	SecurityLevel   string   `json:"security_level"`
}

// CheckRequest is the input for a permission check
type CheckRequest struct {
	AgentID    string                 `json:"agent_id"`
	AgentRole  string                 `json:"agent_role"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// CheckResult is the output of a permission check
type CheckResult struct {
	Allowed    bool     `json:"allowed"`
	Reason     string   `json:"reason,omitempty"`
	MatchedRules []string `json:"matched_rules,omitempty"`
	DeniedBy   string   `json:"denied_by,omitempty"`
}

// Engine is the ABAC+RBAC permission evaluation engine
type Engine struct {
	mu    sync.RWMutex
	rules map[string]*Rule
	roles map[string]*RoleDefinition
}

// NewEngine creates a new permission engine
func NewEngine() *Engine {
	return &Engine{
		rules: make(map[string]*Rule),
		roles: make(map[string]*RoleDefinition),
	}
}

// AddRule adds a permission rule
func (e *Engine) AddRule(rule *Rule) {
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

// Check evaluates whether an action is permitted
func (e *Engine) Check(ctx context.Context, req *CheckRequest) (*CheckResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Collect all applicable rules sorted by priority
	var applicable []*Rule
	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}
		if e.ruleMatches(rule, req) {
			applicable = append(applicable, rule)
		}
	}

	// Sort by priority (higher priority first)
	for i := 0; i < len(applicable); i++ {
		for j := i + 1; j < len(applicable); j++ {
			if applicable[j].Priority > applicable[i].Priority {
				applicable[i], applicable[j] = applicable[j], applicable[i]
			}
		}
	}

	// If no rules match, default deny
	if len(applicable) == 0 {
		return &CheckResult{
			Allowed:  false,
			Reason:   "no matching rule found, default deny",
		}, nil
	}

	// First deny rule wins (deny takes precedence)
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

	// If all matching rules are allow, permit
	return &CheckResult{
		Allowed:      true,
		Reason:       "allowed by matching rules",
		MatchedRules: matchedIDs,
	}, nil
}

// ruleMatches checks if a rule applies to the given request
func (e *Engine) ruleMatches(rule *Rule, req *CheckRequest) bool {
	// Check agent role match
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

	// Check action match
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

	// Check resource match
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

// resourceGlobMatch does simple glob matching for resource patterns
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
func (e *Engine) ListRules() []*Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]*Rule, 0, len(e.rules))
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
