package coll

import (
	"math"
)

// LineCircleOverlap checks if the infinite line defined by the given raySeg intersects with a circle.
// Return true if the infinite line passes through the circle
// This function treats 'raySeg' as a line extending infinitely in both directions
//
// Parameters:
//   - raySeg: Defines the trajectory (slope and position) of the infinite line.
//   - c: The Circle struct to test against.
//   - result: If non-nil and an intersection occurs, this is populated with the
//     two intersection points along the infinite line.
func LineCircleOverlap(raySeg *Segment, c *Circle, result *Segment) bool {

	dp := raySeg.B.Sub(raySeg.A)
	dAPos := raySeg.A.Sub(c.Pos)
	a := dp.MagSq()
	b := 2 * dp.Dot(dAPos)
	cr := dAPos.MagSq() - c.Radius*c.Radius
	bb4ac := b*b - 4*a*cr

	if math.Abs(a) < Epsilon || bb4ac < 0 {
		return false
	}

	if result != nil {
		sqrtBB4AC := math.Sqrt(bb4ac)
		invA2 := 1.0 / (2 * a)
		negB := -b
		result.A = raySeg.A.Add(raySeg.B.Sub(raySeg.A).Scale((negB + sqrtBB4AC) * invA2))
		result.B = raySeg.A.Add(raySeg.B.Sub(raySeg.A).Scale((negB - sqrtBB4AC) * invA2))
	}
	return true
}
