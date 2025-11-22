package coll

import (
	"math"
)

// LineCircleOverlap checks if the infinite line defined by the given segment intersects with a circle.
//
// This function treats 'raySeg' as a line extending infinitely in both directions
// (passing through points A and B), rather than a finite segment.
//
// Parameters:
//   - raySeg: Defines the trajectory (slope and position) of the infinite line.
//   - circ: The Circle struct to test against.
//   - overlapSeg: If non-nil and an intersection occurs, this is populated with the
//     two intersection points along the infinite line.
//
// Returns:
//   - overlap: True if the infinite line passes through the circle.
func LineCircleOverlap(raySeg *Segment, circ *Circle, overlapSeg *Segment) bool {

	dp := raySeg.B.Sub(raySeg.A)
	dAPos := raySeg.A.Sub(circ.Pos)
	a := dp.MagSq()
	b := 2 * dp.Dot(dAPos)
	c := dAPos.MagSq() - circ.Radius*circ.Radius
	bb4ac := b*b - 4*a*c

	if math.Abs(a) < Epsilon || bb4ac < 0 {
		return false
	}

	if overlapSeg != nil {
		sqrtBB4AC := math.Sqrt(bb4ac)
		invA2 := 1.0 / (2 * a)
		negB := -b
		overlapSeg.A = raySeg.A.Add(raySeg.B.Sub(raySeg.A).Scale((negB + sqrtBB4AC) * invA2))
		overlapSeg.B = raySeg.A.Add(raySeg.B.Sub(raySeg.A).Scale((negB - sqrtBB4AC) * invA2))
	}
	return true
}
