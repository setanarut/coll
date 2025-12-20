package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxBoxOverlap checks whether a and b overlap.
//
// If h is not nil, the function fills it with for box b:
//   - Normal: Collision surface normal for box b
//   - Data: the penetration depth for box b (overlap distance)
//
// This method can behave poorly for moving objects.
//
// For continuous motion, BoxBoxSweep1() of BoxBoxSweep2() should be used instead.
func BoxBoxOverlap(a, b *AABB, h *Hit) bool {
	d := b.Pos.Sub(a.Pos)
	absD := d.Abs()
	hSum := a.Half.Add(b.Half)

	px := hSum.X - absD.X

	if px <= 0 {
		return false
	}

	py := hSum.Y - absD.Y

	if py <= 0 {
		return false
	}

	if h == nil {
		return true
	}

	if px < py {
		sx := math.Copysign(1, d.X)
		h.Normal = v.Vec{X: sx, Y: 0}
		h.Data = px
	} else {
		sy := math.Copysign(1, d.Y)
		h.Normal = v.Vec{X: 0, Y: sy}
		h.Data = py
	}

	return true
}
