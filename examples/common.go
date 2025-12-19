package examples

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/v"
)

var im = ebiten.NewImage(1, 1)
var dio = &colorm.DrawImageOptions{}
var clrm = colorm.ColorM{}

func init() {
	im.Fill(color.White)
}

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

func StrokeBox(dst *ebiten.Image, box *coll.AABB, clr color.Color) {
	vector.StrokeRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), 1, clr, false)
}
func FillOBB(screen *ebiten.Image, obox *coll.OBB, c color.Color) {
	clrm.Reset()
	clrm.ScaleWithColor(c)
	dio.GeoM.Reset()
	dio.GeoM.Scale(obox.Half.X*2, obox.Half.Y*2)
	dio.GeoM.Translate(-obox.Half.X, -obox.Half.Y)
	dio.GeoM.Rotate(obox.Angle)
	dio.GeoM.Translate(obox.Pos.X, obox.Pos.Y)
	colorm.DrawImage(screen, im, clrm, dio)
}
func StrokeOBB(screen *ebiten.Image, box *coll.OBB, clr color.Color) {

	bAx := v.FromAngle(box.Angle)
	bAy := v.Vec{X: -bAx.Y, Y: bAx.X}

	axisX := bAx.Scale(box.Half.X)
	axisY := bAy.Scale(box.Half.Y)

	a := box.Pos.Sub(axisX).Sub(axisY) // left-bottom
	b := box.Pos.Add(axisX).Sub(axisY) // right-bottom
	c := box.Pos.Add(axisX).Add(axisY) // right-top
	d := box.Pos.Sub(axisX).Add(axisY) // left-top

	DrawLine(screen, a, b, clr)
	DrawLine(screen, b, c, clr)
	DrawLine(screen, c, d, clr)
	DrawLine(screen, d, a, clr)
}

func StrokeBoxAt(dst *ebiten.Image, pos, half v.Vec, clr color.Color) {
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
func FillBox(dst *ebiten.Image, box *coll.AABB, clr color.Color) {
	vector.FillRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), clr, false)
}

// func DrawHitNormal(dst *ebiten.Image, hit *coll.Hit, pos, vel v.Vec, clr color.Color, arrow bool) {
// 	hitpos := pos.Add(vel.Scale(hit.Time))
// 	DrawRay(dst, hitpos, hit.Normal, 12, clr, arrow)
// }

func CursorPos() v.Vec {
	curX, curY := ebiten.CursorPosition()
	return v.Vec{float64(curX), float64(curY)}
}

func PrintHitInfoAt(dst *ebiten.Image, hit *coll.Hit, x, y int) {
	ebitenutil.DebugPrintAt(
		dst,
		fmt.Sprintf("Normal: %v\nTime: %v", hit.Normal, hit.Time),
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
		arrowLen := 12.0
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
