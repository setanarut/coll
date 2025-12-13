package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxOrientedBoxOverlap tests if an AABB and OBB are currently intersecting.
// For moving objects, use BoxOrientedBoxSweep2 to prevent tunneling.
func BoxOrientedBoxOverlap(a *AABB, b *OBB) bool {
	d := b.Pos.Sub(a.Pos)

	// Precompute axes and their absolute values
	bAxisX := v.FromAngle(b.Angle)
	bAxisY := v.Vec{X: -bAxisX.Y, Y: bAxisX.X}
	bAxisXAbs := bAxisX.Abs()
	bAxisYAbs := bAxisY.Abs()

	// Check AABB axes
	projBOnAx := bAxisXAbs.X*b.Half.X + bAxisYAbs.X*b.Half.Y
	if math.Abs(d.X) > a.Half.X+projBOnAx {
		return false
	}

	projBOnAy := bAxisXAbs.Y*b.Half.X + bAxisYAbs.Y*b.Half.Y
	if math.Abs(d.Y) > a.Half.Y+projBOnAy {
		return false
	}

	// Check OBB axes
	distOnObbX := math.Abs(d.Dot(bAxisX))
	projAOnObbX := bAxisXAbs.X*a.Half.X + bAxisXAbs.Y*a.Half.Y
	if distOnObbX > b.Half.X+projAOnObbX {
		return false
	}

	distOnObbY := math.Abs(d.Dot(bAxisY))
	projAOnObbY := bAxisYAbs.X*a.Half.X + bAxisYAbs.Y*a.Half.Y
	return distOnObbY <= b.Half.Y+projAOnObbY
}
