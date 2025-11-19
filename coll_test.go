package coll

import (
	"math"
	"testing"

	"github.com/setanarut/v"
)

func TestAABBPointOverlap(t *testing.T) {
	box := &AABB{
		Pos:  v.Vec{X: 0, Y: 0},
		Half: v.Vec{X: 10, Y: 5},
	}

	t.Run("Point outside returns false", func(t *testing.T) {
		if AABBPointOverlap(box, v.Vec{20, 0}, nil) {
			t.Fatalf("expected no collision")
		}
	})

	t.Run("Point inside returns true", func(t *testing.T) {
		if !AABBPointOverlap(box, v.Vec{1, 1}, nil) {
			t.Fatalf("expected collision")
		}
	})

	t.Run("HitInfo X-axis response", func(t *testing.T) {
		var h HitInfo
		AABBPointOverlap(box, v.Vec{9, 0}, &h)

		if math.Abs(h.Delta.X) == 0 {
			t.Fatalf("expected X delta")
		}
		if h.Normal.X == 0 {
			t.Fatalf("expected X normal")
		}
		if h.Pos.Y != 0 {
			t.Fatalf("expected hit pos Y == point Y")
		}
	})

	t.Run("HitInfo Y-axis response", func(t *testing.T) {
		var h HitInfo
		AABBPointOverlap(box, v.Vec{0, 4.5}, &h)

		if math.Abs(h.Delta.Y) == 0 {
			t.Fatalf("expected Y delta")
		}
		if h.Normal.Y == 0 {
			t.Fatalf("expected Y normal")
		}
		if h.Pos.X != 0 {
			t.Fatalf("expected hit pos X == point X")
		}
	})
}
