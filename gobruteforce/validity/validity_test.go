package validity_test

import (
	"testing"

	"github.com/kevin51jiang/msci-435-go/validity"
)

func TestCheckForValidity(t *testing.T) {
	tests := []struct {
		name                   string
		potentialAvailability  uint64
		universityAvailability []bool
		want                   bool
	}{
		{
			name:                   "Test 0: All university availabilities are true",
			potentialAvailability:  0b0000,
			universityAvailability: []bool{true, true, true, true},
			want:                   true,
		},
		{
			name:                   "Test 1: All university availabilities are false",
			potentialAvailability:  0b0000,
			universityAvailability: []bool{false, false, false, false},
			want:                   true,
		},
		{
			name:                   "Test 2: Mixed university availabilities",
			potentialAvailability:  0b1010,
			universityAvailability: []bool{true, false, true, false},
			want:                   false,
		},
		{
			name:                   "Test 3: Blocked off university availability",
			potentialAvailability:  0b1,
			universityAvailability: []bool{true},
			want:                   false,
		},
		{
			name:                   "Test 4: Partial blocked off university availability (SMALL)",
			potentialAvailability:  0b10,
			universityAvailability: []bool{true, false},
			want:                   false,
		},
		{
			name:                   "Test 5: Partial blocked off university availability (LARGE)",
			potentialAvailability:  0b0010,
			universityAvailability: []bool{true, false, true, false},
			want:                   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validity.CheckForValidity(tt.potentialAvailability, tt.universityAvailability); got != tt.want {
				t.Errorf("checkForValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}
