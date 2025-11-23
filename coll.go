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

func SegmentNormal(pos1, pos2 v.Vec) (normal v.Vec) {
	d := pos2.Sub(pos1)
	if d.IsZero() {
		return v.Vec{}
	}
	normal = v.Vec{X: d.Y, Y: -d.X}
	return normal.Unit()
}

// CollideAndSlide updates the AABB position by applying the calculated slide velocity parallel to the surface.
func CollideAndSlide(box *AABB, vel v.Vec, hitInfo *HitInfo) {
	box.Pos = box.Pos.Add(CalculateSlideVelocity(vel, hitInfo))
}

// CalculateSlideVelocity preserves the magnitude of the sliding velocity parallel to the surface, whether flat or inclined.
func CalculateSlideVelocity(vel v.Vec, hitInfo *HitInfo) (slideVel v.Vec) {
	movementToHit := vel.Scale(hitInfo.Time)
	remainingVel := vel.Sub(movementToHit)
	originalSpeed := remainingVel.Mag()
	slideDirection := remainingVel.Sub(hitInfo.Normal.Scale(remainingVel.Dot(hitInfo.Normal)))
	if slideDirection.MagSq() < 1e-6 {
		return movementToHit
	}
	slideDirectionUnit := slideDirection.Unit()
	scaledSlideDirection := slideDirectionUnit.Scale(originalSpeed)
	return movementToHit.Add(scaledSlideDirection)
}
