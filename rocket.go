package server

//RocketTTL ...
const RocketTTL = 30

//RocketVelocity ...
const RocketVelocity = 4.0

//Rocket ...
type Rocket struct {
	Pos    *Vector `json:"pos"`
	target RocketTarget
	player *Player
	timer  int
}

//RocketTarget ...
type RocketTarget interface {
	Position() *Vector
	Alive() bool
	DealDamage()
}

//NewRocket ...
func NewRocket(pos *Vector, player *Player, target RocketTarget) *Rocket {
	return &Rocket{
		Pos:    pos,
		target: target,
		player: player,
	}
}

//Tick ...
func (r *Rocket) Tick() {
	if r.timer > RocketTTL || !r.target.Alive() {
		r.player.RemoveRocket(r)
		return
	}
	r.move()
	r.checkForCollisions()
}

func (r *Rocket) move() {
	dir := r.target.Position().Copy().Sub(r.Pos)
	if dir.Length() > 0 {
		dir.ToLength(RocketVelocity)
		r.Pos.Add(dir)
	}
}

func (r *Rocket) checkForCollisions() {
	for _, otherPlayer := range OtherPlayers(r.player) {
		rocketTargets := otherPlayer.RocketTargets()

		for _, target := range rocketTargets {
			if target.Position().Dist(r.Pos) < RocketVelocity {
				target.DealDamage()
				r.player.RemoveRocket(r)
				r.player.AddExplosion(r.Pos)
				return
			}
		}
	}
}
