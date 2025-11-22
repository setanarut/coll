package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxPointOverlap returns true if the point is in the box, false otherwise.
//
// If hitInfo is not nil, it will be filled with collision details.
//
// If a point is behind all of the edges of the box, it’s colliding.
// https://noonat.github.io/intersect/#aabb-vs-point
//
// hitInfo.Pos and hitInfo.Delta will be set to the nearest edge of the box.
//
// This code first finds the overlap on the X and Y axis. If the overlap is less than zero for either,
// a collision is not possible. Otherwise, we find the axis with the smallest overlap and use that to
// create an intersection point on the edge of the box.
func BoxPointOverlap(box *AABB, point v.Vec, hitInfo *HitInfo) bool {
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

	if hitInfo != nil {
		if px < py {
			sx := math.Copysign(1, dx)
			hitInfo.Delta.X = px * sx
			hitInfo.Normal.X = sx
			hitInfo.Pos.X = box.Pos.X + (box.Half.X * sx)
			hitInfo.Pos.Y = point.Y
		} else {
			sy := math.Copysign(1, dy)
			hitInfo.Delta.Y = py * sy
			hitInfo.Normal.Y = sy
			hitInfo.Pos.X = point.X
			hitInfo.Pos.Y = box.Pos.Y + (box.Half.Y * sy)
		}
	}
	return true
}
