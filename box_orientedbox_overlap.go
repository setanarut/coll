package coll

import (
	"math"
)

// BoxOrientedBoxOverlap tests if an AABB and OBB are currently intersecting.
// For moving objects, use BoxOrientedBoxSweep2 to prevent tunneling.
func BoxOrientedBoxOverlap(a *AABB, b *OBB) bool {
	T := b.Pos.Sub(a.Pos)
	Ax := b.AxisX()
	Ay := b.AxisY()
	projBOnAx := math.Abs(Ax.X)*b.Half.X + math.Abs(Ay.X)*b.Half.Y
	if math.Abs(T.X) > a.Half.X+projBOnAx {
		return false
	}
	projBOnAy := math.Abs(Ax.Y)*b.Half.X + math.Abs(Ay.Y)*b.Half.Y
	if math.Abs(T.Y) > a.Half.Y+projBOnAy {
		return false
	}
	distOnObbX := math.Abs(T.Dot(Ax))
	projAOnObbX := math.Abs(Ax.X)*a.Half.X + math.Abs(Ax.Y)*a.Half.Y
	if distOnObbX > b.Half.X+projAOnObbX {
		return false
	}
	distOnObbY := math.Abs(T.Dot(Ay))
	projAOnObbY := math.Abs(Ay.X)*a.Half.X + math.Abs(Ay.Y)*a.Half.Y
	return distOnObbY <= b.Half.Y+projAOnObbY
}
