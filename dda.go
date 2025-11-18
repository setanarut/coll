package coll

import (
	"image"
	"math"

	"github.com/setanarut/v"
)

// HitRayInfo structure stores the result of the DDA (Digital Differential Analysis) raycast operation
type HitRayInfo struct {
	// The distance from origin to the hit point
	Distance float64
	// The exact coordinates of where the ray hit
	Point v.Vec
	// The grid cell coordinates that were hit
	Cell image.Point
	// The surface normal at the hit point
	Normal v.Vec
}

// RaycastDDA performs the DDA (Digital Differential Analysis) algorithm to find intersections with a tile map
//
// youtube.com/watch?v=NbSee-XM7WA
//
// Parameters:
//   - start: Starting position of the ray
//   - dir: Direction unit vector of the ray (should be normalized)
//   - length: Maximum distance the ray can travel
//   - tileMap: 2D grid of cells where any non-zero value represents a wall/obstacle
func RaycastDDA(start, dir v.Vec, length float64, tileMap [][]uint8, cellSize float64, hit *HitRayInfo) bool {
	// Bitiş noktasını hesapla
	end := start.Add(dir.Scale(length))

	// DDA için delta değerlerini hesapla
	delta := end.Sub(start)
	steps := int(max(math.Abs(delta.X), math.Abs(delta.Y)))

	// Her adımdaki artış miktarı
	var inc v.Vec
	if steps != 0 {
		inc = delta.Scale(1.0 / float64(steps))
	}

	// Başlangıç noktası
	current := start

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
				hit.Point = current
				hit.Cell = cell
				hit.Distance = current.Sub(start).Mag()

				// Yüzey normalini hesapla
				// Hangi yüzeye çarptığını belirle (X veya Y yüzeyi)
				cellCenterX := float64(cell.X)*cellSize + cellSize/2
				cellCenterY := float64(cell.Y)*cellSize + cellSize/2

				// Çarpışma noktasının hücre merkezine göre pozisyonunu kullanarak
				// hangi yüzeye çarptığını belirle
				diffX := math.Abs(current.X - cellCenterX)
				diffY := math.Abs(current.Y - cellCenterY)

				if diffX > diffY {
					// X yüzeyine çarpış
					hit.Normal = v.Vec{X: -math.Copysign(1, dir.X), Y: 0}
				} else {
					// Y yüzeyine çarpış
					hit.Normal = v.Vec{X: 0, Y: -math.Copysign(1, dir.Y)}
				}

				return true
			}
		}

		// Bir sonraki piksele geç
		current = current.Add(inc)
	}

	// Çarpışma bulunamadı, maksimum bitiş noktasını döndür
	hit.Point = end
	// hit.Cell = image.Point{
	// 	X: int(end.X / cellSize),
	// 	Y: int(end.Y / cellSize),
	// }
	hit.Distance = length
	hit.Normal = v.Vec{} // Sıfır vektör
	return false
}
