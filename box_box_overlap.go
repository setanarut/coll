package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxBoxOverlap checks whether a and b overlap.
// Hit describes how to move a out of b.
//
// If h is not nil, the function fills it with:
//   - Time: the penetration depth (overlap distance), where pos.Add(normal.Scale(time))
//     gives the pushback distance needed to resolve the overlap
//   - Normal: the direction in which boxA is pushed
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
		h.Time = px
	} else {
		sy := math.Copysign(1, d.Y)
		h.Normal = v.Vec{X: 0, Y: sy}
		h.Time = py
	}

	return true
}
