package server

import (
	"math/rand"
)

//ShipRelativeTargetDist ...
const ShipRelativeTargetDist = 100

// ShipHealth ...
const ShipHealth = 20

//Ship ...
type Ship struct {
	MoverWithAcceleration

	player             *Player
	moveTarget         *Vector
	moveRelativeTarget *Vector
	launcher           *Launcher
	damageReceived     int
}

//NewShip ...
func NewShip(p *Player) *Ship {
	pos := p.Mothership.Pos.Copy()
	s := &Ship{
		player:   p,
		launcher: NewLauncher(pos, 100, p),
		MoverWithAcceleration: MoverWithAcceleration{
			Pos:          pos,
			Vel:          &Vector{},
			Acceleration: 0.025,
			Braking:      0.005,
			MaxVelocity:  3,
		},
	}
	return s
}

//Tick ...
func (s *Ship) Tick() {
	s.move()
	s.launcher.Tick()
}

func (s *Ship) move() {
	if s.moveTarget == nil {
		s.moveTarget = s.player.Mothership.Pos
	}

	if s.moveRelativeTarget == nil || s.Pos.Dist(s.destination()) < 20 {
		s.moveRelativeTarget = s.generateRandomRelativeTargetPos()
	}

	s.accelerateTo(*s.destination())
	s.MoverWithAcceleration.move()
}

func (s *Ship) destination() *Vector {
	return s.moveTarget.Copy().Add(s.moveRelativeTarget)
}

//Position for RocketTargetinterface
func (s *Ship) Position() *Vector {
	return s.Pos
}

// Alive for RocketTargetInterface
func (s *Ship) Alive() bool {
	return s.damageReceived < ShipHealth
}

//DealDamage ...
func (s *Ship) DealDamage() {
	s.damageReceived++
	if !s.Alive() {
		s.player.RemoveShip(s)
	}
}

func (s *Ship) generateRandomRelativeTargetPos() *Vector {
	v := &Vector{
		X: rand.Float64()*2 - 1.0,
		Y: rand.Float64()*2 - 1.0,
	}
	return v.ToLength(ShipRelativeTargetDist)
}
