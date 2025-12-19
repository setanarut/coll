package coll

import "github.com/setanarut/v"

// BoxBoxSweep2 Returns true if collision occurs during movement. Fills hit info h for b if not nil.
func BoxBoxSweep2(a, b *AABB, velA, velB v.Vec, h *Hit) bool {
	delta := velB.Sub(velA)
	return BoxBoxSweep1(a, b, delta, h)
}
