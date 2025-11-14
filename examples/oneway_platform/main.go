package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/v"
)

type AABB = coll.AABB

const (
	screenWidth  = 854
	screenHeight = 480
	gravity      = 1.0
	friction     = 0.9
)

var (
	box = &AABB{
		Pos:  v.Vec{100, 100},
		Half: v.Vec{16, 35},
	}
	platform = &AABB{
		Pos:  v.Vec{400, 250},
		Half: v.Vec{100, 10},
	}

	boxVelocity      = v.Vec{0, 0}
	platformVelocity = v.Vec{3, 0}
	hit              = &coll.HitInfo2{}
)

type Game struct{}

func (g *Game) Update() error {
	axis := Axis()

	// Jump control
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		boxVelocity.Y = -60 // Fixed jump velocity
	}

	boxVelocity.X = axis.X * 10

	// Apply gravity and friction
	boxVelocity.Y += gravity
	boxVelocity.X *= friction
	boxVelocity.Y *= friction

	// Update position
	box.Pos.X += boxVelocity.X
	box.Pos.Y += boxVelocity.Y

	platform.Pos.X += platformVelocity.X
	platform.Pos.Y += platformVelocity.Y

	// Collision check
	*hit = coll.HitInfo2{}
	if coll.AABBPlatform(box, platform, boxVelocity, platformVelocity, hit) {

		if !hit.Top {
			// Adjust position on collision
			box.Pos.X += hit.Delta.X
			box.Pos.Y += hit.Delta.Y
		}

		if hit.Top {
			// boxVelocity.Y = 0 // Reset vertical velocity
		}
		if hit.Bottom {
		}
		if hit.Right {
		}
		if hit.Left {
		}
	}

	// Ground collision check
	if box.Bottom() > screenHeight {
		box.SetBottom(screenHeight)
		boxVelocity.Y = 0
	}

	// Platform boundary check
	if platform.Right() > screenWidth || platform.Left() < 0.0 {
		platformVelocity.X *= -1
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	DrawAABB(screen, box)      // player
	DrawAABB(screen, platform) // platform

}

// Layout returns the size of the game window.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hitbox (Ported to Go with Ebiten)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(fmt.Errorf("error running game: %w", err))
	}
}

func DrawAABB(dst *ebiten.Image, box *AABB) {
	vector.StrokeRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), 1, color.Gray{200}, false)
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
