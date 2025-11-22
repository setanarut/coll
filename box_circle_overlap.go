package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxCircleOverlap checks whether box and circle overlap.
// Any collision information written to hitInfo always describes how to move circle out of box.
//
// It uses a separating-axis test: if the circle do not overlap on either X or Y,
// there is no collision and the function returns false.
//
// If hitInfo is not nil, the function fills it with:
//   - Delta: the minimum vector needed to push boxA out of boxB
//   - Normal: the direction in which boxA is pushed
//   - Pos: an approximate contact position on the collision side
//
// This method can behave poorly for moving objects. For continuous motion,
// AABBCircleSweep2 should be used instead.
//
// If you only need to know whether a collision occurred, pass nil for hitInfo
// to skip generating collision details.
func BoxCircleOverlap(box *AABB, circle *Circle, hitInfo *HitInfo) bool {
	if !boxCircleIntersect(box, circle) {
		return false
	}

	if hitInfo == nil {
		return true
	}
	d := circle.Pos.Sub(box.Pos)
	closest := v.Vec{
		X: math.Max(-box.Half.X, math.Min(d.X, box.Half.X)),
		Y: math.Max(-box.Half.Y, math.Min(d.Y, box.Half.Y)),
	}
	inside := d.Equals(closest)

	if !inside {
		normal := d.Sub(closest)
		distSq := normal.MagSq()
		dist := math.Sqrt(distSq)

		// Normalize the normal vector
		hitInfo.Normal = normal.DivS(dist)

		// Penetration amount: Radius - Distance from center to surface
		penetration := circle.Radius - dist
		hitInfo.Delta = hitInfo.Normal.Scale(penetration)

		// Contact point is the closest point on the box (converted to world coordinates)
		hitInfo.Pos = box.Pos.Add(closest)

	} else {
		absD := d.Abs()
		px := box.Half.X - absD.X
		py := box.Half.Y - absD.Y

		if px < py {
			sx := math.Copysign(1, d.X)
			pushDistance := px + circle.Radius
			hitInfo.Delta = v.Vec{X: pushDistance * sx, Y: 0}
			hitInfo.Normal = v.Vec{X: sx, Y: 0}
			hitInfo.Pos = v.Vec{
				X: box.Pos.X + box.Half.X*sx,
				Y: circle.Pos.Y,
			}
		} else {
			sy := math.Copysign(1, d.Y)

			pushDistance := py + circle.Radius

			hitInfo.Delta = v.Vec{X: 0, Y: pushDistance * sy}
			hitInfo.Normal = v.Vec{X: 0, Y: sy}

			hitInfo.Pos = v.Vec{
				X: circle.Pos.X,
				Y: box.Pos.Y + box.Half.Y*sy,
			}
		}
	}
	return true
}
