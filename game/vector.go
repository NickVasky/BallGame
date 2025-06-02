package game

import (
	"fmt"
	"math"
)

type Vec2 struct {
	x float64
	y float64
}

func (v *Vec2) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vec2) Normalize() {
	vLen := v.Len()
	v.x /= vLen
	v.y /= vLen
}

func (v *Vec2) Add(vec Vec2) {
	v.x += vec.x
	v.y += vec.y
}

func (v *Vec2) MultByScalar(mult float64) Vec2 {
	return Vec2{v.x * mult, v.y * mult}
}

func (v Vec2) String() string {
	return fmt.Sprintf("[x: %.2f, y: %.2f]", v.x, v.y)
}
