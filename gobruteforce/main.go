package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/bits"
	"os"
	"runtime"
	"sort"

	"github.com/cheggaaa/pb/v3"
	"github.com/kevin51jiang/msci-435-go/validity"
)

var topEntriesKept = 1000000

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Unmarshal allAvails.json
type Availability struct {
	M  []bool `json:"M"`
	T  []bool `json:"T"`
	W  []bool `json:"W"`
	Th []bool `json:"Th"`
	F  []bool `json:"F"`
}

type University struct {
	University   string       `json:"university"`
	IsAvailable  Availability `json:"isAvailable"`
	ComboIndexes []int        `json:"comboIndexes"`
}

type Universities map[string]University

// Pretty print JSON
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func readUniversities() Universities {
	data, err := os.ReadFile("allAvails.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal the JSON data into a Universities object
	var unis Universities
	err = json.Unmarshal(data, &unis)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return unis
}

func printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
	fmt.Printf("TotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
	fmt.Printf("Sys = %v MiB\n", bToMb(m.Sys))
	fmt.Printf("NumGC = %v\n", m.NumGC)
}

func generateCombos() {
	// count of iterations
	count := int64(1 << 34)

	// create and start new bar
	bar := pb.Full.Start64(count)

	// Generate every 64bit uint from 0 to 2^34-1
	for i := int64(0); i < count; i++ {
		bar.Increment()
	}

	// finish bar
	bar.Finish()
}

func main() {

	// validity.CheckForValidity(int64(0b00001), []bool{true, false, true, false, true})

	universities := readUniversities()

	validCombos := map[string][]uint64{}
	validComboMinimums := map[string]int{}

	// count of iterations
	count := uint64(1 << 34)

	// create and start new bar
	bar := pb.Full.Start64(int64(count))

	// Generate every 64bit uint from 0 to 2^34-1
	// for i := count - 1; i >= 0; i-- {
	for i := uint64(0); i < count; i++ {
		bar.Increment()

		// For each day in the universities
		for uniName, uni := range universities {
			// check if this combo exceeds the minimum number of meetings, and is valid
			if bits.OnesCount64(i) >= validComboMinimums[uniName] {
				schedule := uni.IsAvailable.M
				if validity.CheckForValidity(i, schedule) {
					validCombos[uniName] = append(validCombos[uniName], uint64(i))
				}
			}
		}

		if i%uint64(topEntriesKept*10) == 0 {
			// in valid combos, keep the top 10000000 combos with the maximum amount of bits for each university
			for ind, _ := range universities {

				if len(validCombos[ind]) > topEntriesKept {
					// sort the valid combos
					sort.Slice(validCombos[ind], func(i, j int) bool {
						return bits.OnesCount64(validCombos[ind][i]) > bits.OnesCount64(validCombos[ind][j])
					})

					validCombos[ind] = validCombos[ind][:topEntriesKept]

					// keep track of the minimum amount of bits
					validComboMinimums[ind] = bits.OnesCount64(validCombos[ind][topEntriesKept-1])
				}
			}
			fmt.Println("M1 Minimum: ", validComboMinimums["M1"])
		}

		if (i % 10000000) == 0 {
			printMemoryUsage()
		}

	}

	// finish bar
	bar.Finish()

}
