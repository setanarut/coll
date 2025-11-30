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
	absT := T.Abs()
	absAx := Ax.Abs()
	absAy := Ay.Abs()
	projBOnAx := (absAx.X * b.Half.X) + (absAy.X * b.Half.Y)
	if absT.X > a.Half.X+projBOnAx {
		return false
	}
	projBOnAy := (absAx.Y * b.Half.X) + (absAy.Y * b.Half.Y)
	if absT.Y > a.Half.Y+projBOnAy {
		return false
	}
	distOnObbX := math.Abs(T.Dot(Ax))
	projAOnObbX := (absAx.X * a.Half.X) + (absAx.Y * a.Half.Y)
	if distOnObbX > b.Half.X+projAOnObbX {
		return false
	}
	distOnObbY := math.Abs(T.Dot(Ay))
	projAOnObbY := (absAy.X * a.Half.X) + (absAy.Y * a.Half.Y)
	return distOnObbY <= b.Half.Y+projAOnObbY
}
