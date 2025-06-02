package game

import "github.com/hajimehoshi/ebiten/v2"

type Ball struct {
	pos    Vec2
	vel    Vec2
	accl   Vec2
	radius int
	decay  float64
	sprite *ebiten.Image
}

func NewBall(sprite *ebiten.Image) Ball {
	return Ball{
		pos:    Vec2{WorldWidth / 2, WorldHeight / 2},
		radius: sprite.Bounds().Dx() / 2,
		decay:  0.3,
		sprite: sprite,
	}
}

func (b *Ball) Reset() {
	b.pos = Vec2{WorldWidth / 2, WorldHeight / 2}
	b.vel = b.vel.MultByScalar(0)
}
