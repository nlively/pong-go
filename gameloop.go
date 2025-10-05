package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
		g.Game.Initialize()
		g.State = GameStateWaitForBallLaunch
	} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// quit game
		os.Exit(0)
	}
}

func (g *GameLoop) Update_WaitForBallLaunch() {
	// handle keys to launch ball
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// position ball and set game state to running
		g.Game.CenterPaddle()
		g.Game.PutBallInMiddle()
		g.Game.SetRandomUpwardBallVector()
		g.State = GameStateRunning
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
	} else {
		g.Game.MovePaddle(DirectionNone)
	}

	// Progress the ball along its trajectory
	collisionType := g.Game.MoveBallAlongTrajectory()
	if collisionType == CollisionTypePaddle {
		// Player gets a point
		g.Game.Score++
		if g.Game.Score%10 == 0 {
			g.Game.Level++
		}
		if g.Game.Level >= 10 {
			g.State = GameStateOver
		}
	}

	// See if the ball has gone past the paddle
	if int(g.Game.BallPosition.Y+BallDiameter) >= g.Game.GridHeight {
		g.Game.Lives--
		if g.Game.Lives == 0 {
			g.State = GameStateOver
		} else {
			g.State = GameStateWaitForBallLaunch
		}
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

func (g *GameLoop) drawAngle(image *ebiten.Image) {
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("Ball Angle: %0.2f", g.Game.BallAngle), 300, 0)
}

func (g *GameLoop) drawGameStats(image *ebiten.Image) {
	g.drawScore(image)
	g.drawLives(image)
	g.drawLevel(image)
	g.drawAngle(image)
}

func (g *GameLoop) drawBall(image *ebiten.Image) {
	radius := BallDiameter / 2.0
	cx := float32(g.Game.BallPosition.X) + float32(radius)
	cy := float32(g.Game.BallPosition.Y) + float32(radius)
	vector.DrawFilledCircle(image, cx, cy, BallDiameter/2, color.RGBA{0, 255, 255, 255}, true)
}

func (g *GameLoop) drawPaddle(image *ebiten.Image) {
	x := float32(g.Game.PaddleX)
	y := float32(g.Game.PaddleY)
	vector.DrawFilledRect(image, x, y, PaddleWidth, PaddleHeight, color.RGBA{0, 255, 255, 255}, true)
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

	x := 0
	y := g.Game.GridHeight / 2

	if g.Game.Lives == 0 {
		ebitenutil.DebugPrintAt(image, "Game over. You lost.", x, y)
	} else if g.Game.Level >= 10 {
		ebitenutil.DebugPrintAt(image, "Game over. You won!", x, y)
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
	return 800, 600
}
