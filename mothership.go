package server

// MothershipVelocity ...
const MothershipVelocity = 1

//MothershipMiningTime ...
const MothershipMiningTime = 300

//MothershipMiningRange ...
const MothershipMiningRange = 100

//MothershipTotalHealth ...
const MothershipTotalHealth = 70

//Mothership ...
type Mothership struct {
	Pos          *Vector `json:"pos"`
	TargetPos    *Vector `json:"target_pos"`
	player       *Player
	MiningTarget *Asteroid `json:"mining_target,omitempty"`
	Health       int       `json:"health"`
	miningTimer  int
	launcher     *Launcher
	ore          int
}

//NewMothership ...
func NewMothership(pos *Vector, player *Player) *Mothership {
	return &Mothership{
		Pos:       pos,
		TargetPos: pos,
		player:    player,
		Health:    MothershipTotalHealth,
		launcher:  NewLauncher(pos, 100, player),
	}
}

//Tick one update
func (m *Mothership) Tick() {
	if !m.Alive() {
		return
	}
	m.spawn()
	m.move()
	m.mine()
	m.launcher.Tick()
}

//Position for RocketTarget interface
func (m *Mothership) Position() *Vector {
	return m.Pos
}

//Alive ˆ
func (m *Mothership) Alive() bool {
	return m.Health > 0
}

//DealDamage ...
func (m *Mothership) DealDamage() {
	m.Health--
}

func (m *Mothership) move() {
	dir := m.TargetPos.Copy().Sub(m.Pos)

	if dir.Length() > 0 {

		m.Pos.Add(dir.ToLength(MothershipVelocity))
	}
}

func (m *Mothership) mine() {
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

func (m *Mothership) spawn() {
	if m.ore >= 3 {
		m.player.AddShip()
		m.player.AddMiner()
		m.ore -= 3
	}
}

func (m *Mothership) findMiningTarget() *Asteroid {
	for _, asteroid := range m.player.game.Asteroids {
		if m.Pos.Dist(asteroid.Pos) < MothershipMiningRange && asteroid.Capacity > 0 {
			return asteroid
		}
	}
	return nil
}
