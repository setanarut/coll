package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

var (
	circle1 = coll.NewCircle(250, 250, 30)
	cVel1   = v.Vec{4, 0}
)
var (
	circle2 = coll.NewCircle(250, 250, 10)
	cVel2   = v.Vec{-4, 0}
)

var collided bool

func main() {
	g := &Game{ScreenWidth: 500, H: 500}
	ebiten.SetWindowSize(int(g.ScreenWidth), int(g.H))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	ScreenWidth, H float64
}

func (g *Game) Update() error {

	cVel2 = examples.CursorPos().Sub(circle2.Pos)
	collided = coll.CircleCircleSweep2(
		circle1,
		circle2,
		cVel1,
		cVel2,
	)

	circle2.Pos = circle2.Pos.Add(cVel2)
	circle1.Pos = circle1.Pos.Add(cVel1)

	cx := circle1.Pos.X
	if cx < 0 || cx > g.ScreenWidth {
		cVel1.X *= -1
	}

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeCircle(s, circle1, colornames.Gray)
	if collided {
		examples.StrokeCircle(s, circle2, colornames.Yellow)
	} else {
		examples.StrokeCircle(s, circle2, colornames.Gray)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return g.ScreenWidth, g.H
}
