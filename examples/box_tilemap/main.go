package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var rect = coll.AABB{
	Pos:  v.Vec{10, 10},
	Half: v.Vec{8, 8},
}

type Tile struct {
	ID uint8
}

func (t *Tile) IsSolid() bool {
	return t.ID != 0
}

var (
	tiles = [][]uint8{
		{0, 0, 0, 1, 1, 0, 9, 1},
		{0, 0, 0, 0, 0, 6, 0, 1},
		{4, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 3, 1},
		{2, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 4, 0, 0, 0, 0, 0},
		{1, 4, 2, 8, 1, 3, 3, 1},
	}

	tileMap  [][]coll.Tile
	collider *coll.TileCollider
)

func init() {
	tileMap = make([][]coll.Tile, len(tiles))

	for y := range tiles {
		tileMap[y] = make([]coll.Tile, len(tiles[y]))
		for x, id := range tiles[y] {
			tileMap[y][x] = &Tile{ID: id}
		}
	}

	collider = &coll.TileCollider{
		TileMap:  tileMap,
		CellSize: image.Point{64, 64},
	}
}

type Game struct {
}

func (g *Game) Update() error {

	// Get input axis
	delta := examples.Axis()
	delta.Y *= 6
	delta.X *= 6

	// Collide with tiles
	allowedDelta := collider.Collide(rect, delta, nil)

	// Update player position
	rect.Pos.X += allowedDelta.X
	rect.Pos.Y += allowedDelta.Y

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw tiles
	for y := range len(tileMap) {
		for x := range len(tileMap[y]) {
			if collider.TileMap[y][x].IsSolid() {
				vector.FillRect(screen,
					float32(x*collider.CellSize.X),
					float32(y*collider.CellSize.Y),
					float32(collider.CellSize.X),
					float32(collider.CellSize.Y),
					color.Gray{Y: 128},
					true)
			}
		}
	}

	// Draw player
	vector.FillRect(screen,
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
				collider.TileMap[c.TileCoords.X][c.TileCoords.Y],
				c.TileCoords,
				c.Normal,
			), 20, 20+(i*20))
	}

	ebitenutil.DebugPrint(screen, "Controls: WASD")

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
