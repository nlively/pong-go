package main_test

import (
	"testing"

	main "github.com/nlively/pong-go"
)

func Test_Polar_ToVector(t *testing.T) {
	tests := []struct {
		name  string
		polar main.Polar
		want  main.Vector
	}{
		{"30 degrees at 5 speed", main.Polar{30, 5}, main.Vector{4.330127018922194, 2.4999999999999996}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.polar.ToVector()
			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}
