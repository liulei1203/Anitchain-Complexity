package main

import (
	"fmt"
	"strings"
	//	"bytes"

	//	"math"

	"os"
	//	"runtime"
	//	"runtime/pprof"
	//	"strconv"
	//	"sync"
	//	"time"
	"bufio"
	"log"
)

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

// find all set which's complexity is 7
func getMaxComSet(srow string) (bool, []string) {
	set := make([]string, 0)
	result := strings.Split(srow, ",")
	r := false
	if strings.Contains(result[1], "6") {
		//start state is 0, then 6 is final state for DFA with complexity 7
		maxComString := result[0]
		set = map2Set(maxComString)
		r = true
	}
	return r, set
}

// map back to set
func map2Set(s string) []string {
	set := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			set = append(set, strings3[i])
			//fmt.Print(strings3[i] + ", ")
		}
	}
	return set
	//%fmt.Println()

}

//check 3-conditions
func check(set []string) bool {
	flag1 := false // condition 1: exist 2 strings with different 1-digit prefixes
	flag2 := false // condition 2: exist 2 strings with commom 2-digit prefix
	flag3 := false // condition 3: exist 2 strings with different 1-digit sufixes and
	// they are have no common 2-digit prefix

	f1 := false //indicator for break
	//pref1 1-digit prefix, eg: 001's pref1 is 0
	//pref2 2-digit prefix, eg: 011's pref1 is 01
	pref1 := set[0][0]
	pref2 := set[0][0:2]
	n := len(set)
	var currStr string // current string
	var currf2 string  /// current strind's 2-digit prefix,
	//var nexrf2 []byte
	newSet := make([]string, 0)
	for i := 1; i < n; i++ {
		currStr = set[i]

		//check condition 1: if exist one string which has a diferent 1-digit prefix
		//                  from pref1, then condition 1 is true
		if f1 == false && pref1 != currStr[0] {
			flag1 = true

			f1 = true // if condition 1 is true, it is not nessesary to check more strins
		}

		//check condition 2: we need to check whether the current string has common 2-digit
		// prefix with the previous one
		currf2 = currStr[0:2]          //string
		if i == 1 && pref2 != currf2 { //for set[0], as set[0] has no previous string
			newSet = append(newSet, set[i-1])
			//new set contains the strings which have no same 2-digit prefix as
			//another string in the set

		}

		if pref2 == currf2 { // if condition 2 is true
			flag2 = true
		} else {
			if i != n-1 {
				if currf2 != set[i+1][0:2] {
					newSet = append(newSet, set[i])
				}
			} else { // for the very last string as no next string
				newSet = append(newSet, set[i])
			}

		}
		pref2 = currf2[0:2]
	}
	// check condition 3
	if len(newSet) == 0 { // condition 3 is false as all strings are paired
		return false
	}

	suf1 := newSet[0][2] // 1-digit suffix byte, eg : 001's suf1 is 1

	f3 := false //indicator for break
	for i := 1; i < len(newSet); i++ {
		if f3 == false && suf1 != newSet[i][2] {
			flag3 = true
			f3 = true
		}
	}

	if flag1 && flag2 && flag3 {
		return true
	}
	return false

}

func main() {
	file, err := os.Open("./mapNew.csv") // path to file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	fmt.Println("Set , true/false ")
	for scanner.Scan() {
		r, set := getMaxComSet(scanner.Text())
		if r == true {
			count++
			flag := check(set)
			fmt.Printf("%v , %v\n", set, flag)
		}
		//if check(map2Set(result[0])){
	}

	fmt.Printf("\ncount =%d\n", count)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
