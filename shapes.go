package coll

import "github.com/setanarut/v" // External package for the Vector (Vec) type.

// AABB represents an Axis-Aligned Bounding Box.
type AABB struct {
	Pos  v.Vec // Center position of the box.
	Half v.Vec // Half-extents (half dimensions) from the center.
}

// Returns the left x-coordinate of the box.
func (a *AABB) Left() float64 { return a.Pos.X - a.Half.X }

// Returns the right x-coordinate of the box.
func (a *AABB) Right() float64 { return a.Pos.X + a.Half.X }

// Returns the top y-coordinate of the box (usually smaller y).
func (a *AABB) Top() float64 { return a.Pos.Y - a.Half.Y }

// Returns the bottom y-coordinate of the box (usually larger y).
func (a *AABB) Bottom() float64 { return a.Pos.Y + a.Half.Y }

// Sets the left x-coordinate and updates Pos.X accordingly.
func (a *AABB) SetLeft(l float64) { a.Pos.X = l + a.Half.X }

// Sets the right x-coordinate and updates Pos.X accordingly.
func (a *AABB) SetRight(r float64) { a.Pos.X = r - a.Half.X }

// Sets the top y-coordinate and updates Pos.Y accordingly.
func (a *AABB) SetTop(t float64) { a.Pos.Y = t + a.Half.Y }

// Sets the bottom y-coordinate and updates Pos.Y accordingly.
func (a *AABB) SetBottom(b float64) { a.Pos.Y = b - a.Half.Y }

// Returns the total width (X dimension) of the box.
func (a *AABB) Width() float64 { return a.Half.X * 2 }

// Returns the total height (Y dimension) of the box.
func (a *AABB) Height() float64 { return a.Half.Y * 2 }

// Circle represents a circular bounding volume.
type Circle struct {
	Pos    v.Vec   // Center position of the circle.
	Radius float64 // The radius of the circle.
}
