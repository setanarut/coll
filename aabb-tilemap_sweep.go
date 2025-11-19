package coll

import (
	"image"
	"math"

	"github.com/setanarut/v"
)

// TileHitInfo stores information about a collision with a tile
type TileHitInfo struct {
	TileCoords image.Point // X,Y coordinates of the tile in the tilemap
	Normal     v.Vec       // Normal vector of the collision (-1/0/1)
}

// TileCollider handles collision detection between AABB and [][]uint8 2D tilemap
type TileCollider struct {
	Collisions     []TileHitInfo // List of collisions from last check
	CellSize       image.Point   // Width and height of tiles
	TileMap        [][]uint8     // 2D grid of tile IDs
	NonSolidTileID uint8         // Sets the ID of non-solid tiles. Defaults to 0.
}

// NewTileCollider creates a new tile collider with the given tilemap and tile dimensions
func NewTileCollider(tileMap [][]uint8, tileWidth, tileHeight int) *TileCollider {
	return &TileCollider{
		TileMap:  tileMap,
		CellSize: image.Point{tileWidth, tileHeight},
	}
}

// TileCollisionCallback is called when collisions occur, receiving collision info and final movement
type TileCollisionCallback func([]TileHitInfo, float64, float64)

// Collide checks for collisions when a moving aabb and returns the allowed movement
func (c *TileCollider) Collide(box AABB, delta v.Vec, onCollide TileCollisionCallback) v.Vec {
	c.Collisions = c.Collisions[:0]

	if delta.X == 0 && delta.Y == 0 {
		return delta
	}

	if math.Abs(delta.X) > math.Abs(delta.Y) {
		if delta.X != 0 {
			delta.X = c.CollideX(&box, delta.X)
		}
		if delta.Y != 0 {
			box.Pos.X += delta.X
			delta.Y = c.CollideY(&box, delta.Y)
		}
	} else {
		if delta.Y != 0 {
			delta.Y = c.CollideY(&box, delta.Y)
		}
		if delta.X != 0 {

			box.Pos.Y += delta.Y
			delta.X = c.CollideX(&box, delta.X)
		}
	}

	if onCollide != nil {
		onCollide(c.Collisions, delta.X, delta.Y)
	}

	return delta
}

// CollideX checks for collisions along the X axis and returns the allowed X movement
func (c *TileCollider) CollideX(aabb *AABB, deltaX float64) float64 {
	checkLimit := max(1, int(math.Ceil(math.Abs(deltaX)/float64(c.CellSize.Y)))+1)

	rectTop := aabb.Top()
	rectBottom := aabb.Bottom()

	rectTileTopCoord := int(math.Floor(rectTop / float64(c.CellSize.Y)))
	rectTileBottomCoord := int(math.Ceil((rectBottom)/float64(c.CellSize.Y))) - 1

	if deltaX > 0 {
		startRightX := int(math.Floor((aabb.Pos.X + aabb.Half.X) / float64(c.CellSize.X)))
		endX := startRightX + checkLimit
		endX = min(endX, len(c.TileMap[0]))

		for y := rectTileTopCoord; y <= rectTileBottomCoord; y++ {
			if y < 0 || y >= len(c.TileMap) {
				continue
			}
			for x := startRightX; x < endX; x++ {
				if x < 0 || x >= len(c.TileMap[0]) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileLeft := float64(x * c.CellSize.X)
					collision := tileLeft - (aabb.Pos.X + aabb.Half.X)
					if collision <= deltaX {
						deltaX = collision
						c.Collisions = append(c.Collisions, TileHitInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Left,
						})
					}
				}
			}
		}
	}

	if deltaX < 0 {
		rectLeft := aabb.Left()

		endX := int(math.Floor(rectLeft / float64(c.CellSize.X)))
		startX := endX - checkLimit
		startX = max(startX, 0)

		for y := rectTileTopCoord; y <= rectTileBottomCoord; y++ {
			if y < 0 || y >= len(c.TileMap) {
				continue
			}
			for x := startX; x <= endX; x++ {
				if x < 0 || x >= len(c.TileMap[0]) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileRight := float64((x + 1) * c.CellSize.X)
					collision := tileRight - rectLeft
					if collision >= deltaX {
						deltaX = collision
						c.Collisions = append(c.Collisions, TileHitInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Right,
						})
					}
				}
			}
		}
	}

	return deltaX
}

// CollideY checks for collisions along the Y axis and returns the allowed Y movement
func (c *TileCollider) CollideY(rect *AABB, deltaY float64) float64 {

	checkLimit := max(1, int(math.Ceil(math.Abs(deltaY)/float64(c.CellSize.Y)))+1)

	rectLeft := rect.Left()
	rectRight := rect.Right()

	rectTileLeftCoord := int(math.Floor(rectLeft / float64(c.CellSize.X)))
	rectTileRightCoord := int(math.Ceil(rectRight/float64(c.CellSize.X))) - 1

	if deltaY > 0 {
		rectBottom := rect.Pos.Y + rect.Half.Y
		startBottomY := int(math.Floor(rectBottom / float64(c.CellSize.Y)))
		endY := startBottomY + checkLimit
		endY = min(endY, len(c.TileMap))

		for x := rectTileLeftCoord; x <= rectTileRightCoord; x++ {
			if x < 0 || x >= len(c.TileMap[0]) {
				continue
			}
			for y := startBottomY; y < endY; y++ {
				if y < 0 || y >= len(c.TileMap) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileTop := float64(y * c.CellSize.Y)
					collision := tileTop - rectBottom
					if collision <= deltaY {
						deltaY = collision
						c.Collisions = append(c.Collisions, TileHitInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Up,
						})
					}
				}
			}
		}
	}

	if deltaY < 0 {
		rectTop := rect.Top()
		endY := int(math.Floor(rectTop / float64(c.CellSize.Y)))
		startY := endY - checkLimit
		startY = max(startY, 0)

		for x := rectTileLeftCoord; x <= rectTileRightCoord; x++ {
			if x < 0 || x >= len(c.TileMap[0]) {
				continue
			}
			for y := startY; y <= endY; y++ {
				if y < 0 || y >= len(c.TileMap) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileBottom := float64((y + 1) * c.CellSize.Y)
					collision := tileBottom - rectTop
					if collision >= deltaY {
						deltaY = collision
						c.Collisions = append(c.Collisions, TileHitInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Down,
						})
					}
				}
			}
		}
	}
	return deltaY
}
