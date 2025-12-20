package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxCircleOverlap checks whether box and circle overlap.
//
// If h is not nil, the function fills it with:
//   - Normal: the surface normal of box
//   - Data: penetration depth (distance to move the box)
//
// To resolve the overlap: newBoxPos = box.Pos.Add(hit.Normal.Neg().Scale(hit.Data))
func BoxCircleOverlap(box *AABB, circle *Circle, h *Hit) bool {

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

	if h == nil {
		return true
	}

	inside := diff.Equals(clamped)

	if !inside {
		// Box is outside circle - push along the line from closest point to circle center
		normal := diff.Sub(clamped)
		distSq := normal.MagSq()
		dist := math.Sqrt(distSq)

		// Normalize the normal vector (direction to push box)
		h.Normal = normal.DivS(dist) // Box pushed away from circle

		// Penetration depth (how far to move)
		penetration := circle.Radius - dist
		h.Data = penetration

	} else {
		// Box center is inside circle - push along shortest axis
		absD := diff.Abs()
		px := box.Half.X - absD.X
		py := box.Half.Y - absD.Y

		if px < py {
			// Push horizontally
			sx := math.Copysign(1, diff.X)
			pushDistance := px + circle.Radius
			h.Normal = v.Vec{X: sx, Y: 0} // Box pushed away from circle
			h.Data = pushDistance
		} else {
			// Push vertically
			sy := math.Copysign(1, diff.Y)
			pushDistance := py + circle.Radius
			h.Normal = v.Vec{X: 0, Y: sy} // Box pushed away from circle
			h.Data = pushDistance
		}
	}
	return true
}
