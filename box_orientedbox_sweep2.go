package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxOrientedBoxSweep2 tests if a moving AABB and moving OBB intersect during their motion.
// Uses swept volume testing to prevent tunneling for fast-moving objects.
func BoxOrientedBoxSweep2(a *AABB, b *OBB, va v.Vec, vb v.Vec) bool {
	V_rel := vb.Sub(va)
	T := b.Pos.Sub(a.Pos)
	Ax := b.AxisX()
	Ay := b.AxisY()
	absT := T.Abs()
	absAx := Ax.Abs()
	absAy := Ay.Abs()
	absVRelOnAx := math.Abs(V_rel.Dot(Ax))
	absVRelOnAy := math.Abs(V_rel.Dot(Ay))
	projBOnAx := (absAx.X * b.Half.X) + (absAy.X * b.Half.Y)
	projVRelOnGlobalX := V_rel.AbsX()
	distX := T.X
	if (distX > 0 && V_rel.X > 0) || (distX < 0 && V_rel.X < 0) {
		if absT.X > a.Half.X+projBOnAx {
			return false
		}
	} else {
		if absT.X > a.Half.X+projBOnAx+projVRelOnGlobalX {
			return false
		}
	}
	projBOnAy := (absAx.Y * b.Half.X) + (absAy.Y * b.Half.Y)
	projVRelOnGlobalY := V_rel.AbsY()
	distY := T.Y
	if (distY > 0 && V_rel.Y > 0) || (distY < 0 && V_rel.Y < 0) {
		if absT.Y > a.Half.Y+projBOnAy {
			return false
		}
	} else {
		if absT.Y > a.Half.Y+projBOnAy+projVRelOnGlobalY {
			return false
		}
	}
	distOnObbX := math.Abs(T.Dot(Ax))
	projAOnObbX := (absAx.X * a.Half.X) + (absAx.Y * a.Half.Y)
	dotTAx := T.Dot(Ax)
	dotVAx := V_rel.Dot(Ax)
	if (dotTAx > 0 && dotVAx > 0) || (dotTAx < 0 && dotVAx < 0) {
		if distOnObbX > b.Half.X+projAOnObbX {
			return false
		}
	} else {
		if distOnObbX > b.Half.X+projAOnObbX+absVRelOnAx {
			return false
		}
	}
	distOnObbY := math.Abs(T.Dot(Ay))
	projAOnObbY := (absAy.X * a.Half.X) + (absAy.Y * a.Half.Y)
	dotTAy := T.Dot(Ay)
	dotVAy := V_rel.Dot(Ay)
	if (dotTAy > 0 && dotVAy > 0) || (dotTAy < 0 && dotVAy < 0) {
		return distOnObbY <= b.Half.Y+projAOnObbY
	} else {
		return distOnObbY <= b.Half.Y+projAOnObbY+absVRelOnAy
	}
}
