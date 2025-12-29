package coll

import (
	"github.com/setanarut/v"
)

// BoxSegmentsSweep1Indexed returns the index of the colliding segment, or -1 if no collision was detected.
//
// Performs a sweep test of an box against a slice of segments.
//
// To determine the earliest point of impact along a movement vector,
// It iterates through the provided segments and finds the collision that occurs at the minimum time value.
// If h is not nil and a collision is detected, it will be populated with:
//   - Normal: Collision surface normal for the box
//   - Data: Normalized time of impact (0.0 to 1.0) along the movement path
func BoxSegmentsSweep1Indexed(s []*Segment, a *AABB, deltaA v.Vec, h *Hit) (index int) {
	colliderIndex := -1
	var resHitTime float64
	var tmpHitInfo Hit

	for i, line := range s {
		if BoxSegmentSweep1(line, a, deltaA, &tmpHitInfo) {
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
