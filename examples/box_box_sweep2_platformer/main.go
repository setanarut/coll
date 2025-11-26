package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/coll"
	"github.com/setanarut/coll/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

const (
	ScreenWidth          = 854
	ScreenHeight         = 480
	Gravity      float64 = 1.6
	JumpForce    float64 = -30
	PlayerSpeed  float64 = 5
	PlatSpeed    float64 = 1.5
)

var (
	box         = coll.NewAABB(425, 250, 16, 35)
	boxVel      = v.Vec{0, 0}
	hitInfoBoxB = &coll.HitInfo{}
)
var (
	platform = coll.NewAABB(400, 300, 100, 10)
	platVel  = v.Vec{0, -PlatSpeed}
)

type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		switch ebiten.TPS() {
		case 60:
			ebiten.SetTPS(4)
		case 4:
			ebiten.SetTPS(60)
		}
	}
	axis := examples.Axis().Unit()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && boxVel.Y == 0 {
		boxVel.Y = JumpForce
	}
	boxVel.X = axis.X * PlayerSpeed
	boxVel.Y += Gravity
	hitInfoBoxB.Reset()
	relVel := boxVel.Sub(platVel)
	if coll.BoxBoxSweep2(platform, box, platVel, boxVel, hitInfoBoxB) {
		moveRel := slide(relVel, hitInfoBoxB)
		totalMove := moveRel.Add(platVel)
		box.Pos = box.Pos.Add(totalMove)
		if hitInfoBoxB.Normal.Y == -1 {
			boxVel.Y = 0
		}
		if hitInfoBoxB.Normal.Y == 1 {
			boxVel.Y = 0
		}
	} else {
		box.Pos = box.Pos.Add(boxVel)
	}
	platform.Pos = platform.Pos.Add(platVel)
	if platform.Top() < 0 || platform.Bottom() > ScreenHeight {
		platVel.Y *= -1
	}
	if box.Bottom() > ScreenHeight {
		box.SetBottom(ScreenHeight)
		boxVel.Y = 0
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	examples.StrokeBox(screen, box, colornames.Green)
	examples.StrokeBox(screen, platform, colornames.Gray)
	ebitenutil.DebugPrintAt(screen, "Controls: WASD / Space", 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Player Velocity: %.2f", boxVel.Y), 10, 30)
	examples.PrintHitInfoAt(screen, hitInfoBoxB, 10, 100)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Relative Velocity Platformer Example")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(fmt.Errorf("error running game: %w", err))
	}
}
func slide(vel v.Vec, hitInfo *coll.HitInfo) (slideVel v.Vec) {
	movementToHit := vel.Scale(hitInfo.Time)
	remainingVel := vel.Sub(movementToHit)
	dot := remainingVel.Dot(hitInfo.Normal)
	if dot > 0 {
		return movementToHit.Add(remainingVel)
	}
	slideDirection := remainingVel.Sub(hitInfo.Normal.Scale(dot))
	return movementToHit.Add(slideDirection)
}
