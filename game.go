package main

import (
	"fmt"
	"math/rand/v2"
)

const PaddleWidth = 70
const PaddleHeight = 15
const BallDiameter = 20
const PaddleRelativePosition = 0.8 // 80% of grid height

type Game struct {
	GridWidth            int    // grid width in pixels
	GridHeight           int    // grid height in pixels
	Lives                int    // how many lives the player has left
	BallPosition         Vector // X,Y of ball on grid
	ballAngleRad         float64
	BallAngle            float64
	BallSpeed            float64
	PaddleX              int    // current position of the paddle horizontally
	PaddleY              int    // fixed position of the paddle vertically
	paddleMovementVector Vector // when the paddle is moving rather than still, this tracks paddle's vector from the last frame
	Score                int    // Player's current score
	Level                int    // Current level
}

func NewGame(gridWidth, gridHeight int) *Game {
	return &Game{
		GridWidth:    gridWidth,
		GridHeight:   gridHeight,
		BallAngle:    45,
		BallSpeed:    2.3,
		ballAngleRad: deg2Rad(45),
		Lives:        3,
		Score:        0,
		Level:        1,
		PaddleY:      int(float64(gridHeight) * PaddleRelativePosition),
	}
}

func (g *Game) Initialize() {
	g.Lives = 3
	g.Score = 0
	g.Level = 1
	g.SetRandomUpwardBallVector()
}

func (g *Game) CenterPaddle() {
	g.PaddleX = (g.GridWidth / 2) - (PaddleWidth / 2)
}

func (g *Game) MovePaddle(direction Direction) {
	moveAmount := 5 // number of pixels to move

	switch direction {
	case DirectionLeft:
		g.PaddleX -= moveAmount
		if g.PaddleX < 0 {
			g.PaddleX = 0
		}
		g.paddleMovementVector = Vector{-float64(moveAmount), 0}
	case DirectionRight:
		g.PaddleX += moveAmount
		if g.PaddleX+PaddleWidth > g.GridWidth {
			g.PaddleX = g.GridWidth - PaddleWidth
		}
		g.paddleMovementVector = Vector{float64(moveAmount), 0}
	case DirectionNone:
		g.paddleMovementVector = Vector{0, 0}
	default:
		panic("unrecognized direction")
	}
}

func (g *Game) PutBallInMiddle() {
	x := (g.GridWidth / 2) + (BallDiameter / 2)
	y := (g.PaddleY / 3) + (PaddleHeight / 2)
	g.BallPosition = Vector{float64(x), float64(y)}
}

func (g *Game) BounceBall(surfaceVector Vector, movementVector Vector) {
	v := (&Polar{Angle: g.BallAngle, Speed: g.BallSpeed}).ToVector()
	v = v.ReflectWithTangentialImpulse(surfaceVector, movementVector, 0.5)
	p := v.ToPolar()
	g.BallAngle = p.Angle
}

func (g *Game) SetRandomUpwardBallVector() {
	g.BallAngle = float64(rand.IntN(45) + 30)
	g.ballAngleRad = deg2Rad(g.BallAngle)
}

func (g *Game) MoveBallAlongTrajectory() CollisionType {
	collisionType := CollisionTypeNone
	horizontalSurfaceVector := Vector{1, 0}
	verticalSurfaceVector := Vector{0, 1}

	speed := g.BallSpeed              // assume the value = the number of pixels moved per tick.  we prob need to adjust this
	speed += float64(g.Level-1) * 0.3 // make the ball move faster with every level. again, we prob need to adjust this

	v := (&Polar{Angle: g.BallAngle, Speed: speed}).ToVector()

	// do the math to calculate the new position based on the ball's vector
	delta := v.Multiply(speed)

	newVector := g.BallPosition.Add(delta)

	// correct based on collisions
	if newVector.X < 0 {
		newVector.X = 0
		collisionType = CollisionTypeWall
		fmt.Println("Detected collision with left wall!")
		g.BounceBall(verticalSurfaceVector, Vector{})
	} else if int(newVector.X)+BallDiameter > g.GridWidth {
		newVector.X = float64(g.GridWidth - BallDiameter)
		fmt.Println("Detected collision with right wall!")
		collisionType = CollisionTypeWall
		g.BounceBall(verticalSurfaceVector, Vector{})
	}

	if newVector.Y < 0 {
		newVector.Y = 0
		collisionType = CollisionTypeWall
		fmt.Println("Detected collision with top wall!")
		g.BounceBall(horizontalSurfaceVector, Vector{})
	} else if int(newVector.Y)+BallDiameter > g.GridWidth {
		newVector.Y = float64(g.GridHeight - BallDiameter)
	}

	// There's a collision with the paddle if:
	// the right edge of the ball's bounding box is at or past the left edge of the paddle and
	// the left edge of the ball's bounding box is at or before the right edge of the paddle and
	// the bottom edge of the ball is at or lower than the paddle's top edge

	ballLeftEdge := newVector.X
	ballRightEdge := newVector.X + BallDiameter
	ballBottom := newVector.Y + BallDiameter

	paddleLeftEdge := float64(g.PaddleX)
	paddleRightEdge := paddleLeftEdge + PaddleWidth
	paddleTop := float64(g.PaddleY)
	paddleBottom := float64(g.PaddleY) + PaddleHeight

	if ballRightEdge >= paddleLeftEdge && ballLeftEdge <= paddleRightEdge && ballBottom >= paddleTop && ballBottom <= paddleBottom {
		newVector.Y = float64(g.PaddleY - BallDiameter) // normalize the ball's Y position
		collisionType = CollisionTypePaddle
		fmt.Println("Detected collision with paddle!")
		g.BounceBall(horizontalSurfaceVector, g.paddleMovementVector)
	}

	g.BallPosition = newVector

	return collisionType
}
