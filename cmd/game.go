package main

import (
	"log"

	"github.com/NickVasky/BallGame/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}
