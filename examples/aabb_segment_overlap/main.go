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
	box               = coll.NewAABB(250, 250, 16, 16)
	hit               = &coll.HitInfo{}
	collided          bool
	pos1, pos2, delta v.Vec
	angle             float64
)

// var center = v.Vec{250, 250}

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

	pos1.X = box.Pos.X + math.Cos(angle)*64
	pos1.Y = box.Pos.Y + math.Sin(angle)*64

	pos2.X = box.Pos.X + math.Sin(angle)*32
	pos2.Y = box.Pos.Y + math.Cos(angle)*32

	delta = pos2.Sub(pos1)

	hit.Reset()
	collided = coll.AABBSegmentOverlap(
		box,
		pos1,
		delta,
		v.Vec{},
		hit,
	)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeAABB(s, box, colornames.Gray)

	if collided {
		examples.DrawRay(s, pos1, delta.Unit(), delta.Mag(), colornames.Red, true)
		examples.DrawSegment(s, pos1, hit.Pos, colornames.Yellow)
		examples.DrawHitNormal(s, hit, colornames.Yellow, false)
	} else {
		examples.DrawRay(s, pos1, delta.Unit(), delta.Mag(), colornames.Green, true)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
