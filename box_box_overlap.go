package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxBoxOverlap checks whether boxA and boxB overlap.
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
func BoxBoxOverlap(boxA, boxB *AABB, hitInfo *HitInfo) bool {
	d := boxB.Pos.Sub(boxA.Pos)
	absD := d.Abs()
	hSum := boxA.Half.Add(boxB.Half)

	px := hSum.X - absD.X

	if px <= 0 {
		return false
	}

	py := hSum.Y - absD.Y

	if py <= 0 {
		return false
	}

	if hitInfo == nil {
		return true
	}

	if px < py {
		sx := math.Copysign(1, d.X)

		hitInfo.Delta = v.Vec{X: px * sx, Y: 0}
		hitInfo.Normal = v.Vec{X: sx, Y: 0}

		hitInfo.Pos = v.Vec{
			X: boxA.Pos.X + boxA.Half.X*sx,
			Y: boxB.Pos.Y,
		}

	} else {
		sy := math.Copysign(1, d.Y)

		hitInfo.Delta = v.Vec{X: 0, Y: py * sy}
		hitInfo.Normal = v.Vec{X: 0, Y: sy}

		hitInfo.Pos = v.Vec{
			X: boxB.Pos.X,
			Y: boxA.Pos.Y + boxA.Half.Y*sy,
		}
	}
	return true
}
