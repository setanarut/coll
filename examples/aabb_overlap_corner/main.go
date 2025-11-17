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

var (
	box  = coll.AABB{Half: v.Vec{16, 16}}
	wall = coll.AABB{Half: v.Vec{16, 16}, Pos: v.Vec{250, 250}}
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
	coll.AABBOverlap(&wall, &box, hit)

	box.Pos = box.Pos.Add(hit.Delta)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {

	colour := colornames.Green
	if collided {
		colour = colornames.Yellow
	}

	vector.StrokeRect(s, float32(wall.Pos.X-wall.Half.X), float32(wall.Pos.Y-wall.Half.Y), float32(wall.Half.X*2), float32(wall.Half.Y*2), 1, colornames.Gray, false)
	vector.StrokeRect(s, float32(box.Pos.X-box.Half.X), float32(box.Pos.Y-box.Half.Y), float32(box.Half.X*2), float32(box.Half.Y*2), 1, colour, false)
	ebitenutil.DebugPrint(s, fmt.Sprintf("Pos: %v Delta: %v Normal: %v Time: %v ", hit.Pos, hit.Delta, hit.Normal, hit.Time))
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
