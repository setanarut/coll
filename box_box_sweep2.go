package coll

import "github.com/setanarut/v"

// BoxBoxSweep2 fills hit info for boxB if not nil. Returns true if collision occurs during movement, false otherwise.
//
// Sweep two moving AABBs to see if and when they first and last were overlapping.
// https://www.gamedeveloper.com/disciplines/simple-intersection-tests-for-games
//
// Params:
//   - boxA - previous state of boxA
//   - boxB - previous state of boxB
//   - boxAVel - displacment vector of boxA
//   - boxBVel - displacement vector of boxB
//   - hitInfo - hit info for boxB. Filled if collision occurs, can be set to nil for performance
func BoxBoxSweep2(boxA, boxB *AABB, boxAVel, boxBVel v.Vec, hitInfo *HitInfo) bool {
	delta := boxBVel.Sub(boxAVel)
	isCollide := BoxBoxSweep1(boxA, boxB, delta, hitInfo)
	if isCollide {
		hitInfo.Pos = hitInfo.Pos.Add(boxAVel.Scale(hitInfo.Time))
		if hitInfo.Normal.X != 0 {
			hitInfo.Pos.X = boxB.Pos.X + (boxBVel.X * hitInfo.Time) - (hitInfo.Normal.X * boxB.Half.X)
		} else {
			hitInfo.Pos.Y = boxB.Pos.Y + (boxBVel.Y * hitInfo.Time) - (hitInfo.Normal.Y * boxB.Half.Y)
		}
	}
	return isCollide
}
