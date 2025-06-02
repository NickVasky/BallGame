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
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*
var assets embed.FS

const (
	ScreenWidth     = 1280
	ScreenHeight    = 720
	GlobalScale     = 0.5
	WorldWidth      = ScreenWidth * GlobalScale
	WorldHeight     = ScreenHeight * GlobalScale
	FloorRatio      = 0.2
	FloorHeight     = FloorRatio * WorldHeight
	FloorLevel      = WorldHeight - FloorHeight
	PhysicsDelta    = time.Second / 60
	Gravity         = 9.8
	ForceScale      = 100
	AudioSampleRate = 44100
)

type Game struct {
	audioCtx             *audio.Context
	audioPlayer          *audio.Player
	ballSprite           *ebiten.Image
	playerControlsVector Vec2
	lastUpdateTime       time.Time
	accumulator          time.Duration
	keys                 []ebiten.Key
	ball                 Ball
}

func NewGame() *Game {
	b := Ball{
		pos:    Vec2{WorldWidth / 2, WorldHeight / 2},
		accl:   Vec2{0, Gravity * ForceScale},
		radius: 16,
		mass:   10,
		decay:  0.3,
	}

	bs := LoadKolobok() //drawCircle(b.radius, color.RGBA{R: 210, G: 100, B: 30, A: 255})
	ctx := audio.NewContext(AudioSampleRate)
	sound := LoadKolobokSound()

	return &Game{
		audioCtx:    ctx,
		audioPlayer: ctx.NewPlayerFromBytes(sound),
		ballSprite:  bs,
		ball:        b,
	}
}

func LoadKolobok() *ebiten.Image {
	data, err := assets.ReadFile("assets/kolobok.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _ := png.Decode(bytes.NewReader(data))
	return ebiten.NewImageFromImage(img)
}

func LoadKolobokSound() []byte {
	data, err := assets.ReadFile("assets/kolobok_jump.wav")
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (g *Game) GetControlVector(keys []ebiten.Key, speed float64) Vec2 {
	v := Vec2{}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		v.x--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		v.x++
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.checkCollision() {
		v.y = -(0.5 * Gravity * ForceScale)
		if !g.audioPlayer.IsPlaying() {
			g.audioPlayer.Rewind() // go back to start
			g.audioPlayer.Play()
		}
	}

	v.x *= speed
	return v
}

func (g *Game) Update() error {
	now := time.Now()
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if g.lastUpdateTime.IsZero() {
		g.lastUpdateTime = now
		return nil
	}

	delta := now.Sub(g.lastUpdateTime)
	g.lastUpdateTime = now
	g.accumulator += delta

	g.playerControlsVector = g.GetControlVector(g.keys, 150)
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.ball.pos = g.ball.pos.MultByScalar(0)
		g.ball.vel = g.ball.vel.MultByScalar(0)
	}

	for g.accumulator >= PhysicsDelta {
		g.physicsStep(PhysicsDelta)
		g.accumulator -= PhysicsDelta
	}

	return nil
}

func (g *Game) physicsStep(dt time.Duration) {
	if g.checkCollision() {
		g.ball.vel.y *= -g.ball.decay
	}
	g.ball.vel.Add(g.ball.accl.MultByScalar(dt.Seconds()))
	g.ball.vel.Add(g.playerControlsVector)
	g.ball.vel.x *= 0.8
	g.ball.pos.Add(g.ball.vel.MultByScalar(dt.Seconds()))

}

func (g *Game) checkCollision() bool {
	// check floor collision
	if (g.ball.pos.y + float64(g.ball.radius)) >= FloorLevel {
		g.ball.pos.y = float64(FloorLevel - g.ball.radius)
		return true
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	// bg
	bgImg := ebiten.NewImage(WorldWidth, WorldHeight)
	bgImg.Fill(color.RGBA{R: 90, G: 165, B: 200, A: 255})
	screen.DrawImage(bgImg, op)

	// ground
	groundImg := ebiten.NewImage(WorldWidth, WorldHeight)
	groundImg.Fill(color.RGBA{R: 65, G: 45, B: 25, A: 255})

	op.GeoM.Translate(0, FloorLevel)

	screen.DrawImage(groundImg, op)

	op.GeoM.Reset()
	op.GeoM.Translate(-float64(g.ballSprite.Bounds().Dx())/2, -float64(g.ballSprite.Bounds().Dy())/2)
	op.GeoM.Rotate(g.ball.pos.x / float64(g.ball.radius))
	op.GeoM.Translate(g.ball.pos.x, g.ball.pos.y)
	screen.DrawImage(g.ballSprite, op)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Press [LEFT/RIGHT] to move,[SPACE] to jump!\n Ball Coords: %s\n Key pressed: %v",
			g.ball.pos,
			g.keys,
		), 5, 5)
}

// func drawCircle(radius int, clr color.Color) *ebiten.Image {
// 	diameter := radius * 2
// 	img := ebiten.NewImage(diameter, diameter)

// 	// r, g, b, a := clr.RGBA()
// 	// fill := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
// 	// stroke := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a / 2)}

// 	for y := 0; y < diameter; y++ {
// 		for x := 0; x < diameter; x++ {
// 			dx := float64(x - radius)
// 			dy := float64(y - radius)
// 			if float64(radius*radius) >= (dx*dx + dy*dy) {
// 				img.Set(x, y, clr)
// 			}
// 		}
// 	}

// 	return img
// }

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WorldWidth, WorldHeight
}
