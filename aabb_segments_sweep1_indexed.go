package coll

import (
	"github.com/setanarut/v"
)

// AABBSegmentSweep1Indexed performs a sweep test of an AABB against a slice of line segments
// to determine the earliest point of impact along a movement vector.
//
// It iterates through the provided segments and finds the collision that occurs
// at the minimum time value. If a collision is found, the provided hitInfo struct
// is updated with the details of that closest intersection.
//
// Returns true if a collision occurred, along with the index of the colliding segment.
// Returns false and -1 if no collision was detected.
func AABBSegmentSweep1Indexed(lines []*Segment, aabb *AABB, delta v.Vec, hitInfo *HitInfo) (bool, int) {
	colliderIndex := -1
	var resHitTime float64
	var tmpHitInfo HitInfo

	for i, line := range lines {
		if AABBSegmentSweep1(line, aabb, delta, &tmpHitInfo) {
			if colliderIndex == -1 || tmpHitInfo.Time < resHitTime {
				colliderIndex = i
				resHitTime = tmpHitInfo.Time
				*hitInfo = tmpHitInfo
			}
		}
	}

	return colliderIndex >= 0, colliderIndex
}
