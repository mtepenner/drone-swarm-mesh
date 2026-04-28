package behavior

import (
	"testing"

	fc "github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/flight_controller"
)

func TestComputeBoidsRespondsToNeighbors(t *testing.T) {
	state := fc.DroneState{Position: fc.Vector3{X: 0, Y: 10, Z: 0}, Velocity: fc.Vector3{}}
	neighbors := []Neighbor{
		{AgentID: "left", Position: fc.Vector3{X: -4, Y: 10, Z: 0}, Velocity: fc.Vector3{X: 1}},
		{AgentID: "right", Position: fc.Vector3{X: 4, Y: 10, Z: 0}, Velocity: fc.Vector3{X: 1}},
	}

	desired := ComputeBoids(state, neighbors)
	if desired.Magnitude() == 0 {
		t.Fatal("expected a non-zero boids steering vector")
	}
}
