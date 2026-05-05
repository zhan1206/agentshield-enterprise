package permission

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// PermissionLevel represents the hierarchical level of a permission rule.
// It defines the scope at which a permission operates: from broad environment-wide
// settings down to fine-grained data-level controls.
type PermissionLevel int

const (
	// LevelEnvironment represents broad, system-wide permission scope.
	// Permissions at this level affect all resources and operations.
	LevelEnvironment PermissionLevel = iota

	// LevelResource represents resource-group level permission scope.
	// Permissions at this level apply to a specific category of resources.
	LevelResource

	// LevelOperation represents operation-level permission scope.
	// Permissions at this level apply to specific actions on resources.
	LevelOperation

	// LevelData represents fine-grained data-level permission scope.
	// Permissions at this level control access to specific data attributes.
	LevelData
)

// String returns the string representation of a PermissionLevel.
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

// Effect defines whether a permission rule allows or denies an action.
type Effect int

const (
	// EffectAllow grants permission when matched.
	EffectAllow Effect = iota

	// EffectDeny revokes permission when matched.
	EffectDeny
)

// String returns the string representation of an Effect.
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

// PermissionRule defines a single permission rule within the ABAC+RBAC hybrid model.
// It combines role-based constraints with attribute-based conditions to provide
// fine-grained access control across multiple dimensions.
type PermissionRule struct {
	// ID is the unique identifier for this permission rule.
	ID string `json:"id"`

	// Name is a human-readable name for this rule.
	Name string `json:"name"`

	// Description provides detailed explanation of what this rule controls.
	Description string `json:"description"`

	// Effect determines whether matching this rule allows or denies access.
	Effect Effect `json:"effect"`

	// Level defines the hierarchical scope at which this rule operates.
	Level PermissionLevel `json:"level"`

	// AgentRoles specifies which agent roles this rule applies to.
	// An empty slice means the rule applies to all roles.
	AgentRoles []string `json:"agentRoles"`

	// Actions defines the specific operations this rule controls.
	// Examples: "read", "write", "execute", "delete".
	Actions []string `json:"actions"`

	// Resources specifies the resource patterns this rule matches.
	// Supports wildcard patterns like "sandbox:*" or "data:customer:*".
	Resources []string `json:"resources"`

	// Conditions defines attribute-based constraints that must be satisfied.
	// Each key is a condition name, and the value is the expected value.
	// Example: {"securityLevel": "high", "teamId": "engineering"}.
	Conditions map[string]string `json:"conditions"`

	// Priority determines evaluation order when multiple rules match.
	// Higher priority rules are evaluated first.
	Priority int `json:"priority"`

	// Enabled indicates whether this rule is actively enforced.
	// Disabled rules are ignored during evaluation.
	Enabled bool `json:"enabled"`
}

// RoleDefinition defines a named role with optional permission inheritance.
// Roles provide a way to group permissions and assign them to agents collectively.
type RoleDefinition struct {
	// ID uniquely identifies this role within the system.
	ID string `json:"id"`

	// Name is a human-readable role name.
	Name string `json:"name"`

	// Description explains the purpose of this role.
	Description string `json:"description"`

	// BasePermissions lists the permission rule IDs that constitute this role.
	// These permissions are granted to any agent assigned to this role.
	BasePermissions []string `json:"basePermissions"`

	// InheritFrom specifies a parent role from which this role inherits permissions.
	// Empty string means no inheritance (top-level role).
	InheritFrom string `json:"inheritFrom"`
}

// PermissionDecision represents the outcome of a permission evaluation.
// It indicates whether an action on a resource is permitted and provides
// the reasoning behind the decision for audit and debugging purposes.
type PermissionDecision struct {
	// Allowed indicates whether the requested action on the resource is permitted.
	Allowed bool `json:"allowed"`

	// MatchedRule identifies the rule that produced this decision.
	// Empty if no rule matched.
	MatchedRule string `json:"matchedRule"`

	// Reason provides human-readable explanation of the decision.
	// Useful for debugging and audit trails.
	Reason string `json:"reason"`
}

// AgentContext encapsulates all the contextual information required for
// permission evaluation. It includes the agent's identity, assigned roles,
// current task characteristics, and environmental factors that influence
// access control decisions.
type AgentContext struct {
	// AgentID uniquely identifies the agent requesting access.
	AgentID string `json:"agentId"`

	// Roles contains the list of role IDs assigned to the agent.
	Roles []string `json:"roles"`

	// TaskType describes the nature of the task being performed.
	// Examples: "data_processing", "model_training", "file_operations".
	TaskType string `json:"taskType"`

	// ExecutionPhase indicates the current phase within the agent's lifecycle.
	// Examples: "initialization", "execution", "cleanup", "suspended".
	ExecutionPhase string `json:"executionPhase"`

	// SecurityLevel represents the trust classification of the agent.
	// Higher security levels receive broader permissions.
	SecurityLevel string `json:"securityLevel"`

	// TeamID identifies the team to which the agent belongs.
	// Team membership can be used for team-scoped permissions.
	TeamID string `json:"teamId"`

	// Environment describes the execution context.
	// Examples: "production", "staging", "development", "sandbox".
	Environment string `json:"environment"`
}

// String returns a string representation of the AgentContext for debugging.
func (c *AgentContext) String() string {
	return fmt.Sprintf("AgentContext{AgentID:%s, Roles:%v, TaskType:%s, Phase:%s, SecurityLevel:%s, Team:%s, Env:%s}",
		c.AgentID, c.Roles, c.TaskType, c.ExecutionPhase, c.SecurityLevel, c.TeamID, c.Environment)
}

// NewAgentContext creates a new AgentContext with default values.
func NewAgentContext(agentID string) *AgentContext {
	return &AgentContext{
		AgentID:        agentID,
		Roles:          []string{},
		TaskType:       "",
		ExecutionPhase: "execution",
		SecurityLevel:  "Standard",
		TeamID:         "",
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

// RuleSet is a thread-safe collection of permission rules.
type RuleSet struct {
	rules map[string]*PermissionRule
	mu    sync.RWMutex
}

// NewRuleSet creates a new empty rule set.
func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make(map[string]*PermissionRule),
	}
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

// EnabledRules returns only enabled rules, sorted by priority.
func (rs *RuleSet) EnabledRules() []*PermissionRule {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	result := make([]*PermissionRule, 0)
	for _, rule := range rs.rules {
		if rule.Enabled {
			result = append(result, rule)
		}
	}
	// Sort by priority descending
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Priority < result[j].Priority {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

// Grant represents a temporary permission grant with expiration.
type Grant struct {
	GrantID    string    `json:"grantId"`
	AgentID    string    `json:"agentId"`
	PermissionID string  `json:"permissionId"`
	GrantedAt  time.Time `json:"grantedAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
	AutoRevoke bool      `json:"autoRevoke"`
}

// IsExpired checks if the grant has expired.
func (g *Grant) IsExpired() bool {
	return time.Now().After(g.ExpiresAt)
}

// RemainingDuration returns the time until expiration.
func (g *Grant) RemainingDuration() time.Duration {
	return time.Until(g.ExpiresAt)
}