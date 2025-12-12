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
// BoxCircleSweep2() should be used instead.
//
// If you only need to know whether a collision occurred, pass nil for hitInfo
// to skip generating collision details.
func BoxCircleOverlap(box *AABB, circle *Circle, hitInfo *HitInfo) bool {

	// intersection test
	diff := circle.Pos.Sub(box.Pos)
	clamped := v.Vec{
		X: max(-box.Half.X, min(diff.X, box.Half.X)),
		Y: max(-box.Half.Y, min(diff.Y, box.Half.Y)),
	}
	closest := box.Pos.Add(clamped)
	if !(circle.Pos.DistSq(closest) <= circle.Radius*circle.Radius) {
		return false
	}

	if hitInfo == nil {
		return true
	}

	inside := diff.Equals(clamped)

	if !inside {
		normal := diff.Sub(clamped)
		distSq := normal.MagSq()
		dist := math.Sqrt(distSq)

		// Normalize the normal vector
		hitInfo.Normal = normal.DivS(dist)

		// Penetration amount: Radius - Distance from center to surface
		penetration := circle.Radius - dist
		hitInfo.Delta = hitInfo.Normal.Scale(penetration)

		// Contact point is the closest point on the box (converted to world coordinates)
		hitInfo.Pos = box.Pos.Add(clamped)

	} else {
		absD := diff.Abs()
		px := box.Half.X - absD.X
		py := box.Half.Y - absD.Y

		if px < py {
			sx := math.Copysign(1, diff.X)
			pushDistance := px + circle.Radius
			hitInfo.Delta = v.Vec{X: pushDistance * sx, Y: 0}
			hitInfo.Normal = v.Vec{X: sx, Y: 0}
			hitInfo.Pos = v.Vec{
				X: box.Pos.X + box.Half.X*sx,
				Y: circle.Pos.Y,
			}
		} else {
			sy := math.Copysign(1, diff.Y)

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
