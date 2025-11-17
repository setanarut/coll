package examples

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

const WindowWidth, WindowHeight int = 500, 500

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

func DrawHitInfo(dst *ebiten.Image, hit *coll.HitInfo) {
	px, py := float32(hit.Pos.X), float32(hit.Pos.Y)
	nx, ny := px+(float32(hit.Normal.X)*8), py+(float32(hit.Normal.Y)*8)
	vector.FillCircle(dst, px, py, 2, colornames.Yellow, true)
	vector.StrokeLine(dst, px, py, nx, ny, 1, colornames.Yellow, false)
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
