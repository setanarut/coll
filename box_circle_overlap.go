package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxCircleOverlap checks whether box and circle overlap.
// Any collision information written to hitInfo always describes how to move box out of circle.
//
// It uses a separating-axis test: if the circle do not overlap on either X or Y,
// there is no collision and the function returns false.
//
// If hitInfo is not nil, the function fills it with:
//   - Normal: the direction in which the box is pushed
//   - Time: penetration depth normalized (0-1 range based on box dimensions)
//
// This method can behave poorly for moving objects.
//
// If you only need to know whether a overlap occurred, pass nil for hitInfo
// to skip generating overlap details.
func BoxCircleOverlap(box *AABB, circle *Circle, hitInfo *Hit) bool {

	// intersection test
	diff := circle.Pos.Sub(box.Pos)
	clamped := v.Vec{
		X: max(-box.Half.X, min(diff.X, box.Half.X)),
		Y: max(-box.Half.Y, min(diff.Y, box.Half.Y)),
	}
	closest := box.Pos.Add(clamped)
	if !(circle.Pos.DistSq(closest) <= circle.Radius*circle.Radius) {
		return false
	}

	if hitInfo == nil {
		return true
	}

	inside := diff.Equals(clamped)

	if !inside {
		normal := diff.Sub(clamped)
		distSq := normal.MagSq()
		dist := math.Sqrt(distSq)

		// Normalize the normal vector
		hitInfo.Normal = normal.DivS(dist).Scale(-1) // Inverse for box

		// Penetration amount
		penetration := circle.Radius - dist
		// Normalize based on the axis with maximum box dimension
		maxBoxDim := math.Max(box.Half.X, box.Half.Y) * 2
		hitInfo.Time = 1.0 - (penetration / maxBoxDim)

	} else {
		absD := diff.Abs()
		px := box.Half.X - absD.X
		py := box.Half.Y - absD.Y

		if px < py {
			sx := math.Copysign(1, diff.X)
			pushDistance := px + circle.Radius
			hitInfo.Normal = v.Vec{X: -sx, Y: 0} // Inverse for box
			hitInfo.Time = 1.0 - (pushDistance / (box.Half.X * 2))
		} else {
			sy := math.Copysign(1, diff.Y)
			pushDistance := py + circle.Radius
			hitInfo.Normal = v.Vec{X: 0, Y: -sy} // Inverse for box
			hitInfo.Time = 1.0 - (pushDistance / (box.Half.Y * 2))
		}
	}
	return true
}

// ResolvePosition calculates the resolved position by applying the collision response.
// It takes the original position, hit normal, and hit time (penetration distance),
// and returns the new position that separates the objects.
func ResolvePosition(pos v.Vec, normal v.Vec, time float64) v.Vec {
	return pos.Add(normal.Scale(time))
}
