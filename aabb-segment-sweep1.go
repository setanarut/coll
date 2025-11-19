package coll

import (
	"math"

	"github.com/setanarut/v"
)

// AABBSegmentSweep1 sweep a moving aabb against a line segment.
//
// Return true when the box hits the segment.
//
// Params:
//
//   - line - non-moving line segment
//   - aabb - moving box
//   - delta - delta movement vector of the aabb
//   - hitInfo - optional contact data structure filled on hit, can be nil
func AABBSegmentSweep1(line *Segment, aabb *AABB, delta v.Vec, hitInfo *HitInfo) bool {

	var lineMin, lineMax v.Vec

	aabbCenter := aabb.Pos
	aabbMin := aabb.Min()
	aabbMax := aabb.Max()

	normalizedDelta := delta.Unit()

	// calculate line bounds
	lineDir := line.B.Sub(line.A)
	if lineDir.X > 0.0 {
		lineMin.X = line.A.X
		lineMax.X = line.B.X
	} else {
		lineMin.X = line.B.X
		lineMax.X = line.A.X
	}
	if lineDir.Y > 0.0 {
		lineMin.Y = line.A.Y
		lineMax.Y = line.B.Y
	} else {
		lineMin.Y = line.B.Y
		lineMax.Y = line.A.Y
	}

	// get aabb's center to line.A distance
	lineAabbDist := line.A.Sub(aabbCenter)

	// get the line's normal
	// if the dot product of it and the delta is larger than 0,
	// it means the line's normal is facing away from the sweep
	lineNormal := segmentNormal(line.A, line.B)
	hitNormal := lineNormal

	hitTime := 0.0
	outTime := 1.0

	// calculate the radius of the box in respect to the line normal
	r := aabb.Half.X*math.Abs(lineNormal.X) + aabb.Half.Y*math.Abs(lineNormal.Y)

	// distance from box to line in respect to the line normal
	boxProj := lineAabbDist.Dot(lineNormal)

	// velocity, projected on the line normal
	velProj := delta.Dot(lineNormal)

	// inverse the radius if required
	if velProj < 0 {
		r *= -1
	}

	// calculate first and last overlap times,
	// as if we're dealing with a line rather than a segment
	hitTime = math.Max((boxProj-r)/velProj, hitTime)
	outTime = math.Min((boxProj+r)/velProj, outTime)

	// run standard AABBvsAABB sweep
	// against an AABB constructed from the extents of the line segment
	// X axis overlap
	if delta.X < 0 {
		// sweeping left
		if aabbMax.X < lineMin.X {
			return false
		}

		hit := (lineMax.X - aabbMin.X) / delta.X
		out := (lineMin.X - aabbMax.X) / delta.X
		outTime = min(out, outTime)
		if hit >= hitTime && hit <= outTime {
			// box is hitting the line on its end:
			// adjust the normal accordingly
			hitNormal = v.Right
		}
		hitTime = max(hit, hitTime)

	} else if delta.X > 0 {
		// sweeping right
		if aabbMin.X > lineMax.X {
			return false
		}

		hit := (lineMin.X - aabbMax.X) / delta.X
		out := (lineMax.X - aabbMin.X) / delta.X
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
	if delta.Y < 0 {
		// sweeping up
		if aabbMax.Y < lineMin.Y {
			return false
		}

		hit := (lineMax.Y - aabbMin.Y) / delta.Y
		out := (lineMin.Y - aabbMax.Y) / delta.Y
		outTime = min(out, outTime)
		if hit >= hitTime && hit <= outTime {
			hitNormal = v.Down
		}

		hitTime = max(hit, hitTime)

	} else if delta.Y > 0 {
		// sweeping down
		if aabbMin.Y > lineMax.Y {
			return false
		}

		hit := (lineMin.Y - aabbMax.Y) / delta.Y
		out := (lineMax.Y - aabbMin.Y) / delta.Y
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

	if hitInfo != nil {
		hitInfo.Delta.X = delta.X*hitTime + (Padding * hitNormal.X)
		hitInfo.Delta.Y = delta.Y*hitTime + (Padding * hitNormal.Y)
		hitInfo.Pos = aabb.Pos.Add(hitInfo.Delta)
		hitInfo.Normal = hitNormal
		hitInfo.Time = hitTime
		// hitInfo.collider = line
	}
	return true
}

func segmentNormal(pos1, pos2 v.Vec) (out v.Vec) {
	d := pos2.Sub(pos1)
	if d.IsZero() {
		return v.Vec{}
	}
	// match JS: perpendicular = [dy, -dx]
	out = v.Vec{X: d.Y, Y: -d.X}
	return out.Unit()
}
