package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/kevin51jiang/msci-435-go/solver"
)

func solveChairs(day int) {

	const numParticipants = 10
	numIts := 0
	const minIts = 10000

	// Read in and parse data/maxMeetings.csv
	maxMeetings, err := solver.ParseMaxMeetingsCSV("../data/maxMeetings.csv")
	if err != nil {
		log.Fatalf("Failed to parse maxMeetings.csv: %v", err)
	}
	fmt.Printf("Max Meetings: %v\n", maxMeetings)

	// // Read in and parse data/combinations-chairs-0.tsv
	chairCombos, err := solver.ParseChairCombos(fmt.Sprintf("../data/combinations-chairs-%v.tsv", day), maxMeetings)
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

	visited := make(map[string]bool)
	visited[solver.GetSolutionKey(solution)] = true
	bestEquitability := solver.CalculateEquitability(solution, numParticipants)
	bestEquitabilityLog := make([]float64, 0)
	bestEquitabilityLog = append(bestEquitabilityLog, bestEquitability)

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
				bestEquitabilityLog = append(bestEquitabilityLog, bestEquitability)
			}
			visited[solver.GetSolutionKey(entry)] = true
		}

		numIts++

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
		if numIts%100 == 0 {
			fmt.Println("Iteration: ", numIts, " Best Equitability: ", bestEquitability, " Solution length: ", len(solution))
		}

		if numIts > minIts && areLast20EntriesSame(bestEquitabilityLog) {
			break
		}
	}

	// write out bestEquitabilityLog to a file
	writeToLog(day, bestEquitabilityLog)
	writeResults(day, solution, numParticipants)

	fmt.Println("Final Solution: ", solution)
	fmt.Println()

	fmt.Println("Initial Meetings: ", solver.GetSolutionDistribution(initialSolution, numParticipants))
	fmt.Println("Final Meetings: ", solver.GetSolutionDistribution(solution, numParticipants))

	fmt.Println()

	fmt.Println(solver.DisplaySolution(solution, numParticipants, "Chair"))

}

func writeToLog(day int, eqLog []float64) {
	filename := fmt.Sprintf("../data/equitabilityLog-%v.csv", day)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	for _, eq := range eqLog {
		file.WriteString(fmt.Sprintf("%v\n", eq))
	}
}

type DisplayStyle string

const (
	CSV      DisplayStyle = "csv"
	HTML     DisplayStyle = "html"
	Markdown DisplayStyle = "markdown"
	Console  DisplayStyle = "console"
)

func DisplaySolution(solution []solver.ParticipantsCombination, numParticipants int8, title string, displayStyle DisplayStyle) string {

	t := table.NewWriter()
	// t.SetCaption("Suggested Schedule (Tabu Search)")

	header := table.Row{}
	header = append(header, "Time")
	for i := 0; i < int(numParticipants); i++ {
		header = append(header, title+" "+strconv.Itoa(i+1))
	}
	t.AppendHeader(header)

	for _, combo := range solution {
		participSchedule := table.Row{}
		participSchedule = append(participSchedule, strconv.Itoa(int(combo.Time)))

		isParticipating := make([]bool, numParticipants)
		for _, participant := range combo.Participants {
			isParticipating[participant] = true
		}
		for _, isParticipating := range isParticipating {
			if isParticipating {
				participSchedule = append(participSchedule, "X")
			} else {
				participSchedule = append(participSchedule, "")
			}
		}

		t.AppendRow(participSchedule)
	}

	if displayStyle == CSV {
		return t.RenderCSV()
	} else if displayStyle == HTML {
		return t.RenderHTML()
	} else if displayStyle == Markdown {
		return t.RenderMarkdown()
	} else {
		return t.Render()
	}

}

func writeResults(day int, solution []solver.ParticipantsCombination, numParticipants int8) {
	filename := fmt.Sprintf("../data/solution-%v.csv", day)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	file.WriteString(DisplaySolution(solution, numParticipants, "Member", CSV))
}

func areLast20EntriesSame(arr []float64) bool {
	if len(arr) < 20 {
		return false
	}

	last20 := arr[len(arr)-20:]
	for i := 1; i < len(last20); i++ {
		if last20[i] != last20[0] {
			return false
		}
	}

	return true
}

// func solveMembers(day int) {
// 	const numParticipants = 30

// 	// availGrid := validity.LoadAvailabilityGrid("../data/allAvailabilities.tsv")

// 	numIts := 0
// 	const maxIts = 50000

// 	// Read in and parse data/maxMeetings.csv
// 	maxMeetings, err := solver.ParseMaxMeetingsCSV("../data/maxMeetings.csv")
// 	if err != nil {
// 		log.Fatalf("Failed to parse maxMeetings.csv: %v", err)
// 	}
// 	fmt.Printf("Max Meetings: %v\n", maxMeetings)

// 	// // Read in and parse data/combinations-members-0.tsv
// 	filename := fmt.Sprintf("../data/combinations-members-%v.tsv", day)
// 	participCombos, err := solver.ParseMemberCombos(filename, maxMeetings)
// 	if err != nil {
// 		log.Fatalf("Failed to parse %v: %v", filename, err)
// 	}
// 	fmt.Printf("Member Combos: len(%v)\n", len(participCombos))

// 	// initial solution is the first entry for each distinct day and time
// 	initialSolution := make([]solver.ParticipantsCombination, 0)

// 	// // write all the keys of participCombos to a file
// 	// keysFile, err := os.Create("../data/participCombosKeys.csv")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create file: %v", err)
// 	// }
// 	// defer keysFile.Close()

// 	// for k := range participCombos {
// 	// 	keysFile.WriteString(fmt.Sprintf("%v\n", k))
// 	// }

// 	for time := 1; time < 34+1; time++ {
// 		hasFound := false
// 		for k, combo := range participCombos {
// 			if combo.Time == int8(time) {
// 				hasFound = true
// 				fmt.Println("Adding combo", k, " ", combo)
// 				selectedCombo := solver.ParticipantsCombination{}
// 				selectedCombo.Time = combo.Time
// 				selectedCombo.Participants = make([]int8, len(combo.Participants))
// 				copy(selectedCombo.Participants, combo.Participants)

// 				initialSolution = append(initialSolution, selectedCombo)
// 				break
// 			}
// 			continue
// 		}
// 		if !hasFound {
// 			log.Panicf("Error: Time not found in participant combos: %v", time)
// 		}
// 	}

// 	fmt.Println("Initial Solution: ", initialSolution, " Length: ", len(initialSolution))

// 	solution := make([]solver.ParticipantsCombination, len(initialSolution))
// 	copy(solution, initialSolution)
// 	// prevSol := initialSolution

// 	visited := make(map[string]bool)
// 	visited[solver.GetSolutionKey(solution)] = true
// 	bestEquitability := solver.CalculateEquitability(solution, numParticipants)
// 	bestEquitabilityLog := make([]float64, 0)
// 	bestEquitabilityLog = append(bestEquitabilityLog, bestEquitability)

// 	// Tabu Search
// 	for {
// 		// Get neighborhood
// 		neighborhood := solver.GetNeighborhood(solution, participCombos, numParticipants)

// 		// Find the best entry in the neighborhood
// 		for _, entry := range neighborhood {
// 			// fmt.Println("Entry: ", solver.GetSolutionDistribution(entry, numParticipants))
// 			newEq := solver.CalculateEquitability(entry, numParticipants)
// 			// fmt.Println("New Eq: ", newEq, " Best Eq: ", bestEquitability)
// 			if newEq <= bestEquitability && !visited[solver.GetSolutionKey(entry)] {
// 				copy(solution, entry)
// 				bestEquitability = solver.CalculateEquitability(entry, numParticipants)
// 				bestEquitabilityLog = append(bestEquitabilityLog, bestEquitability)
// 			}
// 			visited[solver.GetSolutionKey(entry)] = true
// 		}
// 		numIts++

// 		if numIts%100 == 0 {
// 			fmt.Println("Iteration: ", numIts, " Best Equitability: ", bestEquitability, " Solution length: ", len(solution))
// 		}

// 		// if len(solution) == len(prevSol) {
// 		// 	equal := true
// 		// 	for i := range solution {
// 		// 		if !reflect.DeepEqual(solution[i], prevSol[i]) {
// 		// 			equal = false
// 		// 			break
// 		// 		}
// 		// 	}
// 		// 	if equal {
// 		// 		break
// 		// 	}
// 		// }

// 		// check if bestEquitability has not changed in the last n iterations
// 		// if it hasn't, break

// 		if numIts > maxIts && areLast20EntriesSame(bestEquitabilityLog) {
// 			break
// 		}
// 	}

// 	// write out bestEquitabilityLog to a file
// 	writeToLog(day, bestEquitabilityLog)
// 	writeResults(day, solution, numParticipants)

// 	fmt.Println("Final Solution: ", solution)
// 	fmt.Println()

// 	fmt.Println("Initial Meetings: ", solver.GetSolutionDistribution(initialSolution, numParticipants))
// 	fmt.Println("Final Meetings: ", solver.GetSolutionDistribution(solution, numParticipants))

// 	fmt.Println()

// 	fmt.Println(solver.DisplaySolution(solution, numParticipants, "Member"))

// }

func main() {

	// fmt.Println()
	solveChairs(0)

	// solveMembers(0)
	// solveMembers(1)
	// solveMembers(2)
	// solveMembers(3)
	// solveMembers(4)

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
