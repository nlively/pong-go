package main

import (
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func (v *Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v *Vector) Subtract(v2 Vector) Vector {
	return Vector{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v *Vector) Multiply(scale float64) Vector {
	return Vector{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func (v *Vector) Dot(v2 Vector) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v *Vector) Reflect(surfaceDirection Vector, surfaceVelocity Vector) Vector {
	n := surfaceDirection.Perpendicular()
	nHat := n.Normalize()

	vRel := v.Subtract(surfaceVelocity)

	vnMag := vRel.Dot(nHat)

	vRelReflected := vRel.Subtract(nHat.Multiply(2 * vnMag))

	return vRelReflected.Add(surfaceVelocity)
}

// k: tangential impulse gain (0 = none; try 0.5 for Pong-like steering)
func (v *Vector) ReflectWithTangentialImpulse(surfaceDirection Vector, surfaceVelocity Vector, k float64) Vector {
	// unit normal
	n := surfaceDirection.Perpendicular()
	n = n.Normalize()

	t := Vector{-n.Y, n.X} // unit tangent

	vRel := v.Subtract(surfaceVelocity)

	vn := n.Multiply(vRel.Dot(n)) // normal component
	vt := t.Multiply(vRel.Dot(t)) // tangential component

	// reflect normal; add tangential impulse from paddle motion
	vnReflected := vn.Multiply(-1)

	// paddle tangential speed (in world), projected onto t
	u_t := surfaceVelocity.Dot(t)
	vtReflected := vt.Add(t.Multiply(k * u_t))

	combined := vnReflected.Add(vtReflected)
	return combined.Add(surfaceVelocity)
}

func (v *Vector) Perpendicular() Vector {
	return Vector{
		X: v.Y,
		Y: v.X,
	}
}

func (v *Vector) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v *Vector) Normalize() Vector {
	l := v.Length()
	if l == 0 {
		return Vector{0, 0}
	}
	return Vector{v.X / l, v.Y / l}
}

func (v *Vector) ToPolar() Polar {
	speed := math.Hypot(v.X, v.Y) // sqrt(vx^2 + vy^2)
	angleDeg := math.Atan2(v.Y, v.X) * 180.0 / math.Pi

	return Polar{
		Speed: speed,
		Angle: angleDeg,
	}
}
