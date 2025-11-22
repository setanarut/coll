package coll

import "github.com/setanarut/v"

// boxCircleIntersect checks for intersection between AABB and a Circle.
// This method can behave poorly for moving objects. For continuous motion,
// AABBCircleSweep2 should be used instead.
func boxCircleIntersect(box *AABB, circle *Circle) bool {
	diff := circle.Pos.Sub(box.Pos)
	clamped := v.Vec{
		X: max(-box.Half.X, min(diff.X, box.Half.X)),
		Y: max(-box.Half.Y, min(diff.Y, box.Half.Y)),
	}
	closest := box.Pos.Add(clamped)
	return circle.Pos.DistSq(closest) <= circle.Radius*circle.Radius
}
