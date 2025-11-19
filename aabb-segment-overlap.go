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
	signX := math.Copysign(1, scale.X)
	signY := math.Copysign(1, scale.Y)
	nearTimeX := (box.Pos.X - signX*(box.Half.X+padding.X) - start.X) * scale.X
	nearTimeY := (box.Pos.Y - signY*(box.Half.Y+padding.Y) - start.Y) * scale.Y
	farTimeX := (box.Pos.X + signX*(box.Half.X+padding.X) - start.X) * scale.X
	farTimeY := (box.Pos.Y + signY*(box.Half.Y+padding.Y) - start.Y) * scale.Y
	if math.IsNaN(nearTimeY) {
		nearTimeY = math.Inf(1)
	}
	if math.IsNaN(farTimeY) {
		farTimeY = math.Inf(1)
	}
	if nearTimeX > farTimeY || nearTimeY > farTimeX {
		return false
	}
	nearTime := max(nearTimeX, nearTimeY)
	farTime := min(farTimeX, farTimeY)
	if nearTime >= 1 || farTime <= 0 {
		return false
	}
	if hitInfo == nil {
		return true
	}
	hitInfo.Time = max(0, min(1, nearTime))

	if nearTimeX > nearTimeY {
		hitInfo.Normal.X = -signX
		hitInfo.Normal.Y = 0
	} else {
		hitInfo.Normal.X = 0
		hitInfo.Normal.Y = -signY
	}

	hitInfo.Delta.X = (1.0 - hitInfo.Time) * -delta.X
	hitInfo.Delta.Y = (1.0 - hitInfo.Time) * -delta.Y

	hitInfo.Pos = start.Add(delta.Scale(hitInfo.Time))
	return true
}
