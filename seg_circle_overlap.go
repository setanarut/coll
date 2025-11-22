package coll

import (
	"math"

	"github.com/setanarut/v"
)

// SegmentCircleOverlap returns the intersection points of a Segment and a Circle.
//
// If the returned slice is nil, there is no intersection.
//
// The length can be 1 or 2. It can be queried with the len() method.
func SegmentCircleOverlap(seg *Segment, c *Circle) []v.Vec {
	dp := seg.B.Sub(seg.A)
	dAPos := seg.A.Sub(c.Pos)
	a := dp.MagSq()
	b := 2 * dp.Dot(dAPos)
	bb4ac := b*b - 4*a*(dAPos.MagSq()-c.Radius*c.Radius)

	if math.Abs(a) < Epsilon || bb4ac < 0 {
		return nil
	}

	hitPoints := make([]v.Vec, 0, 2)

	sqrtBB4AC := math.Sqrt(bb4ac)
	invA2 := 1.0 / (2 * a)
	negB := -b

	mu1 := (negB + sqrtBB4AC) * invA2
	mu2 := (negB - sqrtBB4AC) * invA2

	if mu1 >= 0 && mu1 <= 1 {
		hitPoints = append(hitPoints, seg.A.Add(seg.B.Sub(seg.A).Scale(mu1)))
	}
	if mu2 >= 0 && mu2 <= 1 {
		hitPoints = append(hitPoints, seg.A.Add(seg.B.Sub(seg.A).Scale(mu2)))
	}
	return hitPoints
}
