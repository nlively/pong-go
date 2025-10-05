package main_test

import (
	"testing"

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
			got := tt.object.Reflect(tt.surface)
			if got != tt.want {
				t.Errorf("%v reflect against surface %v. want %v, got %v", tt.object, tt.surface, tt.want, got)
			}
		})
	}
}
