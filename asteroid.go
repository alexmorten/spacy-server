package server

import (
	"math/rand"
)

//Asteroid ...
type Asteroid struct {
	Pos      *Vector `json:"pos"`
	Capacity int     `json:"capacity"`
}

//NewAsteroid ...
func NewAsteroid(pos *Vector) *Asteroid {
	return &Asteroid{
		Pos:      pos,
		Capacity: rand.Intn(8) + 1,
	}
}
