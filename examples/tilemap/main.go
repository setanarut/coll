package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
)

type Vec = coll.Vec
type AABB = coll.AABB

const (
	screenWidth  = 640
	screenHeight = 480
)

var rect = AABB{
	Pos:  Vec{140, 130},
	Half: Vec{8, 8},
}

var (
	TileMap = [][]uint8{
		{0, 0, 0, 0, 0, 0, 9, 1},
		{0, 0, 0, 0, 0, 6, 0, 1},
		{4, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 8, 0, 8, 3, 1},
		{2, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 4, 0, 5, 0, 0, 0},
		{1, 4, 2, 8, 1, 88, 13, 1},
	}

	collider = coll.NewCollider(TileMap, screenWidth/8, screenHeight/8)
)

type Game struct {
}

func (g *Game) Update() error {

	// Get input axis
	vel := Axis()
	vel.Y *= 6
	vel.X *= 6

	// Collide with tiles
	delta := collider.Collide(rect, vel, nil)

	// Update player position
	rect.Pos.X += delta.X
	rect.Pos.Y += delta.Y

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw tiles
	for y := 0; y < len(TileMap); y++ {
		for x := 0; x < len(TileMap[y]); x++ {
			if TileMap[y][x] != 0 {
				vector.DrawFilledRect(screen,
					float32(x*collider.TileSize.X),
					float32(y*collider.TileSize.Y),
					float32(collider.TileSize.X),
					float32(collider.TileSize.Y),
					color.Gray{Y: 128},
					true)
			}
		}
	}

	// Draw player
	vector.DrawFilledRect(screen,
		float32(rect.Pos.X-rect.Half.X),
		float32(rect.Pos.Y-rect.Half.Y),
		float32(2*rect.Half.X),
		float32(2*rect.Half.Y),
		color.RGBA{47, 36, 254, 255},
		false)

	// Print collisions to the screen
	for i, c := range collider.Collisions {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf(
				"Tile ID: %d, Tile Coords: %v, Collision Normal: %v",
				c.TileID,
				c.TileCoords,
				c.Normal,
			), 20, 20+(i*20))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
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
