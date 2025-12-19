package coll

import (
	"math"
	"testing"

	"github.com/setanarut/v"
)

func TestBoxPointOverlap(t *testing.T) {
	box := &AABB{
		Pos:  v.Vec{X: 0, Y: 0},
		Half: v.Vec{X: 10, Y: 5},
	}
	tests := []struct {
		name     string
		point    v.Vec
		expected bool
		checkHit bool
	}{
		{
			name:     "Point outside returns false",
			point:    v.Vec{X: 20, Y: 0},
			expected: false,
		},
		{
			name:     "Point inside returns true",
			point:    v.Vec{X: 1, Y: 1},
			expected: true,
		},
		{
			name:     "HitInfo X-axis response (Collision)",
			point:    v.Vec{X: 9, Y: 0},
			expected: true,
			checkHit: true,
		},
		{
			name:     "HitInfo Y-axis response (Collision)",
			point:    v.Vec{X: 0, Y: 4.5},
			expected: true,
			checkHit: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h HitInfo
			hit := BoxPointOverlap(box, tt.point, &h)
			if hit != tt.expected {
				t.Errorf("BoxPointOverlap(%v) = %v, expected %v", tt.point, hit, tt.expected)
				return
			}
			if hit && tt.checkHit {
				if tt.name == "HitInfo X-axis response (Collision)" {
					if math.Abs(h.Delta.X) == 0 {
						t.Errorf("expected non-zero X delta, got %v", h.Delta.X)
					}
					if h.Normal.X == 0 {
						t.Errorf("expected X normal, got %v", h.Normal)
					}
					if h.Pos.Y != tt.point.Y {
						t.Errorf("expected hit pos Y == point Y (%v), got %v", tt.point.Y, h.Pos.Y)
					}
				}
				if tt.name == "HitInfo Y-axis response (Collision)" {
					if math.Abs(h.Delta.Y) == 0 {
						t.Errorf("expected non-zero Y delta, got %v", h.Delta.Y)
					}
					if h.Normal.Y == 0 {
						t.Errorf("expected Y normal, got %v", h.Normal)
					}
					if h.Pos.X != tt.point.X {
						t.Errorf("expected hit pos X == point X (%v), got %v", tt.point.X, h.Pos.X)
					}
				}
			}
		})
	}
}
func TestBoxOrientedBoxSweep2(t *testing.T) {
	tests := []struct {
		name     string
		aabb     *AABB
		obb      *OBB
		va       v.Vec
		vb       v.Vec
		expected bool
		desc     string
	}{
		{
			name:     "High speed tunneling",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 0}, Half: v.Vec{X: 1, Y: 1}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -20, Y: 0},
			expected: true,
			desc:     "Fast moving OBB should be detected even if it would tunnel through AABB",
		},
		{
			name:     "Extreme speed tunneling - vertical",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 0, Y: 15}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: 0, Y: -30},
			expected: true,
			desc:     "Vertical high-speed tunneling should be detected",
		},
		{
			name:     "Both objects moving towards each other",
			aabb:     &AABB{Pos: v.Vec{X: -5, Y: 0}, Half: v.Vec{X: 0.5, Y: 0.5}},
			obb:      &OBB{Pos: v.Vec{X: 5, Y: 0}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 15, Y: 0},
			vb:       v.Vec{X: -15, Y: 0},
			expected: true,
			desc:     "Both objects moving towards each other at high speed should collide",
		},
		{
			name:     "Diagonal high-speed tunneling",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 10}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -25, Y: -25},
			expected: true,
			desc:     "Diagonal high-speed movement should be detected",
		},
		{
			name:     "Rotated OBB high-speed tunneling",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 8, Y: 0}, Half: v.Vec{X: 2, Y: 0.5}, Angle: math.Pi / 4},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -20, Y: 0},
			expected: true,
			desc:     "Rotated OBB with high speed should still be detected",
		},
		{
			name:     "Small object high-speed through large object",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 5, Y: 5}},
			obb:      &OBB{Pos: v.Vec{X: 20, Y: 0}, Half: v.Vec{X: 0.1, Y: 0.1}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -50, Y: 0},
			expected: true,
			desc:     "Small fast object tunneling through large object should be detected",
		},
		{
			name:     "Parallel high-speed miss",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 0, Y: 10}, Half: v.Vec{X: 1, Y: 1}, Angle: 0},
			va:       v.Vec{X: 50, Y: 0},
			vb:       v.Vec{X: 50, Y: 0},
			expected: false,
			desc:     "Parallel movement should not collide",
		},
		{
			name:     "Near-miss high-speed",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 2.1}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -30, Y: 0},
			expected: false,
			desc:     "Near-miss should not collide",
		},
		{
			name:     "Bullet-like speed test",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 2, Y: 2}},
			obb:      &OBB{Pos: v.Vec{X: 100, Y: 0}, Half: v.Vec{X: 0.2, Y: 0.2}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -200, Y: 0},
			expected: true,
			desc:     "Bullet-like speed should still be detected",
		},
		{
			name:     "Rotated OBB spinning and moving fast",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 8, Y: 0}, Half: v.Vec{X: 3, Y: 0.5}, Angle: math.Pi / 3},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -25, Y: 0},
			expected: true,
			desc:     "Rotated elongated OBB should be detected",
		},
		{
			name:     "Edge case - exactly swept volume boundary",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 5, Y: 0}, Half: v.Vec{X: 1, Y: 1}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -3, Y: 0},
			expected: true,
			desc:     "Boundary case should be handled correctly",
		},
		{
			name:     "Moving away at high speed",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 0}, Half: v.Vec{X: 1, Y: 1}, Angle: 0},
			va:       v.Vec{X: -50, Y: 0},
			vb:       v.Vec{X: 50, Y: 0},
			expected: false,
			desc:     "Objects moving away from each other should not collide",
		},
		{
			name:     "Extreme: Speed of light simulation",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 1000, Y: 0}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -10000, Y: 0},
			expected: true,
			desc:     "Very high speed collision",
		},
		{
			name:     "Extreme: Clear miss with large gap",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 5}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -100, Y: 0},
			expected: false,
			desc:     "Objects should clearly miss each other",
		},
		{
			name:     "Extreme: Grazing hit - top edge",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 1.5}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 0, Y: 0},
			vb:       v.Vec{X: -20, Y: 0},
			expected: true,
			desc:     "A near miss/grazing hit should be detected",
		},
		{
			name:     "Extreme: Zero relative velocity (Miss)",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 3, Y: 0}, Half: v.Vec{X: 1, Y: 1}, Angle: 0},
			va:       v.Vec{X: 5, Y: 0},
			vb:       v.Vec{X: 5, Y: 0},
			expected: false,
			desc:     "Moving parallel with same speed should not collide if they start separated",
		},
		{
			name:     "Extreme: Perpendicular high-speed pass (Miss)",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 0, Y: 10}, Half: v.Vec{X: 0.5, Y: 0.5}, Angle: 0},
			va:       v.Vec{X: 100, Y: 0},
			vb:       v.Vec{X: 0, Y: -5},
			expected: false,
			desc:     "Objects passing perpendicularly without overlap should not collide",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoxOrientedBoxSweep2(tt.aabb, tt.obb, tt.va, tt.vb)
			if result != tt.expected {
				t.Errorf("%s: Expected %v, got %v\nAABB(P:%v, H:%v, V:%v), OBB(P:%v, H:%v, R:%v, V:%v)",
					tt.desc, tt.expected, result,
					tt.aabb.Pos, tt.aabb.Half, tt.va,
					tt.obb.Pos, tt.obb.Half, tt.obb.Angle, tt.vb)
			}
		})
	}
}
func TestBoxOrientedBoxOverlap(t *testing.T) {
	tests := []struct {
		name     string
		aabb     *AABB
		obb      *OBB
		expected bool
	}{
		{
			name:     "1. Overlap - Simple Alignment",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 2, Y: 2}},
			obb:      &OBB{Pos: v.Vec{X: 1, Y: 1}, Half: v.Vec{X: 2, Y: 2}, Angle: 0},
			expected: true,
		},
		{
			name:     "2. No Overlap - Clear Miss",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 5, Y: 5}, Half: v.Vec{X: 1, Y: 1}, Angle: 0},
			expected: false,
		},
		{
			name:     "3. Overlap - OBB Rotated",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 2, Y: 2}},
			obb:      &OBB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}, Angle: math.Pi / 4},
			expected: true,
		},
		{
			name:     "4. No Overlap - OBB Rotated Miss",
			aabb:     &AABB{Pos: v.Vec{X: 0, Y: 0}, Half: v.Vec{X: 1, Y: 1}},
			obb:      &OBB{Pos: v.Vec{X: 10, Y: 0}, Half: v.Vec{X: 1, Y: 1}, Angle: math.Pi / 4},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := BoxOrientedBoxOverlap(tt.aabb, tt.obb)
			if actual != tt.expected {
				t.Errorf("BoxOrientedBoxOverlap() = %v, expected %v\nAABB: %+v, OBB: %+v", actual, tt.expected, tt.aabb, tt.obb)
			}
		})
	}
}
func TestBoxBoxSweep2(t *testing.T) {
	test := struct {
		name     string
		boxA     *AABB
		boxB     *AABB
		va       v.Vec
		vb       v.Vec
		expected bool
	}{
		name: "High-speed tunneling collision detection",
		boxA: &AABB{
			Pos:  v.Vec{X: 0, Y: 0},
			Half: v.Vec{X: 1, Y: 1},
		},
		boxB: &AABB{
			Pos:  v.Vec{X: 5, Y: 0},
			Half: v.Vec{X: 0.5, Y: 0.5},
		},
		va:       v.Vec{X: 0, Y: 0},
		vb:       v.Vec{X: -100, Y: 0},
		expected: true,
	}
	t.Run(test.name, func(t *testing.T) {
		var h HitInfo
		isCollide := BoxBoxSweep2(test.boxA, test.boxB, test.va, test.vb, &h)
		if isCollide != test.expected {
			t.Errorf("BoxBoxSweep2() = %v, expected %v", isCollide, test.expected)
			return
		}
		if isCollide {
			if h.Time <= 0 || h.Time >= 1.0 {
				t.Errorf("Expected collision time (h.Time) between 0 and 1, got %v", h.Time)
			}
			if h.Pos.X > 1.5 {
				t.Errorf("Hit position X seems too high for collision near AABB edge: %v", h.Pos.X)
			}
		}
	})
}
func TestBoxCircleSweep2(t *testing.T) {
	test := struct {
		name      string
		box       *AABB
		circle    *Circle
		boxVel    v.Vec
		circleVel v.Vec
		expected  bool
	}{
		name: "High-speed Circle tunneling through static AABB",
		box: &AABB{
			Pos:  v.Vec{X: 0, Y: 0},
			Half: v.Vec{X: 1, Y: 1},
		},
		circle: &Circle{
			Pos:    v.Vec{X: 10, Y: 0},
			Radius: 0.5,
		},
		boxVel:    v.Vec{X: 0, Y: 0},
		circleVel: v.Vec{X: -20, Y: 0},
		expected:  true,
	}
	t.Run(test.name, func(t *testing.T) {
		isCollide := BoxCircleSweep2(test.box, test.circle, test.boxVel, test.circleVel)
		if isCollide != test.expected {
			t.Errorf("BoxCircleSweep2() = %v, expected %v\nBox: %+v, Circle: %+v, Vb: %v, Vc: %v",
				isCollide, test.expected,
				test.box, test.circle, test.boxVel, test.circleVel)
		}
	})
}
