package coll

import (
	"math"

	"github.com/setanarut/mathutils"
	"github.com/setanarut/v"
)

// AABBCircleOverlap checks whether box and circle overlap.
// Any collision information written to hitInfo always describes how to move circle out of box.
//
// It uses a closest-point on AABB test.
func AABBCircleOverlap(box *AABB, circle *Circle, hitInfo *HitInfo) bool {
	// 1. Daire merkezi ile Kutu merkezi arasındaki vektör (Local Space)
	d := circle.Pos.Sub(box.Pos)

	// 2. Kutu sınırları (Half extents) içinde d vektörünü "kelepçele" (Clamp).
	// Bu işlem bize kutu üzerinde daire merkezine en yakın noktayı verir.
	closest := v.Vec{
		X: mathutils.Clamp(d.X, -box.Half.X, box.Half.X),
		Y: mathutils.Clamp(d.Y, -box.Half.Y, box.Half.Y),
	}

	// 3. Daire merkezinin kutunun içinde olup olmadığını kontrol et.
	// Eğer en yakın nokta ile orijinal mesafe vektörü aynıysa, daire merkezi kutunun içindedir.
	inside := false
	if d.Equals(closest) {
		inside = true
	}

	// --- Dışarıda Olma Durumu (Genel Çarpışma) ---
	if !inside {
		// closest noktası kutu yüzeyindedir.
		// normal: closest noktasından daire merkezine giden vektör.
		normal := d.Sub(closest)
		distSq := normal.MagSq()
		radiusSq := circle.Radius * circle.Radius

		// Mesafe yarıçaptan büyükse çarpışma yoktur.
		if distSq > radiusSq {
			return false
		}

		// Çarpışma var ancak detay istenmiyor.
		if hitInfo == nil {
			return true
		}

		dist := math.Sqrt(distSq)

		// Normali normalize et
		hitInfo.Normal = normal.DivS(dist)

		// Penetrasyon miktarı: Yarıçap - Merkezden yüzeye olan uzaklık
		penetration := circle.Radius - dist
		hitInfo.Delta = hitInfo.Normal.Scale(penetration)

		// Temas noktası kutu üzerindeki closest noktasıdır (Dünya koordinatlarına çevrilir)
		hitInfo.Pos = box.Pos.Add(closest)

	} else {
		// --- İçeride Olma Durumu (Merkez Kutunun İçinde) ---
		// Daire merkezi kutunun tamamen içinde.
		// AABBAABBOverlap mantığıyla en yakın kenardan dışarı itiyoruz.

		if hitInfo == nil {
			return true
		}

		// Burası AABBAABBOverlap mantığı ile birebir aynıdır.
		// Hangi eksene (X veya Y) daha yakınsak o taraftan dışarı iteriz.
		absD := d.Abs()
		px := box.Half.X - absD.X
		py := box.Half.Y - absD.Y

		if px < py {
			sx := math.Copysign(1, d.X)

			hitInfo.Delta = v.Vec{X: px * sx, Y: 0}
			hitInfo.Normal = v.Vec{X: sx, Y: 0}

			// X eksenindeki kenar noktası
			hitInfo.Pos = v.Vec{
				X: box.Pos.X + box.Half.X*sx,
				Y: circle.Pos.Y,
			}
		} else {
			sy := math.Copysign(1, d.Y)

			hitInfo.Delta = v.Vec{X: 0, Y: py * sy}
			hitInfo.Normal = v.Vec{X: 0, Y: sy}

			// Y eksenindeki kenar noktası
			hitInfo.Pos = v.Vec{
				X: circle.Pos.X,
				Y: box.Pos.Y + box.Half.Y*sy,
			}
		}
	}

	return true
}
