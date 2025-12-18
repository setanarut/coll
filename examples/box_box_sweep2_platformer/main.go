package main

import (
	"fmt"
	"log"
	"math"

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
	platform       = coll.NewAABB(400, 300, 100, 10)
	platVel        = v.Vec{0, 0}
	platformCenter = v.Vec{X: ScreenWidth / 2, Y: 300}
	platformRadius = 150.0
	platformAngle  = 0.0
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

	// 1. Calculate platform movement for this frame
	platformAngle += 0.02
	newPlatCenterX := platformCenter.X + math.Cos(platformAngle)*platformRadius
	newPlatCenterY := platformCenter.Y + math.Sin(platformAngle)*platformRadius
	newPlatPos := v.Vec{X: newPlatCenterX - platform.Half.X, Y: newPlatCenterY - platform.Half.Y}
	platVel = newPlatPos.Sub(platform.Pos)

	// 2. Apply gravity and check for ground collision
	boxVel.Y += Gravity

	hitInfoBoxB.Reset()
	relVel := boxVel.Sub(platVel)

	onGround := false
	if coll.BoxBoxSweep2(platform, box, platVel, boxVel, hitInfoBoxB) {
		if hitInfoBoxB.Normal.Y == -1 {
			onGround = true
			// When on ground, vertical velocity should be based on the platform's
			boxVel = platVel
		}

		moveRel := slide(relVel, hitInfoBoxB)
		totalMove := moveRel.Add(platVel)
		box.Pos = box.Pos.Add(totalMove)

		if hitInfoBoxB.Normal.Y == 1 { // Hit bottom of platform
			boxVel.Y = 0
		}
	} else {
		// In air
		box.Pos = box.Pos.Add(boxVel)
	}

	// 3. Handle player input
	axis := examples.Axis().Unit()
	if onGround {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// Jump relative to the platform's current velocity
			boxVel.Y += JumpForce
		}
		// Add player's input speed to the platform's velocity
		boxVel.X = platVel.X + axis.X*PlayerSpeed
	} else {
		// Air control
		boxVel.X = axis.X * PlayerSpeed
	}

	// 4. Update platform position
	platform.Pos = platform.Pos.Add(platVel)

	// 5. Final ground check (floor)
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
