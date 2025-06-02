package game

import (
	"fmt"
	"image/color"
	"math"

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
)

type Vector struct {
	x float64
	y float64
}

func (v *Vector) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector) Normalize() {
	vLen := v.Len()
	v.x /= vLen
	v.y /= vLen
}

func (v *Vector) Add(vec Vector) {
	v.x += vec.x
	v.y += vec.y
}

func (v *Vector) MultByScalar(mult float64) {
	v.x *= mult
	v.y *= mult
}

func (v Vector) String() string {
	return fmt.Sprintf("[x: %.2f, y: %.2f]", v.x, v.y)
}

type Game struct {
	keys     []ebiten.Key
	msgPoint Vector
}

func GetSpeedVector(keys []ebiten.Key, speed float64) Vector {
	v := Vector{}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		v.x--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		v.x++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		v.y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		v.y--
	}
	if v.Len() > 0 {
		v.Normalize()
	}
	v.MultByScalar(speed)
	return v
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	v := GetSpeedVector(g.keys, 3)
	g.msgPoint.Add(v)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	groundImg := ebiten.NewImage(ScreenWidth, FloorHeight)
	groundImg.Fill(color.RGBA{R: 128, G: 32, B: 16, A: 23})

	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(float64(g.msgPoint.x), float64(h*j))
	op.GeoM.Translate(g.msgPoint.x, FloorLevel)

	screen.DrawImage(groundImg, op)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hello, World!\nCoords: %s", g.msgPoint), int(g.msgPoint.x), int(g.msgPoint.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
