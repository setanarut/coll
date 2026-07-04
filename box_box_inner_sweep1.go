package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxBoxInnerSweep1 checks if moving box 'a' hits the internal bounds of static container box 'b'.
//
// Parameters:
//   - a: The dynamic (moving) box.
//   - b: The static (stationary) container box.
//   - delta: The movement vector of box 'a'.
//   - h: Hit information.
//
// If a collision occurs and 'h' is not nil, it fills 'h' with:
//   - Normal: The surface normal of the inner edge of 'b' that 'a' collided with.
//   - Data: The time of impact (0.0 to 1.0) along the movement path.
func BoxBoxInnerSweep1(a, b *AABB, delta v.Vec, h *Hit) bool {
	if delta.IsZero() {
		return false
	}

	tX, tY := math.Inf(1), math.Inf(1)
	var nX, nY v.Vec

	if delta.X > 0 {
		tX = (b.Pos.X + b.Half.X - a.Half.X - a.Pos.X) / delta.X
		nX = v.Vec{X: -1, Y: 0}
	} else if delta.X < 0 {
		tX = (b.Pos.X - b.Half.X + a.Half.X - a.Pos.X) / delta.X
		nX = v.Vec{X: 1, Y: 0}
	}

	if delta.Y > 0 {
		tY = (b.Pos.Y + b.Half.Y - a.Half.Y - a.Pos.Y) / delta.Y
		nY = v.Vec{X: 0, Y: -1}
	} else if delta.Y < 0 {
		tY = (b.Pos.Y - b.Half.Y + a.Half.Y - a.Pos.Y) / delta.Y
		nY = v.Vec{X: 0, Y: 1}
	}

	minT := math.Inf(1)
	var normal v.Vec

	if tX >= 0 && tX <= 1 {
		minT = tX
		normal = nX
	}
	if tY >= 0 && tY <= 1 && tY < minT {
		minT = tY
		normal = nY
	}

	if minT <= 1 {
		if h != nil {
			h.Data = max(0, minT-Epsilon)
			h.Normal = normal
		}
		return true
	}
	return false
}
