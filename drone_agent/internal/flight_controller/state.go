package flightcontroller

import "math"

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (vector Vector3) Add(other Vector3) Vector3 {
	return Vector3{X: vector.X + other.X, Y: vector.Y + other.Y, Z: vector.Z + other.Z}
}

func (vector Vector3) Sub(other Vector3) Vector3 {
	return Vector3{X: vector.X - other.X, Y: vector.Y - other.Y, Z: vector.Z - other.Z}
}

func (vector Vector3) Scale(multiplier float64) Vector3 {
	return Vector3{X: vector.X * multiplier, Y: vector.Y * multiplier, Z: vector.Z * multiplier}
}

func (vector Vector3) Magnitude() float64 {
	return math.Sqrt(vector.X*vector.X + vector.Y*vector.Y + vector.Z*vector.Z)
}

func (vector Vector3) Normalize() Vector3 {
	length := vector.Magnitude()
	if length == 0 {
		return Vector3{}
	}
	return vector.Scale(1 / length)
}

func (vector Vector3) ClampMagnitude(max float64) Vector3 {
	length := vector.Magnitude()
	if length <= max || length == 0 {
		return vector
	}
	return vector.Scale(max / length)
}

type DroneState struct {
	AgentID      string
	Position     Vector3
	Velocity     Vector3
	BatteryLevel float64
}

type Controller struct {
	KP       float64
	KD       float64
	MaxAccel float64
	MaxSpeed float64
}

func NewController() Controller {
	return Controller{
		KP:       0.9,
		KD:       0.25,
		MaxAccel: 6.0,
		MaxSpeed: 7.5,
	}
}

func (controller Controller) Step(state DroneState, desiredVelocity Vector3, dt float64) DroneState {
	error := desiredVelocity.Sub(state.Velocity)
	acceleration := error.Scale(controller.KP).Sub(state.Velocity.Scale(controller.KD)).ClampMagnitude(controller.MaxAccel)
	state.Velocity = state.Velocity.Add(acceleration.Scale(dt)).ClampMagnitude(controller.MaxSpeed)
	state.Position = state.Position.Add(state.Velocity.Scale(dt))
	state.BatteryLevel = math.Max(10.0, state.BatteryLevel-0.015)
	return state
}
