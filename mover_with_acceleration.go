package server

//MoverWithAcceleration ...
type MoverWithAcceleration struct {
	Pos          *Vector `json:"pos"`
	Vel          *Vector `json:"vel"`
	Acceleration float64
	MaxVelocity  float64
	Braking      float64
}

func (m *MoverWithAcceleration) moveTo(pos Vector) {
	force := (&pos).Sub(m.Pos)
	if force.Length() > 0 {
		force.ToLength(m.Acceleration)
		m.Vel.Add(force)
		brakingForce := m.Vel.Copy().Mul(-1).ToLength(m.Braking)
		m.Vel.Add(brakingForce)
		if m.Vel.Length() > m.MaxVelocity {
			m.Vel.ToLength(m.MaxVelocity)
		}
		m.Pos.Add(m.Vel)
	}

}
