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
	width  = 1024 * 2
	height = 768 * 2
)

func CheckerBoard(pos Vector) (float64, Hit) {
	d := -(10 - pos.Y)

	return d, Hit{
		Ambient: color.Gray{(uint8(math.Mod(pos.X/4, 2)) + uint8(math.Mod(pos.Z/4, 2))) % 2 * 255},
		Normal:  Vector{0, 1, 0},
	}
}

func main() {
	fmt.Println("Hello World!")
	t := tinyfb.New("raytrace", int32(width), int32(height))

	go func() {
		i := image.NewRGBA(image.Rect(0, 0, width, height))
		start := Vector{0, 0, 0}
		geom := Union(
			CheckerBoard,
			Translate(Vector{0, 0, 10},
				Specular(color.NRGBA{255, 255, 255, 200}, Sphere(3)),
			),
			Translate(Vector{8, 0, 10},
				Specular(color.NRGBA{255, 255, 255, 200}, Sphere(3)),
			),
			Translate(Vector{-8, 0, 10},
				Specular(color.NRGBA{255, 255, 255, 200}, Sphere(3)),
			),
			Translate(Vector{4, -8, 10},
				Specular(color.NRGBA{255, 255, 255, 200}, Sphere(3)),
			),
			Translate(Vector{-4, -8, 10},
				Specular(color.NRGBA{255, 255, 255, 200}, Sphere(3)),
			),

		)
		x, y := 0, 0

		for {
			startTime := time.Now()
			if y < height {
				for {
					for c := 0; c < 100; c++ {
						i.Set(x, y,
							Trace(start, Vector{float64(x - width/2), float64(y - height/2), 500}, geom, 10),
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
					if y >= width {
						break
					}
				}
			}
			fmt.Println(y)
			t.Update(i)
			time.Sleep(100*time.Millisecond - time.Now().Sub(startTime))
		}
	}()

	t.Run()
}

func Trace(start Vector, dir Vector, geom HitTest, bounces int) color.Color {
	pos := start
	dir = dir.Normalize()
	for count := 0; count < 500; count++ {
		d, h := geom(pos)
		d = -d
		//fmt.Println(d)
		if d > 1000000 {
			break
		}
		if d < 0.01 {
			r, g, b, _ := h.Ambient.RGBA()
			if bounces > 0 {
				for _, d := range h.Diffusions {
					dr, dg, db := uint32(0), uint32(0), uint32(0)
					if d.Scatter == 0 {
						newDir := dir.Add(h.Normal.Scale(h.Normal.Dot(dir) * -2))
						dC := Trace(pos.Add(newDir.Normalize().Scale(0.1)), newDir, geom, bounces-1)
						dr, dg, db, _ = dC.RGBA()
					} else {
						
						
						
					}
					mr, mg, mb, _ := d.Color.RGBA()
					r += (dr * mr / 0xFFFF)
					g += (dg * mg / 0xFFFF)
					b += (db * mb / 0xFFFF)
				}
			}
			return color.RGBA64{uint16(r), uint16(g), uint16(b), 0xFFFF}
		}
		pos = pos.Add(dir.Scale(d))
	}
	return color.Gray{128}
}
