package main

import (
	"math"
)

type Vector struct {
	X, Y, Z float64
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vector) Sub(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vector) Mult(v2 Vector) Vector {
	return Vector{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

func (v Vector) Scale(s float64) Vector {
	return Vector{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector) Dot(v2 Vector) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v Vector) Normalize() Vector {
	return v.Scale(1.0 / math.Sqrt(v.Dot(v)))
}
