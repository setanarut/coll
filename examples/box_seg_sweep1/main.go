package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const boxSpeed = 6

var (
	staticLine = &coll.Segment{
		v.Vec{0, 250},
		v.Vec{500, 250},
	}
	box      = coll.NewAABB(0, 0, 20, 30)
	hit      = &coll.Hit{}
	boxDelta   = v.FromAngle(math.Pi / 4).Scale(boxSpeed)
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

	hit.Reset()
	collided = coll.BoxSegmentSweep1(staticLine, box, boxDelta, hit)

	if collided {
		box.Pos = box.Pos.Add(slide(boxDelta, hit))
	} else {
		box.Pos = box.Pos.Add(boxDelta)
	}

	if box.Pos.X > 500 {
		box.Pos = v.Vec{}
		staticLine.B.Y = rand.Float64() * 500
	}

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.DrawSegment(s, staticLine, colornames.Gray)
	examples.StrokeBox(s, box, colornames.Green)
	examples.PrintHitInfoAt(s, hit, 10, 10, false)
	if collided {
		examples.DrawRay(s, box.Pos, hit.Normal, 30, colornames.Yellow, true)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func slide(delta v.Vec, hit *coll.Hit) (slideDelta v.Vec) {
	movementToHit := delta.Scale(hit.Data)
	remaining := delta.Sub(movementToHit)
	originalSpeed := remaining.Mag()
	slideDirection := remaining.Sub(hit.Normal.Scale(remaining.Dot(hit.Normal)))
	if slideDirection.MagSq() < coll.Epsilon {
		return movementToHit
	}
	slideDirectionUnit := slideDirection.Unit()
	scaledSlideDirection := slideDirectionUnit.Scale(originalSpeed)
	return movementToHit.Add(scaledSlideDirection)
}
