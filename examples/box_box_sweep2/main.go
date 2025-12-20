package main

import (
	"fmt"
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

var (
	box          = coll.NewAABB(250, 100, 16, 16)
	wall         = coll.NewAABB(250, 250, 16*4, 16)
	hit          = &coll.Hit{}
	wallVelocity = v.Vec{X: 5}
	boxVelocity  = v.Vec{}
	collided     bool
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
	boxVelocity = examples.Axis().Unit().Scale(5)
	hit.Reset()
	collided = coll.BoxBoxSweep2(wall, box, wallVelocity, boxVelocity, hit)
	if collided {

		box.Pos = box.Pos.Add(boxVelocity.Scale(hit.Data))
		box.Pos = box.Pos.Add(wallVelocity)

		if hit.Normal.Y == -1 {
			boxVelocity.Y = 0
		}
		if hit.Normal.X != 0 {
			boxVelocity.X = 0
		}
	} else {
		box.Pos = box.Pos.Add(boxVelocity)
	}

	if wall.Right() > 400 {
		wallVelocity = wallVelocity.NegX()
	}
	if wall.Left() < 100 {
		wallVelocity = wallVelocity.NegX()
	}

	wall.Pos = wall.Pos.Add(wallVelocity)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeBox(s, wall, colornames.Gray)
	if collided {
		examples.StrokeBox(s, box, colornames.Yellow)
		// Draw normal
		examples.DrawRay(s, box.Pos, hit.Normal, 30, colornames.White, true)
	} else {
		examples.StrokeBox(s, box, colornames.Green)
	}
	ebitenutil.DebugPrintAt(
		s,
		fmt.Sprintf(
			"WASD = Move\naVel: %v\nBoxPos: %v",
			boxVelocity,
			box.Pos,
		),
		10,
		10,
	)
	examples.PrintHitInfoAt(s, hit, 10, 100, false)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
