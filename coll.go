// 2D Collision functions
package coll

import (
	"github.com/setanarut/v"
)

const (
	Epsilon float64 = 1e-8
	Padding float64 = 0.005
)

// Hit holds the information about a collision or contact event.
type Hit struct {
	Normal v.Vec   // The normal vector of the surface hit.
	Time   float64 // The time (0.0 to 1.0) along the movement path when the collision happened.
}

// Resets the zero values.
func (h *Hit) Reset() {
	*h = Hit{} // Reinitializes all fields of the struct to their zero values (nil, 0, false, etc.).
}

func SegmentNormal(pos1, pos2 v.Vec) (normal v.Vec) {
	d := pos2.Sub(pos1)
	if d.IsZero() {
		return v.Vec{}
	}
	normal = v.Vec{X: d.Y, Y: -d.X}
	return normal.Unit()
}
