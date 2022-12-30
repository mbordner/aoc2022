package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
)

var (
	reLine = regexp.MustCompile(`^Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.$`)
)

type Resource int

const (
	Ore Resource = iota
	Clay
	Obsidian
	Geode
)

var (
	memo = make(map[int]map[Inventory]*Inventory)
)

func main() {
	bps := getBlueprints("../data.txt")
	if len(bps) > 0 {
		geodes := make([]int, len(bps), len(bps))
		for i, bp := range bps[0:3] {
			geodes[i] = run(bp, 32, NewInventory()).resource[Geode]
			memo = make(map[int]map[Inventory]*Inventory)
			fmt.Println(bp.id, geodes[i])
		}
		fmt.Println(geodes[0] * geodes[1] * geodes[2])
	}
}

func min(a ...int) int {
	m := a[0]
	for _, b := range a[1:] {
		if b < m {
			m = b
		}
	}
	return m
}

func max(a ...int) int {
	m := a[0]
	for _, b := range a[1:] {
		if b > m {
			m = b
		}
	}
	return m
}

type Inventory struct {
	resource [4]int
	robot    [4]int
	build    [4]int
}

func NewInventory() *Inventory {
	i := new(Inventory)
	i.robot[Ore] = 1
	return i
}

func (inv *Inventory) Clone() *Inventory {
	ni := Inventory{}
	for i := 0; i < 4; i++ {
		ni.resource[i] = inv.resource[i]
		ni.robot[i] = inv.robot[i]
		ni.build[i] = inv.build[i]
	}
	return &ni
}

func (inv *Inventory) ProduceResources() {
	for i := 0; i < 4; i++ {
		inv.resource[i] += inv.robot[i]
	}
}

func (inv *Inventory) CompleteBuilds() *Inventory {
	ni := inv.Clone()
	for i := 0; i < 4; i++ {
		ni.robot[i] += ni.build[i]
		ni.build[i] = 0
	}
	return ni
}

func (inv *Inventory) CanBuildCounts(bp *Blueprint) [4]int {
	var counts [4]int

	// num of ore robots that can be built
	counts[Ore] = inv.resource[Ore] / bp.oreRobotOreCost
	// num of clay robots that can be built
	counts[Clay] = inv.resource[Ore] / bp.clayRobotOreCost
	// num of obsidian robots that can be built
	counts[Obsidian] = min(inv.resource[Ore]/bp.obsidianRobotOreCost, inv.resource[Clay]/bp.obsidianRobotClayCost)
	// num of geode cracking robots that can be built
	counts[Geode] = min(inv.resource[Ore]/bp.geodeRobotOreCost, inv.resource[Obsidian]/bp.geodeRobotObsidianCost)

	return counts
}

func (inv *Inventory) Build(bp *Blueprint, r Resource, count int) *Inventory {
	ni := inv.Clone()
	ni.build[r] += count
	switch r {
	case Ore:
		ni.resource[Ore] -= bp.oreRobotOreCost * count
	case Clay:
		ni.resource[Ore] -= bp.clayRobotOreCost * count
	case Obsidian:
		ni.resource[Ore] -= bp.obsidianRobotOreCost * count
		ni.resource[Clay] -= bp.obsidianRobotClayCost * count
	case Geode:
		ni.resource[Ore] -= bp.geodeRobotOreCost * count
		ni.resource[Obsidian] -= bp.geodeRobotObsidianCost * count
	}
	return ni
}

func run(bp *Blueprint, minsLeft int, inventory *Inventory) *Inventory {
	if minsLeft == 0 {
		return inventory
	}

	if _, e := memo[minsLeft]; !e {
		memo[minsLeft] = make(map[Inventory]*Inventory)
	}
	if i, e := memo[minsLeft][*inventory]; e {
		return i
	}

	// at start of min, the builds from the previous minute should complete
	ni := inventory.CompleteBuilds()

	nextTicks := make([]*Inventory, 0, 5)

	buildMaxCounts := ni.CanBuildCounts(bp)

	// if we can build all, there is no reason not to keep saving
	//if buildMaxCounts[Ore] == 0 || buildMaxCounts[Clay] == 0 || buildMaxCounts[Obsidian] == 0 || buildMaxCounts[Geode] == 0 {
	//nextTicks = append(nextTicks, ni)
	//}

	maxOreNeeds := max(bp.oreRobotOreCost, bp.clayRobotOreCost, bp.obsidianRobotOreCost, bp.geodeRobotOreCost)
	if ni.resource[Ore] <= maxOreNeeds {
		nextTicks = append(nextTicks, ni)
	}

	ni.ProduceResources()

	if buildMaxCounts[Ore] > 0 && ni.robot[Ore] < maxOreNeeds {
		nextTicks = append(nextTicks, ni.Build(bp, Ore, 1))
	}
	if buildMaxCounts[Clay] > 0 && ni.robot[Clay] < bp.obsidianRobotClayCost {
		nextTicks = append(nextTicks, ni.Build(bp, Clay, 1))
	}
	if buildMaxCounts[Obsidian] > 0 && ni.robot[Obsidian] < bp.geodeRobotObsidianCost {
		nextTicks = append(nextTicks, ni.Build(bp, Obsidian, 1))
	}
	if buildMaxCounts[Geode] > 0 {
		nextTicks = append(nextTicks, ni.Build(bp, Geode, 1))
	}

	maxGeodes := -1
	var maxInventory *Inventory

	for i := range nextTicks {
		ti := run(bp, minsLeft-1, nextTicks[i])
		if ti.resource[Geode] > maxGeodes {
			maxGeodes = ti.resource[Geode]
			maxInventory = ti
		}
	}

	memo[minsLeft][*inventory] = maxInventory

	return maxInventory
}

type Blueprint struct {
	id                     int
	oreRobotOreCost        int
	clayRobotOreCost       int
	obsidianRobotOreCost   int
	obsidianRobotClayCost  int
	geodeRobotOreCost      int
	geodeRobotObsidianCost int
}

func getBlueprints(path string) []*Blueprint {
	lines, _ := file.GetLines(path)

	bps := make([]*Blueprint, len(lines), len(lines))
	for i, line := range lines {
		if reLine.MatchString(line) {

			bp := Blueprint{}

			matches := reLine.FindStringSubmatch(line)
			if len(matches) == 8 {
				bp.id, _ = strconv.Atoi(matches[1])
				bp.oreRobotOreCost, _ = strconv.Atoi(matches[2])
				bp.clayRobotOreCost, _ = strconv.Atoi(matches[3])
				bp.obsidianRobotOreCost, _ = strconv.Atoi(matches[4])
				bp.obsidianRobotClayCost, _ = strconv.Atoi(matches[5])
				bp.geodeRobotOreCost, _ = strconv.Atoi(matches[6])
				bp.geodeRobotObsidianCost, _ = strconv.Atoi(matches[7])
			}

			bps[i] = &bp
		}
	}

	return bps
}
