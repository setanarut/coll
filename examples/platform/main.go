package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
)

type AABB = coll.AABB
type Vec = coll.Vec

const (
	screenWidth  = 800
	screenHeight = 600
	gravity      = 1.0
	friction     = 0.9
)

var (
	rect1 = &AABB{
		Pos:  Vec{100, 100},
		Half: Vec{40, 40},
	}
	plat = &AABB{
		Pos:  Vec{400, 250},
		Half: Vec{20, 100},
	}

	rect1Vel = Vec{0, 0}
	platVel  = Vec{1, 0}
	hit      = &coll.HitInfo2{}
	jumping  bool
)

type Game struct{}

func (g *Game) Update() error {
	// Get current mouse position.
	mx, _ := ebiten.CursorPosition()
	mouseX := float64(mx)

	// Jump control
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if !jumping {
			jumping = true
			rect1Vel.Y = -80 // Fixed jump velocity
		}
	}

	// Update horizontal velocity of the player based on mouse position.
	rect1Vel.X += (mouseX - rect1.Left() - (rect1.Half.X*2)*0.5) * 0.01

	// Apply gravity and friction
	rect1Vel.Y += gravity
	rect1Vel.X *= friction
	rect1Vel.Y *= friction

	// Update position
	rect1.Pos.X += rect1Vel.X
	rect1.Pos.Y += rect1Vel.Y

	plat.Pos.X += platVel.X
	plat.Pos.Y += platVel.Y

	// Collision check
	*hit = coll.HitInfo2{}
	if coll.AABBPlatform(rect1, plat, rect1Vel, platVel, hit) {
		// Adjust position on collision
		rect1.Pos.X += hit.Delta.X
		rect1.Pos.Y += hit.Delta.Y

		if hit.Top {
			jumping = false
			rect1Vel.Y = 0 // Reset vertical velocity
			fmt.Println("Top")
		}
		if hit.Bottom {
			fmt.Println("Bottom")
		}
		if hit.Right {
			fmt.Println("Right")
		}
		if hit.Left {
			fmt.Println("Left")
		}
	}

	// Ground collision check
	if rect1.Bottom() > screenHeight {
		rect1.SetBottom(screenHeight)
		rect1Vel.Y = 0
		jumping = false
	}

	// Platform boundary check
	if plat.Pos.X > screenWidth {
		plat.Pos.X = 0
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	DrawRect(screen, rect1.Left(), rect1.Top(), rect1.Half.X*2, rect1.Half.Y*2) //player
	DrawRect(screen, plat.Left(), plat.Top(), plat.Half.X*2, plat.Half.Y*2)     // platform

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

func DrawRect(dst *ebiten.Image, x, y, width, height float64) {
	vector.StrokeRect(dst, float32(x), float32(y), float32(width), float32(height), 1, color.Gray{200}, false)
}
