package mesh

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"
)

type Peer struct {
	AgentID    string    `json:"agent_id"`
	ListenAddr string    `json:"listen_addr"`
	LastSeen   time.Time `json:"-"`
}

type Directory struct {
	selfID string
	mu     sync.RWMutex
	peers  map[string]Peer
}

func NewDirectory(selfID string) *Directory {
	return &Directory{selfID: selfID, peers: map[string]Peer{}}
}

func (directory *Directory) Observe(peer Peer) {
	if peer.AgentID == directory.selfID || peer.AgentID == "" {
		return
	}
	peer.LastSeen = time.Now()
	directory.mu.Lock()
	defer directory.mu.Unlock()
	directory.peers[peer.AgentID] = peer
}

func (directory *Directory) Snapshot() []Peer {
	directory.mu.RLock()
	defer directory.mu.RUnlock()
	peers := make([]Peer, 0, len(directory.peers))
	for _, peer := range directory.peers {
		if time.Since(peer.LastSeen) <= 5*time.Second {
			peers = append(peers, peer)
		}
	}
	return peers
}

func Listen(ctx context.Context, port int, directory *Directory) error {
	addr := net.UDPAddr{IP: net.IPv4zero, Port: port}
	conn, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 512)
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

		var peer Peer
		if err := json.Unmarshal(buffer[:n], &peer); err != nil {
			continue
		}
		directory.Observe(peer)
	}
}

func Broadcast(peer Peer, port int) error {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.IPv4bcast, Port: port})
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetWriteBuffer(1024); err != nil {
		return err
	}

	payload, err := json.Marshal(peer)
	if err != nil {
		return err
	}
	_, err = conn.Write(payload)
	return err
}
