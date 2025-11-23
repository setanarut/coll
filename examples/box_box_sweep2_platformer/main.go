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

// Moving platformer example

const (
	ScreenWidth          = 854
	ScreenHeight         = 480
	Gravity      float64 = 0.3
	JumpForce    float64 = -13
	PlayerSpeed  float64 = 4
)

var (
	box         = coll.NewAABB(425, 250, 16, 35)
	boxVelocity = v.Vec{0, 0}
	hitInfoBoxB = &coll.HitInfo{}
)

var (
	platform         = coll.NewAABB(400, 300, 100, 10)
	platformVelocity = v.Vec{3, 0}
)

type Game struct{}

func (g *Game) Update() error {
	axis := examples.Axis()

	// Jump control
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		boxVelocity.Y = JumpForce
	}

	boxVelocity.X = axis.X * PlayerSpeed // WASD Movement
	boxVelocity.Y += Gravity
	platform.Pos = platform.Pos.Add(platformVelocity)

	// Collision check
	hitInfoBoxB.Reset()
	if coll.BoxBoxSweep2(platform, box, platformVelocity, boxVelocity, hitInfoBoxB) {
		coll.CollideAndSlide(box, boxVelocity, hitInfoBoxB)
		// vel := boxVelocity.Add(boxVelocity)
		// box.Pos = box.Pos.Add(vel)
		box.Pos = box.Pos.Add(platformVelocity)

		if hitInfoBoxB.Normal.Y != 0 {
			boxVelocity.Y = 0 // Reset vertical velocity
		}
	} else {
		box.Pos = box.Pos.Add(boxVelocity)
	}

	// Ground collision check
	if box.Bottom() > ScreenHeight {
		box.SetBottom(ScreenHeight)
		boxVelocity.Y = 0
	}

	// Platform boundary check
	if platform.Right() > ScreenWidth || platform.Left() < 0.0 {
		platformVelocity.X *= -1
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	examples.StrokeBox(screen, box, colornames.Green)     // player
	examples.StrokeBox(screen, platform, colornames.Gray) // platform
	ebitenutil.DebugPrintAt(screen, "Controls: WASD / Space)", 10, 10)
	examples.PrintHitInfoAt(screen, hitInfoBoxB, 10, 100)

}

// Layout returns the size of the game window.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(fmt.Errorf("error running game: %w", err))
	}
}
