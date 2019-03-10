package server

import (
	"fmt"
	"sync"
)

// Player ...
type Player struct {
	ID              int `json:"id"`
	stream          *Stream
	game            *Game
	Mothership      *Mothership  `json:"mothership"`
	Ships           []*Ship      `json:"ships"`
	Miners          []*Miner     `json:"miners"`
	Rockets         []*Rocket    `json:"rockets"`
	Explosions      []*explosion `json:"explosions"`
	mutex           sync.Mutex
	shutdownChannel chan struct{}
}

type action struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type moveAction struct {
	Type string  `json:"type"`
	Data *Vector `json:"data"`
}

type explosion struct {
	Pos *Vector `json:"pos"`
}

// NewPlayer ...
func NewPlayer(stream *Stream, game *Game, pos *Vector, id int) *Player {
	p := &Player{
		stream:          stream,
		game:            game,
		shutdownChannel: make(chan struct{}),
		Ships:           []*Ship{},
		Miners:          []*Miner{},
		Rockets:         []*Rocket{},
		ID:              id,
	}
	p.Mothership = NewMothership(pos, p)
	for i := 0; i < 3; i++ {
		p.AddShip()
		p.AddMiner()
	}
	go p.readMessages()
	return p
}

//SendToPlayer Sends the byte Array to the player
func (p *Player) SendToPlayer(byteArray []byte) {
	err := p.stream.Send(&State{
		JsonState: byteArray,
	})
	if err != nil {
		fmt.Println(err)
		p.game.Shutdown()
	}
}

//Tick one update
func (p *Player) Tick() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Mothership.Tick()
	for _, ship := range p.Ships {
		ship.Tick()
	}

	for _, miner := range p.Miners {
		miner.Tick()
	}

	for _, rocket := range p.Rockets {
		rocket.Tick()
	}
}

// AddExplosion to player temp list
func (p *Player) AddExplosion(pos *Vector) {
	p.Explosions = append(p.Explosions, &explosion{Pos: pos})
}

//AddShip to Player
func (p *Player) AddShip() {
	s := NewShip(p)
	p.Ships = append(p.Ships, s)
}

//RemoveShip ...
func (p *Player) RemoveShip(ship *Ship) {
	newShips := []*Ship{}
	for _, s := range p.Ships {
		if s != ship {
			newShips = append(newShips, s)
		}
	}
	p.Ships = newShips
}

//AddMiner to Player
func (p *Player) AddMiner() {
	m := NewMiner(p)
	p.Miners = append(p.Miners, m)
}

//RemoveMiner ...
func (p *Player) RemoveMiner(miner *Miner) {
	newMiners := []*Miner{}
	for _, m := range p.Miners {
		if m != miner {
			newMiners = append(newMiners, m)
		}
	}
	p.Miners = newMiners
}

//AddRocket ...
func (p *Player) AddRocket(pos *Vector, target RocketTarget) {
	p.Rockets = append(p.Rockets, NewRocket(pos.Copy(), p, target))
}

//RemoveRocket from Player
func (p *Player) RemoveRocket(r *Rocket) {
	newRockets := []*Rocket{}
	for _, rocket := range p.Rockets {
		if rocket != r {
			newRockets = append(newRockets, rocket)
		}
	}
	p.Rockets = newRockets
}

func (p *Player) readMessages() {
	for {
		action := <-p.stream.ActionC
		if action == nil {
			p.game.Shutdown()
			break
		}
		moveAction := action.GetMove()
		if moveAction != nil {
			p.mutex.Lock()
			p.Mothership.TargetPos = moveAction.Pos
			p.mutex.Unlock()
			continue
		}
		fmt.Println(action.String(), "not supported")

	}
	fmt.Println("Player ", p.ID, " disconnected")
}

//Shutdown player
func (p *Player) Shutdown() {
	p.stream.ShutdownC <- struct{}{}
}

//OtherPlayers ...
func OtherPlayers(p *Player) (players []*Player) {
	for _, player := range p.game.Players {
		if player != p {
			players = append(players, player)
		}
	}
	return
}

//RocketTargets of player
func (p *Player) RocketTargets() []RocketTarget {
	rocketTargets := []RocketTarget{}
	for _, ship := range p.Ships {
		rocketTargets = append(rocketTargets, ship)
	}

	for _, miner := range p.Miners {
		rocketTargets = append(rocketTargets, miner)
	}

	rocketTargets = append(rocketTargets, p.Mothership)
	return rocketTargets
}
