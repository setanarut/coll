package coll

import "github.com/setanarut/v"

// BoxBoxSweep1 Returns true if collision occurs during movement.
// Fills hit info h for dynamicB if not nil.
func BoxBoxSweep1(staticA, dynamicB *AABB, velB v.Vec, h *Hit) bool {
	if velB.IsZero() {
		return BoxBoxOverlap(staticA, dynamicB, h)
	}
	result := BoxSegmentOverlap(staticA, dynamicB.Pos, velB, dynamicB.Half, h)
	if result && h != nil {
		h.Time = max(0, min(1, h.Time-Epsilon))
	}
	return result
}
