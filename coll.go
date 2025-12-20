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
	// The normal vector of the hit.
	Normal v.Vec
	// 1. The time (0.0 to 1.0) along the movement path for moving objects.
	//
	// 2. Penetration depth for overlap tests
	Data float64
}

// Resets the zero values.
func (h *Hit) Reset() {
	*h = Hit{} // Reinitializes all fields of the struct to their zero values (nil, 0, false, etc.).
}

// SegmentNormal returns surface normal of the segment.
func SegmentNormal(pointA, pointB v.Vec) (normal v.Vec) {
	d := pointB.Sub(pointA)
	if d.IsZero() {
		return
	}
	return v.Vec{X: d.Y, Y: -d.X}.Unit()
}
