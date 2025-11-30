package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxOrientedBoxSweep2 tests if a moving AABB and moving OBB intersect during their motion.
// Uses swept volume testing to prevent tunneling for fast-moving objects.
func BoxOrientedBoxSweep2(a *AABB, b *OBB, va v.Vec, vb v.Vec) bool {
	relVx := vb.X - va.X
	relVy := vb.Y - va.Y
	tx := b.Pos.X - a.Pos.X
	ty := b.Pos.Y - a.Pos.Y

	bAx := b.AxisX()
	bAy := b.AxisY()
	absAxX := math.Abs(bAx.X)
	absAxY := math.Abs(bAx.Y)
	absAyX := math.Abs(bAy.X)
	absAyY := math.Abs(bAy.Y)

	projB_on_GlobalX := (absAxX * b.Half.X) + (absAyX * b.Half.Y)
	limitX := a.Half.X + projB_on_GlobalX
	if (tx > 0 && relVx < 0) || (tx < 0 && relVx > 0) {
		limitX += math.Abs(relVx)
	}
	if math.Abs(tx) > limitX {
		return false
	}

	projB_on_GlobalY := (absAxY * b.Half.X) + (absAyY * b.Half.Y)
	limitY := a.Half.Y + projB_on_GlobalY
	if (ty > 0 && relVy < 0) || (ty < 0 && relVy > 0) {
		limitY += math.Abs(relVy)
	}
	if math.Abs(ty) > limitY {
		return false
	}

	dotT_Ax := (tx * bAx.X) + (ty * bAx.Y)
	projA_on_ObbX := (absAxX * a.Half.X) + (absAxY * a.Half.Y)
	dotV_Ax := (relVx * bAx.X) + (relVy * bAx.Y)
	limitObbX := b.Half.X + projA_on_ObbX
	if (dotT_Ax > 0 && dotV_Ax < 0) || (dotT_Ax < 0 && dotV_Ax > 0) {
		limitObbX += math.Abs(dotV_Ax)
	}
	if math.Abs(dotT_Ax) > limitObbX {
		return false
	}

	dotT_Ay := (tx * bAy.X) + (ty * bAy.Y)
	projA_on_ObbY := (absAyX * a.Half.X) + (absAyY * a.Half.Y)
	dotV_Ay := (relVx * bAy.X) + (relVy * bAy.Y)
	limitObbY := b.Half.Y + projA_on_ObbY
	if (dotT_Ay > 0 && dotV_Ay < 0) || (dotT_Ay < 0 && dotV_Ay > 0) {
		limitObbY += math.Abs(dotV_Ay)
	}
	if math.Abs(dotT_Ay) > limitObbY {
		return false
	}

	return true
}
