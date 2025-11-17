package main

import (
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

var (
	box               = coll.NewAABB(250, 250, 50, 50)
	bodHitInfo        = &coll.HitInfo{}
	collided          bool
	segmentStartPoint = v.Vec{50, 50}
	cursor            = v.Vec{}
)

func main() {
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	cursor = examples.CursorPos()
	bodHitInfo.Reset()
	collided = coll.AABBSegmentOverlap(
		box,
		segmentStartPoint,
		cursor,
		v.Vec{},
		bodHitInfo,
	)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	colour := colornames.Green
	if collided {
		colour = colornames.Yellow
	}
	examples.StrokeAABB(s, box, colornames.Gray)
	examples.DrawHitInfo(s, bodHitInfo)
	vector.StrokeLine(s, float32(segmentStartPoint.X), float32(segmentStartPoint.Y), float32(cursor.X), float32(cursor.Y), 2, colour, true)
	examples.PrintHitInfoAt(s, bodHitInfo, 10, 10)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
