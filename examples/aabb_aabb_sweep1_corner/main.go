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
	box            = coll.NewAABB(0, 0, 20, 20)
	wall           = coll.NewAABB(250, 250, 20, 20)
	hit            = &coll.HitInfo{}
	velocity       = v.Vec{3, 3}
	slidingEnabled = true
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
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		slidingEnabled = !slidingEnabled
	}

	collided = coll.AABBAABBSweep1(wall, box, velocity, hit)
	if collided {
		if slidingEnabled {
			velocity = examples.CalculateSlideVelocity(box, velocity, hit)
		} else {
			velocity = velocity.Add(hit.Delta)
			box.Pos = box.Pos.Add(velocity)

		}

	} else {
		box.Pos = box.Pos.Add(velocity)
		hit.Reset()
	}

	if box.Pos.X > 500 || box.Pos.Y > 500 {
		box.Pos = v.Vec{}
		velocity = v.Vec{3, 3}
		slidingEnabled = !slidingEnabled
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
	examples.PrintHitInfoAt(s, hit, 10, 30)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("Sliding Enabled: %v (Press S)", slidingEnabled), 10, 10)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
