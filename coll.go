// 2D Collision functions
package coll

import (
	"math"

	"github.com/setanarut/v"
)

const EPSILON = 1e-8

type HitInfo struct {
	Pos    v.Vec
	Delta  v.Vec
	Normal v.Vec
	Time   float64
}

type HitInfo2 struct {
	Right, Bottom, Left, Top bool
	Delta                    v.Vec
}

func (h *HitInfo) Reset() {
	*h = HitInfo{}
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
			h.Delta.Y = (platform.Pos.Y + combinedHalfH + EPSILON) - a.Pos.Y
			h.Top = true
		} else if a.Pos.Y < platform.Pos.Y && oldYDist >= combinedHalfH {
			h.Delta.Y = (platform.Pos.Y - combinedHalfH - EPSILON) - a.Pos.Y
			h.Bottom = true
		}
	}

	if xDist < combinedHalfW {
		if a.Pos.X > platform.Pos.X && oldXDist >= combinedHalfW {
			h.Delta.X = (platform.Pos.X + combinedHalfW + EPSILON) - a.Pos.X
			h.Left = true
		} else if a.Pos.X < platform.Pos.X && oldXDist >= combinedHalfW {
			h.Delta.X = (platform.Pos.X - combinedHalfW - EPSILON) - a.Pos.X
			h.Right = true
		}
	}

	return true
}

func AABBSegmentOverlap(a *AABB, pos, delta, padding v.Vec, hit *HitInfo) bool {
	scale := v.One.Div(delta)
	signX := math.Copysign(1, scale.X)
	signY := math.Copysign(1, scale.Y)
	nearTimeX := (a.Pos.X - signX*(a.Half.X+padding.X) - pos.X) * scale.X
	nearTimeY := (a.Pos.Y - signY*(a.Half.Y+padding.Y) - pos.Y) * scale.Y
	farTimeX := (a.Pos.X + signX*(a.Half.X+padding.X) - pos.X) * scale.X
	farTimeY := (a.Pos.Y + signY*(a.Half.Y+padding.Y) - pos.Y) * scale.Y
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

// OverlapSweep returns hit info for b
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

// AABBAABBSweep1 returns hit info for dynamicB
func AABBAABBSweep1(staticA, dynamicB *AABB, bDelta v.Vec, hit *HitInfo) bool {
	if bDelta.IsZero() {
		return AABBOverlap(staticA, dynamicB, hit)
	}
	result := AABBSegmentOverlap(staticA, dynamicB.Pos, bDelta, dynamicB.Half, hit)
	if result {
		hit.Time = max(0, min(1, hit.Time))
		direction := bDelta.Unit()
		hit.Pos.X = max(staticA.Pos.X-staticA.Half.X, min(staticA.Pos.X+staticA.Half.X, hit.Pos.X+direction.X*dynamicB.Half.X))
		hit.Pos.Y = max(hit.Pos.Y+direction.Y*dynamicB.Half.Y, min(staticA.Pos.Y-staticA.Half.Y, staticA.Pos.Y+staticA.Half.Y))
	}
	return result
}

// AABBAABBSweep2 returns hit info for b
func AABBAABBSweep2(a, b *AABB, aDelta, bDelta v.Vec, hit *HitInfo) bool {
	delta := bDelta.Sub(aDelta)
	isCollide := AABBAABBSweep1(a, b, delta, hit)
	if isCollide {
		hit.Pos = hit.Pos.Add(aDelta.Scale(hit.Time))
		if hit.Normal.X != 0 {
			hit.Pos.X = b.Pos.X + (bDelta.X * hit.Time) - (hit.Normal.X * b.Half.X)
		} else {
			hit.Pos.Y = b.Pos.Y + (bDelta.Y * hit.Time) - (hit.Normal.Y * b.Half.Y)
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
