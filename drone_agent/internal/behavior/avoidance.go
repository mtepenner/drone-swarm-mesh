package behavior

import fc "github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/flight_controller"

func ComputeAvoidance(state fc.DroneState, neighbors []Neighbor) fc.Vector3 {
	repulsion := fc.Vector3{}
	for _, neighbor := range neighbors {
		offset := state.Position.Sub(neighbor.Position)
		distance := offset.Magnitude()
		if distance == 0 || distance > 6.0 {
			continue
		}
		repulsion = repulsion.Add(offset.Normalize().Scale((6.0 - distance) / 6.0))
	}
	return repulsion.Scale(4.0).ClampMagnitude(6.5)
}
