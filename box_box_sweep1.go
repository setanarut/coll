package coll

import "github.com/setanarut/v"

// BoxBoxSweep1 Returns true if collision occurs during movement.
// Fills hit info h for dynamicB if not nil.
func BoxBoxSweep1(staticA, dynamicB *AABB, deltaB v.Vec, h *Hit) bool {
	if deltaB.IsZero() {
		return BoxBoxOverlap(staticA, dynamicB, h)
	}
	result := BoxSegmentOverlap(staticA, dynamicB.Pos, deltaB, dynamicB.Half, h)
	if result && h != nil {
		h.Data = max(0, min(1, h.Data-Epsilon))
	}
	return result
}
