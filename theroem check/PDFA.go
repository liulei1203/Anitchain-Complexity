// Generate possible DFAs, then check any of them is for the anti-chain sets. Since we generate DFA
//from small number of the states to the large number ones, then the minimal DFA of the
//anti-chain is the very first one to be found.

package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

var stringsL3 = []string{
	"0",
	"1",
	"00",
	"01",
	"10",
	"11",
}

var strings3 = []string{
	"000",
	"001",
	"010",
	"011",
	"100",
	"101",
	"110",
	"111",
}

// getStringCombination generates antichian sets
func getStringCombinations(elements []string) [][]string {
	results := make([][]string, 0)
	for i := 0; i < len(elements); i++ {
		a := []string{elements[i]}
		results = append(results, a)

		otherCombo := getStringCombinations(elements[i+1:])
		for _, c := range otherCombo {
			newCombo := []string{elements[i]}
			newCombo = append(newCombo, c...)
			results = append(results, newCombo)
		}
		// for j := i + 1; j < len(elements); j++ {
		// 	a = append(a, elements[j])
		// 	results = append(results, a)
		// 	//fmt.Println(strings.Join(a, ", "))
		// }
	}
	//fmt.Printf("\ntotalCom: %d\n", len(results))
	return results
}

func generateAndCheckDFA(numState int, inputs []string) chan *DfaResult {

	resultChan := make(chan *DfaResult)

	sizeRow := numState
	sizeCol := 2 * (numState - 1) //final state is sink
	max := int(math.Pow(float64(sizeRow), float64(sizeCol))) - 1

	//parallel computing
	workers := 15
	var wg sync.WaitGroup
	wg.Add(workers)

	steps := max/workers + 1
	for j := 0; j < workers; j++ {
		go func(start, end int) {
			for i := start; i <= end && i <= max; i++ {
				dfa, valid := getDFA(i, numState)
				if valid {
					result := checkDFA(dfa, inputs)
					dfaResult := DfaResult{
						DFA:    string(dfa),
						Result: result,
					}
					resultChan <- &dfaResult
				}
			}
			wg.Done()
		}(j*steps, (j+1)*steps-1)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()
	return resultChan
}

// func generateDFAs(numState int) chan string {
// 	// Init channel
// 	checkChan := make(chan string)
//
// 	sizeRow := numState
// 	sizeCol := 2 * (numState - 1)
// 	max := int(math.Pow(float64(sizeRow), float64(sizeCol))) - 1
//
// 	go func() {
// 		for i := 0; i < max; i++ {
// 			dfa, valid := getDFA(i, numState)
// 			if valid {
// 				checkChan <- dfa
// 			}
// 		}
// 		close(checkChan)
// 	}()
// 	return checkChan
// }
//////////////////////////////////////////////////////////////

// getDFA gives all the transition functions
func getDFA(seed int, numState int) ([]byte, bool) {
	sizeCol := 2 * (numState - 1)

	// s := strconv.FormatUint(uint64(seed), numState)
	// s := make([]byte, 0, sizeCol)
	s := strconv.AppendUint(nil, uint64(seed), numState)
	// fmt.Println(string(s))
	if len(s) < sizeCol {
		diff := sizeCol - len(s)
		zeros := make([]byte, diff)
		for i := 0; i < diff; i++ {
			zeros[i] = '0'
		}
		s = append(zeros, s...)
		// for len(s) < sizeCol {
		// 	s = append([]byte{'0'}, s...)
		// }
	}
	// maxDigit := strconv.Itoa(numState - 1) // biggest digit
	// digits := strings.Split(s, "")
	maxDigit := byte('0' + numState - 1)
	digits := s

	valid := false
	for _, d := range digits {
		if d == maxDigit {
			valid = true
			break
		}
	}
	if !valid {
		return []byte{}, false
	}
	// check len(string) < 3
	for _, a := range stringsL3 {
		if isAccepted(s, a) {
			return []byte{}, false
		}
	}

	// check connectness
	set := make([]byte, 0, numState)
	stack := make([]byte, 0, numState)
	stack = append(stack, '0')
	for len(stack) > 0 && len(set) < numState {
		currentState := stack[0]
		stack = stack[1:]
		if !inByteArray(set, currentState) {
			set = append(set, currentState)
			if currentState == maxDigit {
				// end state has no out going edge
				continue
			}
			// currentStateIndex, _ := strconv.Atoi(currentState)
			currentStateIndex := int(currentState - '0')
			dest0 := digits[currentStateIndex]
			if !inByteArray(set, dest0) {
				stack = append(stack, dest0)
			}
			dest1 := digits[currentStateIndex+numState-1]
			if !inByteArray(set, dest1) {
				stack = append(stack, dest1)
			}
		}
	}
	if len(set) == numState {
		return s, true
	}

	return []byte{}, false
}

///////////////////////////////////////////////////////

func inByteArray(set []byte, s byte) bool {
	for _, d := range set {
		if d == s {
			return true
		}
	}
	return false
}

func inArray(set []string, s string) bool {
	for _, d := range set {
		if d == s {
			return true
		}
	}
	return false
}

//func getAllStrings(maxLength int) []string {
//	results := make([]string, 0)
//	for i := 0; i <= maxLength; i++ {
//		s := generateStrings(i)
//		results = append(results, s...)
//	}
//	return results
//}

func checkSet([][]string) {

}

func isAccepted(states []byte, inputString string) bool {
	// inputs := strings.Split(input, "")
	// states := strings.Split(dfa, "")
	numStates := len(states) / 2
	currentState := 0
	for _, input := range inputString {
		var nextStateString byte
		if currentState == numStates {
			return false
		}
		if input == '0' {
			// 0 -> first half
			nextStateString = states[currentState]
		} else {
			// 1 -> second half
			nextStateString = states[numStates+currentState]
		}
		currentState = int(nextStateString - '0')
	}
	return currentState == numStates
}

// check all inputs's result against the DFA, return array of strings of 0 or 1
// results stores antichain Set which accept by DFA
func checkDFA(dfa []byte, inputs []string) string {

	var result bytes.Buffer
	// dfaStates := strings.Split(dfa, "")
	for _, a := range inputs {
		accepted := isAccepted(dfa, a)
		if accepted {
			result.WriteByte('1')
		} else {
			result.WriteByte('0')
		}
	}
	return result.String()
}

type DfaResult struct {
	DFA    string
	Result string
}

// read DFA from checkChan, then calculate the accepted inputs
// func processDFA(inputs []string, checkChan chan string, resultChan chan *DfaResult) {
// 	for dfa := range checkChan {
// 		result := checkDFA(dfa, inputs)
// 		dfaResult := DfaResult{
// 			DFA:    dfa,
// 			Result: result,
// 		}
// 		resultChan <- &dfaResult
// 	}
// }

func main() {
	// defer profile.Start().Stop()
	// f, _ := os.Create("cpu.prof")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	// args := os.Args[1:]
	//inputLength, _ := strconv.Atoi(args[0])
	dfaMax := 7 // max number of states

	//allStrings := getAllStrings(inputLength)

	solutionsMap := make(map[string]string) //?

	sets := getStringCombinations(strings3)
	//generateSets(inputLength) // the max Length of strings
	for _, set := range sets {
		solution := "" //ID for cantichain set
		for _, s := range strings3 {
			if inArray(set, s) {
				solution = solution + "1"
			} else {
				solution = solution + "0"
			}
		}
		solutionsMap[solution] = ""
	}
	fmt.Printf("total set %d\n", len(solutionsMap))
	for max := 4; max <= dfaMax; max++ {
		// max starts at 2 in general, but A(Z) >=4 for Z in sets
		start := time.Now()

		// Start generate DFA
		// checkChan := generateDFAs(max)n
		//
		// // Start checkDFA worker
		// workers := 7
		// var wg sync.WaitGroup
		// wg.Add(workers)
		// resultChan := make(chan *DfaResult)
		// for i := 0; i < workers; i++ {
		// 	go func() {
		// 		processDFA(allStrings, checkChan, resultChan)
		// 		wg.Done()
		// 	}()
		// }
		// go func() {
		// 	wg.Wait()
		// 	close(resultChan)
		// }()

		resultChan := generateAndCheckDFA(max, strings3)

		// Check Results
		for dfaResult := range resultChan {
			if solutionsMap[dfaResult.Result] == "" {
				solutionsMap[dfaResult.Result] = dfaResult.DFA
			}
		}

		// Finish current length of DFA
		solvedCount := 0
		for _, dfa := range solutionsMap {
			if dfa != "" {
				solvedCount++
			}
		}
		fmt.Printf("[%d] solved %d/%d, elapsed %s\n", max, solvedCount, len(solutionsMap), time.Since(start))

		// Memory profile
		// mf, _ := os.Create("mem.prof")
		// runtime.GC() // get up-to-date statistics
		// pprof.WriteHeapProfile(mf)
		// mf.Close()
	}

	// for solution, dfa := range solutionsMap {
	// 	if dfa != "" {
	// 		fmt.Printf("%s,\t %s\n", solution, dfa)
	// 	}
	// }

	fmt.Printf("\n\n")
	solvedCount := 0
	for _, dfa := range solutionsMap {
		if dfa != "" {
			solvedCount++
		}
	}
	fmt.Printf("Total set %d, solved %d/%d\n", len(solutionsMap), solvedCount, len(solutionsMap))

	for set, dfa := range solutionsMap {
		fmt.Printf(" %q\t%q\n", set, dfa)
	}

}
