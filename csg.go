package main

import (
	"image/color"
	"math"
)

type Hit struct {
	Color  color.Color
	Normal Vector
}

type HitTest func(Vector) (float64, Hit)

var HitDefault = Hit{Color: color.White}

func Sphere(radius float64) HitTest {
	return func(pos Vector) (float64, Hit) {
		h := HitDefault
		d := radius - math.Sqrt(pos.Dot(pos))
		h.Normal = pos.Normalize().Scale(math.Copysign(1, d))
		return d, h
	}
}

// 2 - 3 = -1
// 3 - 2 = 1
// 6 - 4 = 2
// 4 - 6 = -2

// -10 - 2 =

// -5 - 95 = -100
//

func Box(size Vector) HitTest {
	b1 := size.Scale(0.5)
	b2 := size.Scale(-0.5)
	return func(pos Vector) (float64, Hit) {
		p1 := pos.Sub(b2)
		p2 := b1.Sub(pos)
		h := HitDefault

		d := p1.X
		h.Normal = Vector{1, 0, 0}
		if p2.X < d {
			d = p2.X
			h.Normal = Vector{-1, 0, 0}
		}
		if p1.Y < d {
			d = p1.Y
			h.Normal = Vector{0, 1, 0}
		}
		if p2.Y < d {
			d = p2.Y
			h.Normal = Vector{0, -1, 0}
		}
		if p1.Z < d {
			d = p1.Z
			h.Normal = Vector{0, 0, 1}
		}
		if p2.Z < d {
			d = p2.Z
			h.Normal = Vector{0, 0, -1}
		}
		return d, h
	}
}

func Union(hts ...HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := hts[0](pos)
		for _, ht := range hts[1:] {
			newd, newh := ht(pos)
			if newd > d {
				d, h = newd, newh
			}
		}
		return d, h
	}
}

func Intersect(hts ...HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := hts[0](pos)
		for _, ht := range hts[1:] {
			newd, newh := ht(pos)
			if newd < d {
				d, h = newd, newh
			}
		}
		return d, h
	}
}

func Subtract(ht HitTest, hts ...HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := ht(pos)
		c := h.Color
		i := false
		for _, ht := range hts {
			newd, newh := ht(pos)
			newh.Color = c
			if newd < d || !i {
				d, h = newd, newh
				i = true
			}
		}
		return d, h
	}
}

func Color(c color.Color, ht HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := ht(pos)
		h.Color = c
		return d, h
	}
}

func Translate(trans Vector, ht HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		return ht(pos.Sub(trans))
	}
}

func RotateX(rad float64, ht HitTest) HitTest {
	c := math.Cos(rad)
	s := math.Sin(rad)
	nc := math.Cos(-rad)
	ns := math.Sin(-rad)
	return func(pos Vector) (float64, Hit) {
		p := Vector{pos.X, pos.Y*c - pos.Z*s, pos.Y*s + pos.Z*c}
		d, h := ht(p)
		n := h.Normal
		h.Normal = Vector{n.X, n.Y*nc - n.Z*ns, n.Y*ns + n.Z*nc}
		return d, h
	}
}

func RotateY(rad float64, ht HitTest) HitTest {
	c := math.Cos(rad)
	s := math.Sin(rad)
	nc := math.Cos(-rad)
	ns := math.Sin(-rad)
	return func(pos Vector) (float64, Hit) {
		p := Vector{pos.X*c + pos.Z*s, pos.Y, -pos.X*s + pos.Z*c}
		d, h := ht(p)
		n := h.Normal
		h.Normal = Vector{n.X*nc + n.Z*ns, n.Y, -n.X*ns + n.Z*nc}
		return d, h
	}
}

func RotateZ(rad float64, ht HitTest) HitTest {
	c := math.Cos(rad)
	s := math.Sin(rad)
	nc := math.Cos(-rad)
	ns := math.Sin(-rad)
	return func(pos Vector) (float64, Hit) {
		p := Vector{pos.X*c - pos.Y*s, pos.X*s + pos.Y*c, pos.Z}
		d, h := ht(p)
		n := h.Normal
		h.Normal = Vector{n.X*nc - n.Y*ns, n.X*ns + n.Y*nc, n.Z}
		return d, h
	}
}
