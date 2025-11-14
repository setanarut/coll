package main

import (
	"fmt"
	"log"

	"github.com/setanarut/coll"
	"github.com/setanarut/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

var platformBox = &coll.AABB{
	Pos:  v.Vec{100, 250},
	Half: v.Vec{60, 5},
}
var bulletCircle = &coll.Circle{
	Pos:    v.Vec{100, 100},
	Radius: 20,
}

var bulletColor = colornames.Gray
var overlap bool
var platformVel = v.Vec{2, 0}
var bulletVel = v.Vec{2, 0}
var maxBulletVelXRecord = 0.0

func main() {

	g := &Game{ScreenWidth: 800, H: 500}
	ebiten.SetWindowSize(int(g.ScreenWidth), int(g.H))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	ScreenWidth, H float64
}

func (g *Game) Update() error {

	curX, curY := ebiten.CursorPosition()
	cursorPos := v.Vec{float64(curX), float64(curY)}

	bulletVel = cursorPos.Sub(bulletCircle.Pos)

	overlap = coll.AABBCircleSweep(platformBox, bulletCircle, platformVel, bulletVel)
	if overlap {
		if bulletVel.X > maxBulletVelXRecord {
			maxBulletVelXRecord = bulletVel.X
		}
		bulletColor = colornames.Yellow
	} else {
		bulletColor = colornames.Gray
	}

	platformBox.Pos = platformBox.Pos.Add(platformVel)

	if platformBox.Left() < 0 || platformBox.Right() > g.ScreenWidth {
		platformVel.X *= -1
	}

	bulletCircle.Pos = bulletCircle.Pos.Add(bulletVel)
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	vector.FillRect(
		s,
		float32(platformBox.Pos.X-platformBox.Half.X),
		float32(platformBox.Pos.Y-platformBox.Half.Y),
		float32(platformBox.Half.X*2),
		float32(platformBox.Half.Y*2),
		colornames.Gray,
		true,
	)
	// Bullet
	vector.FillCircle(
		s,
		float32(bulletCircle.Pos.X),
		float32(bulletCircle.Pos.Y),
		float32(bulletCircle.Radius),
		bulletColor,
		true,
	)

	ebitenutil.DebugPrint(s, fmt.Sprintf("tunnelling test max X speed: %v", maxBulletVelXRecord))
}

func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return g.ScreenWidth, g.H
}
