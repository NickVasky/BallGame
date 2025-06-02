package game

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*
var assets embed.FS

const (
	ScreenWidth     = 1280
	ScreenHeight    = 720
	GlobalScale     = 2
	WorldWidth      = ScreenWidth / GlobalScale
	WorldHeight     = ScreenHeight / GlobalScale
	FloorRatio      = 0.2
	FloorHeight     = int(FloorRatio * WorldHeight)
	FloorLevel      = int(WorldHeight - FloorHeight)
	PhysicsDelta    = time.Second / 60
	Gravity         = 9.8
	ForceScale      = 100
	AudioSampleRate = 44100
)

type Game struct {
	audio    Audio
	ball     Ball
	physics  Physics
	player   PlayerController
	drawOpts *ebiten.DrawImageOptions
	keys     []ebiten.Key
}

type Physics struct {
	lastUpdateTime time.Time
	accumulator    time.Duration
}

func NewPhysics() Physics {
	return Physics{
		lastUpdateTime: time.Now(),
	}
}

func NewGame() *Game {
	jumpSound := LoadSound("assets/kolobok_jump.wav")
	playerSprite := LoadImage("assets/kolobok.png")

	ball := NewBall(playerSprite)
	audio := NewAudio()
	physics := NewPhysics()
	PlayerController := NewPlayerController(150, 0.3)

	audio.jumpPlayer = audio.audioCtx.NewPlayerFromBytes(jumpSound)

	return &Game{
		audio:    audio,
		ball:     ball,
		physics:  physics,
		player:   PlayerController,
		drawOpts: &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
	}
}

func LoadImage(path string) *ebiten.Image {
	data, err := assets.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	img, _ := png.Decode(bytes.NewReader(data))
	return ebiten.NewImageFromImage(img)
}

func LoadSound(path string) []byte {
	data, err := assets.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (g *Game) ProcessKeys(keys []ebiten.Key, speed float64) {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.ball.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.controlVector.x = -speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.controlVector.x = speed
	} else {
		g.player.controlVector.x = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.checkCollision() {
		g.player.controlVector.y = -(0.5 * Gravity * ForceScale)

		if !g.audio.jumpPlayer.IsPlaying() {
			g.audio.jumpPlayer.Rewind()
			g.audio.jumpPlayer.Play()
		}
	} else {
		g.player.controlVector.y = 0
	}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	now := time.Now()
	delta := now.Sub(g.physics.lastUpdateTime)
	g.physics.lastUpdateTime = now
	g.physics.accumulator += delta

	g.ProcessKeys(g.keys, 150)

	if g.ball.pos.x >= float64(WorldWidth+g.ball.radius)+0.1 {
		g.ball.pos.x = float64(-g.ball.radius)
	}
	if g.ball.pos.x <= -float64(g.ball.radius)-0.1 {
		g.ball.pos.x = float64(WorldWidth + g.ball.radius)
	}

	for g.physics.accumulator >= PhysicsDelta {
		g.physicsStep(PhysicsDelta)
		g.physics.accumulator -= PhysicsDelta
	}

	return nil
}

func (g *Game) physicsStep(dt time.Duration) {
	if g.checkCollision() {
		g.ball.pos.y = float64(FloorLevel - g.ball.radius)
		g.ball.vel.y *= -g.ball.decay
	}

	g.ball.accl = Vec2{0, Gravity * ForceScale}

	g.ball.vel.Add(g.ball.accl.MultByScalar(dt.Seconds()))
	g.ball.vel.Add(g.player.controlVector)
	g.ball.vel.x *= 0.8

	g.ball.pos.Add(g.ball.vel.MultByScalar(dt.Seconds()))
}

func (g *Game) checkCollision() bool {
	if int(g.ball.pos.y+float64(g.ball.radius)) >= FloorLevel {
		return true
	} else {
		return false
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBg(screen, g.drawOpts)
	g.drawGround(screen, g.drawOpts)
	g.drawBall(screen, g.drawOpts)
	g.drawDebugInfo(screen)
}

func (g *Game) drawBall(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(g.ball.sprite.Bounds().Dx())/2, -float64(g.ball.sprite.Bounds().Dy())/2)
	op.GeoM.Rotate(g.ball.pos.x / float64(g.ball.radius))
	op.GeoM.Translate(g.ball.pos.x, g.ball.pos.y)
	screen.DrawImage(g.ball.sprite, op)
}

func (g *Game) drawGround(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	groundImg := ebiten.NewImage(WorldWidth, WorldHeight)
	groundImg.Fill(color.RGBA{R: 65, G: 45, B: 25, A: 255})
	op.GeoM.Reset()
	op.GeoM.Translate(0, float64(FloorLevel))
	screen.DrawImage(groundImg, op)
}

func (g *Game) drawBg(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	bgImg := ebiten.NewImage(WorldWidth, WorldHeight)
	bgImg.Fill(color.RGBA{R: 90, G: 165, B: 200, A: 255})
	op.GeoM.Reset()
	screen.DrawImage(bgImg, op)
}

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Press [LEFT/RIGHT] to move,[SPACE] to jump! Key pressed: %v\n\tBall's Pos: %s, Vel: %s, Acc: %s\n\tControl vector: %s",
			g.keys,
			g.ball.pos,
			g.ball.vel,
			g.ball.accl,
			g.player.controlVector,
		), 5, 5)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WorldWidth, WorldHeight
}
