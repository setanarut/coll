package main

import (
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

var (
	staticBox = coll.NewAABB(250, 250, 100, 100)
	cursorBox = coll.NewAABB(250, 250, 32, 32)
)

var IsContain bool

type Game struct{}

func (g *Game) Update() error {

	cursorBox.Pos = examples.CursorPos()
	IsContain = coll.BoxBoxContain(staticBox, cursorBox)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	examples.StrokeBox(s, staticBox, colornames.Gray)
	if IsContain {
		examples.FillBox(s, cursorBox, colornames.Green)
		ebitenutil.DebugPrint(s, "inside")
	} else {
		ebitenutil.DebugPrint(s, "outside")
		examples.FillBox(s, cursorBox, colornames.Red)
	}

}
func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func main() {
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
