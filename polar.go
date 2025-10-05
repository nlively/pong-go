package main

import "math"

type Polar struct {
	Angle float64 // degrees
	Speed float64
}

func (p *Polar) ToVector() Vector {
	rad := deg2Rad(p.Angle)
	x := p.Speed * math.Cos(rad)
	y := p.Speed * math.Sin(rad)

	return Vector{
		X: x,
		Y: y,
	}
}
