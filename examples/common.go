package examples

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/v"
)

const WindowWidth, WindowHeight int = 500, 500

func StrokeCircleAt(dst *ebiten.Image, pos v.Vec, r float64, clr color.Color) {
	vector.StrokeCircle(dst, float32(pos.X), float32(pos.Y), float32(r), 2, clr, true)
}
func StrokeCircle(dst *ebiten.Image, c *coll.Circle, clr color.Color) {
	vector.StrokeCircle(dst, float32(c.Pos.X), float32(c.Pos.Y), float32(c.Radius), 2, clr, true)
}

func FillCircle(dst *ebiten.Image, c *coll.Circle, clr color.Color) {
	vector.FillCircle(dst, float32(c.Pos.X), float32(c.Pos.Y), float32(c.Radius), clr, true)
}
func FillCircleAt(dst *ebiten.Image, origin v.Vec, radius float64, clr color.Color) {
	vector.FillCircle(dst, float32(origin.X), float32(origin.Y), float32(radius), clr, true)
}

func StrokeAABB(dst *ebiten.Image, box *coll.AABB, clr color.Color) {
	vector.StrokeRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), 1, clr, false)
}

func StrokeAABBAt(dst *ebiten.Image, pos, half v.Vec, clr color.Color) {
	vector.StrokeRect(
		dst,
		float32(pos.X-half.X),
		float32(pos.Y-half.Y),
		float32(half.X*2),
		float32(half.Y*2),
		1,
		clr,
		false,
	)
}
func FillAABB(dst *ebiten.Image, box *coll.AABB, clr color.Color) {
	vector.FillRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), clr, false)
}

func DrawHitNormal(dst *ebiten.Image, hit *coll.HitInfo, clr color.Color, arrow bool) {
	DrawRay(dst, hit.Pos, hit.Normal, 12, clr, arrow)
}

func CursorPos() v.Vec {
	curX, curY := ebiten.CursorPosition()
	return v.Vec{float64(curX), float64(curY)}
}

func PrintHitInfoAt(dst *ebiten.Image, hit *coll.HitInfo, x, y int) {
	ebitenutil.DebugPrintAt(
		dst,
		fmt.Sprintf("Pos: %v\nDelta: %v\nNormal: %v\nTime: %v", hit.Pos, hit.Delta, hit.Normal, hit.Time),
		x,
		y,
	)
}

func Axis() (axis v.Vec) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axis.X += 1
	}
	return
}

func DrawRay(s *ebiten.Image, pos, dir v.Vec, length float64, clr color.Color, arrow bool) {
	end := pos.Add(dir.Unit().Scale(length))
	DrawLine(s, pos, end, clr)

	if arrow {
		arrowLen := 6.0
		arrowAngle := math.Pi / 7

		unitDir := dir.Unit()

		left := unitDir.Rotate(math.Pi - arrowAngle).Scale(arrowLen)
		right := unitDir.Rotate(-(math.Pi - arrowAngle)).Scale(arrowLen)

		DrawLine(s, end, end.Add(left), clr)
		DrawLine(s, end, end.Add(right), clr)
	}
}

func DrawLine(s *ebiten.Image, start, end v.Vec, clr color.Color) {
	vector.StrokeLine(s, float32(start.X), float32(start.Y), float32(end.X), float32(end.Y), 1.5, clr, true)
}
func DrawSegment(s *ebiten.Image, seg *coll.Segment, clr color.Color) {
	vector.StrokeLine(s, float32(seg.A.X), float32(seg.A.Y), float32(seg.B.X), float32(seg.B.Y), 1.5, clr, true)
}
