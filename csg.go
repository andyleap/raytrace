package main

import (
	"image/color"
	"math"
)

type Hit struct {
	Ambient color.Color

	Diffusions []Diffusion

	Normal Vector
}

type Diffusion struct {
	Scatter float64
	Color   color.Color
}

type HitTest func(Vector) (float64, Hit)

var HitDefault = Hit{Ambient: color.RGBA{0, 0, 0, 0}}

func Uniform() HitTest {
	return func(pos Vector) (float64, Hit) {
		return 1, HitDefault
	}
}

func Sphere(radius float64) HitTest {
	return func(pos Vector) (float64, Hit) {
		h := HitDefault
		d := radius - math.Sqrt(pos.Dot(pos))
		h.Normal = pos.Normalize().Scale(-math.Copysign(1, d))
		return d, h
	}
}

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

func Cylinder(radius float64, height float64) HitTest {
	return func(pos Vector) (float64, Hit) {
		d := height/2 - math.Abs(pos.Y)
		h := HitDefault
		h.Normal = Vector{math.Copysign(1, pos.Y), 0, 0}
		dr := radius - math.Sqrt(pos.X*pos.X+pos.Z*pos.Z)
		if dr < d {
			d = dr
			h.Normal = Vector{pos.X, 0, pos.Z}.Normalize()
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
		for _, ht := range hts {
			newd, newh := ht(pos)
			if newd > 0 {
				newd = -newd
				if newd < d {
					d = newd
					h.Normal = newh.Normal.Scale(-1)
				}
			}
		}
		return d, h
	}
}

func Specular(c color.Color, ht HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := ht(pos)
		h.Diffusions = append(h.Diffusions, Diffusion{Scatter: 0, Color: c})
		return d, h
	}
}

func Diffuse(c color.Color, scatter float64, ht HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := ht(pos)
		h.Diffusions = append(h.Diffusions, Diffusion{Scatter: scatter, Color: c})
		return d, h
	}
}

func Ambient(c color.Color, ht HitTest) HitTest {
	return func(pos Vector) (float64, Hit) {
		d, h := ht(pos)
		h.Ambient = c
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
