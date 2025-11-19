package main

import (
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	box         = coll.NewAABB(100, 250, 10, 100)
	boxVelocity = v.Vec{3, 0}
)
var (
	circle         = coll.NewCircle(100, 100, 20)
	circleVelocity = v.Vec{0, 0}
)

var hitInfo = &coll.HitInfo{}

var sweep bool

func main() {
	g := &Game{ScreenWidth: 800, H: 500}
	ebiten.SetWindowSize(int(g.ScreenWidth), int(g.H))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	ScreenWidth, H float64
}

func (g *Game) Update() error {

	circleVelocity = examples.CursorPos().Sub(circle.Pos)

	sweep = coll.AABBCircleSweep2(
		box,
		circle,
		boxVelocity,
		circleVelocity,
	)

	box.Pos = box.Pos.Add(boxVelocity)
	circle.Pos = circle.Pos.Add(circleVelocity)

	if box.Left() < 0 || box.Right() > g.ScreenWidth {
		boxVelocity.X *= -1
	}
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	if sweep {
		examples.StrokeCircle(s, circle, colornames.Yellow)
		examples.DrawHitNormal(s, hitInfo, colornames.Yellow, true)
	} else {
		examples.StrokeCircle(s, circle, colornames.Gray)
	}
	examples.StrokeAABB(s, box, colornames.Gray)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return g.ScreenWidth, g.H
}
