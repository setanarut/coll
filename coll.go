// 2D Collision functions
package coll

import (
	"github.com/setanarut/v"
)

const (
	Epsilon float64 = 1e-8
	Padding float64 = 0.005
)

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

// CalculateSlideVelocity computes the total movement: movement until collision plus sliding along the surface normal.
func CalculateSlideVelocity(vel v.Vec, hitInfo *HitInfo) (slideVel v.Vec) {
	remainingVel := vel.Scale(1.0 - hitInfo.Time)
	slideVel = remainingVel.Sub(hitInfo.Normal.Scale(remainingVel.Dot(hitInfo.Normal)))
	movementToHit := vel.Scale(hitInfo.Time)
	return movementToHit.Add(slideVel)
}

// ApplySlideVelocity updates the AABB position by applying the calculated slide velocity.
func ApplySlideVelocity(aabb *AABB, vel v.Vec, hitInfo *HitInfo) {
	aabb.Pos = aabb.Pos.Add(CalculateSlideVelocity(vel, hitInfo))
}
