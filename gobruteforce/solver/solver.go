package solver

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/maolinc/copier"
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

type Meeting struct {
	Day         int
	Time        int
	NumMeetings int
	Chairs      []int
	Members     []int
}

type MaxMeeting struct {
	Day         int8
	Time        int8
	MaxMeetings int8
}

func ParseMaxMeetingsCSV(filename string) ([]MaxMeeting, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var maxMeetings []MaxMeeting
	for _, line := range lines[1:] { // Skip header
		day, _ := strconv.Atoi(line[0])
		time, _ := strconv.Atoi(line[1])
		maxMeetingsStr := line[2]
		maxMeetingsFloat, _ := strconv.ParseFloat(maxMeetingsStr, 64)
		maxMeetingsInt := int(maxMeetingsFloat)

		maxMeetings = append(maxMeetings, MaxMeeting{Day: int8(day), Time: int8(time - 24 + 1), MaxMeetings: int8(maxMeetingsInt)})
	}

	return maxMeetings, nil
}

type ParticipantsCombination struct {
	Day          int8
	Time         int8
	Participants []int8
}

func (p ParticipantsCombination) GetKey() string {
	participants := make([]string, len(p.Participants))
	for i, participant := range p.Participants {
		participants[i] = strconv.Itoa(int(participant))
	}
	return strings.Join(participants, "-") + "--" + strconv.Itoa(int(p.Time)) + "--" + strconv.Itoa(int(p.Day))
}

func ParseChairCombos(filename string, maxMeetings []MaxMeeting) ([]ParticipantsCombination, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 3

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var chairCombinations []ParticipantsCombination
	for _, line := range lines[1:] { // Skip header
		day, _ := strconv.Atoi(line[0])
		time, _ := strconv.Atoi(line[1])

		chairsStr := strings.Trim(line[2], "[]")
		chairsStrSplit := strings.Split(chairsStr, ",")
		chairs := make([]int8, len(chairsStrSplit))
		for i, chairStr := range chairsStrSplit {
			chair, _ := strconv.Atoi(strings.TrimSpace(chairStr))
			chairs[i] = int8(chair)
		}

		maxMeetingsAtTime := maxMeetings[day*5+time]
		if int(maxMeetingsAtTime.MaxMeetings) != len(chairs) {
			continue
		}

		if len(chairCombinations) == 2476 {
			fmt.Println("raw", line)
		}

		chairCombinations = append(chairCombinations, ParticipantsCombination{Day: int8(day), Time: int8(time), Participants: chairs})
	}

	return chairCombinations, nil
}

func ParseMemberCombos(filename string, maxMeetings []MaxMeeting) (map[string]ParticipantsCombination, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 3

	_, _ = reader.Read() // Skip header

	memberCombinations := make(map[string]ParticipantsCombination)
	ind := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		day, _ := strconv.Atoi(line[0])
		time, _ := strconv.Atoi(line[1])

		membersStr := strings.Trim(line[2], "[]")
		membersStrSplit := strings.Split(membersStr, ",")
		members := make([]int8, len(membersStrSplit))
		for i, memebersStr := range membersStrSplit {
			member, _ := strconv.Atoi(strings.TrimSpace(memebersStr))
			members[i] = int8(member)
		}

		if ind%1000000 == 0 {
			fmt.Printf("Parsed %v member combinations\n", ind)
		}

		maxMeetingsAtTime := maxMeetings[day*5+time]
		ind++
		if int(maxMeetingsAtTime.MaxMeetings)*5 != len(members) {
			// fmt.Println("Max Meetings at Time: ", maxMeetingsAtTime.MaxMeetings, ' ', len(members))
			continue
		}

		participCombo := ParticipantsCombination{Day: int8(day), Time: int8(time), Participants: members}

		memberCombinations[participCombo.GetKey()] = participCombo

	}

	return memberCombinations, nil
}

func CalculateEquitability(combos []ParticipantsCombination, numParticipants int8) float64 {
	participantMeetings := make([]int8, numParticipants)
	for _, combo := range combos {
		for _, participant := range combo.Participants {
			participantMeetings[participant]++
		}
	}
	// fmt.Println("Participant Meetings: ", participantMeetings)

	// find MSE of participantMeetings
	expectedNumMeetings := 0.0
	for _, numMeetings := range participantMeetings {
		expectedNumMeetings += float64(numMeetings)
	}
	expectedNumMeetings /= float64(numParticipants)
	// fmt.Println("Expected Num Meetings: ", expectedNumMeetings)

	if expectedNumMeetings == 0 {
		return 0
	}

	// Calculate the penalty score based on the difference between the expected number of meetings and the actual number of meetings
	penalty := 0.0
	for _, numMeetings := range participantMeetings {
		delta := expectedNumMeetings - float64(numMeetings)
		// fmt.Println("Expected Num Meetings: ", expectedNumMeetings, "Num Meetings: ", numMeetings, "Delta: ", delta)
		penalty += math.Abs(delta * delta * delta)
		// fmt.Println("Penalty: ", penalty)
	}

	return penalty
}

func GetNeighborhood(combos []ParticipantsCombination, allCombos map[string]ParticipantsCombination, numParticipants int8) [][]ParticipantsCombination {
	neighborhood := make([][]ParticipantsCombination, 0)

	for comboInd, combo := range combos {

		// modify each participant down by 1
		for participantInd, participant := range combo.Participants {
			comboCopy := combo
			comboCopy.Participants = make([]int8, len(combo.Participants))
			copier.Copy(&comboCopy.Participants, &combo.Participants)

			// fmt.Println("Combo Copy: ", comboCopy)

			// modify the participant
			modifiedId := (participant - 1 + numParticipants) % numParticipants
			comboCopy.Participants[participantInd] = modifiedId

			// check to see if the modified participant is in allCombos
			if _, ok := allCombos[combo.GetKey()]; !ok {
				continue
			}
			if len(comboCopy.Participants) == 0 {
				fmt.Println("Zero length combocopy", comboCopy.Participants)
			}

			neighborhoodEntry := make([]ParticipantsCombination, 34)
			copier.Copy(&neighborhoodEntry, &combos)

			neighborhoodEntry[comboInd] = comboCopy
			neighborhood = append(neighborhood, neighborhoodEntry)
		}

		// down by half
		for participantInd, participant := range combo.Participants {
			comboCopy := combo
			comboCopy.Participants = make([]int8, len(combo.Participants))
			copier.Copy(&comboCopy.Participants, &combo.Participants)

			// modify the participant
			modifiedId := (participant - 3 + numParticipants) % numParticipants
			comboCopy.Participants[participantInd] = modifiedId

			// check to see if the modified participant is in allCombos
			if _, ok := allCombos[combo.GetKey()]; !ok {
				continue
			}
			neighborhoodEntry := make([]ParticipantsCombination, 34)
			copier.Copy(&neighborhoodEntry, &combos)
			neighborhoodEntry[comboInd] = comboCopy
			neighborhood = append(neighborhood, neighborhoodEntry)
		}

		// modify each participant up by 1
		for participantInd, participant := range combo.Participants {
			comboCopy := combo
			comboCopy.Participants = make([]int8, len(combo.Participants))
			copier.Copy(&comboCopy.Participants, &combo.Participants)

			// modify the participant
			modifiedId := (participant + 1) % numParticipants
			comboCopy.Participants[participantInd] = modifiedId

			// check to see if the modified participant is in allCombos
			if _, ok := allCombos[combo.GetKey()]; !ok {
				continue
			}

			if len(comboCopy.Participants) == 0 {
				fmt.Println("Zero length combocopy", comboCopy.Participants)
			}

			neighborhoodEntry := make([]ParticipantsCombination, 34)
			copier.Copy(&neighborhoodEntry, &combos)
			neighborhoodEntry[comboInd] = comboCopy
			neighborhood = append(neighborhood, neighborhoodEntry)
		}

		// modify each participant up by half
		for participantInd, participant := range combo.Participants {
			comboCopy := combo
			comboCopy.Participants = make([]int8, len(combo.Participants))
			copier.Copy(&comboCopy.Participants, &combo.Participants)

			// modify the participant
			modifiedId := (participant + 7) % numParticipants
			comboCopy.Participants[participantInd] = modifiedId

			// check to see if the modified participant is in allCombos
			if _, ok := allCombos[combo.GetKey()]; !ok {
				continue
			}
			neighborhoodEntry := make([]ParticipantsCombination, 34)
			copier.Copy(&neighborhoodEntry, &combos)
			neighborhoodEntry[comboInd] = comboCopy
			neighborhood = append(neighborhood, neighborhoodEntry)
		}
	}
	return neighborhood
}

func GetSolutionKey(solution []ParticipantsCombination) string {
	keys := make([]string, len(solution))
	for i, combo := range solution {
		keys[i] = combo.GetKey()
	}
	return strings.Join(keys, "---")
}

func GetSolutionDistribution(solution []ParticipantsCombination, numParticipants int8) []int8 {
	participantMeetings := make([]int8, numParticipants)
	for _, combo := range solution {
		for _, participant := range combo.Participants {
			participantMeetings[participant]++
		}
	}
	return participantMeetings
}

func DisplaySolution(solution []ParticipantsCombination, numParticipants int8, title string) string {

	t := table.NewWriter()
	t.SetCaption("Suggested Schedule (Tabu Search)")

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

	return t.RenderCSV()
}
