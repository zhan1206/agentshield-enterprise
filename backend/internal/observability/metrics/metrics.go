package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Collector collects and exposes system metrics
type Collector struct {
	activeSandboxes  int64
	runningAgents    int64
	securityEvents   int64
	blockedOps       int64
	threatAlerts     int64
	startTime        time.Time
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		startTime: time.Now(),
	}
}

// GinMiddleware returns a Gin middleware for metrics collection
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// GetActiveSandboxes returns the number of active sandboxes
func (c *Collector) GetActiveSandboxes() int64 { return c.activeSandboxes }

// GetRunningAgents returns the number of running agents
func (c *Collector) GetRunningAgents() int64 { return c.runningAgents }

// GetSecurityEvents returns the count of security events
func (c *Collector) GetSecurityEvents() int64 { return c.securityEvents }

// GetBlockedOps returns the count of blocked operations
func (c *Collector) GetBlockedOps() int64 { return c.blockedOps }

// IncrementSandboxes increments the sandbox count
func (c *Collector) IncrementSandboxes() { c.activeSandboxes++ }

// DecrementSandboxes decrements the sandbox count
func (c *Collector) DecrementSandboxes() { c.activeSandboxes-- }

// IncrementAgents increments the agent count
func (c *Collector) IncrementAgents() { c.runningAgents++ }

// DecrementAgents decrements the agent count
func (c *Collector) DecrementAgents() { c.runningAgents-- }

// RecordSecurityEvent records a security event
func (c *Collector) RecordSecurityEvent() { c.securityEvents++ }

// RecordBlockedOp records a blocked operation
func (c *Collector) RecordBlockedOp() { c.blockedOps++ }

// RecordThreatAlert records a threat alert
func (c *Collector) RecordThreatAlert() { c.threatAlerts++ }

// Snapshot returns a snapshot of all metrics
func (c *Collector) Snapshot() map[string]interface{} {
	return map[string]interface{}{
		"active_sandboxes":    c.activeSandboxes,
		"running_agents":      c.runningAgents,
		"security_events":     c.securityEvents,
		"blocked_operations":  c.blockedOps,
		"threat_alerts":       c.threatAlerts,
		"uptime_seconds":      time.Since(c.startTime).Seconds(),
	}
}
