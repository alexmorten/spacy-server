package server

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Game ...
type Game struct {
	Players       []*Player   `json:"players"`
	Asteroids     []*Asteroid `json:"asteroids"`
	shutdownChann chan struct{}
	mutex         sync.Mutex
	running       bool
}

//NewGame ...
func NewGame() *Game {
	return &Game{
		Players:       []*Player{},
		Asteroids:     generateRandomNAsteroids(20),
		shutdownChann: make(chan struct{}),
		mutex:         sync.Mutex{},
		running:       true,
	}
}

//Start the game
func (g *Game) Start() {
	go g.loop()
}

//Tick one game Tick
func (g *Game) Tick() {
	for _, player := range g.Players {
		player.Tick()
	}

	state := g.constructState()

	for _, player := range g.Players {
		player.SendToPlayer(state)

		player.Explosions = []*explosion{}
	}
}

func (g *Game) loop() {
	ticker := time.NewTicker(16 * time.Millisecond)
loop:
	for {
		select {
		case <-g.shutdownChann:
			break loop
		case <-ticker.C:
			g.Tick()
		}
	}
}

//Shutdown the game
func (g *Game) Shutdown() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.running {
		g.shutdownChann <- struct{}{}
		g.running = false
		for _, player := range g.Players {
			player.Shutdown()
		}
	}
}

func (g *Game) constructState() []byte {
	byteArr, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err)
	}
	return byteArr
}

func generateRandomNAsteroids(n int) (asteroids []*Asteroid) {
	for i := 0; i < n; i++ {
		a := NewAsteroid(randomPos())
		asteroids = append(asteroids, a)
	}
	return
}
