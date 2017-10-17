package main

import (
	"strconv"
	"strings"
)

// Patial DFA, final state is sink and unique final state
type State struct {
	Act0    *State
	Act1    *State
	Name    string
	Dist    int // record distance to start
	Parents []string
	// PreState *State // previous state for shortest path current state to Start
	// From0 *State
	// From1 *State
}

func (state *State) String() string {
	s := "[" + state.Name + "]"
	if state.Act0 != nil {
		s = s + " 0->(" + state.Act0.Name + ")"
	}
	if state.Act1 != nil {
		s = s + " 1->(" + state.Act1.Name + ")"
	}
	return "<" + s + ">"
}

func (state *State) setOutput(s string, output *State) {
	if s == "0" {
		state.Act0 = output
	}
	if s == "1" {
		state.Act1 = output
	}
	output.Parents = append(output.Parents, state.Name)
}
func (state *State) getOutput(s string) *State {
	if s == "0" {
		return state.Act0
	}
	if s == "1" {
		return state.Act1
	}
	return nil
}

// is state equivalent to s
func (state *State) IsEqual(s *State) bool {
	if state.Act0 == s.Act0 && state.Act1 == s.Act1 {
		return true
	}
	if state.Act0 == state && s.Act0 == s && state.Act1 == s.Act1 {
		return true
	}
	if state.Act1 == state && s.Act1 == s && state.Act0 == s.Act0 {
		return true
	}
	return false
}

//generate DFA
func generateDFA(set []string, maxLength int) map[string]*State {
	//output: start state, middle states, finish state
	empty := &State{Name: "E"}
	empty.Act0 = empty
	empty.Act1 = empty

	start := &State{Act0: empty, Act1: empty, Name: "S", Parents: []string{}}
	finish := &State{Name: "F", Parents: []string{}}
	// states := make([]*State, 0)
	stateMap := make(map[string]*State)
	stateMap[start.Name] = start
	stateMap[finish.Name] = finish

	// [1] generate raw DFA
	for _, s := range set {
		x := strings.Split(s, "")
		curState := start
		for i := 0; i < len(x); i++ {
			if i == len(x)-1 {
				curState.setOutput(x[i], finish)
			} else {
				nextState := curState.getOutput(x[i])
				if nextState == nil || nextState == empty {
					newState := &State{Act0: empty, Act1: empty, Parents: []string{}} // create a new state
					newState.Name = strconv.Itoa(len(stateMap) - 1)                   // find a name
					stateMap[newState.Name] = newState                                // add to map
					curState.setOutput(x[i], newState)                                // link new state
					newState.Dist = curState.Dist + 1
					curState = newState
				} else {
					// update distance
					if nextState.Dist > curState.Dist+1 {
						nextState.Dist = curState.Dist + 1
					}
					curState = nextState
				}
			}
		}
	}

	reduceStates(stateMap, maxLength, empty)

	checkEmpty(stateMap, empty)

	return stateMap
	// // [end] convert map to array
	// states := make([]*State, 0)
	// for _, state := range stateMap {
	// 	states = append(states, state)
	// }
	//
	// return start, states, finish
}

// DOES NOT WORK !!!
func reverseGenerateDFA(set []string, maxLength int) map[string]*State {
	// DOES NOT WORK !!!
	//output: start state, middle states, finish state
	empty := &State{Name: "E"}
	empty.Act0 = empty
	empty.Act1 = empty

	start := &State{Act0: empty, Act1: empty, Name: "S", Parents: []string{}}
	finish := &State{Name: "F", Parents: []string{}}
	// states := make([]*State, 0)
	stateMap := make(map[string]*State)
	stateMap[start.Name] = start
	stateMap[finish.Name] = finish

	// [1] generate raw DFA
	for _, s := range set {
		x := strings.Split(s, "")
		curState := finish
		for i := len(x) - 1; i >= 0; i-- {
			if i == 0 {
				start.setOutput(x[i], curState)
				break
			}
			if curState == start {
				start.setOutput(x[i], start)
				break
			}
			found := false
			for _, p := range curState.Parents {
				s := stateMap[p]
				if s.getOutput(x[i]) == curState {
					curState = s
					found = true
					break
				}
			}
			if !found {
				newState := &State{Act0: empty, Act1: empty, Parents: []string{}}
				newState.Name = strconv.Itoa(len(stateMap) - 1)
				stateMap[newState.Name] = newState
				newState.setOutput(x[i], curState)
				curState = newState
			}
		}
	}

	findDeadStates(stateMap)

	reduceStates(stateMap, maxLength, empty)

	checkEmpty(stateMap, empty)

	return stateMap
}

func findDeadStates(stateMap map[string]*State) {
	visited := make(map[string]bool)
	start := stateMap["S"]
	stateQueue := make([]*State, 0)
	stateQueue = append(stateQueue, start)
	for len(stateQueue) > 0 {
		current := stateQueue[0]
		visited[current.Name] = true
		if current.Act0 != nil && !visited[current.Act0.Name] {
			stateQueue = append(stateQueue, current.Act0)
		}
		if current.Act1 != nil && !visited[current.Act1.Name] {
			stateQueue = append(stateQueue, current.Act1)
		}
		stateQueue = stateQueue[1:]
	}
	for name, _ := range stateMap {
		if !visited[name] {
			delete(stateMap, name)
		}
	}
}

func updateParents(stateMap map[string]*State) {
	// reset parents
	for _, s := range stateMap {
		s.Parents = make([]string, 0)
	}

	for _, s := range stateMap {
		if s.Act0 != nil {
			s.Act0.Parents = append(s.Act0.Parents, s.Name)
		}
		if s.Act1 != nil {
			s.Act1.Parents = append(s.Act1.Parents, s.Name)
		}
	}

}

func reduceStates(stateMap map[string]*State, maxLength int, empty *State) {
	finish := stateMap["F"]
	// [2] minimize
	// [2.1] Remove links to empty state
	// find distance >= 2, change link to empty state to self loop
	for _, state := range stateMap {
		if state.Dist >= maxLength-1 {
			if state.Act0 == empty {
				state.Act0 = state
			}
			if state.Act1 == empty {
				state.Act1 = state
			}
		}
	}
	// [2.2] find duplicate states, start from F
	queue := make([]string, len(finish.Parents))
	copy(queue, finish.Parents)
	curName := ""
	for len(queue) > 0 {
		// pop from queue
		curName = queue[0]
		queue = queue[1:]
		// add parents to queue
		curState, ok := stateMap[curName]
		if ok {
			for _, p := range curState.Parents {
				if p != curName {
					queue = append(queue, p)
				}
			}
			if len(queue) > 0 && queue[0] != curName {
				nextName := queue[0]
				nextState, ok := stateMap[nextName]
				if ok {
					// test equal
					if curState.IsEqual(nextState) {
						// merge
						// change curState's parent's link to curState to nextState
						for _, name := range curState.Parents {
							s, ok := stateMap[name]
							if ok {
								if s.Act0 == curState {
									s.Act0 = nextState
								}
								if s.Act1 == curState {
									s.Act1 = nextState
								}
							}
						}
						// delete curState
						delete(stateMap, curName)
						// update parents
						updateParents(stateMap)
					}
				}
			}
		}
	}
}

func checkEmpty(stateMap map[string]*State, empty *State) {
	for _, s := range stateMap {
		if s.Act0 == empty {
			stateMap[empty.Name] = empty
			return
		}
		if s.Act1 == empty {
			stateMap[empty.Name] = empty
			return
		}
	}
}

// func main() {
// 	testSets := [][]string{
// 		[]string{"0", "10", "110"},
// 		[]string{"000"},
// 		[]string{"11", "000", "010"},
// 		[]string{"11", "011", "100", "101"},
// 	}
// 	for _, testSet := range testSets {
// 		fmt.Println(testSet)
// 		states := generateDFA(testSet, 3)
// 		fmt.Printf("complexity=%d, %s\n", len(states), states)
// 		states = reverseGenerateDFA(testSet, 3)
// 		fmt.Printf("complexity=%d, %s\n", len(states), states)
// 	}
// }

//========================================================== first draft ======
// minimal partial DFA
// func miniDFA(set []string) (*State, []*State, *State) {
// 	start, states, finish := generateDFA(set)
// 	fmt.Printf("%d\t", len(states))
// 	empty := 0
// 	if start.Act0.Name == "E" && len(set[0]) > 1 {
// 		start.Act0.Act0 = start
// 		start.Act0.Act1 = start
// 		states = append(states, start.Act0)
// 		empty = 1
// 		fmt.Printf("%d, %d\t", len(states), empty)
// 	} else if start.Act1.Name == "E" && len(set[0]) > 1 {
// 		start.Act1.Act0 = start
// 		start.Act1.Act1 = start
// 		states = append(states, start.Act1)
// 		empty = 1
// 		fmt.Printf("%d, %d\t", len(states), empty)
// 	} else {
//
// 		flag := false
// 		for i := len(states) - 1; i >= 0; i-- {
//
// 			if states[i].Act0.Name == "E" && flag == false {
// 				//states[i].Act0.Act0 = states[i]
// 				states[i].Act0.Act1 = states[i]
// 				states = append(states, states[i].Act0)
// 				flag = true
// 				empty = 1
// 			} else if states[i].Act1.Name == "E" && flag == false {
// 				states[i].Act1.Act0 = states[i]
// 				states = append(states, states[i].Act1)
//
// 				//states[i].Act1.Act1 = states[i]
// 				flag = true
// 				empty = 1
// 			}
// 		}
// 	}
// 	// if distantce S to E > 2, then E can be deleted.
// 	if empty == 1 {
// 		curState := states[len(states)-1]
// 		if curState.DistToSt > 2 {
// 			//remove "E"
// 			preState := curState.PreState
// 			if preState.Act0 == curState {
// 				preState.Act0 = preState
// 			}
// 			preState.Act1 = preState
// 		}
//
// 		//# of states -1
//
// 	}
//
// 	// find homo group, state P ~ Q if (P, 0)= (Q,0) and (P,1) = (Q,1)
//
// 	if empty == 0 {
// 		fmt.Printf("%d, %d\t", len(states), empty)
// 		return start, states, finish
//
// 	}
//
// 	return start, states, finish
// }
//===============================================================================================
