package sctp

import (
	"4g-proxy/internal/inspector"
	"4g-proxy/internal/models"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ishidawataru/sctp"
)

const (
	// S1AP SCTP port
	DefaultS1APPort = 36412

	// SCTP payload protocol identifier for S1AP
	S1APPPID = 18

	// Buffer size for SCTP messages
	BufferSize = 65535
)

// Proxy represents the S1AP SCTP proxy
type Proxy struct {
	// Configuration
	listenAddr string
	mmeAddr    string

	// State
	state     *models.ProxyState
	inspector *inspector.Inspector

	// Connection management
	listener *sctp.SCTPListener
	mu       sync.RWMutex
	sessions map[string]*Session

	// Control
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// Session represents a proxy session between eNB and MME
type Session struct {
	id       string
	enbConn  *sctp.SCTPConn
	mmeConn  *sctp.SCTPConn
	proxy    *Proxy
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewProxy creates a new SCTP proxy
func NewProxy(listenAddr, mmeAddr string, state *models.ProxyState) *Proxy {
	return &Proxy{
		listenAddr: listenAddr,
		mmeAddr:    mmeAddr,
		state:      state,
		inspector:  inspector.NewInspector(true),
		sessions:   make(map[string]*Session),
		stopCh:     make(chan struct{}),
	}
}

// Start starts the proxy
func (p *Proxy) Start() error {
	addr, err := sctp.ResolveSCTPAddr("sctp", p.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve listen address: %w", err)
	}

	listener, err := sctp.ListenSCTP("sctp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on SCTP: %w", err)
	}

	p.listener = listener
	log.Printf("S1AP Proxy listening on %s", p.listenAddr)
	log.Printf("Forwarding to MME at %s", p.mmeAddr)

	p.wg.Add(1)
	go p.acceptLoop()

	return nil
}

// Stop stops the proxy
func (p *Proxy) Stop() {
	close(p.stopCh)

	if p.listener != nil {
		p.listener.Close()
	}

	// Close all sessions
	p.mu.Lock()
	for _, session := range p.sessions {
		session.Close()
	}
	p.mu.Unlock()

	p.wg.Wait()
	log.Println("Proxy stopped")
}

func (p *Proxy) acceptLoop() {
	defer p.wg.Done()

	// Create a channel to receive accept results
	type acceptResult struct {
		conn *sctp.SCTPConn
		err  error
	}
	acceptCh := make(chan acceptResult, 1)

	for {
		select {
		case <-p.stopCh:
			return
		default:
		}

		// Run accept in a goroutine so we can check stop channel
		go func() {
			conn, err := p.listener.AcceptSCTP()
			select {
			case acceptCh <- acceptResult{conn, err}:
			default:
				// Channel full, close connection if we got one
				if conn != nil {
					conn.Close()
				}
			}
		}()

		// Wait for accept or stop signal
		select {
		case <-p.stopCh:
			return
		case result := <-acceptCh:
			if result.err != nil {
				select {
				case <-p.stopCh:
					return
				default:
					log.Printf("Accept error: %v", result.err)
					continue
				}
			}

			session, err := p.createSession(result.conn)
			if err != nil {
				log.Printf("Failed to create session: %v", err)
				result.conn.Close()
				continue
			}

			p.mu.Lock()
			p.sessions[session.id] = session
			p.mu.Unlock()

			go session.Start()
		}
	}
}

func (p *Proxy) createSession(enbConn *sctp.SCTPConn) (*Session, error) {
	// Connect to MME
	mmeAddr, err := sctp.ResolveSCTPAddr("sctp", p.mmeAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve MME address: %w", err)
	}

	mmeConn, err := sctp.DialSCTP("sctp", nil, mmeAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MME: %w", err)
	}

	remoteAddr := enbConn.RemoteAddr()
	sessionID := remoteAddr.String()

	log.Printf("New session from eNB %s", sessionID)

	return &Session{
		id:      sessionID,
		enbConn: enbConn,
		mmeConn: mmeConn,
		proxy:   p,
		stopCh:  make(chan struct{}),
	}, nil
}

func (p *Proxy) removeSession(id string) {
	p.mu.Lock()
	delete(p.sessions, id)
	p.mu.Unlock()
}

// GetSessionCount returns the number of active sessions
func (p *Proxy) GetSessionCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.sessions)
}

// Start starts the session proxy
func (s *Session) Start() {
	s.wg.Add(2)
	go s.proxyUplink()
	go s.proxyDownlink()

	s.wg.Wait()
	s.cleanup()
}

// Close closes the session
func (s *Session) Close() {
	select {
	case <-s.stopCh:
		// Already closed
	default:
		close(s.stopCh)
	}
}

func (s *Session) cleanup() {
	if s.enbConn != nil {
		s.enbConn.Close()
	}
	if s.mmeConn != nil {
		s.mmeConn.Close()
	}
	s.proxy.removeSession(s.id)
	log.Printf("Session %s closed", s.id)
}

func (s *Session) proxyUplink() {
	defer s.wg.Done()
	s.proxyData(s.enbConn, s.mmeConn, inspector.DirectionUplink)
}

func (s *Session) proxyDownlink() {
	defer s.wg.Done()
	s.proxyData(s.mmeConn, s.enbConn, inspector.DirectionDownlink)
}

func (s *Session) proxyData(src, dst *sctp.SCTPConn, dir inspector.Direction) {
	buf := make([]byte, BufferSize)

	for {
		select {
		case <-s.stopCh:
			return
		default:
		}

		// Set read deadline
		src.SetReadDeadline(time.Now().Add(1 * time.Second))

		n, info, err := src.SCTPRead(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if err == io.EOF {
				log.Printf("Connection closed by peer (%s)", dir)
				s.Close()
				return
			}
			select {
			case <-s.stopCh:
				return
			default:
				log.Printf("Read error (%s): %v", dir, err)
				s.Close()
				return
			}
		}

		if n == 0 {
			continue
		}

		data := buf[:n]

		// Inspect the message
		result, _ := s.proxy.inspector.InspectAndLog(data, dir)

		// Check if message should be dropped
		if result != nil && s.proxy.state.DropFlags.ShouldDrop(result.SignalType, dir) {
			log.Printf("DROPPED: %s", result.Summary)
			s.proxy.state.IncrementDropped()
			continue
		}

		// Check if message should be delayed
		var delay time.Duration
		if result != nil {
			delay = s.proxy.state.DelayConfig.GetDelayWithHandoverSubType(result.SignalType, result.HandoverSubType, dir)
		}

		if delay > 0 {
			log.Printf("DELAYING %v: %s", delay, result.Summary)
			s.proxy.state.IncrementDelayed()

			// Apply delay with cancellation support
			select {
			case <-s.stopCh:
				return
			case <-time.After(delay):
				// Continue after delay
			}
		}

		// Update statistics
		if dir == inspector.DirectionUplink {
			s.proxy.state.IncrementUplink(n)
		} else {
			s.proxy.state.IncrementDownlink(n)
		}

		// Forward the message
		sndInfo := &sctp.SndRcvInfo{
			PPID: S1APPPID,
		}
		if info != nil {
			sndInfo.Stream = info.Stream
			sndInfo.Context = info.Context
		}
		_, err = dst.SCTPWrite(data, sndInfo)
		if err != nil {
			log.Printf("Write error (%s): %v", dir, err)
			s.Close()
			return
		}
	}
}
