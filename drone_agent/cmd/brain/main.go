package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/behavior"
	fc "github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/flight_controller"
	"github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/mesh"
)

type telemetrySample struct {
	AgentID      string   `json:"agent_id"`
	X            float64  `json:"x"`
	Y            float64  `json:"y"`
	Z            float64  `json:"z"`
	VX           float64  `json:"vx"`
	VY           float64  `json:"vy"`
	VZ           float64  `json:"vz"`
	BatteryLevel float64  `json:"battery_level"`
	PeerIDs      []string `json:"peer_ids"`
}

func main() {
	agentID := getenv("DRONE_AGENT_ID", hostnameFallback())
	discoveryPort := getenvInt("DISCOVERY_PORT", 10010)
	gossipPort := getenvInt("GOSSIP_PORT", 10020)
	simulatorAddr := getenv("SWARM_SIMULATOR_UDP_ADDR", "127.0.0.1:9010")
	tickMillis := getenvInt("TICK_MS", 200)
	index := getenvInt("DRONE_INDEX", 1)
	advertiseHost := resolveAdvertiseHost(simulatorAddr)

	directory := mesh.NewDirectory(agentID)
	bus := mesh.NewBus(gossipPort)
	controller := fc.NewController()
	state := fc.DroneState{
		AgentID:      agentID,
		Position:     seededPosition(index),
		Velocity:     fc.Vector3{},
		BatteryLevel: 100,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := mesh.Listen(ctx, discoveryPort, directory); err != nil {
			log.Printf("discovery listener stopped: %v", err)
		}
	}()

	go func() {
		if err := bus.Listen(ctx); err != nil {
			log.Printf("gossip listener stopped: %v", err)
		}
	}()

	ticker := time.NewTicker(time.Duration(tickMillis) * time.Millisecond)
	defer ticker.Stop()

	for tick := 0; ; tick++ {
		peerAddresses := directory.Snapshot()
		neighbors := make([]behavior.Neighbor, 0, len(peerAddresses))
		peerIDs := make([]string, 0, len(peerAddresses))
		for _, peer := range bus.Snapshot() {
			if peer.AgentID == agentID {
				continue
			}
			neighbors = append(neighbors, behavior.Neighbor{AgentID: peer.AgentID, Position: peer.Position, Velocity: peer.Velocity})
			peerIDs = append(peerIDs, peer.AgentID)
		}

		desired := behavior.ComputeBoids(state, neighbors).Add(behavior.ComputeAvoidance(state, neighbors)).Add(orbitBias(index, tick)).ClampMagnitude(6.0)
		state = controller.Step(state, desired, float64(tickMillis)/1000.0)

		if tick%5 == 0 {
			err := mesh.Broadcast(mesh.Peer{AgentID: agentID, ListenAddr: net.JoinHostPort(advertiseHost, strconv.Itoa(gossipPort))}, discoveryPort)
			if err != nil {
				log.Printf("discovery broadcast failed: %v", err)
			}
		}

		snapshot := mesh.PeerSnapshot{AgentID: agentID, Position: state.Position, Velocity: state.Velocity, Timestamp: time.Now().UnixMilli()}
		if err := bus.Broadcast(snapshot, peerAddresses); err != nil {
			log.Printf("gossip broadcast failed: %v", err)
		}

		if err := sendTelemetry(simulatorAddr, telemetrySample{
			AgentID:      agentID,
			X:            state.Position.X,
			Y:            state.Position.Y,
			Z:            state.Position.Z,
			VX:           state.Velocity.X,
			VY:           state.Velocity.Y,
			VZ:           state.Velocity.Z,
			BatteryLevel: state.BatteryLevel,
			PeerIDs:      peerIDs,
		}); err != nil {
			log.Printf("simulator uplink failed: %v", err)
		}

		<-ticker.C
	}
}

func sendTelemetry(address string, payload telemetrySample) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = conn.Write(encoded)
	return err
}

func orbitBias(index int, tick int) fc.Vector3 {
	angle := float64(index)*0.35 + float64(tick)*0.08
	return fc.Vector3{X: math.Cos(angle) * 0.6, Y: math.Sin(angle*0.3) * 0.15, Z: math.Sin(angle) * 0.6}
}

func seededPosition(index int) fc.Vector3 {
	angle := float64(index) * 0.7
	radius := 10.0 + float64(index%6)*2.5
	return fc.Vector3{X: math.Cos(angle) * radius, Y: 8 + float64(index%5), Z: math.Sin(angle) * radius}
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func hostnameFallback() string {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		return "drone-agent"
	}
	return hostname
}

func resolveAdvertiseHost(simulatorAddr string) string {
	conn, err := net.Dial("udp", simulatorAddr)
	if err == nil {
		defer conn.Close()
		if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok && addr.IP != nil {
			return addr.IP.String()
		}
	}
	return hostnameFallback()
}
