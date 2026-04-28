package behavior

import (
	"testing"

	fc "github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/flight_controller"
)

func TestComputeAvoidancePushesAwayFromCloseNeighbor(t *testing.T) {
	state := fc.DroneState{Position: fc.Vector3{X: 0, Y: 8, Z: 0}}
	neighbors := []Neighbor{{AgentID: "close", Position: fc.Vector3{X: 1, Y: 8, Z: 0}}}

	force := ComputeAvoidance(state, neighbors)
	if force.X >= 0 {
		t.Fatalf("expected avoidance force to push left, got %+v", force)
	}
}
