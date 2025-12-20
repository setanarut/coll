package coll

import (
	"math"

	"github.com/setanarut/v"
)

// CircleCircleSweep2 checks for collision between two moving circles to prevent tunneling.
//
// Returns true if collision occurs during movement or if they are already overlapping.
//
// If h is not nil and a collision is detected, it will be populated with:
//   - Normal: Collision surface normal for c2
//   - Data: Normalized time of impact (0.0 to 1.0) along the movement path
func CircleCircleSweep2(c1, c2 *Circle, c1Vel, c2Vel v.Vec, h *Hit) bool {
	rSum := c1.Radius + c2.Radius
	rSumSq := rSum * rSum
	pDiff := c2.Pos.Sub(c1.Pos)
	if pDiff.MagSq() <= rSumSq {
		if h != nil {
			h.Data = 0.0
			h.Normal = pDiff.Unit().Neg()
		}
		return true
	}
	relVel := c2Vel.Sub(c1Vel)
	a := relVel.MagSq()
	if a < Epsilon {
		return false
	}
	b := 2.0 * pDiff.Dot(relVel)
	c := pDiff.MagSq() - rSumSq
	discriminant := b*b - 4.0*a*c
	if discriminant < 0.0 {
		return false
	}
	sqrtDisc := math.Sqrt(discriminant)
	t1 := (-b - sqrtDisc) / (2.0 * a)
	if t1 >= 0.0 && t1 <= 1.0 {
		if h != nil {
			h.Data = t1
			h.Normal = pDiff.Add(relVel.Scale(t1)).Unit().Neg()
		}
		return true
	}

	return false
}
