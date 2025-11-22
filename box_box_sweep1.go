package coll

import "github.com/setanarut/v"

// BoxBoxSweep1 fills hit info for boxB if not nil. Returns true if collision occurs during movement.
//
// https://noonat.github.io/intersect/#aabb-vs-swept-aabb
//
// returns bool true if the two AABBs collide, false otherwise. If hitInfo is not nil, the function fills it.
//
// Params:
//   - staticBoxA - The static box
//   - boxB - The moving box
//   - boxBVel - The displacement vector of boxB
//   - hitInfo - The contact object. Filled if collision occurs
func BoxBoxSweep1(staticBoxA, boxB *AABB, boxBVel v.Vec, hitInfo *HitInfo) bool {
	if boxBVel.IsZero() {
		return BoxBoxOverlap(staticBoxA, boxB, hitInfo)
	}
	result := BoxSegmentOverlap(staticBoxA, boxB.Pos, boxBVel, boxB.Half, hitInfo)
	if result {
		hitInfo.Time = max(0, min(1, hitInfo.Time-Epsilon))
		direction := boxBVel.Unit()
		hitInfo.Pos.X = max(staticBoxA.Left(), min(staticBoxA.Right(), hitInfo.Pos.X+direction.X*boxB.Half.X))
		hitInfo.Pos.Y = max(hitInfo.Pos.Y+direction.Y*boxB.Half.Y, min(staticBoxA.Top(), staticBoxA.Bottom()))
	}
	return result
}
