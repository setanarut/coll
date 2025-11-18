package main

import (
	"log"
	"math"

	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	angle, factor float64
	staticBox     = coll.NewAABB(250, 250, 16*7, 16)
	sweepBoxes    = [2]*coll.AABB{
		coll.NewAABB(98, 274, 16, 16),
		coll.NewAABB(128+250, -48+250, 16, 16),
	}
	tempBoxes = [2]*coll.AABB{
		coll.NewAABB(250, 250, 16, 16),
		coll.NewAABB(250, 250, 16, 16),
	}
	sweepDeltas = [2]v.Vec{{64, -12}, {-32, 96}}
	deltas      = [2]v.Vec{}

	hitInfos = [2]*coll.HitInfo{{}, {}}
	sweeps   [2]bool
)

func main() {
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
}

func (g *Game) Update() error {
	angle += 0.5 * math.Pi * 0.02
	factor = math.Max((math.Cos(angle)+1)*0.5, 1e-8)

	for i, box := range sweepBoxes {
		deltas[i] = sweepDeltas[i].Scale(factor)
		hitInfos[i].Reset()
		sweeps[i] = coll.AABBAABBSweep1(staticBox, box, deltas[i], hitInfos[i])
	}
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {

	examples.StrokeAABB(s, staticBox, colornames.Gray)

	for i, box := range sweepBoxes {
		examples.StrokeAABB(s, box, colornames.Gray)
		if sweeps[i] {
			// Draw a red box at the point where it was trying to move to
			examples.DrawRay(s, box.Pos, deltas[i].Unit(), deltas[i].Mag(), colornames.Red, true)
			tempBoxes[i].Pos = box.Pos.Add(deltas[i])

			examples.StrokeAABB(s, tempBoxes[i], colornames.Red)

			// Draw a yellow box at the point it actually got to
			tempBoxes[i].Pos = box.Pos.Add(deltas[i].Scale(hitInfos[i].Time))
			examples.StrokeAABB(s, tempBoxes[i], colornames.Yellow)
			examples.DrawHitNormal(s, hitInfos[i], colornames.Yellow, false)

		} else {
			tempBoxes[i].Pos = box.Pos.Add(deltas[i])
			examples.StrokeAABB(s, tempBoxes[i], colornames.Green)
			examples.DrawRay(s, box.Pos, deltas[i].Unit(), deltas[i].Mag(), colornames.Green, true)
		}

	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
