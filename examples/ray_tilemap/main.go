package main

import (
	"fmt"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

var screen = v.Vec{400, 180}

var (
	hitInfo = &coll.HitInfo{}
	IsHit   bool
	start   = screen.Scale(0.5)
	angle   float64
	TileMap = [][]uint8{
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 6, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 3, 1},
		{0, 1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 4, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 0, 0, 0},
		{4, 2, 8, 1, 1, 1, 1, 0, 1, 1, 1, 1, 88, 13, 1},
	}
	cellSize = 25
	dir      v.Vec
	coords   image.Point
	collided bool
)

type Game struct{}

func (g *Game) Update() error {
	angle += 0.01
	if angle >= 2*math.Pi {
		angle = 0
	}
	dir = v.FromAngle(angle)

	collided, coords = coll.RaycastDDA(start, dir, 400, TileMap, float64(cellSize), hitInfo)

	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	// draw tiles
	for y := range TileMap {
		for x := range TileMap[y] {
			if TileMap[y][x] != 0 {
				vector.StrokeRect(s,
					float32(x*cellSize),
					float32(y*cellSize),
					float32(cellSize),
					float32(cellSize),
					1,
					colornames.Gray,
					false)
			}
		}
	}

	examples.DrawSegment(s, start, hitInfo.Pos, colornames.Green)
	examples.DrawHitNormal(s, hitInfo, colornames.Yellow, true)

	// collision info
	examples.PrintHitInfoAt(s, hitInfo, 10, 10)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("Coords: %v", coords), 10, 100)
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
