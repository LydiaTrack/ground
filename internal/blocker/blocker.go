package blocker

import (
	"sync"
	"time"
)

// Block struct stores blocking details
type Block struct {
	ExpiresAt time.Time
	Method    string
	Endpoint  string
}

// EndpointBlocker struct stores blocks by IP
type EndpointBlocker struct {
	blocks map[string][]Block
	mutex  sync.Mutex
}

// GlobalBlocker is the global instance of EndpointBlocker
var GlobalBlocker *EndpointBlocker

// Initialize the global instance
func Initialize() {
	GlobalBlocker = &EndpointBlocker{
		blocks: make(map[string][]Block),
	}
}

// Add method adds a new block for an IP
func (eb *EndpointBlocker) Add(ip, method, endpoint string, duration time.Duration) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	block := Block{
		ExpiresAt: time.Now().Add(duration),
		Method:    method,
		Endpoint:  endpoint,
	}
	eb.blocks[ip] = append(eb.blocks[ip], block)
}

// RemoveExpired method removes expired blocks
func (eb *EndpointBlocker) RemoveExpired() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	now := time.Now()
	for ip, blocks := range eb.blocks {
		var validBlocks []Block
		for _, block := range blocks {
			if block.ExpiresAt.After(now) {
				validBlocks = append(validBlocks, block)
			}
		}
		if len(validBlocks) > 0 {
			eb.blocks[ip] = validBlocks
		} else {
			delete(eb.blocks, ip)
		}
	}
}

// IsBlocked method checks if a given IP, method, and endpoint are blocked
func (eb *EndpointBlocker) IsBlocked(ip, method, endpoint string) bool {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	blocks, exists := eb.blocks[ip]
	if !exists {
		return false
	}
	now := time.Now()
	for _, block := range blocks {
		if block.Method == method && block.Endpoint == endpoint {
			if block.ExpiresAt.After(now) {
				return true
			}
		}
	}
	return false
}
