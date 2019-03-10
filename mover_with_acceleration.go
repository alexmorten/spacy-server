package server

//MoverWithAcceleration ...
type MoverWithAcceleration struct {
	Pos          *Vector `json:"pos"`
	Vel          *Vector `json:"vel"`
	Acceleration float64
	MaxVelocity  float64
	Braking      float64
}

func (m *MoverWithAcceleration) accelerateTo(pos Vector) {
	force := (&pos).Sub(m.Pos)
	if force.Length() > 0 {
		force.ToLength(m.Acceleration)
		m.Vel.Add(force)
		brakingForce := m.Vel.Copy().Mul(-1).ToLength(m.Braking)
		m.Vel.Add(brakingForce)
		if m.Vel.Length() > m.MaxVelocity {
			m.Vel.ToLength(m.MaxVelocity)
		}
	}
}

func (m *MoverWithAcceleration) move() {
	m.Pos.Add(m.Vel)
}

func (m *MoverWithAcceleration) breakWith(strength float64) {
	if m.Vel.Length() > 0 {
		brakingForce := m.Vel.Copy().Mul(-1).ToLength(strength)
		m.Vel.Add(brakingForce)
		if m.Vel.Length() > m.MaxVelocity {
			m.Vel.ToLength(m.MaxVelocity)
		}
	}
}
