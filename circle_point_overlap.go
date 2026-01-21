package coll

import (
	"math"

	"github.com/setanarut/v"
)

// CirclePointOverlap returns true if the point is in the circle, false otherwise.
// If h is not nil, the function fills it with for point:
//   - Normal: Collision surface normal for circle
//   - Data: the penetration depth for point
func CirclePointOverlap(circle *Circle, point v.Vec, h *Hit) bool {
	dx := point.X - circle.Pos.X
	dy := point.Y - circle.Pos.Y

	distSq := dx*dx + dy*dy
	radiusSq := circle.Radius * circle.Radius

	if distSq > radiusSq {
		return false
	}

	if h != nil {
		if distSq == 0 {
			h.Normal.X = 1
			h.Normal.Y = 0
			h.Data = circle.Radius
		} else {
			dist := math.Sqrt(distSq)
			h.Data = circle.Radius - dist

			invDist := 1.0 / dist
			h.Normal.X = dx * invDist
			h.Normal.Y = dy * invDist
		}
	}
	return true
}
