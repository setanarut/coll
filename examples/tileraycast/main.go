package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

type Vec = v.Vec

var screen = Vec{600, 600}

var (
	hitRayInfo = &coll.HitRayInfo{}
	IsHit      bool
	start      = screen.Scale(0.6)
	angle      float64
	TileMap    = [][]uint8{
		{0, 1, 1, 1, 6, 0, 1},
		{0, 0, 1, 0, 1, 0, 0},
		{0, 0, 0, 0, 8, 3, 1},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0},
		{0, 4, 0, 5, 0, 0, 0},
		{4, 2, 8, 1, 88, 13, 1},
	}
	cellSize = 100
)

type Game struct{}

func (g *Game) Update() error {
	angle += 0.02
	if angle >= 2*math.Pi {
		angle = 0
	}
	dir := v.FromAngle(angle)

	coll.RaycastDDA(start, dir, 400, TileMap, float64(cellSize), hitRayInfo) // length'i 600 yaptık
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	for y := range TileMap {
		for x := range TileMap[y] {
			if TileMap[y][x] != 0 {
				vector.StrokeRect(s,
					float32(x*cellSize),
					float32(y*cellSize),
					float32(cellSize),
					float32(cellSize),
					1,
					colornames.Blue,
					false)
			}
		}
	}

	// ray
	vector.StrokeLine(
		s,
		float32(start.X),
		float32(start.Y),
		float32(hitRayInfo.Point.X),
		float32(hitRayInfo.Point.Y),
		2,
		colornames.Yellow,
		true)

	// dot
	vector.DrawFilledCircle(
		s,
		float32(hitRayInfo.Point.X),
		float32(hitRayInfo.Point.Y),
		5,
		colornames.Red,
		true,
	)
	//normal
	vector.StrokeLine(
		s,
		float32(hitRayInfo.Point.X),
		float32(hitRayInfo.Point.Y),
		float32(hitRayInfo.Point.X+hitRayInfo.Normal.X*30),
		float32(hitRayInfo.Point.Y+hitRayInfo.Normal.Y*30),
		2,
		colornames.Yellow,
		true)

	ebitenutil.DebugPrint(s, fmt.Sprintf("%v", hitRayInfo.Cell))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 0, 0
}
func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (float64, float64) {
	return screen.X, screen.Y
}

func main() {
	ebiten.SetWindowSize(int(screen.X), int(screen.Y))
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
