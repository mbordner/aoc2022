package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
	"strings"
)

var (
	reInput        = regexp.MustCompile(`^Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? (.*)`)
	memo           = make(map[int]map[int]map[uint64]map[uint64]int) // map[ time/mins remaining int ]map[ int id of valve ]map[ bit map of valves ] int score
	vIds           = NewValveIDMap()
	allOpen        OpenValveMap // maps all the valves that can open (rates > 0)
	possibleToOpen = make(ValveMap)
)

type ValveIDMap struct {
	intToStr map[uint64]string
	strToInt map[string]uint64
}

func NewValveIDMap() *ValveIDMap {
	vm := new(ValveIDMap)
	vm.intToStr = make(map[uint64]string)
	vm.strToInt = make(map[string]uint64)
	return vm
}

func (vm *ValveIDMap) addStr(strId string) uint64 {
	if v, exists := vm.strToInt[strId]; exists {
		return v
	}
	intId := uint64(len(vm.strToInt))
	intId = 1 << intId
	vm.strToInt[strId] = intId
	vm.intToStr[intId] = strId
	return intId
}

func (vm *ValveIDMap) getStr(id uint64) string {
	return vm.intToStr[id]
}

func (vm *ValveIDMap) getInt(id string) uint64 {
	return vm.strToInt[id]
}

type OpenValveMap uint64

func (o OpenValveMap) isOpen(v uint64) bool {
	if (uint64(o) & v) == v {
		return true
	}
	return false
}

func (o OpenValveMap) open(v uint64) OpenValveMap {
	no := uint64(o) | v
	return OpenValveMap(no)
}

type ValveMap map[uint64]*Valve

type Valve struct {
	id      uint64
	strId   string
	rate    int
	tunnels []uint64
	to      ValveMap
	from    ValveMap
}

func main() {

	valvesMap := getValves("../data.txt")

	var opened OpenValveMap
	pressure := getPressureFrom(valvesMap, opened, vIds.getInt("AA"), 26, 1)

	fmt.Println(pressure)

}

// too high 5498
func getPressureFrom(valvesMap ValveMap, opened OpenValveMap, valveId uint64, timeRemaining int, actor int) int {
	if timeRemaining < 1 {
		if actor == 0 {
			return 0
		}
		return getPressureFrom(valvesMap, opened, vIds.getInt("AA"), 26, 0)
	}

	if _, e := memo[actor]; !e {
		memo[actor] = make(map[int]map[uint64]map[uint64]int)
	}
	if _, e := memo[actor][timeRemaining]; !e {
		memo[actor][timeRemaining] = make(map[uint64]map[uint64]int)
	}
	if _, e := memo[actor][timeRemaining][valveId]; !e {
		memo[actor][timeRemaining][valveId] = make(map[uint64]int)
	}

	if v, e := memo[actor][timeRemaining][valveId][uint64(opened)]; e {
		return v
	}

	maxPressure := 0

	if !opened.isOpen(valveId) && valvesMap[valveId].rate > 0 {
		openedWithThis := opened.open(valveId)
		maxPressure = (valvesMap[valveId].rate * (timeRemaining - 1)) + getPressureFrom(valvesMap, openedWithThis, valveId, timeRemaining-1, actor)
	}

	for vId := range valvesMap[valveId].to {
		tmpPressure := getPressureFrom(valvesMap, opened, vId, timeRemaining-1, actor)
		if tmpPressure > maxPressure {
			maxPressure = tmpPressure
		}
	}

	memo[actor][timeRemaining][valveId][uint64(opened)] = maxPressure

	return maxPressure
}

func getValves(path string) ValveMap {
	vm := make(ValveMap)
	lines, _ := file.GetLines(path)

	for _, line := range lines {
		if reInput.MatchString(line) {
			matches := reInput.FindStringSubmatch(line)
			if len(matches) == 4 {
				v := &Valve{}
				v.strId = matches[1]
				v.id = vIds.addStr(matches[1])
				v.rate, _ = strconv.Atoi(matches[2])
				tunnelIds := strings.Split(matches[3], ", ")
				v.tunnels = make([]uint64, len(tunnelIds), len(tunnelIds))
				for i := range tunnelIds {
					v.tunnels[i] = vIds.addStr(tunnelIds[i])
				}
				v.to = make(ValveMap)
				v.from = make(ValveMap)
				vm[v.id] = v

				if v.rate > 0 {
					allOpen.open(v.id)       // building up what all open look like
					possibleToOpen[v.id] = v // saving which ids could have been opened from the start
				}
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
