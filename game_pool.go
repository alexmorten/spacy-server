package server

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

const mapWidth = 1200.0
const mapHeight = 800.0

//GamePool handles the games
type GamePool struct {
	gamelessConnections []*Stream
	mutex               sync.Mutex
	playerCountPerGame  int
	games               []*Game
	streamWraps         map[string]*Stream
}

//NewGamePool creates a new connection pool
func NewGamePool() *GamePool {
	return &GamePool{
		gamelessConnections: []*Stream{},
		mutex:               sync.Mutex{},
		playerCountPerGame:  2,
		streamWraps:         map[string]*Stream{},
	}
}

//NewCredentials that point to a StreamWrapper
func (p *GamePool) NewCredentials() *Credentials {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	id := uuid.New().String()
	streamWrap := &Stream{}
	p.streamWraps[id] = streamWrap
	return &Credentials{
		Id: id,
	}
}

//HandleAction ...
func (p *GamePool) HandleAction(action *Action) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	streamWrap := p.streamWraps[action.Credentials.Id]
	if streamWrap == nil {
		return errors.New("credentials not valid")
	}
	streamWrap.ActionC <- action

	return nil
}

//AddConnection to the pool thread safely
func (p *GamePool) AddConnection(credentials *Credentials, stream SpacyServer_GetUpdatesServer) (*Stream, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	streamWrap := p.streamWraps[credentials.Id]
	if streamWrap == nil {
		return nil, errors.New("credentials not valid")
	}
	streamWrap.SpacyServer_GetUpdatesServer = stream
	streamWrap.ShutdownC = make(chan struct{})
	streamWrap.ActionC = make(chan *Action)

	p.gamelessConnections = append(p.gamelessConnections, streamWrap)
	if len(p.gamelessConnections) == p.playerCountPerGame {
		players := []*Player{}
		g := NewGame()
		for index, c := range p.gamelessConnections {
			player := NewPlayer(c, g, randomPos(), index)
			players = append(players, player)
		}
		p.gamelessConnections = []*Stream{}

		g.Players = players
		g.Start()
	}
	return streamWrap, nil
}

func randomPos() *Vector {
	return &Vector{
		X: rand.Float64() * mapWidth,
		Y: rand.Float64() * mapHeight,
	}
}
