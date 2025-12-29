package main

import (
	"image/color"
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

var box = coll.NewAABB(250, 250, 40, 40)
var circle = coll.NewCircle(200, 200, 25)

var hit = &coll.Hit{}

var screenBox = coll.NewAABB(250, 250, 220, 220)

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
	delta = cursor.Sub(circle.Pos)
	circle.Pos = circle.Pos.Add(delta)

	hit.Reset()
	collided = coll.BoxCircleOverlap(box, circle, hit)

	box.Pos = box.Pos.Add(hit.Normal.Neg().Scale(hit.Data))

	if !coll.BoxBoxContain(screenBox, box) {
		box.Pos.X = 250
		box.Pos.Y = 250
	}
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeBox(s, box, colornames.Gray)
	if collided {
		examples.StrokeCircleAt(s, cursor, circle.Radius, colornames.Red)
		examples.StrokeCircle(s, circle, colornames.Yellow)
		examples.DrawRay(s, box.Pos, hit.Normal, 30, color.White, true)
	} else {
		examples.StrokeCircle(s, circle, colornames.Gray)
	}

	examples.StrokeBox(s, screenBox, colornames.Red)

	examples.PrintHitInfoAt(s, hit, 40, 40, true)
	ebitenutil.DebugPrintAt(s, "Push the box with the cursor.", 160, 100)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
