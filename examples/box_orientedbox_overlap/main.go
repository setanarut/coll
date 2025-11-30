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
	orbitalAngle, radius = 0., 42.5
	origin               = v.Vec{250, 250}
	collided             bool

	bullet = &coll.OBB{
		Pos:   v.Vec{250, 250},
		Half:  v.Vec{37, 8},
		Angle: 0,
	}
	player = coll.NewAABB(250, 250, 16, 33)
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
			ebiten.SetTPS(10)
		case 10:
			ebiten.SetTPS(60)

		}
	}

	bullet.Pos.X = origin.X + radius*math.Cos(orbitalAngle)
	bullet.Pos.Y = origin.Y + radius*math.Sin(orbitalAngle)
	bullet.Angle = orbitalAngle + math.Pi/2

	collided = coll.BoxOrientedBoxOverlap(player, bullet)

	orbitalAngle += 0.03
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
