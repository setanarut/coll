package coll

import (
	"math"

	"github.com/setanarut/mathutils"
	"github.com/setanarut/v"
)

// AABBCircleOverlap checks whether box and circle overlap.
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
// AABBCircleSweep2() should be used instead.
//
// If you only need to know whether a collision occurred, pass nil for hitInfo
// to skip generating collision details.
func AABBCircleOverlap(box *AABB, circle *Circle, hitInfo *HitInfo) bool {
	// Use AABBCircleIntersect for the initial intersection check
	if !AABBCircleIntersect(box, circle) {
		return false
	}

	// If we only need a boolean result, return early
	if hitInfo == nil {
		return true
	}

	// Calculate detailed collision information
	// 1. Vector between circle center and box center (Local Space)
	d := circle.Pos.Sub(box.Pos)

	// 2. Clamp the d vector within the box boundaries (Half extents)
	// This gives us the closest point on the box to the circle center
	closest := v.Vec{
		X: mathutils.Clamp(d.X, -box.Half.X, box.Half.X),
		Y: mathutils.Clamp(d.Y, -box.Half.Y, box.Half.Y),
	}

	// 3. Check if the circle center is inside the box
	// If the closest point equals the original distance vector, the circle center is inside
	inside := d.Equals(closest)

	// --- Outside Case (General Collision) ---
	if !inside {
		// The closest point is on the box surface
		// normal: vector from closest point to circle center
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
		// --- Inside Case (Center Inside Box) ---
		// Circle center is completely inside the box
		// Use AABB-AABB overlap logic to push out from the closest edge

		// This logic is exactly the same as in AABBAABBOverlap
		// We push out from the edge we're closest to
		absD := d.Abs()
		px := box.Half.X - absD.X
		py := box.Half.Y - absD.Y

		if px < py {
			sx := math.Copysign(1, d.X)

			hitInfo.Delta = v.Vec{X: px * sx, Y: 0}
			hitInfo.Normal = v.Vec{X: sx, Y: 0}

			// Edge point on X axis
			hitInfo.Pos = v.Vec{
				X: box.Pos.X + box.Half.X*sx,
				Y: circle.Pos.Y,
			}
		} else {
			sy := math.Copysign(1, d.Y)

			hitInfo.Delta = v.Vec{X: 0, Y: py * sy}
			hitInfo.Normal = v.Vec{X: 0, Y: sy}

			// Edge point on Y axis
			hitInfo.Pos = v.Vec{
				X: circle.Pos.X,
				Y: box.Pos.Y + box.Half.Y*sy,
			}
		}
	}

	return true
}
