package coll

import (
	"github.com/setanarut/v"
)

// BoxSegmentsSweep1Indexed performs a sweep test of an AABB against a slice of line segments
// to determine the earliest point of impact along a movement vector.
//
// It iterates through the provided segments and finds the collision that occurs
// at the minimum time value. If a collision is found and hitInfo is not nil,
// the provided hitInfo struct is updated with the details of that closest intersection.
//
// Parameters:
//   - lines: Slice of line segments to test against
//   - aabb: The axis-aligned bounding box
//   - delta: Movement vector
//   - hitInfo: Optional pointer to HitInfo struct (can be nil)
//
// Returns the index of the colliding segment, or -1 if no collision was detected.
func BoxSegmentsSweep1Indexed(lines []*Segment, aabb *AABB, delta v.Vec, hitInfo *HitInfo) (index int) {
	colliderIndex := -1
	var resHitTime float64
	var tmpHitInfo HitInfo

	for i, line := range lines {
		if BoxSegmentSweep1(line, aabb, delta, &tmpHitInfo) {
			if colliderIndex == -1 || tmpHitInfo.Time < resHitTime {
				colliderIndex = i
				resHitTime = tmpHitInfo.Time
				// hitInfo nil değilse güncelle
				if hitInfo != nil {
					*hitInfo = tmpHitInfo
				}
			}
		}
	}
	return colliderIndex
}
