package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxPointOverlap returns true if the point is in the box, false otherwise.
// If h is not nil, the function fills it with for point:
//   - Normal: Collision surface normal for box
//   - Data: the penetration depth for point
func BoxPointOverlap(box *AABB, point v.Vec, h *Hit) bool {
	dx := point.X - box.Pos.X
	px := box.Half.X - math.Abs(dx)

	if px <= 0 {
		return false
	}

	dy := point.Y - box.Pos.Y
	py := box.Half.Y - math.Abs(dy)

	if py <= 0 {
		return false
	}

	if h != nil {
		if px < py {
			h.Normal.X = math.Copysign(1, dx)
			h.Data = px
		} else {
			h.Normal.Y = math.Copysign(1, dy)
			h.Data = py
		}
	}
	return true
}
