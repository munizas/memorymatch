package main

import (
	"log"

	memorymatch "github.com/munizas/memorymatch/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game, err := memorymatch.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(memorymatch.ScreenWidth, memorymatch.ScreenHeight)
	ebiten.SetWindowTitle("Memory Match!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
