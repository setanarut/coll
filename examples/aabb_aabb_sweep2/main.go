package main

import (
	"fmt"
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

var (
	box            = coll.NewAABB(0, 0, 16, 16)
	wall           = coll.NewAABB(250, 250, 16*4, 16)
	hit            = &coll.HitInfo{}
	wallVelocity   = v.Vec{X: 3}
	boxVelocity    = v.Vec{}
	slidingEnabled bool
	collided       bool
)

func main() {
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
}

func (g *Game) Update() error {
	boxVelocity = examples.Axis().Scale(3)

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		slidingEnabled = !slidingEnabled
	}

	hit.Reset()
	wall.Pos = wall.Pos.Add(wallVelocity)
	collided = coll.AABBAABBSweep2(wall, box, wallVelocity, boxVelocity, hit)
	if collided {
		if slidingEnabled {
			boxVelocity = examples.CalculateSlideVelocity(box, boxVelocity, hit)
		} else {
			boxVelocity = boxVelocity.Add(hit.Delta)
		}
	}

	if wall.Right() > 500 {
		wallVelocity = wallVelocity.NegX()
	}
	if wall.Left() < 0 {
		wallVelocity = wallVelocity.NegX()
	}

	box.Pos = box.Pos.Add(boxVelocity)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeAABB(s, wall, colornames.Gray)
	if collided {
		examples.StrokeAABB(s, box, colornames.Yellow)
		examples.DrawHitInfo(s, hit)
	} else {
		examples.StrokeAABB(s, box, colornames.Green)
	}
	examples.PrintHitInfoAt(s, hit, 10, 100)
	ebitenutil.DebugPrintAt(
		s,
		fmt.Sprintf(
			"Sliding Enabled: %v (Press S)\nVel: %v\nBoxPos: %v",
			slidingEnabled,
			boxVelocity,
			box.Pos,
		),
		10,
		10,
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
