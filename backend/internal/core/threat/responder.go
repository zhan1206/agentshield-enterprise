package threat

import (
	"fmt"
	"time"
)

// ResponseAction defines the type of automated response
type ResponseAction string

const (
	ActionAlert    ResponseAction = "alert"
	ActionBlock    ResponseAction = "block"
	ActionSuspend  ResponseAction = "suspend"
	ActionTerminate ResponseAction = "terminate"
	ActionIsolate  ResponseAction = "isolate"
	ActionRollback ResponseAction = "rollback"
)

// ResponsePolicy defines when to trigger which response
type ResponsePolicy struct {
	Level  ThreatLevel    `json:"level"`
	Actions []ResponseAction `json:"actions"`
	AutoExecute bool       `json:"auto_execute"`
}

// Responder handles automated threat responses
type Responder struct {
	detector *Detector
	policies map[ThreatLevel]*ResponsePolicy
}

// NewResponder creates a new threat responder with default policies
func NewResponder(detector *Detector) *Responder {
	r := &Responder{
		detector: detector,
		policies: make(map[ThreatLevel]*ResponsePolicy),
	}
	r.initDefaultPolicies()
	return r
}

// initDefaultPolicies sets up default response policies
func (r *Responder) initDefaultPolicies() {
	r.policies[LevelCritical] = &ResponsePolicy{
		Level:       LevelCritical,
		Actions:     []ResponseAction{ActionAlert, ActionBlock, ActionTerminate, ActionRollback},
		AutoExecute: true,
	}
	r.policies[LevelHigh] = &ResponsePolicy{
		Level:       LevelHigh,
		Actions:     []ResponseAction{ActionAlert, ActionBlock, ActionSuspend},
		AutoExecute: true,
	}
	r.policies[LevelMedium] = &ResponsePolicy{
		Level:       LevelMedium,
		Actions:     []ResponseAction{ActionAlert, ActionBlock},
		AutoExecute: false,
	}
	r.policies[LevelLow] = &ResponsePolicy{
		Level:       LevelLow,
		Actions:     []ResponseAction{ActionAlert},
		AutoExecute: false,
	}
}

// Respond executes the response policy for a given alert
func (r *Responder) Respond(alertID string, action ResponseAction) error {
	r.detector.mu.RLock()
	alert, ok := r.detector.alerts[alertID]
	r.detector.mu.RUnlock()

	if !ok {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	// Execute the response action
	switch action {
	case ActionAlert:
		fmt.Printf("[ALERT] Threat detected: %s (Agent: %s, Level: %s)\n", alert.Description, alert.AgentID, alert.Level)
	case ActionBlock:
		fmt.Printf("[BLOCK] Blocking agent %s operations\n", alert.AgentID)
	case ActionSuspend:
		fmt.Printf("[SUSPEND] Suspending agent %s\n", alert.AgentID)
	case ActionTerminate:
		fmt.Printf("[TERMINATE] Terminating agent %s sandbox\n", alert.AgentID)
	case ActionIsolate:
		fmt.Printf("[ISOLATE] Isolating agent %s\n", alert.AgentID)
	case ActionRollback:
		fmt.Printf("[ROLLBACK] Rolling back agent %s to last snapshot\n", alert.AgentID)
	}

	// Update alert status
	r.detector.UpdateAlertStatus(alertID, "responded")
	alert.Responses = append(alert.Responses, string(action))

	return nil
}

// AutoRespond checks if auto-response is enabled and executes the policy
func (r *Responder) AutoRespond(alert *Alert) error {
	policy, ok := r.policies[alert.Level]
	if !ok || !policy.AutoExecute {
		return nil
	}

	for _, action := range policy.Actions {
		if err := r.Respond(alert.ID, action); err != nil {
			return err
		}
	}
	return nil
}

// Resolve marks an alert as resolved
func (r *Responder) Resolve(alertID string) error {
	r.detector.UpdateAlertStatus(alertID, "resolved")
	return nil
}

// GetPolicy returns the response policy for a threat level
func (r *Responder) GetPolicy(level ThreatLevel) *ResponsePolicy {
	return r.policies[level]
}

// StartWatcher begins periodic threat checking (placeholder)
func (r *Responder) StartWatcher(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			alerts := r.detector.ListAlerts("active")
			for _, alert := range alerts {
				_ = r.AutoRespond(alert)
			}
		}
	}()
}
