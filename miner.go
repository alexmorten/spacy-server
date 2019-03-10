package server

import (
	"math/rand"
)

//MinerOreCapacity ...
const MinerOreCapacity = 1

//MinerDeliveryRange ...
const MinerDeliveryRange = 40

//MinerMiningRange ...
const MinerMiningRange = 150

//MinerRelativeTargetDist ...
const MinerRelativeTargetDist = 30

// Miner collects resource and brings them back to the mothership
type Miner struct {
	MoverWithAcceleration
	player             *Player
	ore                int
	MiningTarget       *Asteroid `json:"mining_target,omitempty"`
	miningTimer        int
	Health             int
	moveTarget         *Vector
	moveRelativeTarget *Vector
}

// NewMiner ...
func NewMiner(p *Player) *Miner {
	pos := p.Mothership.Pos.Copy()
	return &Miner{
		player: p,
		Health: 2,
		MoverWithAcceleration: MoverWithAcceleration{
			Acceleration: 0.05,
			Braking:      0.005,
			MaxVelocity:  4,
			Pos:          pos,
			Vel:          &Vector{},
		}}
}

//Tick ...
func (m *Miner) Tick() {
	m.move()
	m.mine()
	m.deliver()
}

//Position ...
func (m *Miner) Position() *Vector {
	return m.Pos
}

//Alive ...
func (m *Miner) Alive() bool {
	return m.Health > 0
}

//DealDamage ...
func (m *Miner) DealDamage() {
	m.Health--
	if !m.Alive() {
		m.player.RemoveMiner(m)
	}
}

func (m *Miner) destination() *Vector {
	return m.moveTarget.Copy().Add(m.moveRelativeTarget)
}

func (m *Miner) move() {
	defer m.MoverWithAcceleration.move()

	shouldBreak := false
	if m.ore >= MinerOreCapacity {
		m.moveTarget = m.player.Mothership.Pos
	} else {
		asteroid := m.closestAsteroid()
		if asteroid != nil {
			m.moveTarget = asteroid.Pos
			distToAsteroid := m.Pos.Dist(asteroid.Pos)
			if distToAsteroid < MinerMiningRange {
				shouldBreak = true
			}
		} else {
			m.moveTarget = m.player.Mothership.Pos
		}

	}

	if m.moveRelativeTarget == nil || m.Pos.Dist(m.destination()) < 20 {
		m.moveRelativeTarget = m.generateRandomRelativeTargetPos()
	}

	m.accelerateTo(*m.destination())
	if shouldBreak {
		m.breakWith(0.035)
	}
}

func (m *Miner) mine() {
	if m.MiningTarget != nil {
		if m.Pos.Dist(m.MiningTarget.Pos) > MothershipMiningRange || m.MiningTarget.Capacity <= 0 {
			m.miningTimer = 0
			m.MiningTarget = m.findMiningTarget()
		}

		m.miningTimer++
		if m.miningTimer >= MothershipMiningTime {
			m.ore++
			m.miningTimer = 0
			m.MiningTarget.Capacity--
		}

	} else {
		m.MiningTarget = m.findMiningTarget()
	}
}

func (m *Miner) deliver() {
	if m.ore >= MinerOreCapacity && m.Pos.Dist(m.player.Mothership.Pos) < MinerDeliveryRange {
		m.player.Mothership.ore += m.ore
		m.ore = 0
	}
}

func (m *Miner) closestAsteroid() (asteroid *Asteroid) {
	var shortestDistance *float64
	for _, a := range m.player.game.Asteroids {
		if a.Capacity > 0 {
			dist := m.Pos.Dist(a.Pos)
			if shortestDistance != nil && *shortestDistance < dist {
				continue
			}
			asteroid = a
			shortestDistance = &dist
		}
	}
	return
}

func (m *Miner) findMiningTarget() *Asteroid {
	asteroid := m.closestAsteroid()
	if asteroid != nil && asteroid.Pos.Dist(m.Pos) <= MinerMiningRange {
		return asteroid
	}
	return nil
}

func (m *Miner) generateRandomRelativeTargetPos() *Vector {
	v := &Vector{
		X: rand.Float64()*2 - 1.0,
		Y: rand.Float64()*2 - 1.0,
	}
	return v.ToLength(MinerRelativeTargetDist)
}
