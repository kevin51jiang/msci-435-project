package main

import (
	"testing"

	"github.com/kevin51jiang/msci-435-go/solver"
)

func TestCalculateEquitability(t *testing.T) {

	m1 := []int8{0, 1, 2}

	m2 := []int8{1, 2}

	tests := []struct {
		name            string
		combos          []solver.ParticipantsCombination
		numParticipants int8
		want            float64
	}{
		{
			name: "Test 0: Equal number of meetings",
			combos: []solver.ParticipantsCombination{
				{Participants: m1, Day: 1, Time: 1},
				{Participants: m1, Day: 1, Time: 1},
				{Participants: m1, Day: 1, Time: 1},
				{Participants: m1, Day: 1, Time: 1},
			},
			numParticipants: 3,
			want:            0.0,
		},
		{
			name: "Test 1: Unequal number of meetings",
			combos: []solver.ParticipantsCombination{
				{Participants: m2},
				{Participants: m1},
				{Participants: m1},
				{Participants: m1},
			},
			numParticipants: 3,
			want:            0.6666666666666667,
		},
		{
			name:            "Test 2: No meetings",
			combos:          []solver.ParticipantsCombination{},
			numParticipants: 5,
			want:            0.0,
		},
		{
			name: "Test 3: Quadratic score",
			combos: []solver.ParticipantsCombination{
				{Participants: []int8{1, 2}},
				{Participants: []int8{1, 2}},
				{Participants: []int8{1, 2}},
				{Participants: []int8{0, 1, 2}},
			},
			numParticipants: 3,
			want:            6.0,
		},
		{
			name: "Test 4: Empty universities",
			combos: []solver.ParticipantsCombination{
				{Participants: []int8{0, 1, 2}},
				{Participants: []int8{0, 1, 2, 3}},
				{Participants: []int8{0, 1, 2}},
			},
			numParticipants: 6,
			want:            1.66666666667,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solver.CalculateEquitability(tt.combos, tt.numParticipants); got != tt.want {
				t.Errorf("calculateEquitability() = %v, want %v", got, tt.want)
			}
		})
	}
}
