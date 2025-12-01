package coll

import (
	"math"

	"github.com/setanarut/v"
)

// CircleCircleSweep2 checks for collision between two moving circles to prevent tunneling.
// It treats the movement as a continuous sweep from t=0 to t=1.
//
// Returns true if collision occurs during movement or if they are already overlapping.
func CircleCircleSweep2(c1, c2 *Circle, c1Vel, c2Vel v.Vec) bool {
	rSum := c1.Radius + c2.Radius
	rSumSq := rSum * rSum
	pDiff := c2.Pos.Sub(c1.Pos)
	if pDiff.MagSq() <= rSumSq {
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
		return true
	}
	return false
}
