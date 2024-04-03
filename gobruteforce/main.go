package main

import (
	"fmt"
	"log"

	"github.com/kevin51jiang/msci-435-go/solver"
)

func main() {

	// Read in and parse data/maxMeetings.csv

	maxMeetings, err := solver.ParseMaxMeetingsCSV("../data/maxMeetings.csv")
	if err != nil {
		log.Fatalf("Failed to parse maxMeetings.csv: %v", err)
	}
	fmt.Printf("Max Meetings: %v\n", maxMeetings)

	// // Read in and parse data/combinations-chairs-0.tsv
	chairCombos, err := solver.ParseChairCombos("../data/combinations-chairs-0.tsv", maxMeetings)
	if err != nil {
		log.Fatalf("Failed to parse combinations-chairs-0.tsv: %v", err)
	}
	fmt.Printf("Chair Combos: len(%v)\n", len(chairCombos))

	// Create a map that goes from the participants in a combo to the combo
	chairCombosMap := make(map[string][]solver.ParticipantsCombination)
	for _, combo := range chairCombos {
		chairCombosMap[combo.GetKey()] = append(chairCombosMap[combo.GetKey()], combo)
	}

	// initial solution is the first entry for each distinct day and time
	initialSolution := make([]solver.ParticipantsCombination, 0)

	const numParticipants = 10
	for time := 1; time < 34+1; time++ {
		for c_ind, combo := range chairCombos {
			if combo.Time != int8(time) {
				continue
			}
			fmt.Println("Adding combo", c_ind, " ", combo)
			initialSolution = append(initialSolution, combo)
			break
		}
	}

	fmt.Println("Initial Solution: ", initialSolution, " Length: ", len(initialSolution))

	solution := make([]solver.ParticipantsCombination, len(initialSolution))
	copy(solution, initialSolution)
	// prevSol := initialSolution

	numIts := 0
	const maxIts = 10000

	visited := make(map[string]bool)
	visited[solver.GetSolutionKey(solution)] = true
	bestEquitability := solver.CalculateEquitability(solution, numParticipants)

	// Tabu Search
	for {
		// Get neighborhood
		neighborhood := solver.GetNeighborhood(solution, chairCombosMap, numParticipants)

		// Find the best entry in the neighborhood
		for _, entry := range neighborhood {
			// fmt.Println("Entry: ", solver.GetSolutionDistribution(entry, numParticipants))
			newEq := solver.CalculateEquitability(entry, numParticipants)
			// fmt.Println("New Eq: ", newEq, " Best Eq: ", bestEquitability)
			if newEq <= bestEquitability && !visited[solver.GetSolutionKey(entry)] {
				copy(solution, entry)
				bestEquitability = solver.CalculateEquitability(entry, numParticipants)
			}
			visited[solver.GetSolutionKey(entry)] = true
		}
		numIts++
		fmt.Println("Iteration: ", numIts, " Best Equitability: ", bestEquitability, " Solution length: ", len(solution))

		// if len(solution) == len(prevSol) {
		// 	equal := true
		// 	for i := range solution {
		// 		if !reflect.DeepEqual(solution[i], prevSol[i]) {
		// 			equal = false
		// 			break
		// 		}
		// 	}
		// 	if equal {
		// 		break
		// 	}
		// }

		if numIts > maxIts {
			break
		}
	}

	fmt.Println("Final Solution: ", solution)
	fmt.Println()

	fmt.Println("Initial Meetings: ", solver.GetSolutionDistribution(initialSolution, numParticipants))
	fmt.Println("Final Meetings: ", solver.GetSolutionDistribution(solution, numParticipants))

	fmt.Println()

	fmt.Println(solver.DisplaySolution(solution, numParticipants, "Chair"))

	// // Read in and parse data/combinations-members-0.tsv
	// memberCombos, err := parseMemberCombos("../data/combinations-members-0.tsv", maxMeetings)
	// if err != nil {
	// 	log.Fatalf("Failed to parse combinations-members-0.tsv: %v", err)
	// }
	// fmt.Printf("Member Combos: len(%v)\n", len(memberCombos))

	// // Show the distribution of the number of meetings, based on time for memberCombos
	// meetingsDistribution := make(map[int]int)
	// for _, memberCombo := range memberCombos {
	// 	meetingsDistribution[int(memberCombo.Day)*5+int(memberCombo.Time)]++
	// }
	// fmt.Printf("Meetings Distribution: %v\n", meetingsDistribution)

	// // Calculate the equitability of the chairCombos
	// chairEquitability := calculateEquitability(chairCombos[:34], 10)
	// fmt.Printf("Chair Equitability: %v\n", chairEquitability)

	// fmt.Println("Neighborhood:")
	// // Get neighborhood
	// neighborhood := getNeighborhood(chairCombos[34:34+34], chairCombosMap, 10)
	// for _, entry := range neighborhood {
	// 	fmt.Println(entry)
	// }

	// for ind, entry := range neighborhood {
	// 	fmt.Println(ind, ' ', len(entry), "eq: ", calculateEquitability(entry, 10))
	// }
	// fmt.Printf("Neighborhood: %v\n", neighborhood)

}
