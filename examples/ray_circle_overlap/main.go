package main

import (
	"log"
	"math"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	circle                              = coll.NewCircle(250, 250, 20)
	raySeg, resultSegment, longRayDebug coll.Segment
	hit                                 bool
	angle                               float64
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

	angle += 0.5 * math.Pi * 0.005

	raySeg.A.X = circle.Pos.X + math.Cos(angle)*64
	raySeg.A.Y = circle.Pos.Y + math.Sin(angle)*64
	raySeg.B.X = circle.Pos.X + math.Sin(angle)*32
	raySeg.B.Y = circle.Pos.Y + math.Cos(angle)*32

	hit = coll.RayCircleOverlap(&raySeg, circle, &resultSegment)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	ResizeSegment(&raySeg, &longRayDebug, 1000)
	examples.StrokeCircle(s, circle, colornames.Gray)
	if hit {
		examples.DrawSegment(s, &longRayDebug, colornames.Red)
		examples.DrawSegment(s, &raySeg, colornames.Orange)
		examples.DrawSegment(s, &resultSegment, colornames.Yellow)
		examples.FillCircleAt(s, resultSegment.A, 5, colornames.Yellow)
		examples.FillCircleAt(s, resultSegment.B, 5, colornames.Yellow)
	} else {
		examples.DrawSegment(s, &longRayDebug, colornames.Green)
		examples.DrawSegment(s, &raySeg, colornames.Lightgreen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func ResizeSegment(s *coll.Segment, out *coll.Segment, targetLength float64) {
	currentLen := s.A.Dist(s.B)
	delta := (targetLength - currentLen) / 2.0
	dir := s.B.Sub(s.A).Unit()
	offset := dir.Scale(delta)
	out.A = s.A.Sub(offset)
	out.B = s.B.Add(offset)
}
