package data

import (
	"regexp"
	"strings"
	"sync"

	"github.com/zhan1206/agentshield-enterprise/backend/internal/security/encryption"
)

// SensitivityLevel represents data sensitivity classification
type SensitivityLevel string

const (
	SensitivityPublic    SensitivityLevel = "public"
	SensitivityInternal  SensitivityLevel = "internal"
	SensitivityConfidential SensitivityLevel = "confidential"
	SensitivitySecret    SensitivityLevel = "secret"
)

// Finding represents a detected sensitive data item
type Finding struct {
	Type       string           `json:"type"`
	Category   string           `json:"category"`
	Sensitivity SensitivityLevel `json:"sensitivity"`
	Position   [2]int           `json:"position"` // start, end
	Matched    string           `json:"matched"`
	Masked     string           `json:"masked"`
}

// DetectionRule defines a sensitive data detection pattern
type DetectionRule struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Category    string           `json:"category"`
	Pattern     string           `json:"pattern"`
	Sensitivity SensitivityLevel `json:"sensitivity"`
	Enabled     bool             `json:"enabled"`
	MaskFunc    func(string) string `json:"-"`
}

// Guard protects sensitive data from exposure
type Guard struct {
	mu       sync.RWMutex
	rules    map[string]*DetectionRule
	encEngine *encryption.Engine
}

// NewGuard creates a new data guard with built-in detection rules
func NewGuard(encEngine *encryption.Engine) *Guard {
	g := &Guard{
		rules:     make(map[string]*DetectionRule),
		encEngine: encEngine,
	}
	g.initBuiltinRules()
	return g
}

// initBuiltinRules loads built-in sensitive data detection rules
func (g *Guard) initBuiltinRules() {
	builtins := []*DetectionRule{
		{ID: "dr-001", Name: "中国身份证号", Category: "身份证", Pattern: `\b[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]\b`, Sensitivity: SensitivitySecret, Enabled: true, MaskFunc: maskIDCard},
		{ID: "dr-002", Name: "手机号码", Category: "手机号", Pattern: `\b1[3-9]\d{9}\b`, Sensitivity: SensitivityConfidential, Enabled: true, MaskFunc: maskPhone},
		{ID: "dr-003", Name: "银行卡号", Category: "银行卡", Pattern: `\b\d{16,19}\b`, Sensitivity: SensitivitySecret, Enabled: true, MaskFunc: maskBankCard},
		{ID: "dr-004", Name: "邮箱地址", Category: "邮箱", Pattern: `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`, Sensitivity: SensitivityConfidential, Enabled: true, MaskFunc: maskEmail},
		{ID: "dr-005", Name: "API密钥", Category: "密钥", Pattern: `\b(sk|pk|ghp|gho|AKIA|AIza)[-_]?[A-Za-z0-9]{20,}\b`, Sensitivity: SensitivitySecret, Enabled: true, MaskFunc: maskAPIKey},
		{ID: "dr-006", Name: "AWS密钥", Category: "密钥", Pattern: `\bAKIA[0-9A-Z]{16}\b`, Sensitivity: SensitivitySecret, Enabled: true, MaskFunc: maskAPIKey},
		{ID: "dr-007", Name: "JWT Token", Category: "令牌", Pattern: `\beyJ[A-Za-z0-9-_]+\.eyJ[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\b`, Sensitivity: SensitivitySecret, Enabled: true, MaskFunc: maskAPIKey},
		{ID: "dr-008", Name: "IP地址", Category: "网络", Pattern: `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`, Sensitivity: SensitivityInternal, Enabled: true, MaskFunc: maskIP},
	}
	for _, r := range builtins {
		g.rules[r.ID] = r
	}
}

// ScanContent scans content for sensitive data and returns findings
func (g *Guard) ScanContent(content string) ([]Finding, string) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var findings []Finding
	maskedContent := content

	for _, rule := range g.rules {
		if !rule.Enabled {
			continue
		}
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			continue
		}
		matches := re.FindAllStringIndex(content, -1)
		for _, match := range matches {
			matched := content[match[0]:match[1]]
			masked := "****"
			if rule.MaskFunc != nil {
				masked = rule.MaskFunc(matched)
			}
			findings = append(findings, Finding{
				Type:        rule.Name,
				Category:    rule.Category,
				Sensitivity: rule.Sensitivity,
				Position:    [2]int{match[0], match[1]},
				Matched:     matched,
				Masked:      masked,
			})
			maskedContent = strings.Replace(maskedContent, matched, masked, 1)
		}
	}

	return findings, maskedContent
}

// AddRule adds a custom detection rule
func (g *Guard) AddRule(rule *DetectionRule) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.rules[rule.ID] = rule
}

// ListRules returns all detection rules
func (g *Guard) ListRules() []*DetectionRule {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]*DetectionRule, 0, len(g.rules))
	for _, r := range g.rules {
		result = append(result, r)
	}
	return result
}

// --- Masking Functions ---

func maskIDCard(s string) string {
	if len(s) < 8 {
		return "****"
	}
	return s[:3] + "***********" + s[len(s)-1:]
}

func maskPhone(s string) string {
	if len(s) < 7 {
		return "****"
	}
	return s[:3] + "****" + s[7:]
}

func maskBankCard(s string) string {
	if len(s) < 4 {
		return "****"
	}
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}

func maskEmail(s string) string {
	parts := strings.SplitN(s, "@", 2)
	if len(parts) != 2 {
		return "****@****"
	}
	return string(parts[0][0]) + "***@" + parts[1]
}

func maskAPIKey(s string) string {
	if len(s) < 8 {
		return "****"
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}

func maskIP(s string) string {
	parts := strings.SplitN(s, ".", 4)
	if len(parts) != 4 {
		return "***.***.***.***"
	}
	return parts[0] + ".*.*." + parts[3]
}
