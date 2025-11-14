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

var a = &AABB{
	Pos:  Vec{0, 250},
	Half: Vec{60, 20},
}
var b = &AABB{
	Pos:  Vec{100, 100},
	Half: Vec{12, 12},
}

var hit *coll.HitInfo = &coll.HitInfo{}

var collided bool

func main() {
	g := &Game{W: 800, H: 500}
	ebiten.SetWindowSize(int(g.W), int(g.H))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	W, H float64
}

func (g *Game) Update() error {
	av := v.Vec{2, 0}
	bv := Axis().Scale(6)
	hit.Reset()

	collided = coll.AABBAABBSweep2(a, b, av, bv, hit)
	if collided {

		bv = bv.Add(hit.Delta)

	}

	a.Pos = a.Pos.Add(av)

	if a.Pos.X > g.W {
		a.Pos.X = 0
	}

	b.Pos = b.Pos.Add(bv)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	vector.StrokeRect(
		s,
		float32(a.Pos.X-a.Half.X),
		float32(a.Pos.Y-a.Half.Y),
		float32(a.Half.X*2),
		float32(a.Half.Y*2),
		1,
		colornames.Gray,
		false,
	)
	vector.StrokeRect(
		s,
		float32(b.Pos.X-b.Half.X),
		float32(b.Pos.Y-b.Half.Y),
		float32(b.Half.X*2),
		float32(b.Half.Y*2),
		1,
		colornames.Gray,
		false,
	)

	if collided {
		// contact point
		px, py := float32(hit.Pos.X), float32(hit.Pos.Y)
		nx, ny := px+(float32(hit.Normal.X)*8), py+(float32(hit.Normal.Y)*8)
		vector.FillCircle(s, px, py, 2, colornames.Yellow, true)
		vector.StrokeLine(s, px, py, nx, ny, 1, colornames.Yellow, false)
	}

	ebitenutil.DebugPrint(s, fmt.Sprintf(
		"Delta: %v\nNormal: %v\nTime: %v\nPos: %v",
		hit.Delta,
		hit.Normal,
		hit.Time,
		hit.Pos,
	))
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return g.W, g.H
}

func Axis() (axis Vec) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		axis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		axis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		axis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		axis.X += 1
	}
	return
}
