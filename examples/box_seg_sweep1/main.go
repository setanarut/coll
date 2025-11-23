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

var (
	staticLine = &coll.Segment{
		v.Vec{0, 250},
		v.Vec{500, 250},
	}
	box      = coll.NewAABB(0, 0, 16, 16)
	hit      = &coll.HitInfo{}
	dir      = v.FromAngle(math.Pi / 4).Scale(6)
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

	collided = coll.BoxSegmentSweep1(staticLine, box, dir, hit)

	if collided {
		coll.CollideAndSlide(box, dir, hit)
	} else {
		box.Pos = box.Pos.Add(dir)
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
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
