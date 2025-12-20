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
	rayPos   = screen.Scale(0.5)
	rayDir   v.Vec
	rayMag   = 100.0
	hit      = &coll.Hit{}
	cellSize = 25
	TileMap  = [][]uint8{
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 6, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 3, 1},
		{0, 1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 4, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 0, 0, 0},
		{4, 2, 8, 1, 1, 1, 1, 0, 1, 1, 1, 1, 88, 13, 1},
	}
)

var (
	IsHit    bool
	angle    float64
	coords   image.Point
	collided bool
)

type Game struct{}

func (g *Game) Update() error {
	angle += 0.01
	if angle >= 2*math.Pi {
		angle = 0
	}
	rayDir = v.FromAngle(angle)

	hit.Reset()
	collided, coords = coll.RayTilemapDDA(rayPos, rayDir, rayMag, TileMap, float64(cellSize), hit)

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

	// Draw full ray
	examples.DrawRay(s, rayPos, rayDir, rayMag, colornames.White, true)

	// Draw hit segment
	collisionPoint := rayPos.Add(rayDir.Scale(rayMag * hit.Data))
	examples.DrawLine(s, rayPos, collisionPoint, colornames.Lime)

	// Draw collision point
	examples.FillCircleAt(s, collisionPoint, 3, colornames.Red)

	// collision info
	examples.PrintHitInfoAt(s, hit, 10, 10, false)
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
