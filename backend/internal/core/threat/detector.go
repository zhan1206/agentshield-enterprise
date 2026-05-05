package threat

import (
	"sync"
	"time"
)

// ThreatLevel represents the severity of a threat
type ThreatLevel string

const (
	LevelCritical ThreatLevel = "critical"
	LevelHigh     ThreatLevel = "high"
	LevelMedium   ThreatLevel = "medium"
	LevelLow      ThreatLevel = "low"
)

// ThreatType categorizes the threat
type ThreatType string

const (
	TypeUnauthorizedAccess ThreatType = "unauthorized_access"
	TypeMaliciousCode      ThreatType = "malicious_code"
	TypeDataExfiltration   ThreatType = "data_exfiltration"
	TypeResourceAbuse      ThreatType = "resource_abuse"
	TypePrivilegeEscalation ThreatType = "privilege_escalation"
	TypeAnomalousBehavior  ThreatType = "anomalous_behavior"
)

// Alert represents a detected threat alert
type Alert struct {
	ID          string      `json:"id"`
	Timestamp   time.Time   `json:"timestamp"`
	AgentID     string      `json:"agent_id"`
	ThreatType  ThreatType  `json:"threat_type"`
	Level       ThreatLevel `json:"level"`
	Description string      `json:"description"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Status      string      `json:"status"` // active, responded, resolved
	Responses   []string    `json:"responses,omitempty"`
}

// DetectionRule defines a threat detection rule
type DetectionRule struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        ThreatType  `json:"type"`
	Level       ThreatLevel `json:"level"`
	Description string      `json:"description"`
	Condition   string      `json:"condition"` // rule expression
	Enabled     bool        `json:"enabled"`
}

// Detector detects threats using rules and AI analysis
type Detector struct {
	mu    sync.RWMutex
	rules map[string]*DetectionRule
	alerts map[string]*Alert
}

// NewDetector creates a new threat detector
func NewDetector() *Detector {
	d := &Detector{
		rules:  make(map[string]*DetectionRule),
		alerts: make(map[string]*Alert),
	}
	d.initBuiltinRules()
	return d
}

// initBuiltinRules loads built-in threat detection rules
func (d *Detector) initBuiltinRules() {
	builtins := []*DetectionRule{
		{ID: "rule-001", Name: "越权访问检测", Type: TypeUnauthorizedAccess, Level: LevelHigh, Description: "检测Agent访问未授权资源", Condition: "action=='access' && permission.denied==true", Enabled: true},
		{ID: "rule-002", Name: "恶意代码执行", Type: TypeMaliciousCode, Level: LevelCritical, Description: "检测沙箱中执行恶意代码", Condition: "action=='execute' && content.matches('rm -rf|del /|format c:')", Enabled: true},
		{ID: "rule-003", Name: "数据外泄检测", Type: TypeDataExfiltration, Level: LevelCritical, Description: "检测敏感数据外传行为", Condition: "action=='send' && data.sensitivity=='high'", Enabled: true},
		{ID: "rule-004", Name: "资源滥用检测", Type: TypeResourceAbuse, Level: LevelMedium, Description: "检测CPU/内存资源异常使用", Condition: "cpu_usage>90 || memory_usage>95", Enabled: true},
		{ID: "rule-005", Name: "权限提升检测", Type: TypePrivilegeEscalation, Level: LevelHigh, Description: "检测尝试提权操作", Condition: "action=='chmod' || action=='sudo' || action=='runas'", Enabled: true},
		{ID: "rule-006", Name: "异常行为检测", Type: TypeAnomalousBehavior, Level: LevelMedium, Description: "AI模型检测行为偏差", Condition: "anomaly_score>0.8", Enabled: true},
	}
	for _, r := range builtins {
		d.rules[r.ID] = r
	}
}

// Detect evaluates an event against detection rules and generates alerts
func (d *Detector) Detect(agentID string, event map[string]interface{}) []*Alert {
	d.mu.Lock()
	defer d.mu.Unlock()

	var alerts []*Alert
	for _, rule := range d.rules {
		if !rule.Enabled {
			continue
		}
		if d.evaluateCondition(rule, event) {
			alert := &Alert{
				ID:          "alert-" + rule.ID + "-" + time.Now().Format("20060102150405"),
				Timestamp:   time.Now(),
				AgentID:     agentID,
				ThreatType:  rule.Type,
				Level:       rule.Level,
				Description: rule.Description,
				Details:     event,
				Status:      "active",
			}
			d.alerts[alert.ID] = alert
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

// evaluateCondition evaluates a simple rule condition (placeholder for real engine)
func (d *Detector) evaluateCondition(rule *DetectionRule, event map[string]interface{}) bool {
	// In production, this would use a proper rule engine (e.g., go-ceval)
	// For now, return false as default (no false positives)
	return false
}

// AddRule adds a custom detection rule
func (d *Detector) AddRule(rule *DetectionRule) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rules[rule.ID] = rule
}

// ListAlerts returns all alerts, optionally filtered by status
func (d *Detector) ListAlerts(status string) []*Alert {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var result []*Alert
	for _, a := range d.alerts {
		if status == "" || a.Status == status {
			result = append(result, a)
		}
	}
	return result
}

// UpdateAlertStatus updates the status of an alert
func (d *Detector) UpdateAlertStatus(alertID, status string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if alert, ok := d.alerts[alertID]; ok {
		alert.Status = status
	}
}
