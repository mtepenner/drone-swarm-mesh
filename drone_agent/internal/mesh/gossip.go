package mesh

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"

	fc "github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/flight_controller"
)

type PeerSnapshot struct {
	AgentID   string     `json:"agent_id"`
	Position  fc.Vector3 `json:"position"`
	Velocity  fc.Vector3 `json:"velocity"`
	Timestamp int64      `json:"timestamp"`
}

type Bus struct {
	listenPort int
	mu         sync.RWMutex
	peers      map[string]PeerSnapshot
}

func NewBus(listenPort int) *Bus {
	return &Bus{listenPort: listenPort, peers: map[string]PeerSnapshot{}}
}

func (bus *Bus) Listen(ctx context.Context) error {
	addr := net.UDPAddr{IP: net.IPv4zero, Port: bus.listenPort}
	conn, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(750 * time.Millisecond))
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-ctx.Done():
					return nil
				default:
					continue
				}
			}
			return err
		}

		var peer PeerSnapshot
		if err := json.Unmarshal(buffer[:n], &peer); err != nil {
			continue
		}

		bus.mu.Lock()
		bus.peers[peer.AgentID] = peer
		bus.mu.Unlock()
	}
}

func (bus *Bus) Snapshot() []PeerSnapshot {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	peers := make([]PeerSnapshot, 0, len(bus.peers))
	for _, peer := range bus.peers {
		if time.Since(time.UnixMilli(peer.Timestamp)) <= 5*time.Second {
			peers = append(peers, peer)
		}
	}
	return peers
}

func (bus *Bus) Broadcast(sample PeerSnapshot, peers []Peer) error {
	payload, err := json.Marshal(sample)
	if err != nil {
		return err
	}

	for _, peer := range peers {
		addr, err := net.ResolveUDPAddr("udp4", peer.ListenAddr)
		if err != nil {
			continue
		}
		conn, err := net.DialUDP("udp4", nil, addr)
		if err != nil {
			continue
		}
		_, _ = conn.Write(payload)
		_ = conn.Close()
	}
	return nil
}
