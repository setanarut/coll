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

var wall = coll.NewAABB(250, 250, 100, 100)
var circle = coll.NewCircle(200, 200, 25)

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
	cursor = examples.CursorPos()
	velocity = cursor.Sub(circle.Pos)
	circle.Pos = circle.Pos.Add(velocity)

	hit.Reset()
	collided = coll.BoxCircleOverlap(wall, circle, hit)

	circle.Pos = circle.Pos.Add(hit.Delta)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeBox(s, wall, colornames.Gray)
	if collided {
		examples.StrokeCircleAt(s, cursor, circle.Radius, colornames.Red)
		examples.StrokeCircle(s, circle, colornames.Yellow)
		examples.DrawHitNormal(s, hit, colornames.Yellow, false)
	} else {
		examples.StrokeCircle(s, circle, colornames.Gray)
	}
	examples.PrintHitInfoAt(s, hit, 10, 10)

	ebitenutil.DebugPrintAt(
		s,
		fmt.Sprintf("Vel: %v\nBoxPos: %v", velocity, circle.Pos),
		10,
		100,
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
