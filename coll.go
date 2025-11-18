// 2D Collision functions
package coll

import (
	"math"

	"github.com/setanarut/v"
)

const Epsilon float64 = 1e-8

// HitInfo holds the detailed information about a collision or contact event.
type HitInfo struct {
	Pos    v.Vec   // The position where the collision occurred (contact point).
	Delta  v.Vec   // The remaining movement vector after the collision.
	Normal v.Vec   // The normal vector of the surface hit.
	Time   float64 // The time (0.0 to 1.0) along the movement path when the collision happened.
}

// Resets the HitInfo struct to its zero values.
func (h *HitInfo) Reset() {
	*h = HitInfo{} // Reinitializes all fields of the struct to their zero values (nil, 0, false, etc.).
}

// AABBAABBContain returns true if a fully contains b (b is fully inside of the bounds of a).
func AABBAABBContain(a, b *AABB) bool {
	if b.Left() < a.Left() {
		return false
	}
	if b.Right() > a.Right() {
		return false
	}
	if b.Top() < a.Top() {
		return false
	}
	if b.Bottom() > a.Bottom() {
		return false
	}
	return true
}

// AABBSegmentOverlap returns true if they intersect, false otherwise
//
// Params:
//
//   - box - Bounding box to check
//   - start - Line segment origin/start position
//   - delta - Line segment move/displacement vector
//   - padding - Padding added to the radius of the bounding box
//   - hitInfo - Contact info. Filled when argument isn't nil and a collision occurs
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

// AABBOverlap checks whether boxA and boxB overlap.
// Any collision information written to hitInfo always describes how to move boxA out of boxB.
//
// It uses a separating-axis test: if the boxes do not overlap on either X or Y,
// there is no collision and the function returns false.
//
// If hitInfo is not nil, the function fills it with:
//   - Delta: the minimum vector needed to push boxA out of boxB
//   - Normal: the direction in which boxA is pushed
//   - Pos: an approximate contact position on the collision side
//
// This method can behave poorly for moving objects. For continuous motion,
// sweepAABB should be used instead.
//
// If you only need to know whether a collision occurred, pass nil for hitInfo
// to skip generating collision details.
func AABBOverlap(boxA, boxB *AABB, hitInfo *HitInfo) bool {

	dx := boxB.Pos.X - boxA.Pos.X
	px := boxB.Half.X + boxA.Half.X - math.Abs(dx)

	if px <= 0 {
		return false
	}

	dy := boxB.Pos.Y - boxA.Pos.Y
	py := boxB.Half.Y + boxA.Half.Y - math.Abs(dy)

	if py <= 0 {
		return false
	}

	if hitInfo == nil {
		return true
	}

	// if if hitInfo is not nil, fill
	if px < py {
		sx := math.Copysign(1, dx)
		hitInfo.Delta.X = px * sx
		hitInfo.Normal.X = sx
		hitInfo.Pos.X = boxA.Pos.X + boxA.Half.X*sx
		hitInfo.Pos.Y = boxB.Pos.Y
	} else {
		sy := math.Copysign(1, dy)
		hitInfo.Delta.Y = py * sy
		hitInfo.Normal.Y = sy
		hitInfo.Pos.X = boxB.Pos.X
		hitInfo.Pos.Y = boxA.Pos.Y + boxA.Half.Y*sy
	}
	return true
}

// AABBAABBSweep1 fills hit info for boxB if not nil. Returns true if collision occurs during movement.
//
// https://noonat.github.io/intersect/#aabb-vs-swept-aabb
//
// returns bool true if the two AABBs collide, false otherwise. If hitInfo is not nil, the function fills it.
//
// Params:
//   - staticBoxA - The static box
//   - boxB - The moving box
//   - boxBVel - The displacement vector of boxB
//   - hitInfo - The contact object. Filled if collision occurs
func AABBAABBSweep1(staticBoxA, boxB *AABB, boxBVel v.Vec, hitInfo *HitInfo) bool {
	if boxBVel.IsZero() {
		return AABBOverlap(staticBoxA, boxB, hitInfo)
	}
	result := AABBSegmentOverlap(staticBoxA, boxB.Pos, boxBVel, boxB.Half, hitInfo)
	if result {
		hitInfo.Time = max(0, min(1, hitInfo.Time-Epsilon))
		direction := boxBVel.Unit()
		hitInfo.Pos.X = max(staticBoxA.Left(), min(staticBoxA.Right(), hitInfo.Pos.X+direction.X*boxB.Half.X))
		hitInfo.Pos.Y = max(hitInfo.Pos.Y+direction.Y*boxB.Half.Y, min(staticBoxA.Top(), staticBoxA.Bottom()))
	}
	return result
}

// AABBAABBSweep2 fills hit info for boxB if not nil. Returns true if collision occurs during movement, false otherwise.
//
// Sweep two moving AABBs to see if and when they first and last were overlapping.
// https://www.gamedeveloper.com/disciplines/simple-intersection-tests-for-games
//
// Params:
//   - boxA - previous state of boxA
//   - boxB - previous state of boxB
//   - boxAVel - displacment vector of boxA
//   - boxBVel - displacement vector of boxB
//   - hitInfo - hit info for boxB. Filled if collision occurs, can be set to nil for performance
func AABBAABBSweep2(boxA, boxB *AABB, boxAVel, boxBVel v.Vec, hitInfo *HitInfo) bool {
	delta := boxBVel.Sub(boxAVel)
	isCollide := AABBAABBSweep1(boxA, boxB, delta, hitInfo)
	if isCollide {
		hitInfo.Pos = hitInfo.Pos.Add(boxAVel.Scale(hitInfo.Time))
		if hitInfo.Normal.X != 0 {
			hitInfo.Pos.X = boxB.Pos.X + (boxBVel.X * hitInfo.Time) - (hitInfo.Normal.X * boxB.Half.X)
		} else {
			hitInfo.Pos.Y = boxB.Pos.Y + (boxBVel.Y * hitInfo.Time) - (hitInfo.Normal.Y * boxB.Half.Y)
		}
	}
	return isCollide
}

// AABBCircleSweep checks for collision between a moving AABB and a moving Circle.
// Returns true if collision occurs during movement, false otherwise.
func AABBCircleSweep(box *AABB, circle *Circle, boxVel, CircleVel v.Vec) bool {
	// Calculate circle's movement relative to AABB (AABB becomes stationary reference frame)
	relativeDelta := CircleVel.Sub(boxVel)

	// Use Segment to check if circle (treated as point with radius padding) hits the AABB
	// padding expands AABB by circle's radius to simplify collision detection
	return AABBSegmentOverlap(box, circle.Pos, relativeDelta, v.Vec{X: circle.Radius, Y: circle.Radius}, nil)
}

// CalculateSlideVelocity computes the total movement: movement until collision plus sliding along the surface normal.
func CalculateSlideVelocity(vel v.Vec, hitInfo *HitInfo) (slideVel v.Vec) {
	remainingVel := vel.Scale(1.0 - hitInfo.Time)
	slideVel = remainingVel.Sub(hitInfo.Normal.Scale(remainingVel.Dot(hitInfo.Normal)))
	movementToHit := vel.Scale(hitInfo.Time)
	return movementToHit.Add(slideVel)
}

// ApplySlideVelocity updates the AABB position by applying the calculated slide velocity.
func ApplySlideVelocity(box *AABB, vel v.Vec, hitInfo *HitInfo) {
	box.Pos = box.Pos.Add(CalculateSlideVelocity(vel, hitInfo))
}
