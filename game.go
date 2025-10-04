package main

import "math"

const PaddleWidth = 60
const PaddleHeight = 10
const BallDiameter = 10
const PaddleRelativePosition = 0.8 // 80% of grid height

type Game struct {
	GridWidth    int   // grid width in pixels
	GridHeight   int   // grid height in pixels
	Lives        int   // how many lives the player has left
	BallPosition Point // X,Y of ball on grid
	BallVector   Vector
	PaddleX      int // current position of the paddle horizontally
	PaddleY      int // fixed position of the paddle vertically
	Score        int // Player's current score
	Level        int // Current level
}

func NewGame(gridWidth, gridHeight int) *Game {
	return &Game{
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
		Lives:      3,
		Score:      0,
		Level:      1,
		PaddleY:    int(float64(gridHeight) * PaddleRelativePosition),
	}
}

func (g *Game) Initialize() {
	g.Lives = 3
	g.Score = 0
	g.Level = 1
}

func (g *Game) MovePaddle(direction Direction) {
	moveAmount := 5 // number of pixels to move

	switch direction {
	case DirectionLeft:
		g.PaddleX -= moveAmount
		if g.PaddleX < 0 {
			g.PaddleX = 0
		}
	case DirectionRight:
		g.PaddleX += moveAmount
		if g.PaddleX+PaddleWidth > g.GridWidth {
			g.PaddleX = g.GridWidth - PaddleWidth
		}
	default:
		panic("unrecognized direction")
	}
}

func (g *Game) PutBallInMiddle() {
	x := (g.GridWidth / 2) + (BallDiameter / 2)
	y := (g.PaddleY / 2) + (PaddleHeight / 2)
	g.BallPosition = Point{x, y}
}

func (g *Game) BounceBall(axis Axis) {
	// TODO: implement angle inversion along given axis
}

func (g *Game) SetRandomUpwardBallVector() {
	// TODO: implement this
}

func (g *Game) MoveBallAlongTrajectory() CollisionType {
	collisionType := CollisionTypeNone
	currentX := g.BallPosition.X
	currentY := g.BallPosition.Y

	speed := float64(g.BallVector.Speed) // assume the value = the number of pixels moved per tick.  we prob need to adjust this
	speed += float64((g.Level - 1) * 2)  // make the ball move faster with every level. again, we prob need to adjust this

	// do the math to calculate the new position based on the ball's vector
	dX := speed * math.Cos(g.BallVector.Angle)
	dY := speed * math.Sin(g.BallVector.Angle)

	newX := currentX + int(dX)
	newY := currentY + int(dY)

	// correct based on collisions
	if newX < 0 {
		newX = 0
		collisionType = CollisionTypeWall
		g.BounceBall(XAxis)
	} else if newX+BallDiameter > g.GridWidth {
		newX = g.GridWidth - BallDiameter
	}

	if newY < 0 {
		newY = 0
		collisionType = CollisionTypeWall
		g.BounceBall(YAxis)
	} else if newY+BallDiameter > g.GridWidth {
		newY = g.GridHeight - BallDiameter
	}

	// There's a collision with the paddle if:
	// the right edge of the ball's bounding box is at or past the left edge of the paddle and
	// the left edge of the ball's bounding box is at or before the right edge of the paddle and
	// the bottom edge of the ball is at or lower than the paddle's top edge
	if newX+BallDiameter >= g.PaddleX && newX <= g.PaddleX+PaddleWidth && newY+BallDiameter >= g.PaddleY {
		newY = g.PaddleY - BallDiameter // normalize the ball's Y position
		collisionType = CollisionTypePaddle
		g.BounceBall(YAxis)
	}

	g.BallPosition.X = newX
	g.BallPosition.Y = newY

	return collisionType
}
