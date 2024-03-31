package validity

import (
	"slices"
)

// var numTimePeriod = 34

func CheckForValidity(potentialAvailability uint64, universityAvailability []bool) bool {
	// fmt.Println("====")
	// get the position of the first 1 in the potentialAvailability

	// convert universityAvailability to int64
	var universityAvailabilityInt uint64
	slices.Reverse(universityAvailability)
	defer slices.Reverse(universityAvailability)
	for i, v := range universityAvailability {
		if v {
			universityAvailabilityInt |= 1 << i
		}
	}
	// fmt.Printf("0b%b\n0b%b\n", potentialAvailability, universityAvailabilityInt)

	// check if the "and" is greater than 0
	// aka if there is a time period that is blocked off that is also active for the university in the schedule
	res := potentialAvailability & universityAvailabilityInt
	// fmt.Printf("Res=0b%b, %v\n", res, res)
	return res <= 0
}
