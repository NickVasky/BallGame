package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	FloorRatio   = 0.2
	FloorHeight  = FloorRatio * ScreenHeight
	FloorLevel   = ScreenHeight - FloorHeight
	PhysicsDelta = time.Second / 60
)

type Game struct {
	ballSprite           *ebiten.Image
	playerControlsVector Vec2
	lastUpdateTime       time.Time
	accumulator          time.Duration
	keys                 []ebiten.Key
	ball                 Ball
}

func NewGame() *Game {
	b := Ball{
		pos:    Vec2{ScreenWidth / 2, ScreenHeight / 2},
		accl:   Vec2{0, 9.8 * 100},
		radius: 15,
		mass:   10,
	}
	bs := drawCircle(b.radius, color.RGBA{R: 128, G: 150, B: 32, A: 128})
	return &Game{
		ballSprite: bs,
		ball:       b,
	}
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
		v.y -= 600
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
	//now := time.Now()

	//g.msgPoint.Add(v)
	if g.checkCollision() {
		g.ball.vel.y *= -0.5
	}
	g.ball.vel.Add(g.ball.accl.MultByScalar(dt.Seconds()))
	g.ball.vel.Add(g.playerControlsVector)
	g.ball.vel.x *= 0.8
	g.ball.pos.Add(g.ball.vel.MultByScalar(dt.Seconds()))
	// if now.Before(g.body.AccelEndTime) {
	// 	g.body.Acceleration = g.body.AppliedAcceleration
	// } else {
	// 	g.body.Acceleration = 0
	// }

	// Apply acceleration to velocity, position, etc. here
}

func (g *Game) checkCollision() bool {
	// check floor collision
	if (g.ball.pos.y + float64(g.ball.radius)*2) >= FloorLevel {
		g.ball.pos.y = float64(FloorLevel - g.ball.radius*2)
		return true
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	groundImg := ebiten.NewImage(ScreenWidth, FloorHeight)
	groundImg.Fill(color.RGBA{R: 128, G: 32, B: 16, A: 23})

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(0, FloorLevel)
	screen.DrawImage(groundImg, op)

	op.GeoM.Reset()
	op.GeoM.Translate(g.ball.pos.x, g.ball.pos.y)
	screen.DrawImage(g.ballSprite, op)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Hello, World!\n Ball Coords: %s\n Ball sprite %v\n Key pressed: %v",
			g.ball.pos,
			g.ballSprite.Bounds().Max,
			g.keys,
		), 10, 10)
}

func drawCircle(radius int, clr color.Color) *ebiten.Image {
	diameter := radius * 2
	img := ebiten.NewImage(diameter, diameter)

	// r, g, b, a := clr.RGBA()
	// fill := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	// stroke := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a / 2)}

	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x - radius)
			dy := float64(y - radius)
			if float64(radius*radius) >= (dx*dx + dy*dy) {
				img.Set(x, y, clr)
			}
		}
	}

	return img
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
