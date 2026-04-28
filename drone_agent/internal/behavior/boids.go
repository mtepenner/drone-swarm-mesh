package behavior

import fc "github.com/mtepenner/drone-swarm-mesh/drone_agent/internal/flight_controller"

type Neighbor struct {
	AgentID  string
	Position fc.Vector3
	Velocity fc.Vector3
}

func ComputeBoids(state fc.DroneState, neighbors []Neighbor) fc.Vector3 {
	if len(neighbors) == 0 {
		return fc.Vector3{}
	}

	separation := fc.Vector3{}
	alignment := fc.Vector3{}
	cohesion := fc.Vector3{}

	for _, neighbor := range neighbors {
		offset := state.Position.Sub(neighbor.Position)
		distance := offset.Magnitude()
		if distance > 0 && distance < 12.0 {
			separation = separation.Add(offset.Normalize().Scale(1 / distance))
		}
		alignment = alignment.Add(neighbor.Velocity)
		cohesion = cohesion.Add(neighbor.Position)
	}

	count := float64(len(neighbors))
	alignment = alignment.Scale(1 / count).Sub(state.Velocity)
	cohesion = cohesion.Scale(1 / count).Sub(state.Position)

	return separation.Scale(1.7).Add(alignment.Scale(0.65)).Add(cohesion.Scale(0.35)).ClampMagnitude(5.5)
}
