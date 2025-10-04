package main

type Point struct {
	X float64
	Y float64
}

type Vector struct {
	Angle int // degrees
	Speed int
}
type Direction int

const (
	DirectionLeft  Direction = 0
	DirectionRight Direction = 1
)

type CollisionType int

const (
	CollisionTypeNone   CollisionType = 0
	CollisionTypePaddle CollisionType = 1
	CollisionTypeWall   CollisionType = 2
)

type Axis int

const (
	XAxis Axis = 0
	YAxis Axis = 1
)
