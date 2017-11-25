package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func permunationDFAConcurrent(numState int) []string {
	sizeRow := numState
	sizeCol := 2 * (numState - 1)

	max := int(math.Pow(float64(sizeRow), float64(sizeCol))) - 1
	DFA := make([]string, 0)

	dfas := make(chan []string)
	// separate jobs to this number of pieces, and calculate them concurrently.
	concurrentLimit := 20
	steps := max/concurrentLimit + 1
	for j := 0; j < concurrentLimit; j++ {
		go func(start, end int) {
			results := make([]string, 0)
			for i := start; i <= end && i <= max; i++ {
				s, valid := getDFA(i, numState)
				if valid {
					results = append(results, s)
				}
			}
			dfas <- results
		}(j*steps, (j+1)*steps-1)
	}
	for j := 0; j < concurrentLimit; j++ {
		s := <-dfas
		DFA = append(DFA, s...)
	}
	return DFA
}

func permunationDFA(numState int) []string {
	sizeRow := numState
	sizeCol := 2 * (numState - 1)

	max := int(math.Pow(float64(sizeRow), float64(sizeCol))) - 1
	DFA := make([]string, 0)

	for i := 0; i < max; i++ {
		s, valid := getDFA(i, numState)
		if valid {
			DFA = append(DFA, s)
		}
	}
	return DFA
}

func getDFA(seed int, numState int) (string, bool) {
	sizeCol := 2 * (numState - 1)

	s := strconv.FormatUint(uint64(seed), numState)
	for len(s) < sizeCol {
		s = "0" + s
	}
	maxDigit := strconv.Itoa(numState - 1) // biggest digit
	digits := strings.Split(s, "")
	valid := false
	for _, d := range digits {
		if d == maxDigit {
			valid = true
			break
		}
	}
	if !valid {
		return "", false
	}

	// check connectness
	set := make([]string, 0)
	stack := make([]string, 0)
	stack = append(stack, "0")
	for len(stack) > 0 && len(set) < numState {
		currentState := stack[0]
		stack = stack[1:]
		if !inArray(set, currentState) {
			set = append(set, currentState)
			if currentState == maxDigit {
				// end state has no out going edge
				continue
			}
			currentStateIndex, _ := strconv.Atoi(currentState)
			dest0 := digits[currentStateIndex]
			if !inArray(set, dest0) {
				stack = append(stack, dest0)
			}
			dest1 := digits[currentStateIndex+numState-1]
			if !inArray(set, dest1) {
				stack = append(stack, dest1)
			}
		}
	}
	if len(set) == numState {
		return s, true
	}

	return "", false
}

func inArray(set []string, s string) bool {
	for _, d := range set {
		if d == s {
			return true
		}
	}
	return false
}

func getAllStrings(maxLength int) []string {
	results := make([]string, 0)
	for i := 0; i <= maxLength; i++ {
		s := generateStrings(i)
		results = append(results, s...)
	}
	return results
}

func checkSet([][]string) {

}

func isAccepted(dfa string, input string) bool {
	inputs := strings.Split(input, "")
	states := strings.Split(dfa, "")
	numStates := len(states) / 2
	currentState := 0
	for _, input := range inputs {
		var nextStateString string
		if currentState == numStates {
			return false
		}
		if input == "0" {
			// 0 -> first half
			nextStateString = states[currentState]
		} else {
			// 1 -> second half
			nextStateString = states[numStates+currentState]
		}
		currentState, _ = strconv.Atoi(nextStateString)
	}
	return currentState == numStates
}

// check all inputs's result against the DFA, return array of strings of 0 or 1
func checkDFA(dfa string, inputs []string) []string {
	results := make([]string, 0)
	for _, a := range inputs {
		accepted := isAccepted(dfa, a)
		if accepted {
			results = append(results, "1")
		} else {
			results = append(results, "0")
		}
	}
	return results
}

func checkAllDfa(DFAs []string, inputs []string) [][]string {
	size := len(DFAs)
	allResults := make([][]string, 0)
	resultChan := make(chan [][]string)
	// separate jobs to this number of pieces, and calculate them concurrently.
	concurrentLimit := 20
	steps := size/concurrentLimit + 1
	for j := 0; j < concurrentLimit; j++ {
		go func(start, end int) {
			results := make([][]string, 0)
			for i := start; i <= end && i < size; i++ {
				dfa := DFAs[i]
				result := checkDFA(dfa, inputs)
				result = append([]string{dfa}, result...)
				results = append(results, result)
			}
			resultChan <- results
		}(j*steps, (j+1)*steps-1)
	}
	for j := 0; j < concurrentLimit; j++ {
		s := <-resultChan
		allResults = append(allResults, s...)
	}
	return allResults
}

func getUniqueSetDFA(DFACheckedResults [][]string) map[string][]string {
	solutionMap := make(map[string][]string)
	for _, data := range DFACheckedResults {
		dfa := data[0]
		dfaResult := data[1:]
		resultId := strings.Join(dfaResult, "")
		_, ok := solutionMap[resultId]
		if !ok {
			solutionMap[resultId] = make([]string, 0)
		}
		solutionMap[resultId] = append(solutionMap[resultId], dfa)
	}
	return solutionMap
}

func main() {
	// DFAs := permunationDFAConcurrent(2)
	// fmt.Println(DFAs)

	solutionsMap := make(map[string][]string)

	allStrings := getAllStrings(3)
	sets := generateSets(3)
	for _, set := range sets {
		ids := make([]string, len(allStrings))
		for i, s := range allStrings {
			if inArray(set, s) {
				ids[i] = "1"
			} else {
				ids[i] = "0"
			}
		}
		solution := strings.Join(ids, "")
		solutionsMap[solution] = make([]string, 0)
		// fmt.Println(strings.Join(ids, ""))
	}
	// fmt.Printf("\ntotal: %d\n", len(sets))
	// return

	for max := 2; max <= 6; max++ {

		allStrings := getAllStrings(3)
		// fmt.Println(allStrings)

		start := time.Now()
		DFAs := permunationDFAConcurrent(max)
		fmt.Println(time.Since(start))

		// fmt.Println(len(DFAs))
		// for _, d := range DFAs {
		// 	fmt.Println(d)
		// }

		start = time.Now()
		results := checkAllDfa(DFAs, allStrings)
		fmt.Println(time.Since(start))

		// for _, r := range results {
		// 	fmt.Println(strings.Join(r, ","))
		// }
		start = time.Now() //print all solution with 6 states
		solutionMap := getUniqueSetDFA(results)
		fmt.Println(time.Since(start))
		fmt.Println(strings.Join(allStrings, ", "))
		for solution, dfas := range solutionMap {
			solutionsMap[solution] = append(solutionsMap[solution], dfas...)
		}
	}

	for solution, dfas := range solutionsMap {
		if len(dfas) == 0 {
			continue
		}
		l := len(dfas)
		if l > 5 {
			l = 5
		}
		fmt.Printf("%s, %d, \t %s\n", solution, len(dfas), strings.Join(dfas[:l], ", "))
	}

	fmt.Printf("\n\n")
	for solution, dfas := range solutionsMap {
		if len(dfas) == 0 {
			fmt.Println(solution)
		}
	}
}
