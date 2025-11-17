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

var wall = &coll.AABB{Pos: v.Vec{250, 250}, Half: v.Vec{100, 25}}
var box = &coll.AABB{Pos: v.Vec{200, 200}, Half: v.Vec{25, 25}}

var hit = &coll.HitInfo{}

var collided bool
var velocity = v.Vec{}
var cursor = v.Vec{}

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

	curX, curY := ebiten.CursorPosition()
	cursor = v.Vec{float64(curX), float64(curY)}

	velocity = cursor.Sub(box.Pos)
	box.Pos = box.Pos.Add(velocity)

	hit.Reset()
	collided = coll.AABBOverlap(wall, box, hit)

	box.Pos = box.Pos.Add(hit.Delta)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeAABB(s, wall, colornames.Gray)
	if collided {
		examples.StrokeAABBAt(s, cursor, box.Half, colornames.Red)
		examples.StrokeAABB(s, box, colornames.Yellow)
		examples.DrawHitInfo(s, hit)
	} else {
		examples.StrokeAABB(s, box, colornames.Gray)
	}
	examples.PrintHitInfoAt(s, hit, 10, 10)

	ebitenutil.DebugPrintAt(
		s,
		fmt.Sprintf(
			"Vel: %v\nBoxPos: %v",
			velocity,
			box.Pos,
		),
		10,
		100,
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
