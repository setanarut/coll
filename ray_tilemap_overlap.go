package coll

import (
	"image"
	"math"

	"github.com/setanarut/v"
)

// RaycastDDA performs the DDA (Digital Differential Analysis) algorithm to find intersections with a tile map
//
// youtube.com/watch?v=NbSee-XM7WA
//
// Parameters:
//   - pos: Starting position of the ray
//   - dir: Direction unit vector of the ray (should be normalized)
//   - length: Maximum distance the ray can travel
//   - tileMap: 2D grid of cells where any non-zero value represents a wall/obstacle
//
// Returns:
//   - bool: True if a collision occurred
//   - image.Point: The grid coordinates of the wall that was hit (0,0 if no hit)
func RaycastDDA(pos, dir v.Vec, length float64, tileMap [][]uint8, cellSize float64, hit *HitInfo) (bool, image.Point) {
	// Bitiş noktasını hesapla
	end := pos.Add(dir.Scale(length))

	// DDA için delta değerlerini hesapla
	delta := end.Sub(pos)
	steps := int(max(math.Abs(delta.X), math.Abs(delta.Y)))

	// Her adımdaki artış miktarı
	var inc v.Vec
	if steps != 0 {
		inc = delta.Scale(1.0 / float64(steps))
	}

	// Başlangıç noktası
	current := pos

	// Her piksel için kontrol et
	for i := 0; i <= steps; i++ {
		// Grid hücresini bul
		cell := image.Point{
			X: int(current.X / cellSize),
			Y: int(current.Y / cellSize),
		}

		// Sınırları kontrol et
		if cell.X >= 0 && cell.X < len(tileMap[0]) && cell.Y >= 0 && cell.Y < len(tileMap) {
			if tileMap[cell.Y][cell.X] != 0 {
				// Çarpışma noktasını kaydet
				hit.Pos = current
				hit.Normal = v.Vec{}
				hit.Delta = v.Vec{}
				hit.Time = 0

				// Mesafe ve zaman (0..1) hesapla
				distance := current.Sub(pos).Mag()
				if length > 0 {
					hit.Time = distance / length
				}

				// Kalan hareket vektörü (başlangıçtan hedefe kadar - gittiğimiz mesafe)
				remaining := length - distance
				if remaining < 0 {
					remaining = 0
				}
				hit.Delta = dir.Scale(remaining)

				// Yüzey normalini hesapla
				cellCenterX := float64(cell.X)*cellSize + cellSize/2
				cellCenterY := float64(cell.Y)*cellSize + cellSize/2

				diffX := math.Abs(current.X - cellCenterX)
				diffY := math.Abs(current.Y - cellCenterY)

				if diffX > diffY {
					// X yüzeyine çarpış
					hit.Normal = v.Vec{X: -math.Copysign(1, dir.X), Y: 0}
				} else {
					// Y yüzeyine çarpış
					hit.Normal = v.Vec{X: 0, Y: -math.Copysign(1, dir.Y)}
				}

				// Return true and the cell coordinate
				return true, cell
			}
		}

		// Bir sonraki piksele geç
		current = current.Add(inc)
	}

	// Çarpışma bulunamadı
	hit.Pos = end
	hit.Normal = v.Vec{}
	hit.Delta = v.Vec{}
	if length > 0 {
		hit.Time = 1.0
	} else {
		hit.Time = 0.0
	}

	// Return false and an empty point
	return false, image.Point{}
}
