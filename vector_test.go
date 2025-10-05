package main_test

import (
	"testing"

	"math"

	main "github.com/nlively/pong-go"
)

func LeftWall() main.Vector {
	return main.Vector{0, 1}
}
func RightWall() main.Vector {
	return main.Vector{0, -1}
}
func Ceiling() main.Vector {
	return main.Vector{-1, 0}
}
func Floor() main.Vector {
	return main.Vector{1, 0}
}

func Test_Vector_Perpendicular(t *testing.T) {
	tests := []struct {
		name    string
		surface main.Vector
		want    main.Vector
	}{
		{"ceiling", Ceiling(), main.Vector{0, -1}},
		{"floor", Floor(), main.Vector{0, 1}},
		{"left wall", LeftWall(), main.Vector{1, 0}},
		{"right wall", RightWall(), main.Vector{-1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.surface.Perpendicular()
			if got != tt.want {
				t.Errorf("%v surface normal. want %v, got %v", tt.surface, tt.want, got)
			}
		})
	}
}

func Test_Vector_Reflect(t *testing.T) {
	tests := []struct {
		name    string
		object  main.Vector
		surface main.Vector
		want    main.Vector
	}{
		{"flat horizontal surface, straight drop", main.Vector{0, -5}, Floor(), main.Vector{0, 5}},
		{"hit the ceiling at a 30 degree angle", main.Vector{4.330127018922194, 2.4999999999999996}, Ceiling(), main.Vector{4.330127018922194, -2.4999999999999996}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// surface is static in these tests, so pass zero surface velocity
			got := tt.object.Reflect(tt.surface, main.Vector{0, 0})
			if !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("%v reflect against surface %v. want %v, got %v", tt.object, tt.surface, tt.want, got)
			}
		})
	}
}

func Test_Vector_ReflectWithTangentialImpulse(t *testing.T) {
	tests := []struct {
		name       string
		object     main.Vector
		surface    main.Vector
		surfaceVel main.Vector
		k          float64
		want       main.Vector
	}{
		{"flat horizontal, no surface velocity", main.Vector{0, -5}, Floor(), main.Vector{0, 0}, 0.5, main.Vector{0, 5}},
		{"flat horizontal, paddle moving right", main.Vector{0, -5}, Floor(), main.Vector{2, 0}, 0.5, main.Vector{1, 5}},
		{"flat horizontal, paddle moving left", main.Vector{0, -5}, Floor(), main.Vector{-2, 0}, 0.5, main.Vector{-1, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.object.ReflectWithTangentialImpulse(tt.surface, tt.surfaceVel, tt.k)
			if !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("%v reflect with tangential impulse. want %v, got %v", tt.object, tt.want, got)
			}
		})
	}

	// Additional tests: angled surfaces and varying k
	angledTests := []struct {
		name       string
		object     main.Vector
		surface    main.Vector
		surfaceVel main.Vector
		k          float64
	}{
		{"angled surface static, k=0.5", main.Vector{3, -1}, main.Vector{1, 1}, main.Vector{0, 0}, 0.5},
		{"angled surface moving, k=1.0", main.Vector{3, -1}, main.Vector{1, 1}, main.Vector{2, 0}, 1.0},
		{"steep surface moving left, k=0.2", main.Vector{-2, -3}, main.Vector{0, 1}, main.Vector{-3, 0}, 0.2},
	}

	for _, tt := range angledTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.object.ReflectWithTangentialImpulse(tt.surface, tt.surfaceVel, tt.k)
			want := computeExpectedReflectWithTangentialImpulse(tt.object, tt.surface, tt.surfaceVel, tt.k)
			if !approxEqual(got, want, 1e-9) {
				t.Errorf("%s: got %v, want %v", tt.name, got, want)
			}
		})
	}
}

// approxEqual returns true if two vectors are within eps of each other
func approxEqual(a, b main.Vector, eps float64) bool {
	return math.Abs(a.X-b.X) <= eps && math.Abs(a.Y-b.Y) <= eps
}

// computeExpectedReflectWithTangentialImpulse reproduces the math of ReflectWithTangentialImpulse
// so tests can verify the implementation behaves as intended.
func computeExpectedReflectWithTangentialImpulse(v main.Vector, surfaceDir main.Vector, surfaceVel main.Vector, k float64) main.Vector {
	// unit normal
	n := surfaceDir.Perpendicular()
	n = n.Normalize()

	t := main.Vector{-n.Y, n.X}

	vRel := v.Subtract(surfaceVel)

	vn := n.Multiply(vRel.Dot(n))
	vt := t.Multiply(vRel.Dot(t))

	vnReflected := vn.Multiply(-1)

	u_t := surfaceVel.Dot(t)
	vtReflected := vt.Add(t.Multiply(k * u_t))

	combined := vnReflected.Add(vtReflected)
	return combined.Add(surfaceVel)
}
