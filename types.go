package main

import "math"

type Point struct {
	X float64
	Y float64
}

func deg2Rad(angleDeg float64) float64 {
	return angleDeg * (math.Pi / 180.0)
}

type Direction int

const (
	DirectionLeft  Direction = 1
	DirectionRight Direction = 2
	DirectionNone  Direction = 0
)

type CollisionType int

const (
	CollisionTypeNone   CollisionType = 0
	CollisionTypePaddle CollisionType = 1
	CollisionTypeWall   CollisionType = 2
)
