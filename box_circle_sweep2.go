package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxCircleSweep2 checks for collision between a moving AABB and a moving Circle.
//
// Returns true if collision occurs during movement, false otherwise.
func BoxCircleSweep2(a *AABB, b *Circle, deltaA, deltaB v.Vec, h *Hit) bool {
	relDeltaX := deltaB.X - deltaA.X
	relDeltaY := deltaB.Y - deltaA.Y

	// If no relative movement, check overlap directly
	relDeltaMagSq := relDeltaX*relDeltaX + relDeltaY*relDeltaY
	if relDeltaMagSq < 1e-8 {
		// Find closest point on AABB to circle center
		closestX := math.Max(a.Pos.X-a.Half.X, math.Min(b.Pos.X, a.Pos.X+a.Half.X))
		closestY := math.Max(a.Pos.Y-a.Half.Y, math.Min(b.Pos.Y, a.Pos.Y+a.Half.Y))

		// Vector from closest point to circle center
		distX := b.Pos.X - closestX
		distY := b.Pos.Y - closestY
		distSq := distX*distX + distY*distY

		// Check if overlapping
		if distSq >= b.Radius*b.Radius {
			return false
		}

		if h != nil {
			dist := math.Sqrt(distSq)

			// If circle center is inside box
			if distSq < 1e-8 {
				// Calculate penetration depths for each side
				leftDist := (b.Pos.X - (a.Pos.X - a.Half.X)) + b.Radius
				rightDist := ((a.Pos.X + a.Half.X) - b.Pos.X) + b.Radius
				topDist := (b.Pos.Y - (a.Pos.Y - a.Half.Y)) + b.Radius
				bottomDist := ((a.Pos.Y + a.Half.Y) - b.Pos.Y) + b.Radius

				// Find minimum penetration
				minDist := math.Min(math.Min(leftDist, rightDist), math.Min(topDist, bottomDist))

				if minDist == leftDist {
					h.Normal.X = -1
					h.Normal.Y = 0
				} else if minDist == rightDist {
					h.Normal.X = 1
					h.Normal.Y = 0
				} else if minDist == topDist {
					h.Normal.X = 0
					h.Normal.Y = -1
				} else {
					h.Normal.X = 0
					h.Normal.Y = 1
				}
			} else {
				// Normal from box to circle
				h.Normal.X = distX / dist
				h.Normal.Y = distY / dist
			}

			h.Data = 0
		}

		return true
	}

	// Expanded box bounds (box + circle radius)
	halfX := a.Half.X + b.Radius
	halfY := a.Half.Y + b.Radius
	boxMinX := a.Pos.X - halfX
	boxMaxX := a.Pos.X + halfX
	boxMinY := a.Pos.Y - halfY
	boxMaxY := a.Pos.Y + halfY

	// AABB ray intersection (slab method)
	var tminX, tmaxX, tminY, tmaxY float64
	var hitX, hitY bool

	if math.Abs(relDeltaX) > 1e-8 {
		invDirX := 1.0 / relDeltaX
		t1 := (boxMinX - b.Pos.X) * invDirX
		t2 := (boxMaxX - b.Pos.X) * invDirX
		tminX = min(t1, t2)
		tmaxX = max(t1, t2)
		hitX = true
	} else {
		if b.Pos.X < boxMinX || b.Pos.X > boxMaxX {
			return false
		}
		tminX = -math.MaxFloat64
		tmaxX = math.MaxFloat64
		hitX = false
	}

	if math.Abs(relDeltaY) > 1e-8 {
		invDirY := 1.0 / relDeltaY
		t3 := (boxMinY - b.Pos.Y) * invDirY
		t4 := (boxMaxY - b.Pos.Y) * invDirY
		tminY = min(t3, t4)
		tmaxY = max(t3, t4)
		hitY = true
	} else {
		if b.Pos.Y < boxMinY || b.Pos.Y > boxMaxY {
			return false
		}
		tminY = -math.MaxFloat64
		tmaxY = math.MaxFloat64
		hitY = false
	}

	tmin := max(tminX, tminY)
	tmax := min(tmaxX, tmaxY)

	// No intersection
	if tmax < 0 || tmin > tmax || tmin > 1.0 {
		return false
	}

	// Already overlapping at t=0
	if tmin < 0 {
		// Find closest point on AABB to circle center
		closestX := math.Max(a.Pos.X-a.Half.X, math.Min(b.Pos.X, a.Pos.X+a.Half.X))
		closestY := math.Max(a.Pos.Y-a.Half.Y, math.Min(b.Pos.Y, a.Pos.Y+a.Half.Y))

		// Vector from closest point to circle center
		distX := b.Pos.X - closestX
		distY := b.Pos.Y - closestY
		distSq := distX*distX + distY*distY

		// Check if overlapping
		if distSq >= b.Radius*b.Radius {
			return false
		}

		if h != nil {
			dist := math.Sqrt(distSq)

			// If circle center is inside box
			if distSq < 1e-8 {
				// Calculate penetration depths for each side
				leftDist := (b.Pos.X - (a.Pos.X - a.Half.X)) + b.Radius
				rightDist := ((a.Pos.X + a.Half.X) - b.Pos.X) + b.Radius
				topDist := (b.Pos.Y - (a.Pos.Y - a.Half.Y)) + b.Radius
				bottomDist := ((a.Pos.Y + a.Half.Y) - b.Pos.Y) + b.Radius

				// Find minimum penetration
				minDist := math.Min(math.Min(leftDist, rightDist), math.Min(topDist, bottomDist))

				if minDist == leftDist {
					h.Normal.X = -1
					h.Normal.Y = 0
				} else if minDist == rightDist {
					h.Normal.X = 1
					h.Normal.Y = 0
				} else if minDist == topDist {
					h.Normal.X = 0
					h.Normal.Y = -1
				} else {
					h.Normal.X = 0
					h.Normal.Y = 1
				}
			} else {
				// Normal from box to circle
				h.Normal.X = distX / dist
				h.Normal.Y = distY / dist
			}

			h.Data = 0
		}

		return true
	}

	if h == nil {
		return true
	}

	// Determine which face was hit based on which axis contributed to tmin
	if !hitX {
		// Only Y axis moving
		h.Normal.X = 0
		h.Normal.Y = math.Copysign(1, -relDeltaY)
	} else if !hitY {
		// Only X axis moving
		h.Normal.X = math.Copysign(1, -relDeltaX)
		h.Normal.Y = 0
	} else if tminX > tminY {
		// X axis constrained the hit
		h.Normal.X = math.Copysign(1, -relDeltaX)
		h.Normal.Y = 0
	} else {
		// Y axis constrained the hit
		h.Normal.X = 0
		h.Normal.Y = math.Copysign(1, -relDeltaY)
	}

	h.Data = tmin

	return true
}
