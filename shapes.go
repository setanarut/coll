package coll

import "github.com/setanarut/v"

type AABB struct {
	Pos  v.Vec
	Half v.Vec
}

func (a *AABB) Left() float64       { return a.Pos.X - a.Half.X }
func (a *AABB) Right() float64      { return a.Pos.X + a.Half.X }
func (a *AABB) Top() float64        { return a.Pos.Y - a.Half.Y }
func (a *AABB) Bottom() float64     { return a.Pos.Y + a.Half.Y }
func (a *AABB) SetLeft(l float64)   { a.Pos.X = l + a.Half.X }
func (a *AABB) SetRight(r float64)  { a.Pos.X = r - a.Half.X }
func (a *AABB) SetTop(t float64)    { a.Pos.Y = t + a.Half.Y }
func (a *AABB) SetBottom(b float64) { a.Pos.Y = b - a.Half.Y }

type Circle struct {
	Pos    v.Vec
	Radius float64
}
