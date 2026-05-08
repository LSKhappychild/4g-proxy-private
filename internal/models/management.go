package models

import (
	"4g-proxy/internal/inspector"
	"encoding/json"
	"sync"
	"time"
)

// DropFlags controls which signal types should be dropped
type DropFlags struct {
	mu sync.RWMutex

	// Signal type drop flags
	DropAttach           bool `json:"dropAttach"`
	DropDetach           bool `json:"dropDetach"`
	DropTAU              bool `json:"dropTAU"`
	DropServiceRequest   bool `json:"dropServiceRequest"`
	DropUEContextRelease bool `json:"dropUEContextRelease"`
	DropPDNConnectivity  bool `json:"dropPDNConnectivity"`
	DropHandover         bool `json:"dropHandover"`
	DropReset            bool `json:"dropReset"`
	DropPaging           bool `json:"dropPaging"`

	// Direction-specific drops
	DropUplinkOnly   bool `json:"dropUplinkOnly"`
	DropDownlinkOnly bool `json:"dropDownlinkOnly"`
}

// NewDropFlags creates a new DropFlags with all flags disabled
func NewDropFlags() *DropFlags {
	return &DropFlags{}
}

// ShouldDrop checks if a message with the given signal type should be dropped
func (d *DropFlags) ShouldDrop(signalType inspector.SignalType, direction inspector.Direction) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Check direction filter
	if d.DropUplinkOnly && direction != inspector.DirectionUplink {
		return false
	}
	if d.DropDownlinkOnly && direction != inspector.DirectionDownlink {
		return false
	}

	switch signalType {
	case inspector.SignalTypeAttach:
		return d.DropAttach
	case inspector.SignalTypeDetach:
		return d.DropDetach
	case inspector.SignalTypeTAU:
		return d.DropTAU
	case inspector.SignalTypeServiceRequest:
		return d.DropServiceRequest
	case inspector.SignalTypeUEContextRelease:
		return d.DropUEContextRelease
	case inspector.SignalTypePDNConnectivity:
		return d.DropPDNConnectivity
	case inspector.SignalTypeHandover:
		return d.DropHandover
	case inspector.SignalTypeReset:
		return d.DropReset
	case inspector.SignalTypePaging:
		return d.DropPaging
	default:
		return false
	}
}

// SetDrop sets the drop flag for a signal type
func (d *DropFlags) SetDrop(signalType inspector.SignalType, drop bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch signalType {
	case inspector.SignalTypeAttach:
		d.DropAttach = drop
	case inspector.SignalTypeDetach:
		d.DropDetach = drop
	case inspector.SignalTypeTAU:
		d.DropTAU = drop
	case inspector.SignalTypeServiceRequest:
		d.DropServiceRequest = drop
	case inspector.SignalTypeUEContextRelease:
		d.DropUEContextRelease = drop
	case inspector.SignalTypePDNConnectivity:
		d.DropPDNConnectivity = drop
	case inspector.SignalTypeHandover:
		d.DropHandover = drop
	case inspector.SignalTypeReset:
		d.DropReset = drop
	case inspector.SignalTypePaging:
		d.DropPaging = drop
	}
}

// SetDropByName sets the drop flag by signal type name
func (d *DropFlags) SetDropByName(name string, drop bool) bool {
	signalType := ParseSignalType(name)
	if signalType == inspector.SignalTypeUnknown {
		return false
	}
	d.SetDrop(signalType, drop)
	return true
}

// GetAll returns all drop flags as a map
func (d *DropFlags) GetAll() map[string]bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]bool{
		"attach":           d.DropAttach,
		"detach":           d.DropDetach,
		"tau":              d.DropTAU,
		"serviceRequest":   d.DropServiceRequest,
		"ueContextRelease": d.DropUEContextRelease,
		"pdnConnectivity":  d.DropPDNConnectivity,
		"handover":         d.DropHandover,
		"reset":            d.DropReset,
		"paging":           d.DropPaging,
	}
}

// Reset resets all drop flags to false
func (d *DropFlags) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.DropAttach = false
	d.DropDetach = false
	d.DropTAU = false
	d.DropServiceRequest = false
	d.DropUEContextRelease = false
	d.DropPDNConnectivity = false
	d.DropHandover = false
	d.DropReset = false
	d.DropPaging = false
	d.DropUplinkOnly = false
	d.DropDownlinkOnly = false
}

// ToJSON returns JSON representation of drop flags
func (d *DropFlags) ToJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return json.Marshal(d)
}

// FromJSON updates drop flags from JSON
func (d *DropFlags) FromJSON(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return json.Unmarshal(data, d)
}

// DelayConfig controls delay for each signal type (in milliseconds)
type DelayConfig struct {
	mu sync.RWMutex

	// Signal type delays (in milliseconds)
	DelayAttach           int64 `json:"delayAttach"`
	DelayDetach           int64 `json:"delayDetach"`
	DelayTAU              int64 `json:"delayTAU"`
	DelayServiceRequest   int64 `json:"delayServiceRequest"`
	DelayUEContextRelease int64 `json:"delayUEContextRelease"`
	DelayPDNConnectivity  int64 `json:"delayPDNConnectivity"`
	DelayHandover         int64 `json:"delayHandover"`
	DelayReset            int64 `json:"delayReset"`
	DelayPaging           int64 `json:"delayPaging"`

	// Handover-specific delays (in milliseconds)
	// These take precedence over DelayHandover for specific handover messages
	DelayHandoverRequired int64 `json:"delayHandoverRequired"` // Source eNB -> MME
	DelayHandoverNotify   int64 `json:"delayHandoverNotify"`   // Target eNB -> MME

	// Default delay for unspecified signal types (in milliseconds)
	DelayDefault int64 `json:"delayDefault"`

	// Direction-specific delays
	DelayUplinkOnly   bool `json:"delayUplinkOnly"`
	DelayDownlinkOnly bool `json:"delayDownlinkOnly"`
}

// NewDelayConfig creates a new DelayConfig with zero delays
func NewDelayConfig() *DelayConfig {
	return &DelayConfig{}
}

// GetDelay returns the delay duration for a signal type
func (d *DelayConfig) GetDelay(signalType inspector.SignalType, direction inspector.Direction) time.Duration {
	return d.GetDelayWithHandoverSubType(signalType, inspector.HandoverSubTypeNone, direction)
}

// GetDelayWithHandoverSubType returns the delay duration, with specific handling for handover sub-types
func (d *DelayConfig) GetDelayWithHandoverSubType(signalType inspector.SignalType, hoSubType inspector.HandoverSubType, direction inspector.Direction) time.Duration {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Check direction filter
	if d.DelayUplinkOnly && direction != inspector.DirectionUplink {
		return 0
	}
	if d.DelayDownlinkOnly && direction != inspector.DirectionDownlink {
		return 0
	}

	var delayMs int64
	switch signalType {
	case inspector.SignalTypeAttach:
		delayMs = d.DelayAttach
	case inspector.SignalTypeDetach:
		delayMs = d.DelayDetach
	case inspector.SignalTypeTAU:
		delayMs = d.DelayTAU
	case inspector.SignalTypeServiceRequest:
		delayMs = d.DelayServiceRequest
	case inspector.SignalTypeUEContextRelease:
		delayMs = d.DelayUEContextRelease
	case inspector.SignalTypePDNConnectivity:
		delayMs = d.DelayPDNConnectivity
	case inspector.SignalTypeHandover:
		// Check for specific handover sub-type delays first
		switch hoSubType {
		case inspector.HandoverSubTypeRequired:
			if d.DelayHandoverRequired > 0 {
				delayMs = d.DelayHandoverRequired
			} else {
				delayMs = d.DelayHandover
			}
		case inspector.HandoverSubTypeNotify:
			if d.DelayHandoverNotify > 0 {
				delayMs = d.DelayHandoverNotify
			} else {
				delayMs = d.DelayHandover
			}
		default:
			delayMs = d.DelayHandover
		}
	case inspector.SignalTypeReset:
		delayMs = d.DelayReset
	case inspector.SignalTypePaging:
		delayMs = d.DelayPaging
	default:
		delayMs = d.DelayDefault
	}

	if delayMs <= 0 {
		return 0
	}
	return time.Duration(delayMs) * time.Millisecond
}

// SetHandoverSubTypeDelay sets delay for a specific handover sub-type
func (d *DelayConfig) SetHandoverSubTypeDelay(hoSubType inspector.HandoverSubType, delayMs int64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if delayMs < 0 {
		delayMs = 0
	}

	switch hoSubType {
	case inspector.HandoverSubTypeRequired:
		d.DelayHandoverRequired = delayMs
	case inspector.HandoverSubTypeNotify:
		d.DelayHandoverNotify = delayMs
	}
}

// SetDelay sets the delay for a signal type (in milliseconds)
func (d *DelayConfig) SetDelay(signalType inspector.SignalType, delayMs int64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if delayMs < 0 {
		delayMs = 0
	}

	switch signalType {
	case inspector.SignalTypeAttach:
		d.DelayAttach = delayMs
	case inspector.SignalTypeDetach:
		d.DelayDetach = delayMs
	case inspector.SignalTypeTAU:
		d.DelayTAU = delayMs
	case inspector.SignalTypeServiceRequest:
		d.DelayServiceRequest = delayMs
	case inspector.SignalTypeUEContextRelease:
		d.DelayUEContextRelease = delayMs
	case inspector.SignalTypePDNConnectivity:
		d.DelayPDNConnectivity = delayMs
	case inspector.SignalTypeHandover:
		d.DelayHandover = delayMs
	case inspector.SignalTypeReset:
		d.DelayReset = delayMs
	case inspector.SignalTypePaging:
		d.DelayPaging = delayMs
	case inspector.SignalTypeUnknown:
		d.DelayDefault = delayMs
	}
}

// SetDelayByName sets the delay by signal type name (in milliseconds)
func (d *DelayConfig) SetDelayByName(name string, delayMs int64) bool {
	if name == "default" || name == "Default" {
		d.SetDelay(inspector.SignalTypeUnknown, delayMs)
		return true
	}

	// Check for handover sub-types
	hoSubType := ParseHandoverSubType(name)
	if hoSubType != inspector.HandoverSubTypeNone {
		d.SetHandoverSubTypeDelay(hoSubType, delayMs)
		return true
	}

	signalType := ParseSignalType(name)
	if signalType == inspector.SignalTypeUnknown {
		return false
	}
	d.SetDelay(signalType, delayMs)
	return true
}

// ParseHandoverSubType parses a handover sub-type from string
func ParseHandoverSubType(s string) inspector.HandoverSubType {
	switch s {
	case "handoverRequired", "HandoverRequired", "ho-required":
		return inspector.HandoverSubTypeRequired
	case "handoverNotify", "HandoverNotify", "ho-notify":
		return inspector.HandoverSubTypeNotify
	default:
		return inspector.HandoverSubTypeNone
	}
}

// GetAll returns all delay settings as a map (in milliseconds)
func (d *DelayConfig) GetAll() map[string]int64 {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]int64{
		"attach":           d.DelayAttach,
		"detach":           d.DelayDetach,
		"tau":              d.DelayTAU,
		"serviceRequest":   d.DelayServiceRequest,
		"ueContextRelease": d.DelayUEContextRelease,
		"pdnConnectivity":  d.DelayPDNConnectivity,
		"handover":         d.DelayHandover,
		"handoverRequired": d.DelayHandoverRequired,
		"handoverNotify":   d.DelayHandoverNotify,
		"reset":            d.DelayReset,
		"paging":           d.DelayPaging,
		"default":          d.DelayDefault,
	}
}

// Reset resets all delays to zero
func (d *DelayConfig) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.DelayAttach = 0
	d.DelayDetach = 0
	d.DelayTAU = 0
	d.DelayServiceRequest = 0
	d.DelayUEContextRelease = 0
	d.DelayPDNConnectivity = 0
	d.DelayHandover = 0
	d.DelayHandoverRequired = 0
	d.DelayHandoverNotify = 0
	d.DelayReset = 0
	d.DelayPaging = 0
	d.DelayDefault = 0
	d.DelayUplinkOnly = false
	d.DelayDownlinkOnly = false
}

// ToJSON returns JSON representation of delay config
func (d *DelayConfig) ToJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return json.Marshal(d)
}

// FromJSON updates delay config from JSON
func (d *DelayConfig) FromJSON(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return json.Unmarshal(data, d)
}

// ParseSignalType parses a signal type from string
func ParseSignalType(s string) inspector.SignalType {
	switch s {
	case "attach", "Attach":
		return inspector.SignalTypeAttach
	case "detach", "Detach":
		return inspector.SignalTypeDetach
	case "tau", "TAU", "trackingAreaUpdate":
		return inspector.SignalTypeTAU
	case "serviceRequest", "ServiceRequest", "service":
		return inspector.SignalTypeServiceRequest
	case "ueContextRelease", "UEContextRelease", "release":
		return inspector.SignalTypeUEContextRelease
	case "pdnConnectivity", "PDNConnectivity", "pdn":
		return inspector.SignalTypePDNConnectivity
	case "handover", "Handover", "ho":
		return inspector.SignalTypeHandover
	case "reset", "Reset":
		return inspector.SignalTypeReset
	case "paging", "Paging":
		return inspector.SignalTypePaging
	default:
		return inspector.SignalTypeUnknown
	}
}

// ProxyState represents the overall state of the proxy
type ProxyState struct {
	mu sync.RWMutex

	// Connection counts
	ActiveENBConnections int `json:"activeENBConnections"`
	ActiveMMEConnections int `json:"activeMMEConnections"`
	ActiveUEContexts     int `json:"activeUEContexts"`

	// Statistics
	TotalUplinkPackets   uint64 `json:"totalUplinkPackets"`
	TotalDownlinkPackets uint64 `json:"totalDownlinkPackets"`
	TotalUplinkBytes     uint64 `json:"totalUplinkBytes"`
	TotalDownlinkBytes   uint64 `json:"totalDownlinkBytes"`
	DroppedPackets       uint64 `json:"droppedPackets"`
	DelayedPackets       uint64 `json:"delayedPackets"`

	// Drop flags
	DropFlags *DropFlags `json:"dropFlags"`

	// Delay config
	DelayConfig *DelayConfig `json:"delayConfig"`
}

// NewProxyState creates a new ProxyState
func NewProxyState() *ProxyState {
	return &ProxyState{
		DropFlags:   NewDropFlags(),
		DelayConfig: NewDelayConfig(),
	}
}

// IncrementDelayed increments delayed packet counter
func (s *ProxyState) IncrementDelayed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DelayedPackets++
}

// IncrementUplink increments uplink counters
func (s *ProxyState) IncrementUplink(bytes int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalUplinkPackets++
	s.TotalUplinkBytes += uint64(bytes)
}

// IncrementDownlink increments downlink counters
func (s *ProxyState) IncrementDownlink(bytes int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalDownlinkPackets++
	s.TotalDownlinkBytes += uint64(bytes)
}

// IncrementDropped increments dropped packet counter
func (s *ProxyState) IncrementDropped() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DroppedPackets++
}

// SetActiveConnections sets active connection counts
func (s *ProxyState) SetActiveConnections(enb, mme, ue int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ActiveENBConnections = enb
	s.ActiveMMEConnections = mme
	s.ActiveUEContexts = ue
}

// GetStats returns a copy of the current statistics
func (s *ProxyState) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"activeENBConnections": s.ActiveENBConnections,
		"activeMMEConnections": s.ActiveMMEConnections,
		"activeUEContexts":     s.ActiveUEContexts,
		"totalUplinkPackets":   s.TotalUplinkPackets,
		"totalDownlinkPackets": s.TotalDownlinkPackets,
		"totalUplinkBytes":     s.TotalUplinkBytes,
		"totalDownlinkBytes":   s.TotalDownlinkBytes,
		"droppedPackets":       s.DroppedPackets,
		"delayedPackets":       s.DelayedPackets,
	}
}

// ToJSON returns JSON representation of proxy state
func (s *ProxyState) ToJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s)
}
