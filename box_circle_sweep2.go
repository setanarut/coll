package coll

import (
	"math"

	"github.com/setanarut/v"
)

// BoxCircleSweep2 checks for collision between a moving AABB and a moving Circle.
//
// Returns true if collision occurs during movement, false otherwise.
func BoxCircleSweep2(box *AABB, circle *Circle, boxVel, circleVel v.Vec) bool {

	boxMin := box.Min()
	boxMax := box.Max()
	R := circle.Radius
	R_sq := R * R

	if boxCircleIntersect(box, circle) {
		return true
	}

	V := circleVel.Sub(boxVel)
	t_enter := 0.0
	t_exit := 1.0

	expandedMin := boxMin.Sub(v.Vec{X: R, Y: R})
	expandedMax := boxMax.Add(v.Vec{X: R, Y: R})

	if V.X != 0.0 {
		invD := 1.0 / V.X
		t1 := (expandedMin.X - circle.Pos.X) * invD
		t2 := (expandedMax.X - circle.Pos.X) * invD

		if t1 > t2 {
			t1, t2 = t2, t1
		}

		t_enter = max(t_enter, t1)
		t_exit = min(t_exit, t2)
	} else {
		if circle.Pos.X < expandedMin.X || circle.Pos.X > expandedMax.X {
			return false
		}
	}

	if V.Y != 0.0 {
		invD := 1.0 / V.Y
		t1 := (expandedMin.Y - circle.Pos.Y) * invD
		t2 := (expandedMax.Y - circle.Pos.Y) * invD

		if t1 > t2 {
			t1, t2 = t2, t1
		}

		t_enter = max(t_enter, t1)
		t_exit = min(t_exit, t2)
	} else {
		if circle.Pos.Y < expandedMin.Y || circle.Pos.Y > expandedMax.Y {
			return false
		}
	}

	if t_enter > t_exit || t_enter > 1.0 {
		return false
	}
	corners := []v.Vec{boxMin, {X: boxMax.X, Y: boxMin.Y}, boxMax, {X: boxMin.X, Y: boxMax.Y}}
	min_t_corner := 2.0

	for _, corner := range corners {
		P_diff := circle.Pos.Sub(corner)

		A := V.MagSq()
		B := 2.0 * V.Dot(P_diff)
		C := P_diff.MagSq() - R_sq

		if A < 0.000001 {
			continue
		}

		discriminant := B*B - 4.0*A*C

		if discriminant < 0.0 {
			continue
		}

		sqrt_disc := math.Sqrt(discriminant)
		t1 := (-B - sqrt_disc) / (2.0 * A)
		t2 := (-B + sqrt_disc) / (2.0 * A)

		t_col := min(t1, t2)

		if t_col > 0.0 && t_col <= 1.0 {
			min_t_corner = min(min_t_corner, t_col)
		}
	}

	final_t_col := min(t_enter, min_t_corner)

	return final_t_col > 0.0 && final_t_col <= 1.0
}
