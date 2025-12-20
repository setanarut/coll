package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxSegmentOverlap returns true if they intersect, false otherwise.
//
// Params:
//
//   - box - Bounding box to check
//   - start - Ray segment origin/start position
//   - delta - Ray segment move/displacement vector
//   - padding - Padding added to the radius of the bounding box
//   - h - Contact info for segment. Filled when argument isn't nil and a collision occurs
func BoxSegmentOverlap(box *AABB, start, delta, padding v.Vec, h *Hit) bool {

	scale := v.One.Div(delta)
	signVec := v.Vec{X: math.Copysign(1, scale.X), Y: math.Copysign(1, scale.Y)}
	signedExtent := box.Half.Add(padding).Mul(signVec)
	nearTimes := box.Pos.Sub(signedExtent).Sub(start).Mul(scale)
	farTimes := box.Pos.Add(signedExtent).Sub(start).Mul(scale)

	if math.IsNaN(nearTimes.Y) {
		nearTimes.Y = math.Inf(1)
	}
	if math.IsNaN(farTimes.Y) {
		farTimes.Y = math.Inf(1)
	}
	if nearTimes.X > farTimes.Y || nearTimes.Y > farTimes.X {
		return false
	}

	nearTime := max(nearTimes.X, nearTimes.Y)
	farTime := min(farTimes.X, farTimes.Y)

	if nearTime >= 1 || farTime <= 0 {
		return false
	}

	if h == nil {
		return true
	}

	h.Data = max(0, min(1, nearTime))

	if nearTimes.X > nearTimes.Y {
		h.Normal = v.Vec{X: -signVec.X, Y: 0}
	} else {
		h.Normal = v.Vec{X: 0, Y: -signVec.Y}
	}

	return true
}
