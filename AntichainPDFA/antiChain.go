package main

import (
	"fmt"
	"sort"
	"strings"
)

type chain struct {
	strings   []int
	maxLength int
}

// Generate array of 01 strings with the given length
func generateStrings(length int) []string {
	if length <= 0 {
		return []string{}
	}
	if length == 1 {
		return []string{"0", "1"}
	}
	previousStrings := generateStrings(length - 1)
	newStrings := make([]string, 0)
	for _, str := range previousStrings {
		newStrings = append(newStrings, str+"0", str+"1")
		//print
	}

	return newStrings
}

// given an array of string, make all combinations of them
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

// from the elements get the valid ones according to the given set (no prefix)
func getValidStrings(set []string, elements []string) []string {
	toDelete := make([]int, len(elements))
	for i, e := range elements {
		for _, s := range set {
			if strings.HasPrefix(e, s) {
				toDelete[i] = 1
				break
			}
		}
	}
	newArray := make([]string, 0)
	for i, e := range elements {
		if toDelete[i] == 0 {
			newArray = append(newArray, e)
		}
	}
	return newArray
}

// generate the antiChain sets
func generateSets(maxLength int) [][]string {
	if maxLength <= 0 {
		return [][]string{}
	}
	previousSets := generateSets(maxLength - 1)
	newStrings := generateStrings(maxLength)
	// [1] previous set + combinations of new elements
	newSets := append(previousSets, getStringCombinations(newStrings)...)
	// [2] extend previous set with new elements
	for _, set := range previousSets {
		validStrings := getValidStrings(set, newStrings)
		validStringCombo := getStringCombinations(validStrings)
		for _, newCombo := range validStringCombo {
			// need to make more explicit copy. got problem of new newSet overwrite the previous newSet
			newSet := make([]string, 0)
			newSet = append(newSet, set...)
			newSet = append(newSet, newCombo...)
			newSets = append(newSets, newSet)
		}
	}
	return newSets
}

// func main() {
// 	args := os.Args[1:]
// 	length, err := strconv.Atoi(args[0])
// 	if err != nil {
// 		panic(err)
// 		length = 4
// 	}
//
// 	// max length elements
// 	elements := generateStrings(length)
// 	fmt.Println(strings.Join(elements, ", "))
//
// 	// max length combinations
// 	// combinations := getStringCombinations(elements)
// 	// printSets(combinations)
// 	// fmt.Printf("combinations: %d\n", len(combinations))
//
// 	// all sets final antiChain
// 	// sets := generateSets(length)
// 	// printSets(sets, length)
// 	// fmt.Printf("\ntotal: %d\n", len(sets))
//
// }

//minimal DFA

//str's conjugation
func conjStr(str string) string {
	s := strings.Split(str, "")
	m := len(s)
	conS := make([]string, m)
	for i := 0; i < m; i++ {
		if s[i] == "0" {
			conS[i] = "1"
		} else {
			conS[i] = "0"
		}
	}
	cStr := strings.Join(conS, "")
	return cStr
}

//sort then joint strings to form a id
func id(set []string) string {
	sort.Strings(set)
	id := strings.Join(set, "-")
	return id
}

// find the self-symetric set
func selfSym(set []string) bool {
	conjSet := make([]string, 0)

	for _, str := range set {
		conj := conjStr(str)
		conjSet = append(conjSet, conj)
	}
	// compare congSet and set's ids
	if id(set) == id(conjSet) {
		return true //set is self-symetric
	} else {
		return false
	}

}

func printSets(sets [][]string, length int) {
	count := 0
	for _, set := range sets {
		states := generateDFA(set, length)
		if len(states) < 8 {
			continue
		}
		count = count + 1
		fmt.Println(strings.Join(set, ", "))
		fmt.Printf("complexity=%d, %s\n\n", len(states), states)
		// if selfSym(set) == true {
		// 	fmt.Printf("\n self-symetric set \n")
		// }
	}
	fmt.Println(count)
}

// func (anti *chain) generateStr(maxLength int) []strings {
// 	numEle := 0
// 	//for i := 0; i< maxLength; i++ {
// 		//numEle = numEle + int(Exp2(i+1))
// 	//}
// 	base[]//var strings []int
// 	for i := 0; i< maxLength; i++{
// 		num0 = numEle
// 		numEle = numEle + int(Exp2(i+1))
// 		base[j] = base[j-1] apending 0
// 		s
// 		for j := 1; j< numEle; i++{
// 			strings[num0+j] =
// 		}
// 	}
// 	return b
// }
//
// func updateWeight(weight []float64, sample []float64, c float64) []float64 {
// 	W := make([]float64, len(weight))
// 	for i := 0; i < len(weight); i++ {
// 		W[i] = weight[i] + c*sample[i]
// 	}
// 	return W
// }
