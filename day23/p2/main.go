package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
)

type Direction int

const (
	North Direction = iota
	South
	West
	East
)

type Pos struct {
	X int
	Y int
}

type Positions []Pos

func (p Pos) Transform(x, y int) Pos {
	return Pos{X: p.X + x, Y: p.Y + y}
}

type GroveState struct {
	round int
	elves Positions
	moved int
}

func (gs *GroveState) Print() {
	minX, minY, maxX, maxY := gs.GetExtents()
	groveMap := gs.GetMap()

	l := maxX - minX + 1

	for y := maxY; y >= minY; y-- {
		rs := make([]rune, l, l)
		for x, i := minX, 0; x <= maxX; x, i = x+1, i+1 {
			p := Pos{X: x, Y: y}
			if groveMap.Has(p) {
				rs[i] = '#'
			} else {
				rs[i] = '.'
			}
		}
		fmt.Println(string(rs))
	}
}

func (gs *GroveState) GetExtents() (minX int, minY int, maxX int, maxY int) {
	e := gs.elves[0]
	minX = e.X
	minY = e.Y
	maxX = e.X
	maxY = e.Y
	for _, e = range gs.elves[1:] {
		if e.X < minX {
			minX = e.X
		}
		if e.Y < minY {
			minY = e.Y
		}
		if e.X > maxX {
			maxX = e.X
		}
		if e.Y > maxY {
			maxY = e.Y
		}
	}
	return
}

func (gs *GroveState) GetRoundProposedMovementOrder() [4]Direction {
	switch gs.round % 4 {
	default:
		return [4]Direction{North, South, West, East}
	case 1:
		return [4]Direction{South, West, East, North}
	case 2:
		return [4]Direction{West, East, North, South}
	case 3:
		return [4]Direction{East, North, South, West}
	}
}

func (gs *GroveState) NextState() *GroveState {
	ns := new(GroveState)
	ns.round = gs.round + 1
	ns.elves = make(Positions, len(gs.elves), len(gs.elves))

	pes, pmc := gs.ProposeMoves()
	for i, pe := range pes {
		if pmc.Get(pe) == 1 {
			ns.elves[i] = pe
			ns.moved++
		} else {
			ns.elves[i] = gs.elves[i]
		}
	}

	return ns
}

func (gs *GroveState) NumMoved() int {
	return gs.moved
}

func (gs *GroveState) RoundIndex() int {
	return gs.round
}

func (gm GroveMap) PrintAdjacent(p Pos) {
	chars := make([][]rune, 3, 3)
	for j, y := 0, 1; y >= -1; j, y = j+1, y-1 {
		chars[j] = make([]rune, 3, 3)
		for i, x := 0, -1; x <= 1; i, x = i+1, x+1 {
			if gm.Has(p.Transform(x, y)) {
				chars[j][i] = '#'
			} else {
				chars[j][i] = '.'
			}
		}
		fmt.Println(string(chars[j]))
	}
}

func (gs *GroveState) ProposeMoves() (Positions, ProposedMoveCounts) {
	moveOrder := gs.GetRoundProposedMovementOrder()
	groveMap := gs.GetMap()
	proposedElves := make(Positions, len(gs.elves), len(gs.elves))
	proposedMoveCounts := make(ProposedMoveCounts)

	for i, e := range gs.elves {
		proposedElves[i] = e
		if groveMap.HasAdjacent(e) {
			for _, d := range moveOrder {
				if d == North {
					if !groveMap.HasNorthAdjacent(e) {
						pe := e.Transform(0, 1)
						proposedElves[i] = pe
						proposedMoveCounts.Add(pe)
						break
					}
				} else if d == South {
					if !groveMap.HasSouthAdjacent(e) {
						pe := e.Transform(0, -1)
						proposedElves[i] = pe
						proposedMoveCounts.Add(pe)
						break
					}
				} else if d == West {
					if !groveMap.HasWestAdjacent(e) {
						pe := e.Transform(-1, 0)
						proposedElves[i] = pe
						proposedMoveCounts.Add(pe)
						break
					}
				} else if d == East {
					if !groveMap.HasEastAdjacent(e) {
						pe := e.Transform(1, 0)
						proposedElves[i] = pe
						proposedMoveCounts.Add(pe)
						break
					}
				}
			}
		}
	}

	return proposedElves, proposedMoveCounts
}

type GroveMap map[Pos]int

func (gm GroveMap) HasNorthAdjacent(p Pos) bool {
	n := p.Transform(0, 1)
	if gm.Has(n) || gm.Has(n.Transform(-1, 0)) || gm.Has(n.Transform(1, 0)) {
		return true
	}
	return false
}

func (gm GroveMap) HasSouthAdjacent(p Pos) bool {
	s := p.Transform(0, -1)
	if gm.Has(s) || gm.Has(s.Transform(-1, 0)) || gm.Has(s.Transform(1, 0)) {
		return true
	}
	return false
}

func (gm GroveMap) HasWestAdjacent(p Pos) bool {
	w := p.Transform(-1, 0)
	if gm.Has(w) || gm.Has(w.Transform(0, 1)) || gm.Has(w.Transform(0, -1)) {
		return true
	}
	return false
}

func (gm GroveMap) HasEastAdjacent(p Pos) bool {
	e := p.Transform(1, 0)
	if gm.Has(e) || gm.Has(e.Transform(0, 1)) || gm.Has(e.Transform(0, -1)) {
		return true
	}
	return false
}

func (gm GroveMap) HasAdjacent(p Pos) bool {
	if gm.HasNorthAdjacent(p) {
		return true
	}
	if gm.Has(p.Transform(-1, 0)) || gm.Has(p.Transform(1, 0)) {
		return true
	}
	if gm.HasSouthAdjacent(p) {
		return true
	}
	return false
}

type ProposedMoveCounts map[Pos]int

func (pmc *ProposedMoveCounts) Add(p Pos) {
	(*pmc)[p] = pmc.Get(p) + 1
}

func (pmc *ProposedMoveCounts) Get(p Pos) int {
	if c, exists := (*pmc)[p]; exists {
		return c
	}
	return 0
}

func (gm GroveMap) Has(p Pos) bool {
	if _, exists := gm[p]; exists {
		return true
	}
	return false
}

func (gs *GroveState) GetMap() GroveMap {
	gm := make(GroveMap)
	for i, e := range gs.elves {
		gm[e] = i
	}
	return gm
}

func main() {
	initialGS := getInitialGroveState("../data.txt")
	lastGS := initialGS

	if initialGS != nil {

		for {
			lastGS = lastGS.NextState()
			if lastGS.NumMoved() == 0 {
				break
			}
		}

		fmt.Println(lastGS.RoundIndex())
	}
}

func getInitialGroveState(path string) *GroveState {
	gs := new(GroveState)

	lines, _ := file.GetLines(path)

	elves := make([]Pos, 0, len(lines)*len(lines[0]))

	for j, line := range lines {
		for i, r := range line {
			if r == '#' {
				elves = append(elves, Pos{X: i, Y: -j})
			}
		}
	}

	gs.elves = elves

	return gs
}
