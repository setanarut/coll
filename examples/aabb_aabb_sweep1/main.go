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
	velocity       = v.Vec{}
	cursorPosition = v.Vec{}
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
	cursorPosition = examples.CursorPos()

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		slidingEnabled = !slidingEnabled
	}

	velocity = cursorPosition.Sub(box.Pos)

	hit.Reset()
	collided = coll.AABBAABBSweep1(wall, box, velocity, hit)

	if collided {
		if slidingEnabled {
			newVel := examples.CalculateSlideVelocity(box, velocity, hit) // Sliding correction
			box.Pos = box.Pos.Add(newVel)
		} else {
			newVel := velocity.Scale(hit.Time) // or velocity.Add(hit.Delta)
			box.Pos = box.Pos.Add(newVel)
		}
	} else {
		box.Pos = box.Pos.Add(velocity)
	}
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
			"Sliding Enabled: %v (Press Tab)\nVel: %v\nBoxPos: %v",
			slidingEnabled,
			velocity,
			box.Pos,
		),
		10,
		10,
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
