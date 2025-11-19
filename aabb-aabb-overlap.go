package coll

import "math"

// AABBAABBOverlap checks whether boxA and boxB overlap.
// Any collision information written to hitInfo always describes how to move boxA out of boxB.
//
// It uses a separating-axis test: if the boxes do not overlap on either X or Y,
// there is no collision and the function returns false.
//
// If hitInfo is not nil, the function fills it with:
//   - Delta: the minimum vector needed to push boxA out of boxB
//   - Normal: the direction in which boxA is pushed
//   - Pos: an approximate contact position on the collision side
//
// This method can behave poorly for moving objects. For continuous motion,
// sweepAABB should be used instead.
//
// If you only need to know whether a collision occurred, pass nil for hitInfo
// to skip generating collision details.
func AABBAABBOverlap(boxA, boxB *AABB, hitInfo *HitInfo) bool {

	dx := boxB.Pos.X - boxA.Pos.X
	px := boxB.Half.X + boxA.Half.X - math.Abs(dx)

	if px <= 0 {
		return false
	}

	dy := boxB.Pos.Y - boxA.Pos.Y
	py := boxB.Half.Y + boxA.Half.Y - math.Abs(dy)

	if py <= 0 {
		return false
	}

	if hitInfo == nil {
		return true
	}

	// if if hitInfo is not nil, fill
	if px < py {
		sx := math.Copysign(1, dx)
		hitInfo.Delta.X = px * sx
		hitInfo.Normal.X = sx
		hitInfo.Pos.X = boxA.Pos.X + boxA.Half.X*sx
		hitInfo.Pos.Y = boxB.Pos.Y
	} else {
		sy := math.Copysign(1, dy)
		hitInfo.Delta.Y = py * sy
		hitInfo.Normal.Y = sy
		hitInfo.Pos.X = boxB.Pos.X
		hitInfo.Pos.Y = boxA.Pos.Y + boxA.Half.Y*sy
	}
	return true
}
