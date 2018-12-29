// RayTrace project main.go
package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/andyleap/tinyfb"
)

var (
	width  = 1024 * 3
	height = 768 * 3
)

func CheckerBoard(pos Vector) (float64, Hit) {
	d := -pos.Y

	return d, Hit{
		Ambient: color.Gray{(uint8(math.Mod(pos.X/4, 2)) + uint8(math.Mod(pos.Z/4, 2))) % 2 * 255},
		Normal:  Vector{0, 1, 0},
	}
}

type Light struct {
	Color    color.Color
	Pos      Vector
	Strength float64
}

type World struct {
	Geom   HitTest
	Lights []Light
}

func main() {
	fmt.Println("Hello World!")
	t := tinyfb.New("raytrace", int32(width), int32(height))

	go func() {
		i := image.NewRGBA(image.Rect(0, 0, width, height))
		start := Vector{0, 0, 0}
		geom := Union(
			Translate(Vector{0, -10, 0}, CheckerBoard),
			Translate(Vector{0, 0, 15},
				Subtract(Specular(color.NRGBA{255, 255, 255, 200}, Ambient(color.NRGBA{255, 255, 255, 128}, Sphere(5))),
					Ambient(color.NRGBA{255, 255, 255, 128}, Cylinder(2.5, 11)),
					Ambient(color.NRGBA{255, 255, 255, 128}, RotateX(math.Pi/2, Cylinder(2.5, 11))),
					Ambient(color.NRGBA{255, 255, 255, 128}, RotateZ(math.Pi/2, Cylinder(2.5, 11))),
				),
			),
		)

		w := &World{
			Geom: geom,
			Lights: []Light{
				{color.RGBA{255, 255, 255, 255}, Vector{0, 100, -50}, 150.0},
				{color.RGBA{0, 0, 255, 255}, Vector{0, 2, 15}, 4},
			},
		}
		x, y := 0, 0

		for {
			startTime := time.Now()
			if y < height {
				for {
					for c := 0; c < 100; c++ {
						i.Set(x, y,
							Trace(start, Vector{float64(x-width/2) / float64(height/2), float64(-(y - height/2)) / float64(height/2), 1.4}, w, 10),
						)
						x++
						if x >= width {
							x = 0
							y++
						}
						if y >= height {
							break
						}
					}
					if time.Now().Sub(startTime) > 100*time.Millisecond {
						break
					}
					if y >= height {
						break
					}
				}
			}
			if y < height {
				fmt.Println(y)
			}
			t.Update(i)
			time.Sleep(100*time.Millisecond - time.Now().Sub(startTime))
		}
	}()

	t.Run()
}

func Trace(start Vector, dir Vector, w *World, bounces int) color.Color {
	pos := start
	dir = dir.Normalize()
	for count := 0; count < 500; count++ {
		d, h := w.Geom(pos)
		d = -d
		//fmt.Println(d)
		if d > 1000000 {
			break
		}
		if d < 0.01 {
			r, g, b := uint32(0), uint32(0), uint32(0)
			ar, ag, ab, _ := h.Ambient.RGBA()
			for _, l := range w.Lights {
				if LoSTrace(pos, l.Pos, w) {
					path := l.Pos.Sub(pos)
					cr, cg, cb, _ := l.Color.RGBA()
					scale := h.Normal.Dot(path.Normalize())
					if scale < 0 {
						scale = 0
					}

					scale *= l.Strength * l.Strength
					scale /= path.Dot(path)

					r += uint32(float64(ar*cr/0xFFFF) * scale)
					g += uint32(float64(ag*cg/0xFFFF) * scale)
					b += uint32(float64(ab*cb/0xFFFF) * scale)
				}
			}

			if bounces > 0 {
				for _, d := range h.Diffusions {
					dr, dg, db := uint32(0), uint32(0), uint32(0)
					if d.Scatter == 0 {
						newDir := dir.Add(h.Normal.Scale(h.Normal.Dot(dir) * -2))
						dC := Trace(pos.Add(newDir.Normalize().Scale(0.1)), newDir, w, bounces-1)
						dr, dg, db, _ = dC.RGBA()
					} else {

					}
					mr, mg, mb, _ := d.Color.RGBA()
					r += (dr * mr / 0xFFFF)
					g += (dg * mg / 0xFFFF)
					b += (db * mb / 0xFFFF)
				}
			}
			if r > 0xFFFF {
				r = 0xFFFF
			}
			if g > 0xFFFF {
				g = 0xFFFF
			}
			if b > 0xFFFF {
				b = 0xFFFF
			}
			return color.RGBA64{uint16(r), uint16(g), uint16(b), 0xFFFF}
		}
		pos = pos.Add(dir.Scale(d))
	}
	return color.Black
}

func LoSTrace(start Vector, target Vector, w *World) bool {
	pos := start
	dir := target.Sub(start).Normalize()
	for count := 0; count < 500; count++ {
		d, _ := w.Geom(pos)
		d = -d
		if d < 0 {
			break
		}
		if d < 0.01 && count > 10 {
			break
		}
		if pos.Sub(target).Length() < d {
			return true
		}
		pos = pos.Add(dir.Scale(d))
	}
	return false
}
