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
	circle    = coll.NewCircle(250, 250, 100)
	cursorPos = v.Vec{}
	hit       = &coll.Hit{}
)

var In bool

type Game struct{}

func (g *Game) Update() error {
	cursorPos = examples.CursorPos()

	hit.Reset()
	In = coll.CirclePointOverlap(circle, cursorPos, hit)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeCircle(s, circle, colornames.Gray)

	examples.FillCircleAt(s, cursorPos, 3, colornames.White)
	if In {
		contactPoint := circle.Pos.Add(hit.Normal.Scale(circle.Radius))
		examples.DrawRay(s, contactPoint, hit.Normal, 30, colornames.Yellow, true)
		examples.FillCircleAt(s, contactPoint, 3, colornames.Red)
	}

	if In {
		examples.FillCircleAt(s, cursorPos, 3, colornames.Yellow)
	} else {
		examples.FillCircleAt(s, cursorPos, 3, colornames.Red)
	}
	examples.PrintHitInfoAt(s, hit, 10, 10, true)

}
func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func main() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
