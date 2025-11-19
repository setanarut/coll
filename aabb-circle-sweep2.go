package coll

import "github.com/setanarut/v"

// AABBCircleSweep2 checks for collision between a moving AABB and a moving Circle.
//
// Returns true if collision occurs during movement, false otherwise.
func AABBCircleSweep2(box *AABB, circle *Circle, boxVel, circleVel v.Vec) bool {
	// Calculate circle's movement relative to AABB (AABB becomes stationary reference frame)
	relativeDelta := circleVel.Sub(boxVel)

	// Use Segment to check if circle (treated as point with radius padding) hits the AABB
	// padding expands AABB by circle's radius to simplify collision detection
	return AABBSegmentOverlap(box, circle.Pos, relativeDelta, v.Vec{X: circle.Radius, Y: circle.Radius}, nil)
}
