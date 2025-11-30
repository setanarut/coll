package coll

import (
	"math"
)

// BoxOrientedBoxOverlap tests if an AABB and OBB are currently intersecting.
// For moving objects, use BoxOrientedBoxSweep2 to prevent tunneling.
func BoxOrientedBoxOverlap(a *AABB, b *OBB) bool {
	T := b.Pos.Sub(a.Pos)

	// Precompute axes and their absolute values
	bx, by := b.AxisX(), b.AxisY()
	absBxX, absBxY := math.Abs(bx.X), math.Abs(bx.Y)
	absByX, absByY := math.Abs(by.X), math.Abs(by.Y)

	// Check AABB axes
	projBOnAx := absBxX*b.Half.X + absByX*b.Half.Y
	if math.Abs(T.X) > a.Half.X+projBOnAx {
		return false
	}

	projBOnAy := absBxY*b.Half.X + absByY*b.Half.Y
	if math.Abs(T.Y) > a.Half.Y+projBOnAy {
		return false
	}

	// Check OBB axes
	distOnObbX := math.Abs(T.Dot(bx))
	projAOnObbX := absBxX*a.Half.X + absBxY*a.Half.Y
	if distOnObbX > b.Half.X+projAOnObbX {
		return false
	}

	distOnObbY := math.Abs(T.Dot(by))
	projAOnObbY := absByX*a.Half.X + absByY*a.Half.Y
	return distOnObbY <= b.Half.Y+projAOnObbY
}
