package server

import (
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
)

const mapWidth = 1200.0
const mapHeight = 800.0

//GamePool handles the games
type GamePool struct {
	gamelessConnections []*websocket.Conn
	mutex               sync.Mutex
	playerCountPerGame  int
	games               []*Game
}

//NewGamePool creates a new connection pool
func NewGamePool() *GamePool {
	return &GamePool{
		gamelessConnections: []*websocket.Conn{},
		mutex:               sync.Mutex{},
		playerCountPerGame:  2,
	}
}

//AddConnection to the pool thread safely
func (p *GamePool) AddConnection(conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.gamelessConnections = append(p.gamelessConnections, conn)
	if len(p.gamelessConnections) == p.playerCountPerGame {
		players := []*Player{}
		g := NewGame()
		for index, c := range p.gamelessConnections {
			player := NewPlayer(c, g, randomPos(), index)
			players = append(players, player)
		}
		p.gamelessConnections = []*websocket.Conn{}

		g.Players = players
		g.Start()
	}
}

func randomPos() *Vector {
	return &Vector{
		X: rand.Float64() * mapWidth,
		Y: rand.Float64() * mapHeight,
	}

}
