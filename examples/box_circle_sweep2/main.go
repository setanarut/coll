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

var box = coll.NewAABB(250, 250, 20, 20)
var circle = coll.NewCircle(200, 200, 16)

var hit = &coll.Hit{}

var screenBox = coll.NewAABB(250, 250, 220, 220)

var collided bool
var circleDelta = v.Vec{}
var boxDelta = v.Vec{0, 0}
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
	circleDelta = cursor.Sub(circle.Pos)
	circle.Pos = circle.Pos.Add(circleDelta)

	hit.Reset()
	collided = coll.BoxCircleSweep2(box, circle, boxDelta, circleDelta, hit)

	if collided {
		// Push the box
		box.Pos = box.Pos.Add(circleDelta.Scale(1.0 - hit.Data))
	} else {
		box.Pos = box.Pos.Add(boxDelta)
	}

	// Ekran sınırları kontrolü
	if !coll.BoxBoxContain(screenBox, box) {
		box.Pos.X = 250
		box.Pos.Y = 250
		boxDelta = v.Vec{0, 0}
	}

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeBox(s, box, colornames.Gray)
	if collided {
		examples.StrokeCircleAt(s, cursor, circle.Radius, colornames.Red)
		examples.StrokeCircle(s, circle, colornames.Yellow)
		examples.DrawRay(s, box.Pos, hit.Normal, 20, color.White, true)
	} else {
		examples.StrokeCircle(s, circle, colornames.Gray)
	}

	examples.StrokeBox(s, screenBox, colornames.Red)

	examples.PrintHitInfoAt(s, hit, 40, 40, false)
	ebitenutil.DebugPrintAt(s, "Push the box with the cursor.", 160, 100)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
