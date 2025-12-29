package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxSegmentSweep1 sweep a moving box against a static line segment.
//
// If h is not nil and a collision is detected, it will be populated with:
//   - Normal: Collision surface normal for the box
//   - Data: Normalized time of impact (0.0 to 1.0) along the movement path
func BoxSegmentSweep1(s *Segment, a *AABB, deltaA v.Vec, h *Hit) bool {

	var lineMin, lineMax v.Vec

	aabbCenter := a.Pos
	aabbMin := a.Min()
	aabbMax := a.Max()

	normalizedDelta := deltaA.Unit()

	// calculate line bounds
	lineDir := s.B.Sub(s.A)
	if lineDir.X > 0.0 {
		lineMin.X = s.A.X
		lineMax.X = s.B.X
	} else {
		lineMin.X = s.B.X
		lineMax.X = s.A.X
	}
	if lineDir.Y > 0.0 {
		lineMin.Y = s.A.Y
		lineMax.Y = s.B.Y
	} else {
		lineMin.Y = s.B.Y
		lineMax.Y = s.A.Y
	}

	// get aabb's center to line.A distance
	lineAabbDist := s.A.Sub(aabbCenter)

	// get the line's normal
	// if the dot product of it and the delta is larger than 0,
	// it means the line's normal is facing away from the sweep
	lineNormal := SegmentNormal(s.A, s.B)
	hitNormal := lineNormal

	hitTime := 0.0
	outTime := 1.0

	// calculate the radius of the box in respect to the line normal
	r := a.Half.X*math.Abs(lineNormal.X) + a.Half.Y*math.Abs(lineNormal.Y)

	// distance from box to line in respect to the line normal
	boxProj := lineAabbDist.Dot(lineNormal)

	dProj := deltaA.Dot(lineNormal)

	// inverse the radius if required
	if dProj < 0 {
		r *= -1
	}

	// calculate first and last overlap times,
	// as if we're dealing with a line rather than a segment
	hitTime = math.Max((boxProj-r)/dProj, hitTime)
	outTime = math.Min((boxProj+r)/dProj, outTime)

	// run standard AABBvsAABB sweep
	// against an AABB constructed from the extents of the line segment
	// X axis overlap
	if deltaA.X < 0 {
		// sweeping left
		if aabbMax.X < lineMin.X {
			return false
		}

		hit := (lineMax.X - aabbMin.X) / deltaA.X
		out := (lineMin.X - aabbMax.X) / deltaA.X
		outTime = min(out, outTime)
		if hit >= hitTime && hit <= outTime {
			// box is hitting the line on its end:
			// adjust the normal accordingly
			hitNormal = v.Right
		}
		hitTime = max(hit, hitTime)

	} else if deltaA.X > 0 {
		// sweeping right
		if aabbMin.X > lineMax.X {
			return false
		}

		hit := (lineMin.X - aabbMax.X) / deltaA.X
		out := (lineMax.X - aabbMin.X) / deltaA.X
		outTime = min(out, outTime)
		if hit >= hitTime && hit <= outTime {
			hitNormal = v.Left
		}

		hitTime = max(hit, hitTime)

	} else if lineMin.X > aabbMax.X || lineMax.X < aabbMin.X {
		return false
	}

	if hitTime > outTime {
		return false
	}

	// Y axis overlap
	if deltaA.Y < 0 {
		// sweeping up
		if aabbMax.Y < lineMin.Y {
			return false
		}

		hit := (lineMax.Y - aabbMin.Y) / deltaA.Y
		out := (lineMin.Y - aabbMax.Y) / deltaA.Y
		outTime = min(out, outTime)
		if hit >= hitTime && hit <= outTime {
			hitNormal = v.Down
		}

		hitTime = max(hit, hitTime)

	} else if deltaA.Y > 0 {
		// sweeping down
		if aabbMin.Y > lineMax.Y {
			return false
		}

		hit := (lineMin.Y - aabbMax.Y) / deltaA.Y
		out := (lineMax.Y - aabbMin.Y) / deltaA.Y
		outTime = min(out, outTime)
		if hit >= hitTime && hit <= outTime {
			hitNormal = v.Up
		}

		hitTime = max(hit, hitTime)

	} else if lineMin.Y > aabbMax.Y || lineMax.Y < aabbMin.Y {
		return false
	}

	if hitTime > outTime {
		return false
	}

	// ignore this line if its normal is facing away from the sweep delta
	// check for this only at this point to account for a possibly changed hitNormal
	// from a hit on the line's end
	//
	// also ignore this line its normal is facing away from the adjusted hitNormal
	if normalizedDelta.Dot(hitNormal) > 0 || lineNormal.Dot(hitNormal) < 0 {
		return false
	}

	if h != nil {
		h.Normal = hitNormal
		h.Data = hitTime
	}
	return true
}
