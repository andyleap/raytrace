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

func (v Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Rotate(axis Vector, angle float64) Vector {
	c, s := math.Cos(angle), math.Sin(angle)
	return Vector{
		v.X*(c+axis.X*axis.X*(1-c)) + v.Y*(axis.X*axis.Y*(1-c)-axis.Z*s) + v.Z*(axis.X*axis.Z*(1-c)+axis.Y*s),
		v.X*(axis.Y*axis.X*(1-c)+axis.Z*s) + v.Y*(c+axis.Y*axis.Y*(1-c)) + v.Z*(axis.Y*axis.Z*(1-c)-axis.X*s),
		v.X*(axis.Z*axis.X*(1-c)-axis.Y*s) + v.Y*(axis.Z*axis.Y*(1-c)+axis.X*s) + v.Z*(c+axis.Z*axis.Z*(1-c)),
	}
}
