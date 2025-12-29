package main

import (
	"fmt"
	"image/color"
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
	ScreenWidth   = 854
	ScreenHeight  = 480
	MoveSpeedX    = 6.125
	JumpPower     = -21.46
	Gravity       = 0.86
	PlatformSpeed = 0.03
)

var (
	player        = coll.NewAABB(425, 250, 16, 36)
	playerDelta   = v.Vec{0, 0}
	playerHitInfo = &coll.Hit{}
)
var (
	platform       = coll.NewAABB(400, 300, 64, 16)
	platformDelta  = v.Vec{}
	platformCenter = v.Vec{X: ScreenWidth / 2, Y: 200}
	platformRadius = 150.0
	platformAngle  = 0.0
)

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
	updatePlatformDelta()
	playerDelta.Y += Gravity
	speed := examples.Axis().Unit().Scale(MoveSpeedX)
	playerHitInfo.Reset()
	hit := coll.BoxBoxSweep2(platform, player, platformDelta, playerDelta, playerHitInfo)
	onGround := false
	if hit && playerHitInfo.Normal.Y == -1 {
		playerPosAtHit := player.Pos.Add(playerDelta.Scale(playerHitInfo.Data))
		platformPosAtHit := platform.Pos.Add(platformDelta.Scale(playerHitInfo.Data))
		playerBottomAtHit := playerPosAtHit.Y + player.Half.Y
		platformTopAtHit := platformPosAtHit.Y - platform.Half.Y
		onGround = playerBottomAtHit <= platformTopAtHit
	}
	if onGround {
		player.Pos = player.Pos.Add(playerDelta.Scale(playerHitInfo.Data))
		player.Pos.Y = platform.Pos.Y + platformDelta.Y - player.Half.Y - platform.Half.Y
		player.Pos.X += platformDelta.X + speed.X
		playerDelta.X = platformDelta.X + speed.X
		playerDelta.Y = platformDelta.Y
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			playerDelta.Y = JumpPower
		}
	} else {
		playerDelta.X = speed.X
		player.Pos = player.Pos.Add(playerDelta)
	}
	platform.Pos = platform.Pos.Add(platformDelta)
	if player.Bottom() >= ScreenHeight {
		player.SetBottom(ScreenHeight)
		playerDelta.Y = 0
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			playerDelta.Y = JumpPower
		}
	}
	// inherit only the platform's fractional X offset to avoid sub-pixel jitter.
	if onGround {
		player.Pos.X = math.Floor(player.Pos.X) + Fract(platform.Pos.X)
	}
	return nil
}

func updatePlatformDelta() {
	platformAngle += PlatformSpeed
	newPlatCenterX := platformCenter.X + math.Cos(platformAngle)*platformRadius
	newPlatCenterY := platformCenter.Y + math.Sin(platformAngle)*platformRadius
	newPlatPos := v.Vec{X: newPlatCenterX, Y: newPlatCenterY}
	platformDelta = newPlatPos.Sub(platform.Pos)
}
func (g *Game) Draw(s *ebiten.Image) {
	s.Fill(color.Gray{20})
	examples.StrokeBox(s, player, colornames.Green)
	examples.StrokeBox(s, platform, colornames.Darkgray)
	examples.PrintHitInfoAt(s, playerHitInfo, 30, 30, false)
	ebitenutil.DebugPrintAt(s, "Space - Jump\nA/D - Move", 400, 100)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
func main() {
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Platformer Example")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(fmt.Errorf("error running game: %w", err))
	}

}

// Fract returns the fractional part of x.
func Fract(x float64) float64 {
	if x >= 0 {
		return x - math.Floor(x)
	}
	return x - math.Ceil(x)
}
