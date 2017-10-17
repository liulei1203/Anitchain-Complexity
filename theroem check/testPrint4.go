//generate all sets which contain strings with length 3
// find the sets with complexity 7
// check conditions for those sets

package main

type chain struct {
	strings   []int
	maxLength int
}

var strings4 = []string{
	"0000",
	"0001",
	"0010",
	"0011",
	"0100",
	"0101",
	"0110",
	"0111",
	"1000",
	"1001",
	"1010",
	"1011",
	"1100",
	"1101",
	"1110",
	"1111",
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

//func main() {

//allSets := getStringCombinations(strings3)
//count := 0
//for _, set := range allSets {

//	fmt.Println(strings.Join(set, ", "))
//	if len(set) == 2 {
//count = count + 1
//	}

//}

//fmt.Printf("\ntotalCom: %d\n", len(allSets))
//fmt.Printf("\ncount2: %d\n", count)

//}
