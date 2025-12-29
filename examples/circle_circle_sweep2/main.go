package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

var (
	circle1    = coll.NewCircle(250, 250, 30)
	circle1Delta = v.Vec{2, 0}
)
var (
	circle2    = coll.NewCircle(250, 250, 30)
	circle2Delta = v.Vec{0, 0}
)

var circ2HitInfo = &coll.Hit{}

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

	circle1Delta = examples.CursorPos().Sub(circle1.Pos)

	circ2HitInfo.Reset()
	collided = coll.CircleCircleSweep2(
		circle1,
		circle2,
		circle1Delta,
		circle2Delta,
		circ2HitInfo,
	)

	circle1.Pos = circle1.Pos.Add(circle1Delta)
	circle2.Pos = circle2.Pos.Add(circle2Delta)

	cx := circle2.Pos.X
	if cx < 0 || cx > g.ScreenWidth {
		circle2Delta.X *= -1
	}

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeCircle(s, circle1, colornames.Gray)
	if collided {
		examples.StrokeCircle(s, circle2, colornames.Yellow)
		examples.DrawRay(s, circle2.Pos, circ2HitInfo.Normal, 20, color.White, true)
	} else {
		examples.StrokeCircle(s, circle2, colornames.Gray)
	}

	examples.PrintHitInfoAt(s, circ2HitInfo, 30, 30, false)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return g.ScreenWidth, g.H
}
