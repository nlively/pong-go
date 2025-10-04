package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameState int

const (
	GameStateSplash            GameState = 0
	GameStateRunning           GameState = 1
	GameStateOver              GameState = 2
	GameStateWaitForBallLaunch GameState = 3
)

type GameLoop struct {
	Game  *Game
	State GameState // current state of game
}

func (g *GameLoop) Update_Splash() {
	// handle keys to start game or quit
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// start game
	} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// quit game
		os.Exit(0)
	}
}

func (g *GameLoop) Update_WaitForBallLaunch() {
	// handle keys to launch ball
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// position ball and set game state to running
		g.Game.PutBallInMiddle()
		g.Game.SetRandomUpwardBallVector()
	}
}

func (g *GameLoop) Update_GameOver() {
	// handle keys to restart game or quit
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.Game.Initialize()
		g.State = GameStateWaitForBallLaunch
	} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// quit game
		os.Exit(0)
	}
}

func (g *GameLoop) Update_GameRunning() {
	// Respond to arrow keys and move the paddle
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Game.MovePaddle(DirectionLeft)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Game.MovePaddle(DirectionRight)
	}

	// Progress the ball along its trajectory
	collisionType := g.Game.MoveBallAlongTrajectory()
	if collisionType == CollisionTypePaddle {
		// Player gets a point
		g.Game.Score++
	}

	// See if the ball has gone past the paddle
	if g.Game.BallPosition.Y > g.Game.PaddleY {
		g.Game.Lives--
		g.State = GameStateWaitForBallLaunch
	}

	// Resolve game state
	if g.Game.Lives == 0 {
		g.State = GameStateOver
	}

	if g.Game.Score > 0 && g.Game.Score%50 == 0 {
		g.Game.Level++
	}
	if g.Game.Level >= 10 {
		g.State = GameStateOver
	}
}

// Progresses the game once a tick
func (g *GameLoop) Update() error {
	switch g.State {
	case GameStateSplash:
		g.Update_Splash()
	case GameStateWaitForBallLaunch:
		g.Update_WaitForBallLaunch()
	case GameStateRunning:
		g.Update_GameRunning()
	case GameStateOver:
		g.Update_GameOver()
	default:
		return fmt.Errorf("unrecognized game state")
	}

	return nil
}

func (g *GameLoop) drawScore(image *ebiten.Image) {
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("Score: %d", g.Game.Score), 0, 0)
}

func (g *GameLoop) drawLives(image *ebiten.Image) {
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("Lives: %d", g.Game.Lives), 100, 0)
}

func (g *GameLoop) drawLevel(image *ebiten.Image) {
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("Level: %d", g.Game.Level), 200, 0)
}

func (g *GameLoop) drawGameStats(image *ebiten.Image) {
	g.drawScore(image)
	g.drawLives(image)
	g.drawLevel(image)
}

func (g *GameLoop) drawBall(image *ebiten.Image) {
	// TODO: draw the ball
}

func (g *GameLoop) drawPaddle(image *ebiten.Image) {
	// TODO: draw the paddle
}

func (g *GameLoop) Draw_Splash(image *ebiten.Image) {
	ebitenutil.DebugPrintAt(image, "Press Enter to play", 0, g.Game.GridHeight/2)
}

func (g *GameLoop) Draw_WaitForBallLaunch(image *ebiten.Image) {
	ebitenutil.DebugPrintAt(image, "Press Space to toss the ball", 0, g.Game.GridHeight/2)

	g.drawGameStats(image)
}

func (g *GameLoop) Draw_GameOver(image *ebiten.Image) {
	g.drawGameStats(image)

	if g.Game.Lives == 0 {
		ebitenutil.DebugPrint(image, "Game over. You lost.")
	} else if g.Game.Level >= 10 {
		ebitenutil.DebugPrint(image, "Game over. You won!")
	}

}

func (g *GameLoop) Draw_GameRunning(image *ebiten.Image) {
	g.drawBall(image)
	g.drawPaddle(image)
	g.drawGameStats(image)
}

// Renders the game once a tick
func (g *GameLoop) Draw(image *ebiten.Image) {
	switch g.State {
	case GameStateSplash:
		g.Draw_Splash(image)
	case GameStateWaitForBallLaunch:
		g.Draw_WaitForBallLaunch(image)
	case GameStateRunning:
		g.Draw_GameRunning(image)
	case GameStateOver:
		g.Draw_GameOver(image)
	}
}

func (g *GameLoop) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
