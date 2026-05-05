package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// Entry represents a single audit log entry in the chain
type Entry struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	AgentID    string                 `json:"agent_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	Result     string                 `json:"result"` // success, failed, denied
	Details    map[string]interface{} `json:"details,omitempty"`
	PrevHash   string                 `json:"prev_hash"`
	Hash       string                 `json:"hash"`
}

// Chain implements an immutable audit chain using SHA-256 linking
type Chain struct {
	mu      sync.RWMutex
	entries []*Entry
	lastHash string
}

// NewChain creates a new audit chain with a genesis block
func NewChain() *Chain {
	c := &Chain{
		entries:  make([]*Entry, 0),
		lastHash: "0",
	}
	// Genesis entry
	c.Append(&Entry{
		AgentID:  "system",
		Action:   "chain_init",
		Resource: "audit_chain",
		Result:   "success",
	})
	return c
}

// Append adds a new entry to the audit chain
func (c *Chain) Append(entry *Entry) *Entry {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	entry.ID = "aud-" + hex.EncodeToString(sha256.New().Sum([]byte(entry.Timestamp.String()+entry.AgentID))[:8]
	entry.PrevHash = c.lastHash

	// Compute hash of this entry
	hash := sha256.New()
	hash.Write([]byte(entry.ID))
	hash.Write([]byte(entry.Timestamp.String()))
	hash.Write([]byte(entry.AgentID))
	hash.Write([]byte(entry.Action))
	hash.Write([]byte(entry.Resource))
	hash.Write([]byte(entry.Result))
	hash.Write([]byte(entry.PrevHash))
	entry.Hash = hex.EncodeToString(hash.Sum(nil))

	c.lastHash = entry.Hash
	c.entries = append(c.entries, entry)
	return entry
}

// Verify checks the integrity of the entire chain
func (c *Chain) Verify() (bool, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	prevHash := "0"
	for i, entry := range c.entries {
		if entry.PrevHash != prevHash {
			return false, i
		}
		// Recompute hash
		hash := sha256.New()
		hash.Write([]byte(entry.ID))
		hash.Write([]byte(entry.Timestamp.String()))
		hash.Write([]byte(entry.AgentID))
		hash.Write([]byte(entry.Action))
		hash.Write([]byte(entry.Resource))
		hash.Write([]byte(entry.Result))
		hash.Write([]byte(entry.PrevHash))
		computed := hex.EncodeToString(hash.Sum(nil))
		if computed != entry.Hash {
			return false, i
		}
		prevHash = entry.Hash
	}
	return true, len(c.entries)
}

// List returns paginated audit entries (most recent first)
func (c *Chain) List(limit, offset int) []*Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := len(c.entries)
	if offset >= total {
		return nil
	}
	end := total - offset
	start := end - limit
	if start < 0 {
		start = 0
	}

	result := make([]*Entry, end-start)
	copy(result, c.entries[start:end])
	// Reverse for most recent first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// GetByAgent returns all entries for a specific agent
func (c *Chain) GetByAgent(agentID string) []*Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []*Entry
	for _, entry := range c.entries {
		if entry.AgentID == agentID {
			result = append(result, entry)
		}
	}
	return result
}
