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
// To resolve the overlap:
//
//	newBoxPos = box.Pos.Add(hit.Normal.Neg().Scale(hit.Data))
func BoxCircleOverlap(a *AABB, c *Circle, h *Hit) bool {

	// intersection test
	diff := c.Pos.Sub(a.Pos)
	clamped := v.Vec{
		X: max(-a.Half.X, min(diff.X, a.Half.X)),
		Y: max(-a.Half.Y, min(diff.Y, a.Half.Y)),
	}
	closest := a.Pos.Add(clamped)
	if !(c.Pos.DistSq(closest) <= c.Radius*c.Radius) {
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
		penetration := c.Radius - dist
		h.Data = penetration

	} else {
		// Box center is inside circle - push along shortest axis
		absD := diff.Abs()
		px := a.Half.X - absD.X
		py := a.Half.Y - absD.Y

		if px < py {
			// Push horizontally
			sx := math.Copysign(1, diff.X)
			pushDistance := px + c.Radius
			h.Normal = v.Vec{X: sx, Y: 0} // Box pushed away from circle
			h.Data = pushDistance
		} else {
			// Push vertically
			sy := math.Copysign(1, diff.Y)
			pushDistance := py + c.Radius
			h.Normal = v.Vec{X: 0, Y: sy} // Box pushed away from circle
			h.Data = pushDistance
		}
	}
	return true
}
