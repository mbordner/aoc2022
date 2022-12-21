package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/array"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
	"strings"
)

var (
	reInput = regexp.MustCompile(`^Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? (.*)`)
)

type Timeline string

const (
	chunkSize = 3
	openOp    = 'o'
	travelOp  = 't'
)

type ValveMap map[string]*Valve

func (vm ValveMap) clone() ValveMap {
	m := make(ValveMap)
	for k, v := range vm {
		m[k] = v
	}
	return m
}

func (vm ValveMap) totalRate() int {
	sum := 0
	for _, v := range vm {
		sum += v.rate
	}
	return sum
}

type Valve struct {
	id      string
	rate    int
	tunnels []string
	to      ValveMap
	from    ValveMap
}

// takes 1 min to move to cave
// takes 1 min to open value

// have 30 mins

func (t Timeline) travel(id string) Timeline {
	return Timeline(fmt.Sprintf("%s%s%s", t, string(travelOp), id))
}

func (t Timeline) open(id string) Timeline {
	return Timeline(fmt.Sprintf("%s%s%s", t, string(openOp), id))
}

func (t Timeline) time() int {
	return len(t) / chunkSize
}

func (t Timeline) location() string {
	if len(t) > 0 {
		return string(t[len(t)-2:])
	}
	return "AA"
}

func (t Timeline) getOpened(vm ValveMap, minsBefore int) ValveMap {
	ovm := make(ValveMap)
	for i := 0; i < len(t)-(minsBefore*chunkSize); i += chunkSize {
		if t[i] == openOp {
			openId := string(t[i+1 : i+3])
			ovm[openId] = vm[openId]
		}
	}
	return ovm
}

type memoValues struct {
	openMap ValveMap
	sum     int
}

var (
	memo = make(map[int]map[string]memoValues)
)

func (t Timeline) getReleasedPressure(vm ValveMap, time int) (int, ValveMap) {
	strVal := string(t)
	//fmt.Println("getting pressure for: ", strVal)

	if _, exists := memo[time]; !exists {
		memo[time] = make(map[string]memoValues)
	}

	if v, exists := memo[time][strVal]; exists {
		return v.sum, v.openMap.clone()
	}

	sum := 0
	timeVal := 0
	var openMap ValveMap
	var ops [][]rune

	if t.time() > 2 {
		prevTimeline := Timeline(strVal[0 : len(strVal)-chunkSize])
		curTimeline := Timeline(strVal[len(strVal)-chunkSize:])
		timeVal = prevTimeline.time()
		sum, openMap = prevTimeline.getReleasedPressure(vm, timeVal)
		ops = array.ChunkBy([]rune(curTimeline), chunkSize)
	} else {
		openMap = make(ValveMap)
		ops = array.ChunkBy([]rune(t), chunkSize)
	}

	for i := range ops {
		pressurePerMin := openMap.totalRate()
		sum += pressurePerMin
		op := ops[i][0]
		opTarget := string(ops[i][1:])
		if op == openOp {
			openMap[opTarget] = vm[opTarget]
		}
		timeVal++
		if timeVal == time {
			break
		}
	}
	pressurePerMin := openMap.totalRate()
	for timeVal < time {
		sum += pressurePerMin
		timeVal++
	}

	if t.time() > 3 {
		memo[time][strVal] = memoValues{
			openMap: openMap,
			sum:     sum,
		}
	}

	return sum, openMap.clone()
}

func (t Timeline) getCanOpen(vm ValveMap) ValveMap {
	ovm := t.getOpened(vm, 0)
	covm := make(ValveMap)
	for k, v := range vm {
		if _, exists := ovm[k]; !exists {
			if v.rate > 0 {
				covm[k] = v
			}
		}
	}
	return covm
}

func NewValveOpenTimeline() Timeline {
	return ""
}

func main() {

	valvesMap := getValves("../test.txt")

	/*
		tl := NewValveOpenTimeline()
		tl = tl.travel("DD")
		tl = tl.open("DD")
		tl = tl.travel("CC")
		tl = tl.travel("BB")
		tl = tl.open("BB")
		tl = tl.travel("AA")
		tl = tl.travel("II")
		tl = tl.travel("JJ")
		tl = tl.open("JJ")
		tl = tl.travel("II")
		tl = tl.travel("AA")
		tl = tl.travel("DD")
		tl = tl.travel("EE")
		tl = tl.travel("FF")
		tl = tl.travel("GG")
		tl = tl.travel("HH")
		tl = tl.open("HH")
		tl = tl.travel("GG")
		tl = tl.travel("FF")
		tl = tl.travel("EE")
		tl = tl.open("EE")
		tl = tl.travel("DD")
		tl = tl.travel("CC")
		tl = tl.open("CC")
		fmt.Println(tl.getReleasedPressure(valvesMap, 30))
		memo = make(map[int]map[string]memoValues)
	*/

	//tDDoDDtCCtBBoBB tAAtIItJJoJJtII tAAtDDtEEtFFtGG tHHoHHtGGtFFtEE oEEtDDtCCoCC
	pressure := bestPath(valvesMap, 30, NewValveOpenTimeline())
	fmt.Println(pressure)
}

func (v *Valve) getStartsFor(visited ValveMap, targets ValveMap, startGoalTargets ValveMap) {
	visited[v.id] = v

	for fId, fNode := range v.from {
		if _, exists := startGoalTargets[fId]; exists {
			targets[fId] = fNode
		} else {
			if _, exists := visited[fId]; !exists {
				fNode.getStartsFor(visited, targets, startGoalTargets)
			}
		}
	}
}

func (v *Valve) buildFromMap(visited ValveMap, targets ValveMap) {
	visited[v.id] = v

	for fId, fNode := range v.from {
		if _, exists := targets[fNode.id]; !exists {
			targets[fNode.id] = v
		} else {
			if fId != "" {
				fmt.Println(fId)
			}
		}
	}
}

func bestPath(valves ValveMap, maxTime int, cur Timeline) int {
	curTime := cur.time()
	if curTime == maxTime {
		pressure, _ := cur.getReleasedPressure(valves, maxTime)
		return pressure
	}
	locationValve := valves[cur.location()]
	nextPossible := make([]Timeline, 0, 1+len(locationValve.tunnels))
	canOpenMap := cur.getCanOpen(valves)
	if len(canOpenMap) == 0 {
		pressure, _ := cur.getReleasedPressure(valves, maxTime)
		return pressure
	}
	if _, exists := canOpenMap[locationValve.id]; exists {
		nextPossible = append(nextPossible, cur.open(locationValve.id))
	}
	targets := make(ValveMap)
	for _, tNode := range canOpenMap {
		tNode.getStartsFor(make(ValveMap), targets, locationValve.to)
	}
	for tid := range targets {
		nextTravelTimeline := cur.travel(tid)
		travelLoopDetected := false
		searchString := string(nextTravelTimeline[len(nextTravelTimeline)-chunkSize:])
		for i := len(cur) - chunkSize; i > 0; i -= chunkSize {
			prevOpStr := string(cur[i : i+chunkSize])
			if prevOpStr[0] == openOp {
				break
			}
			if searchString == prevOpStr {
				travelLoopDetected = true
				break
			}
		}
		if !travelLoopDetected {
			nextPossible = append(nextPossible, nextTravelTimeline)
		}
	}
	maxPressure := 0
	if len(nextPossible) > 0 {
		for i := range nextPossible {
			pressure := bestPath(valves, maxTime, nextPossible[i])
			if pressure > maxPressure {
				maxPressure = pressure
			}
		}
		return maxPressure
	}
	pressure, _ := cur.getReleasedPressure(valves, maxTime)
	return pressure
}

func getValves(path string) ValveMap {
	vm := make(ValveMap)
	lines, _ := file.GetLines(path)

	for _, line := range lines {
		if reInput.MatchString(line) {
			matches := reInput.FindStringSubmatch(line)
			if len(matches) == 4 {
				v := &Valve{}
				v.id = matches[1]
				v.rate, _ = strconv.Atoi(matches[2])
				v.tunnels = strings.Split(matches[3], ", ")
				v.to = make(ValveMap)
				v.from = make(ValveMap)
				vm[v.id] = v
			}
		}
	}

	for nodeId, node := range vm {
		for _, tId := range node.tunnels {
			node.to[tId] = vm[tId]
			vm[tId].from[nodeId] = node
		}
	}

	return vm
}
