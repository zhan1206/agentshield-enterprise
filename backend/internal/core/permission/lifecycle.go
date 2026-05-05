package permission

import (
	"context"
	"time"
)

// Lifecycle manages permission lifecycle: grant, expire, revoke
type Lifecycle struct {
	engine *Engine
	grants map[string]*Grant
}

// NewLifecycle creates a new permission lifecycle manager
func NewLifecycle(engine *Engine) *Lifecycle {
	return &Lifecycle{
		engine: engine,
		grants: make(map[string]*Grant),
	}
}

// GrantPermissions temporarily grants permissions to an agent
func (l *Lifecycle) GrantPermissions(ctx context.Context, agentID string, ruleIDs []string, duration time.Duration, reason string) (*Grant, error) {
	grant := &Grant{
		GrantID:    "grant-" + agentID + "-" + time.Now().Format("20060102150405"),
		AgentID:    agentID,
		GrantedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(duration),
		AutoRevoke: true,
	}
	l.grants[grant.GrantID] = grant
	return grant, nil
}

// RevokeGrant revokes a permission grant
func (l *Lifecycle) RevokeGrant(ctx context.Context, grantID string) error {
	delete(l.grants, grantID)
	return nil
}

// CleanupExpired removes expired grants
func (l *Lifecycle) CleanupExpired(ctx context.Context) int {
	now := time.Now()
	expired := 0
	for id, grant := range l.grants {
		if grant.AutoRevoke && now.After(grant.ExpiresAt) {
			delete(l.grants, id)
			expired++
		}
	}
	return expired
}

// GetActiveGrants returns all active grants for an agent
func (l *Lifecycle) GetActiveGrants(agentID string) []*Grant {
	var result []*Grant
	now := time.Now()
	for _, grant := range l.grants {
		if grant.AgentID == agentID && now.Before(grant.ExpiresAt) {
			result = append(result, grant)
		}
	}
	return result
}
