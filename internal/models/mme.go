package models

import (
	"fmt"
	"net"
	"sync"
)

// MME represents an MME backend
type MME struct {
	ID      string `json:"id" yaml:"id"`
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
	Port    int    `json:"port" yaml:"port"`
}

// Endpoint returns the full endpoint address
func (m *MME) Endpoint() string {
	return fmt.Sprintf("%s:%d", m.Address, m.Port)
}

// MMEPool manages a pool of MME backends
type MMEPool struct {
	mu   sync.RWMutex
	mmes map[string]*MME
}

// NewMMEPool creates a new MME pool
func NewMMEPool() *MMEPool {
	return &MMEPool{
		mmes: make(map[string]*MME),
	}
}

// Add adds an MME to the pool
func (p *MMEPool) Add(mme *MME) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mmes[mme.ID] = mme
}

// Remove removes an MME from the pool
func (p *MMEPool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.mmes, id)
}

// Get retrieves an MME by ID
func (p *MMEPool) Get(id string) (*MME, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	mme, ok := p.mmes[id]
	return mme, ok
}

// List returns all MMEs
func (p *MMEPool) List() []*MME {
	p.mu.RLock()
	defer p.mu.RUnlock()
	mmes := make([]*MME, 0, len(p.mmes))
	for _, mme := range p.mmes {
		mmes = append(mmes, mme)
	}
	return mmes
}

// UEContext represents UE context in the proxy
type UEContext struct {
	MMEUEID uint32
	ENBUEID uint32

	// Connection tracking
	ENBConn net.Conn
	MMEConn net.Conn

	// State
	State UEState

	// Statistics
	UplinkPackets   uint64
	DownlinkPackets uint64
	UplinkBytes     uint64
	DownlinkBytes   uint64
}

// UEState represents the state of a UE
type UEState int

const (
	UEStateIdle UEState = iota
	UEStateConnecting
	UEStateConnected
	UEStateReleasing
)

func (s UEState) String() string {
	switch s {
	case UEStateIdle:
		return "Idle"
	case UEStateConnecting:
		return "Connecting"
	case UEStateConnected:
		return "Connected"
	case UEStateReleasing:
		return "Releasing"
	default:
		return "Unknown"
	}
}

// UEContextManager manages UE contexts
type UEContextManager struct {
	mu       sync.RWMutex
	contexts map[uint32]*UEContext // keyed by MME-UE-S1AP-ID
}

// NewUEContextManager creates a new UE context manager
func NewUEContextManager() *UEContextManager {
	return &UEContextManager{
		contexts: make(map[uint32]*UEContext),
	}
}

// Add adds a UE context
func (m *UEContextManager) Add(ctx *UEContext) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.contexts[ctx.MMEUEID] = ctx
}

// Get retrieves a UE context by MME-UE-ID
func (m *UEContextManager) Get(mmeUeId uint32) (*UEContext, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ctx, ok := m.contexts[mmeUeId]
	return ctx, ok
}

// GetByENBUEID retrieves a UE context by eNB-UE-ID
func (m *UEContextManager) GetByENBUEID(enbUeId uint32) (*UEContext, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, ctx := range m.contexts {
		if ctx.ENBUEID == enbUeId {
			return ctx, true
		}
	}
	return nil, false
}

// Remove removes a UE context
func (m *UEContextManager) Remove(mmeUeId uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.contexts, mmeUeId)
}

// List returns all UE contexts
func (m *UEContextManager) List() []*UEContext {
	m.mu.RLock()
	defer m.mu.RUnlock()
	contexts := make([]*UEContext, 0, len(m.contexts))
	for _, ctx := range m.contexts {
		contexts = append(contexts, ctx)
	}
	return contexts
}

// Count returns the number of UE contexts
func (m *UEContextManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.contexts)
}
