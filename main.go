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

func main() {
	fmt.Println("Hello World!")
	t := tinyfb.New("raytrace", 1024, 768)

	go func() {
		i := image.NewRGBA(image.Rect(0, 0, 1024, 768))
		start := Vector{0, 0, 0}
		geom := Translate(Vector{0, 0, 10},
			Union(
				Translate(Vector{-3, 0, 0},
					RotateX(math.Pi/4,
						Color(color.RGBA{0, 0, 255, 0},
							Box(Vector{5, 5, 5}),
						),
					),
				),
				Translate(Vector{3, 0, 0},
					Color(color.RGBA{0, 255, 0, 0},
						Sphere(2.5),
					),
				),
			),
		)
		x, y := 0, 0

		for {
			startTime := time.Now()
			if y < 768 {
				for {
					for c := 0; c < 100; c++ {
						i.Set(x, y, Trace(start, Vector{float64(x - 1024/2), float64(y - 768/2), 500}, geom, 5))
						x++
						if x >= 1024 {
							x = 0
							y++
						}
						if y >= 768 {
							break
						}
					}
					if time.Now().Sub(startTime) > 10*time.Millisecond {
						break
					}
					if y >= 768 {
						break
					}
				}
			}
			t.Update(i)
		}
	}()

	t.Run()
}

func Trace(start Vector, dir Vector, geom HitTest, bounces int) color.Color {
	pos := start
	dir = dir.Normalize()
	for count := 0; count < 100; count++ {
		d, h := geom(pos)
		d = -d
		//fmt.Println(d)
		if d > 1000000 {
			break
		}
		if d < 0.01 {
			if bounces <= 0 {
				return h.Color
			}

			newDir := dir.Add(h.Normal.Scale(h.Normal.Dot(dir) * -2))

			newC := Trace(pos.Add(newDir.Normalize().Scale(0.1)), newDir, geom, bounces-1)

			r, g, b, _ := h.Color.RGBA()
			nr, ng, nb, _ := newC.RGBA()

			return color.RGBA64{uint16((r + nr) / 2), uint16((g + ng) / 2), uint16((b + nb) / 2), 0xFFFF}
		}
		pos = pos.Add(dir.Scale(d))
	}
	return color.Gray{128}
}
