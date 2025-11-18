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
	box  = &coll.AABB{Half: v.Vec{32, 32}}
	wall = &coll.AABB{Half: v.Vec{32, 32}, Pos: v.Vec{250, 250}}
)

var hit = &coll.HitInfo{}
var collided bool
var velocity = v.Vec{2, 2}

func main() {
	g := &Game{W: 500, H: 500}
	ebiten.SetWindowSize(int(g.W), int(g.H))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	W, H float64
}

func (g *Game) Update() error {
	box.Pos = box.Pos.Add(velocity)
	hit.Reset()
	collided = coll.AABBOverlap(wall, box, hit)
	box.Pos = box.Pos.Add(hit.Delta)

	if box.Pos.X > 500 || box.Pos.Y > 500 {
		box.Pos = v.Vec{}
	}
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeAABB(s, wall, colornames.Gray)
	if collided {
		examples.StrokeAABB(s, box, colornames.Yellow)
		examples.DrawHitNormal(s, hit)
	} else {
		examples.StrokeAABB(s, box, colornames.Green)
	}
	examples.PrintHitInfoAt(s, hit, 10, 200)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
