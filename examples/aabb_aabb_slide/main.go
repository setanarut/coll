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
	Gravity      float64 = 1.0
	Damping      float64 = 0.9
)

var (
	box         = coll.NewAABB(100, 100, 16, 35)
	boxVelocity = v.Vec{0, 0}
	hitInfoBoxB = &coll.HitInfo2{}
)

var (
	platform         = coll.NewAABB(400, 250, 100, 10)
	platformVelocity = v.Vec{3, 0}
)

type Game struct{}

func (g *Game) Update() error {
	axis := examples.Axis()

	// Jump control
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		boxVelocity.Y = -60 // Fixed jump velocity
	}

	// WASDMovement
	boxVelocity.X = axis.X * 10

	// Apply gravity and friction
	boxVelocity.Y += Gravity
	boxVelocity.X *= Damping
	boxVelocity.Y *= Damping

	platform.Pos = platform.Pos.Add(platformVelocity)
	box.Pos = box.Pos.Add(boxVelocity)

	// Collision check
	hitInfoBoxB.Reset()
	if coll.AABBAABBSlide(platform, box, platformVelocity, boxVelocity, hitInfoBoxB) {
		box.Pos = box.Pos.Add(hitInfoBoxB.Delta)
		box.Pos = box.Pos.Add(platformVelocity)
		if hitInfoBoxB.Top {
			boxVelocity.Y = 0 // Reset vertical velocity
		}
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
	examples.StrokeAABB(screen, box, colornames.Gray)      // player
	examples.StrokeAABB(screen, platform, colornames.Gray) // platform
	examples.PrintHitInfoAt2(screen, hitInfoBoxB, 10, 10)
	ebitenutil.DebugPrintAt(screen, "Controls: WASD / Space)", 10, 30)

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
