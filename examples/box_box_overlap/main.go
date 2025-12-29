package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

var wall = coll.NewAABB(250, 250, 100, 25)
var box = coll.NewAABB(200, 200, 25, 25)

var hit = &coll.Hit{}

var collided bool
var delta = v.Vec{}
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
	delta = cursor.Sub(box.Pos)
	box.Pos = box.Pos.Add(delta)

	hit.Reset()
	collided = coll.BoxBoxOverlap(wall, box, hit)
	box.Pos = box.Pos.Add(hit.Normal.Scale(hit.Data))
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeBox(s, wall, colornames.Gray)
	if collided {
		examples.StrokeBoxAt(s, cursor, box.Half, colornames.Red)
		examples.StrokeBox(s, box, colornames.Yellow)

		// DrawNormal
		contactEdgePos := box.Pos.Sub(hit.Normal.Mul(box.Half))
		examples.DrawRay(s, contactEdgePos, hit.Normal, 20, color.White, true)

	} else {
		examples.StrokeBox(s, box, colornames.Gray)
	}
	examples.PrintHitInfoAt(s, hit, 10, 10, true)

	ebitenutil.DebugPrintAt(
		s,
		fmt.Sprintf("Delta: %v\nBoxPos: %v", delta, box.Pos),
		10,
		100,
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
