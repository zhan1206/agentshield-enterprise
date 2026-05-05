package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims extends JWT claims with agent-specific fields
type Claims struct {
	AgentID  string `json:"agent_id"`
	Role     string `json:"role"`
	TeamID   string `json:"team_id,omitempty"`
	jwt.RegisteredClaims
}

// Manager handles JWT authentication
type Manager struct {
	secretKey  []byte
	expiry     time.Duration
}

// NewManager creates a new auth manager
func NewManager(secretKey string, expiry time.Duration) *Manager {
	return &Manager{
		secretKey: []byte(secretKey),
		expiry:    expiry,
	}
}

// GenerateToken generates a JWT token for an agent
func (m *Manager) GenerateToken(agentID, role, teamID string) (string, error) {
	claims := &Claims{
		AgentID: agentID,
		Role:    role,
		TeamID:  teamID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "agentshield-enterprise",
			Subject:   agentID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateToken validates a JWT token and returns claims
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// RefreshToken refreshes an existing token
func (m *Manager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return m.GenerateToken(claims.AgentID, claims.Role, claims.TeamID)
}
