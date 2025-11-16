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

	colour := colornames.Green
	if collided {
		colour = colornames.Yellow
	}

	vector.StrokeRect(
		s,
		float32(wall.Pos.X-wall.Half.X),
		float32(wall.Pos.Y-wall.Half.Y),
		float32(wall.Half.X*2),
		float32(wall.Half.Y*2),
		1,
		colornames.Gray,
		false,
	)

	if collided {
		// contact point
		px, py := float32(hit.Pos.X), float32(hit.Pos.Y)
		nx, ny := px+(float32(hit.Normal.X)*8), py+(float32(hit.Normal.Y)*8)

		vector.StrokeRect(s, float32(cursor.X-box.Half.X), float32(cursor.Y-box.Half.Y), float32(box.Half.X*2), float32(box.Half.Y*2), 1, colornames.Red, false)
		vector.FillCircle(s, px, py, 2, colornames.Yellow, true)
		vector.StrokeLine(s, px, py, nx, ny, 1, colornames.Yellow, false)
	}

	vector.StrokeRect(s, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), 1, colour, false)
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
