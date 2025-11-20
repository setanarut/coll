package main

import (
	"log"
	"math"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	angle      float64
	staticLine = &coll.Segment{
		A: v.Vec{300, 200},
		B: v.Vec{230, 250},
	}

	box = coll.NewAABB(250, 310, 16, 16)

	sweepDelta = v.Vec{32, -96}
	delta      v.Vec

	tempBox  = coll.NewAABB(250, 250, 16, 16)
	hit      = &coll.HitInfo{}
	collided bool
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
	factor := max((math.Cos(angle)+1)*0.5, 1e-8)
	delta = sweepDelta.Scale(factor)

	collided = coll.AABBSegmentSweep1(staticLine, box, delta, hit)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.DrawSegment(s, staticLine, colornames.Gray)

	if collided {
		// Draw a red box at the point where it was trying to move to
		examples.DrawRay(s, box.Pos, delta.Unit(), delta.Mag(), colornames.Red, true)
		tempBox.Pos = box.Pos.Add(delta)
		examples.StrokeAABB(s, tempBox, colornames.Red)

		// Draw a yellow box at the point it actually got to
		tempBox.Pos = box.Pos.Add(delta.Scale(hit.Time))
		examples.StrokeAABB(s, tempBox, colornames.Yellow)
		examples.DrawHitNormal(s, hit, colornames.Yellow, false)

	} else {
		tempBox.Pos = box.Pos.Add(delta)
		examples.StrokeAABB(s, tempBox, colornames.Green)
		examples.DrawRay(s, box.Pos, delta.Unit(), delta.Mag(), colornames.Green, true)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
