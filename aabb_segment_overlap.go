package coll

import (
	"math"

	"github.com/setanarut/v"
)

// AABBSegmentOverlap returns true if they intersect, false otherwise.
//
// Params:
//
//   - box - Bounding box to check
//   - start - Line segment origin/start position
//   - delta - Line segment move/displacement vector
//   - padding - Padding added to the radius of the bounding box
//   - hitInfo - Contact info for segment. Filled when argument isn't nil and a collision occurs
func AABBSegmentOverlap(box *AABB, start, delta, padding v.Vec, hitInfo *HitInfo) bool {

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

	if hitInfo == nil {
		return true
	}

	hitInfo.Time = max(0, min(1, nearTime))

	if nearTimes.X > nearTimes.Y {
		hitInfo.Normal = v.Vec{X: -signVec.X, Y: 0}
	} else {
		hitInfo.Normal = v.Vec{X: 0, Y: -signVec.Y}
	}

	hitInfo.Delta = delta.Neg().Scale(1.0 - hitInfo.Time)
	hitInfo.Pos = start.Add(delta.Scale(hitInfo.Time))

	return true
}
