package server

//MinerOreCapacity ...
const MinerOreCapacity = 1

//MinerDeliveryRange ...
const MinerDeliveryRange = 20

//MinerMiningRange ...
const MinerMiningRange = 150

// Miner collects resource and brings them back to the mothership
type Miner struct {
	MoverWithAcceleration
	player       *Player
	ore          int
	MiningTarget *Asteroid `json:"mining_target,omitempty"`
	miningTimer  int
	Health       int
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

func (m *Miner) move() {
	if m.ore >= MinerOreCapacity {
		m.moveTo(*m.player.Mothership.Pos)
		return
	}

	asteroid := m.closestAsteroid()
	if asteroid != nil {
		m.moveTo(*asteroid.Pos)
	} else {
		m.moveTo(*m.player.Mothership.Pos)
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
