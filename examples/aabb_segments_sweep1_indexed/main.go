package main

import (
	"log"
	"math"
	"strconv"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

var (
	delta v.Vec
	angle float64

	staticLines = []*coll.Segment{
		{v.Vec{350, 261}, v.Vec{220, 210}},
		{v.Vec{350, 235}, v.Vec{220, 190}},
		{v.Vec{346, 278}, v.Vec{228, 235}},
	}

	box        = coll.NewAABB(218, 310, 16, 16)
	sweepDelta = v.Vec{32, -96}
	tempBox    = coll.NewAABB(250, 250, 16, 16)

	hit      = &coll.HitInfo{}
	collided bool
	index    int
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

	angle += 0.5 * math.Pi * 0.02
	factor := max((math.Cos(angle)+1)*0.5, 1e-8)
	delta = sweepDelta.Scale(factor)

	index = coll.AABBSegmentSweep1Indexed(staticLines, box, delta, hit)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {

	for _, line := range staticLines {
		examples.DrawLine(s, line.A, line.B, colornames.Gray)
	}

	if index >= 0 {
		// Draw a red box at the point where it was trying to move to
		examples.DrawRay(s, box.Pos, delta.Unit(), delta.Mag(), colornames.Red, true)
		tempBox.Pos = box.Pos.Add(delta)
		examples.StrokeAABB(s, tempBox, colornames.Red)

		// Draw a yellow box at the point it actually got to
		tempBox.Pos = box.Pos.Add(delta.Scale(hit.Time))
		examples.StrokeAABB(s, tempBox, colornames.Yellow)
		examples.DrawHitNormal(s, hit, colornames.Yellow, false)
		ebitenutil.DebugPrint(s, "line index: "+strconv.Itoa(index))

	} else {
		tempBox.Pos = box.Pos.Add(delta)
		examples.StrokeAABB(s, tempBox, colornames.Green)
		examples.DrawRay(s, box.Pos, delta.Unit(), delta.Mag(), colornames.Green, true)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
