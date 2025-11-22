package main

import (
	"log"
	"math"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	circle             = coll.NewCircle(250, 250, 20)
	intersectionPoints []v.Vec
	seg                coll.Segment
	angle              float64
)

func main() {
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {

	angle += 0.5 * math.Pi * 0.02

	seg.A.X = circle.Pos.X + math.Cos(angle)*64
	seg.A.Y = circle.Pos.Y + math.Sin(angle)*64
	seg.B.X = circle.Pos.X + math.Sin(angle)*32
	seg.B.Y = circle.Pos.Y + math.Cos(angle)*32

	intersectionPoints = coll.SegmentCircleOverlap(&seg, circle)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeCircle(s, circle, colornames.Gray)
	if intersectionPoints != nil {
		examples.DrawSegment(s, &seg, colornames.Red)
		for i := range intersectionPoints {
			examples.FillCircleAt(s, intersectionPoints[i], 3, colornames.Yellow)
		}
	} else {
		examples.DrawSegment(s, &seg, colornames.Green)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
