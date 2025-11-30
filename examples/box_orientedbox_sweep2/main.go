package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

var (
	orbitalAngle, radius = 0., 60.
	origin               = v.Vec{250, 250}
	collided             bool

	bullet = &coll.OBB{
		Pos:   v.Vec{250, 250},
		Half:  v.Vec{2, 2},
		Angle: 0,
	}
	player = coll.NewAABB(250, 250, 200, 2)

	bulletOldPos = bullet.Pos
	playerOldPos = player.Pos
)

func main() {
	ebiten.SetScreenClearedEveryFrame(false)
	g := &Game{}
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		switch ebiten.TPS() {
		case 60:
			ebiten.SetTPS(3)
		case 3:
			ebiten.SetTPS(60)

		}
	}

	bullet.Pos.X = origin.X + radius*math.Cos(orbitalAngle)
	bullet.Pos.Y = origin.Y + radius*math.Sin(orbitalAngle)
	bullet.Angle = orbitalAngle + math.Pi/2

	player.Pos.X = origin.X + radius*math.Sin(orbitalAngle)*3
	player.Pos.Y = origin.X + radius*math.Cos(orbitalAngle)*3

	velBullet := bullet.Pos.Sub(bulletOldPos)
	velPlayer := player.Pos.Sub(playerOldPos)

	collided = coll.BoxOrientedBoxSweep2(player, bullet, velPlayer, velBullet)

	orbitalAngle += 0.2

	bulletOldPos = bullet.Pos
	playerOldPos = player.Pos
	return nil
}
func (g *Game) Draw(s *ebiten.Image) {
	s.Fill(color.Gray{32})
	examples.StrokeBox(s, player, colornames.Green)
	if collided {
		examples.StrokeOBB(s, bullet, colornames.Yellow)
	} else {
		examples.StrokeOBB(s, bullet, colornames.Green)
	}
}
func (g *Game) Layout(w, h int) (int, int) {
	return 500, 500
}
