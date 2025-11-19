package coll

import "github.com/setanarut/v"

// AABBCircleSweep checks for collision between a moving AABB and a moving Circle.
// Returns true if collision occurs during movement, false otherwise.
func AABBCircleSweep(aabb *AABB, circle *Circle, boxVel, CircleVel v.Vec) bool {
	// Calculate circle's movement relative to AABB (AABB becomes stationary reference frame)
	relativeDelta := CircleVel.Sub(boxVel)

	// Use Segment to check if circle (treated as point with radius padding) hits the AABB
	// padding expands AABB by circle's radius to simplify collision detection
	return AABBSegmentOverlap(aabb, circle.Pos, relativeDelta, v.Vec{X: circle.Radius, Y: circle.Radius}, nil)
}
