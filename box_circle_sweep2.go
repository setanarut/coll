package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxCircleSweep2 checks for collision between a moving AABB and a moving Circle.
//
// Returns true if collision occurs during movement, false otherwise.
func BoxCircleSweep2(box *AABB, circle *Circle, boxVel, circleVel v.Vec, h *Hit) bool {
	// Relative velocity (circle relative to box)
	relVelX := circleVel.X - boxVel.X
	relVelY := circleVel.Y - boxVel.Y

	// If no relative movement, check overlap directly
	relVelMagSq := relVelX*relVelX + relVelY*relVelY
	if relVelMagSq < 1e-8 {
		// Find closest point on AABB to circle center
		closestX := math.Max(box.Pos.X-box.Half.X, math.Min(circle.Pos.X, box.Pos.X+box.Half.X))
		closestY := math.Max(box.Pos.Y-box.Half.Y, math.Min(circle.Pos.Y, box.Pos.Y+box.Half.Y))

		// Vector from closest point to circle center
		distX := circle.Pos.X - closestX
		distY := circle.Pos.Y - closestY
		distSq := distX*distX + distY*distY

		// Check if overlapping
		if distSq >= circle.Radius*circle.Radius {
			return false
		}

		if h != nil {
			dist := math.Sqrt(distSq)

			// If circle center is inside box
			if distSq < 1e-8 {
				// Calculate penetration depths for each side
				leftDist := (circle.Pos.X - (box.Pos.X - box.Half.X)) + circle.Radius
				rightDist := ((box.Pos.X + box.Half.X) - circle.Pos.X) + circle.Radius
				topDist := (circle.Pos.Y - (box.Pos.Y - box.Half.Y)) + circle.Radius
				bottomDist := ((box.Pos.Y + box.Half.Y) - circle.Pos.Y) + circle.Radius

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
	halfX := box.Half.X + circle.Radius
	halfY := box.Half.Y + circle.Radius
	boxMinX := box.Pos.X - halfX
	boxMaxX := box.Pos.X + halfX
	boxMinY := box.Pos.Y - halfY
	boxMaxY := box.Pos.Y + halfY

	// AABB ray intersection (slab method)
	var tminX, tmaxX, tminY, tmaxY float64
	var hitX, hitY bool

	if math.Abs(relVelX) > 1e-8 {
		invDirX := 1.0 / relVelX
		t1 := (boxMinX - circle.Pos.X) * invDirX
		t2 := (boxMaxX - circle.Pos.X) * invDirX
		tminX = min(t1, t2)
		tmaxX = max(t1, t2)
		hitX = true
	} else {
		if circle.Pos.X < boxMinX || circle.Pos.X > boxMaxX {
			return false
		}
		tminX = -math.MaxFloat64
		tmaxX = math.MaxFloat64
		hitX = false
	}

	if math.Abs(relVelY) > 1e-8 {
		invDirY := 1.0 / relVelY
		t3 := (boxMinY - circle.Pos.Y) * invDirY
		t4 := (boxMaxY - circle.Pos.Y) * invDirY
		tminY = min(t3, t4)
		tmaxY = max(t3, t4)
		hitY = true
	} else {
		if circle.Pos.Y < boxMinY || circle.Pos.Y > boxMaxY {
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
		closestX := math.Max(box.Pos.X-box.Half.X, math.Min(circle.Pos.X, box.Pos.X+box.Half.X))
		closestY := math.Max(box.Pos.Y-box.Half.Y, math.Min(circle.Pos.Y, box.Pos.Y+box.Half.Y))

		// Vector from closest point to circle center
		distX := circle.Pos.X - closestX
		distY := circle.Pos.Y - closestY
		distSq := distX*distX + distY*distY

		// Check if overlapping
		if distSq >= circle.Radius*circle.Radius {
			return false
		}

		if h != nil {
			dist := math.Sqrt(distSq)

			// If circle center is inside box
			if distSq < 1e-8 {
				// Calculate penetration depths for each side
				leftDist := (circle.Pos.X - (box.Pos.X - box.Half.X)) + circle.Radius
				rightDist := ((box.Pos.X + box.Half.X) - circle.Pos.X) + circle.Radius
				topDist := (circle.Pos.Y - (box.Pos.Y - box.Half.Y)) + circle.Radius
				bottomDist := ((box.Pos.Y + box.Half.Y) - circle.Pos.Y) + circle.Radius

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
		h.Normal.Y = math.Copysign(1, -relVelY)
	} else if !hitY {
		// Only X axis moving
		h.Normal.X = math.Copysign(1, -relVelX)
		h.Normal.Y = 0
	} else if tminX > tminY {
		// X axis constrained the hit
		h.Normal.X = math.Copysign(1, -relVelX)
		h.Normal.Y = 0
	} else {
		// Y axis constrained the hit
		h.Normal.X = 0
		h.Normal.Y = math.Copysign(1, -relVelY)
	}

	h.Data = tmin

	return true
}
