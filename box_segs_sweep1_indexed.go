package coll

import (
	"github.com/setanarut/v"
)

// BoxSegmentsSweep1Indexed returns the index of the colliding segment, or -1 if no collision was detected.
//
// Performs a sweep test of an AABB against a slice of line segments.
//
// To determine the earliest point of impact along a movement vector,
// It iterates through the provided segments and finds the collision that occurs at the minimum time value.
//   - Normal: Collision surface normal for box
//   - Data: Normalized time of impact (0.0 to 1.0) along the movement path
//
// Parameters:
//   - lines: Slice of line segments to test against
//   - box: The axis-aligned bounding box
//   - velBox: Movement vector for box
//   - h: Optional pointer to Hit struct (can be nil)
func BoxSegmentsSweep1Indexed(lines []*Segment, box *AABB, velBox v.Vec, h *Hit) (index int) {
	colliderIndex := -1
	var resHitTime float64
	var tmpHitInfo Hit

	for i, line := range lines {
		if BoxSegmentSweep1(line, box, velBox, &tmpHitInfo) {
			if colliderIndex == -1 || tmpHitInfo.Data < resHitTime {
				colliderIndex = i
				resHitTime = tmpHitInfo.Data
				// hitInfo nil değilse güncelle
				if h != nil {
					*h = tmpHitInfo
				}
			}
		}
	}
	return colliderIndex
}
