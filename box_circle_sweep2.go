package coll

import (
	"math"

	"github.com/setanarut/v"
)

func BoxCircleSweep2(a *AABB, b *Circle, deltaA, deltaB v.Vec, h *Hit) bool {
	relDelta := deltaB.Sub(deltaA)

	rad := v.Vec{X: b.Radius, Y: b.Radius}
	boxMin := a.Pos.Sub(a.Half).Sub(rad)
	boxMax := a.Pos.Add(a.Half).Add(rad)

	var tminX, tmaxX, tminY, tmaxY float64
	var hitX, hitY bool

	if math.Abs(relDelta.X) > 1e-8 {
		invX := 1.0 / relDelta.X
		t1, t2 := (boxMin.X-b.Pos.X)*invX, (boxMax.X-b.Pos.X)*invX
		tminX, tmaxX, hitX = min(t1, t2), max(t1, t2), true
	} else if b.Pos.X < boxMin.X || b.Pos.X > boxMax.X {
		return false
	} else {
		tminX, tmaxX = -math.MaxFloat64, math.MaxFloat64
	}

	if math.Abs(relDelta.Y) > 1e-8 {
		invY := 1.0 / relDelta.Y
		t3, t4 := (boxMin.Y-b.Pos.Y)*invY, (boxMax.Y-b.Pos.Y)*invY
		tminY, tmaxY, hitY = min(t3, t4), max(t3, t4), true
	} else if b.Pos.Y < boxMin.Y || b.Pos.Y > boxMax.Y {
		return false
	} else {
		tminY, tmaxY = -math.MaxFloat64, math.MaxFloat64
	}

	tmin := max(tminX, tminY)
	tmax := min(tmaxX, tmaxY)

	if tmax < 0 || tmin > tmax || tmin > 1.0 {
		return false
	}

	if tmin < 0 {
		hit := BoxCircleOverlap(a, b, h)
		if hit && h != nil {
			h.Data = 0
		}
		return hit
	}

	if h == nil {
		return true
	}

	switch {
	case !hitX:
		h.Normal = v.Vec{X: 0, Y: math.Copysign(1, -relDelta.Y)}
	case !hitY:
		h.Normal = v.Vec{X: math.Copysign(1, -relDelta.X), Y: 0}
	case tminX > tminY:
		h.Normal = v.Vec{X: math.Copysign(1, -relDelta.X), Y: 0}
	default:
		h.Normal = v.Vec{X: 0, Y: math.Copysign(1, -relDelta.Y)}
	}
	h.Data = tmin
	return true
}
