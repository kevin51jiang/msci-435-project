package validity

import (
	"encoding/csv"
	"io"
	"os"
	"slices"
)

type University struct {
	Name         string
	Code         string
	Availability [5][]bool
}

// UniversityID, Day, Time
type ParticipantAvailabilityGrid [][5][34]bool

func LoadAvailabilityGrid(csvName string) ParticipantAvailabilityGrid {
	// Load the availability grid from a csv file

	// open file
	file, err := os.Open(csvName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read file
	var availabilityGrid ParticipantAvailabilityGrid

	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	// skip header
	csvReader.Read()
	for {
		// read line
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		availabilityGrid = append(availabilityGrid, [5][34]bool{})
		for i := 0; i < 5; i++ {
			for j := 0; j < 34; j++ {
				availabilityGrid[len(availabilityGrid)-1][i][j] = line[i*34+j] == "True"
			}
		}
	}

	return availabilityGrid
}

func CheckIfUniversityAvailable(participant int8, day int8, time int8, universityAvailability *ParticipantAvailabilityGrid) bool {
	return (*universityAvailability)[participant][day][time]
}

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
