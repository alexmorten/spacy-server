package server

import (
	"math"
)

//Vector for points and movement
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//Copy a Vector
func (v *Vector) Copy() *Vector {
	return &Vector{
		X: v.X,
		Y: v.Y,
	}
}

//Add the supplied vector to the vector add is called on
func (v *Vector) Add(v2 *Vector) *Vector {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

//Sub the supplied vector from the vector sub is called on
func (v *Vector) Sub(v2 *Vector) *Vector {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}

//Mul the vector with r
func (v *Vector) Mul(r float64) *Vector {
	v.X *= r
	v.Y *= r
	return v
}

//Div the vector through r
func (v *Vector) Div(r float64) *Vector {
	v.X /= r
	v.Y /= r
	return v
}

// Dist between v and v2
func (v *Vector) Dist(v2 *Vector) float64 {
	dx := v2.X - v.X
	dy := v2.Y - v.Y
	return math.Sqrt(dx*dx + dy*dy)
}

//Length of v
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

//Normalize 's the Vector to length 1
func (v *Vector) Normalize() *Vector {
	l := v.Length()
	v.X /= l
	v.Y /= l
	return v
}

//ToLength sets the Vector to length r
func (v *Vector) ToLength(r float64) *Vector {
	v.Normalize().Mul(r)
	return v
}
