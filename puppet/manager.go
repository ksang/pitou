package puppet

import (
	"sync"

	"github.com/ksang/pitou/store"
)

// Manager manages puppets
type Manager struct {
	collectors map[Collector]chan error
	// PuppetCount is the number of puppets managed
	PuppetCount int
	mu          sync.RWMutex
	Store       *store.Client
}

func NewManager() *Manager {
	return &Manager{
		collectors: make(map[Collector]chan error),
	}
}

func (m *Manager) Add(c Collector) error {
	m.mu.Lock()
	m.collectors[c] = nil
	errCh := c.Start()
	m.collectors[c] = errCh
	m.PuppetCount += 1
	m.mu.Unlock()
	return nil
}
