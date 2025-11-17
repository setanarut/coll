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

// HitInfo2 tracks which sides of an AABB (Box) were contacted.
type HitInfo2 struct {
	// Flags indicating which side of the bounding box was hit.
	Right, Bottom, Left, Top bool
	// The remaining movement vector after the collision.
	Delta v.Vec
}

// Resets the HitInfo struct to its zero values.
func (h *HitInfo) Reset() {
	*h = HitInfo{} // Reinitializes all fields of the struct to their zero values (nil, 0, false, etc.).
}

// Resets the HitInfo struct to its zero values.
func (h *HitInfo2) Reset() {
	*h = HitInfo2{} // Reinitializes all fields of the struct to their zero values (nil, 0, false, etc.).
}

// AABBPlatform returns hit info for a. Does not prevent collision tunneling, a slides along the platform edges.
func AABBPlatform(a, platform *AABB, aVel, platformVel v.Vec, h *HitInfo2) bool {
	// Calculate old positions using velocities
	aOldPos := v.Vec{a.Pos.X - aVel.X, a.Pos.Y - aVel.Y}
	bOldPos := v.Vec{platform.Pos.X - platformVel.X, platform.Pos.Y - platformVel.Y}

	// Check collision at current positions using half dimensions
	xDist := math.Abs(a.Pos.X - platform.Pos.X)
	yDist := math.Abs(a.Pos.Y - platform.Pos.Y)

	// Combined half widths and heights
	combinedHalfW := a.Half.X + platform.Half.X
	combinedHalfH := a.Half.Y + platform.Half.Y

	// Early exit check
	if xDist > combinedHalfW || yDist > combinedHalfH {
		return false
	}

	// Calculate old distances using calculated old positions
	oldXDist := math.Abs(aOldPos.X - bOldPos.X)
	oldYDist := math.Abs(aOldPos.Y - bOldPos.Y)

	// Check collision direction and calculate position delta
	if yDist < combinedHalfH {
		if a.Pos.Y > platform.Pos.Y && oldYDist >= combinedHalfH {
			h.Delta.Y = (platform.Pos.Y + combinedHalfH + Epsilon) - a.Pos.Y
			h.Top = true
		} else if a.Pos.Y < platform.Pos.Y && oldYDist >= combinedHalfH {
			h.Delta.Y = (platform.Pos.Y - combinedHalfH - Epsilon) - a.Pos.Y
			h.Bottom = true
		}
	}

	if xDist < combinedHalfW {
		if a.Pos.X > platform.Pos.X && oldXDist >= combinedHalfW {
			h.Delta.X = (platform.Pos.X + combinedHalfW + Epsilon) - a.Pos.X
			h.Left = true
		} else if a.Pos.X < platform.Pos.X && oldXDist >= combinedHalfW {
			h.Delta.X = (platform.Pos.X - combinedHalfW - Epsilon) - a.Pos.X
			h.Right = true
		}
	}

	return true
}

// AABBSegmentOverlap returns true if they intersect, false otherwise
//
// Params:
//
//   - box - Bounding box to check
//   - pos - Line segment origin/start position
//   - delta - Line segment move/displacement vector
//   - padding - Padding added to the radius of the bounding box
//   - hit - Physics contact info. Filled when argument isn't nil and a collision occurs
func AABBSegmentOverlap(box *AABB, pos, delta, padding v.Vec, hit *HitInfo) bool {
	scale := v.One.Div(delta)
	signX := math.Copysign(1, scale.X)
	signY := math.Copysign(1, scale.Y)
	nearTimeX := (box.Pos.X - signX*(box.Half.X+padding.X) - pos.X) * scale.X
	nearTimeY := (box.Pos.Y - signY*(box.Half.Y+padding.Y) - pos.Y) * scale.Y
	farTimeX := (box.Pos.X + signX*(box.Half.X+padding.X) - pos.X) * scale.X
	farTimeY := (box.Pos.Y + signY*(box.Half.Y+padding.Y) - pos.Y) * scale.Y
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
	if hit == nil {
		return true
	}
	hit.Time = max(0, min(1, nearTime))

	if nearTimeX > nearTimeY {
		hit.Normal.X = -signX
		hit.Normal.Y = 0
	} else {
		hit.Normal.X = 0
		hit.Normal.Y = -signY
	}
	hit.Delta.X = (1.0 - hit.Time) * -delta.X
	hit.Delta.Y = (1.0 - hit.Time) * -delta.Y

	hit.Pos = pos.Add(delta.Scale(hit.Time))
	return true
}

// AABBOverlap returns true if they intersect, false otherwise
// https://noonat.github.io/intersect/#aabb-vs-aabb
//
// This test uses a separating axis test, which checks for overlaps between the
// two boxes on each axis. If either axis is not overlapping, the boxes aren’t colliding.
//
// The function fills a *HitInfo object, and gives the axis of least overlap as the contact point.
//
// That is, it sets hit.Delta so that the colliding box will be pushed out of the nearest edge
// This can cause weird behavior for moving boxes, so you should use sweepAABB
// instead for moving boxes.
//
// If collision information is not needed, the hit argument can be set to nil for performance.
func AABBOverlap(a, b *AABB, hit *HitInfo) bool {
	dx := b.Pos.X - a.Pos.X
	px := b.Half.X + a.Half.X - math.Abs(dx)
	if px <= 0 {
		return false
	}

	dy := b.Pos.Y - a.Pos.Y
	py := b.Half.Y + a.Half.Y - math.Abs(dy)
	if py <= 0 {
		return false
	}

	if hit == nil {
		return true
	}

	// hit reset here
	if px < py {
		sx := math.Copysign(1, dx)
		hit.Delta.X = px * sx
		hit.Normal.X = sx
		hit.Pos.X = a.Pos.X + a.Half.X*sx
		hit.Pos.Y = b.Pos.Y
	} else {
		sy := math.Copysign(1, dy)
		hit.Delta.Y = py * sy
		hit.Normal.Y = sy
		hit.Pos.X = b.Pos.X
		hit.Pos.Y = a.Pos.Y + a.Half.Y*sy
	}
	return true
}

// AABBAABBSweep1 finds the intersection of this box (a) and another moving box (b),
// where the delta argument is the displacement vector of the moving box (b).
// https://noonat.github.io/intersect/#aabb-vs-swept-aabb
//
// returns bool true if the two AABBs collide, false otherwise
//
// Params:
//   - a - The static box
//   - b - The moving box
//   - delta - The displacement vector of b
//   - hit - The contact object. Filled if collision occurs
//
// If collision information is not needed, the hit argument can be set to nil for performance.
func AABBAABBSweep1(a, b *AABB, delta v.Vec, hit *HitInfo) bool {
	if delta.IsZero() {
		return AABBOverlap(a, b, hit)
	}
	result := AABBSegmentOverlap(a, b.Pos, delta, b.Half, hit)
	if result {
		hit.Time = max(0, min(1, hit.Time))
		direction := delta.Unit()
		hit.Pos.X = max(a.Left(), min(a.Right(), hit.Pos.X+direction.X*b.Half.X))
		hit.Pos.Y = max(hit.Pos.Y+direction.Y*b.Half.Y, min(a.Top(), a.Bottom()))
	}
	return result
}

// AABBAABBSweep2 Sweep two moving AABBs to see if and when they first and last were overlapping.
// https://www.gamedeveloper.com/disciplines/simple-intersection-tests-for-games
//
// Params:
//   - a - previous state of AABB a
//   - b - previous state of AABB b
//   - aV - displacment vector of a
//   - bV - displacement vector of b
//   - hit - The contact object. Filled if collision occurs
//
// If collision information is not needed, the hit argument can be set to nil for performance.
func AABBAABBSweep2(a, b *AABB, aV, bV v.Vec, hit *HitInfo) bool {
	delta := bV.Sub(aV)
	isCollide := AABBAABBSweep1(a, b, delta, hit)
	if isCollide {
		hit.Pos = hit.Pos.Add(aV.Scale(hit.Time))
		if hit.Normal.X != 0 {
			hit.Pos.X = b.Pos.X + (bV.X * hit.Time) - (hit.Normal.X * b.Half.X)
		} else {
			hit.Pos.Y = b.Pos.Y + (bV.Y * hit.Time) - (hit.Normal.Y * b.Half.Y)
		}
	}
	return isCollide
}

// AABBCircleSweep checks for collision between a moving AABB and a moving Circle.
// Returns true if collision occurs during movement, false otherwise.
func AABBCircleSweep(a *AABB, c *Circle, aDelta, cDelta v.Vec) bool {
	// Calculate circle's movement relative to AABB (AABB becomes stationary reference frame)
	relativeDelta := cDelta.Sub(aDelta)

	// Use Segment to check if circle (treated as point with radius padding) hits the AABB
	// padding expands AABB by circle's radius to simplify collision detection
	return AABBSegmentOverlap(a, c.Pos, relativeDelta, v.Vec{X: c.Radius, Y: c.Radius}, nil)
}
