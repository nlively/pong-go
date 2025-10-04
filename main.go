package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := NewGame(800, 600)
	gameLoop := &GameLoop{
		Game:  game,
		State: GameStateSplash,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(gameLoop); err != nil {
		log.Fatal(err)
	}
}
