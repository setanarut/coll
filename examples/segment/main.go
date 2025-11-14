package main

import (
	"fmt"
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Vec = v.Vec
type AABB = coll.AABB

var box1 = AABB{
	Pos:  Vec{200, 200},
	Half: Vec{100, 100},
}

var hit *coll.HitInfo
var collided bool

var pos = Vec{50, 50}
var cursor = Vec{}

func main() {
	g := &Game{W: 900, H: 500}

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
	cursor = Vec{float64(curX), float64(curY)}

	delta := cursor.Sub(pos)

	hit = &coll.HitInfo{}
	collided = coll.Segment(&box1, pos, delta, v.Vec{}, hit)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {

	colour := colornames.Green
	if collided {
		colour = colornames.Yellow
	}

	vector.StrokeRect(
		s,
		float32(box1.Pos.X-box1.Half.X),
		float32(box1.Pos.Y-box1.Half.Y),
		float32(box1.Half.X*2),
		float32(box1.Half.Y*2),
		1,
		colornames.Gray,
		false,
	)
	px, py := float32(hit.Pos.X), float32(hit.Pos.Y)
	nx, ny := px+(float32(hit.Normal.X)*8), py+(float32(hit.Normal.Y)*8)
	vector.FillCircle(s, px, py, 2, colornames.Yellow, true)
	vector.StrokeLine(s, float32(pos.X), float32(pos.Y), float32(cursor.X), float32(cursor.Y), 1, colour, true)
	vector.StrokeLine(s, px, py, nx, ny, 1, colornames.Yellow, true)

	ebitenutil.DebugPrint(s, fmt.Sprintf(
		"Pos: %v Delta: %v Normal: %v Time: %v ",
		hit.Pos,
		hit.Delta,
		hit.Normal,
		hit.Time,
	))
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return g.W, g.H
}
