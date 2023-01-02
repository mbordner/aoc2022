package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
)

func main() {
	m := NewMountain("../data.txt")
	if m != nil {

		path := m.StepsToGoal(m.start, m.goal, 0, 500)
		if path != nil {
			fmt.Println(path)
			//for i := 0; i < len(path); i++ {
			//fmt.Println(i, " :=-=-=-=-=-=-=-=-=-=-=-=-=-=")
			//m.PrintWithPlayerPos(i, path[i])
			//}

			fmt.Println(len(path) - 1)
		}
	}
}

type Pos struct {
	C int // column (x)
	R int // row (y)
}

type Mountain struct {
	start      Pos
	goal       Pos
	rowStreams []*RowStream
	colStreams []*ColumnStream
}

// NewMountain sets up mountain, capture start and goal positions
// sets up column and row streams that translates origin corner 1,1 in
// data file to 0,0 .. i.e. row[0] stream represents the 2nd row in the mountain file data
// set between the start and ending # chars
func NewMountain(path string) *Mountain {
	m := new(Mountain)

	lines, _ := file.GetLines(path)

	var rows, cols int
	rows = len(lines) - 2
	cols = len(lines[0]) - 2

	m.rowStreams = make([]*RowStream, rows, rows)
	m.colStreams = make([]*ColumnStream, cols, cols)

	for r := 0; r < rows; r++ {
		m.rowStreams[r] = NewRowStream(cols)
	}
	for c := 0; c < cols; c++ {
		m.colStreams[c] = NewColumnStream(rows)
	}

	for r := 0; r < len(lines); r++ {
		if r == 0 || r == len(lines)-1 {
			for c := 0; c < len(lines[r]); c++ {
				if lines[r][c] == '.' {
					p := Pos{R: r - 1, C: c - 1}
					if r == 0 {
						m.start = p
					} else {
						m.goal = p
					}
					break
				}
			}
		} else {
			for c, char := range lines[r][1 : len(lines[r])-1] {
				row := r - 1
				col := c

				switch char {
				case '>':
					m.rowStreams[row].PlaceEastward(byte(char), col)
				case '<':
					m.rowStreams[row].PlaceWestward(byte(char), col)
				case 'v':
					m.colStreams[col].PlaceSouthward(byte(char), row)
				case '^':
					m.colStreams[col].PlaceNorthward(byte(char), row)
				}
			}
		}
	}

	return m
}

func (m *Mountain) Print(stepsTaken int) {
	lines := m.GetSteamPrintLines(stepsTaken)
	for _, l := range lines {
		fmt.Println(l)
	}
}

func (m *Mountain) PrintWithPlayerPos(stepsTaken int, p Pos) {
	lines := m.GetSteamPrintLines(stepsTaken)
	j, i := p.R+1, p.C+1
	line := []byte(lines[j])
	line[i] = 'P'
	lines[j] = string(line)
	for _, l := range lines {
		fmt.Println(l)
	}
}

func (m *Mountain) GetSteamPrintLines(stepsTaken int) []string {
	lines := make([]string, 0, len(m.rowStreams)+2)
	line := make([]byte, len(m.colStreams)+2, len(m.colStreams)+2)
	for c := 0; c < len(line); c++ {
		if m.start.C == c-1 {
			line[c] = '.'
		} else {
			line[c] = '#'
		}
	}
	lines = append(lines, string(line))

	line[0] = '#'
	line[len(line)-1] = '#'

	for r := 0; r < len(m.rowStreams); r++ {
		for c, i := 0, 1; c < len(m.colStreams); c, i = c+1, i+1 {
			char, _ := m.GetBlizzardCount(stepsTaken, r, c)
			line[i] = char
		}
		lines = append(lines, string(line))
	}

	for c := 0; c < len(line); c++ {
		if m.goal.C == c-1 {
			line[c] = '.'
		} else {
			line[c] = '#'
		}
	}
	lines = append(lines, string(line))

	return lines
}

func (m *Mountain) HasBlizzard(stepsTaken, r, c int) bool {
	var count int
	_, count = m.colStreams[c].northward.GetBlizzardCount(stepsTaken, r)
	if count > 0 {
		return true
	}
	_, count = m.colStreams[c].southward.GetBlizzardCount(stepsTaken, r)
	if count > 0 {
		return true
	}
	_, count = m.rowStreams[r].westward.GetBlizzardCount(stepsTaken, c)
	if count > 0 {
		return true
	}
	_, count = m.rowStreams[r].eastward.GetBlizzardCount(stepsTaken, c)
	if count > 0 {
		return true
	}
	return false
}

func (m *Mountain) GetBlizzardCount(stepsTaken, r, c int) (char byte, count int) {
	var tmpCount int
	var tmpChar byte
	// northward
	tmpChar, tmpCount = m.colStreams[c].northward.GetBlizzardCount(stepsTaken, r)
	count += tmpCount
	if tmpChar != ' ' {
		char = tmpChar
	}
	// southward
	tmpChar, tmpCount = m.colStreams[c].southward.GetBlizzardCount(stepsTaken, r)
	count += tmpCount
	if tmpChar != ' ' {
		char = tmpChar
	}
	// westward
	tmpChar, tmpCount = m.rowStreams[r].westward.GetBlizzardCount(stepsTaken, c)
	count += tmpCount
	if tmpChar != ' ' {
		char = tmpChar
	}
	// eastward
	tmpChar, tmpCount = m.rowStreams[r].eastward.GetBlizzardCount(stepsTaken, c)
	count += tmpCount
	if tmpChar != ' ' {
		char = tmpChar
	}

	if count > 1 {
		char = '0' + byte(count)
	} else if count == 0 {
		char = '.'
	}
	return
}

var (
	memo = make(map[int]map[Pos][]Pos)
)

func (m *Mountain) StepsToGoal(from Pos, goal Pos, stepsTaken int, maxStepsTaken int) []Pos {

	if _, e := memo[stepsTaken]; !e {
		memo[stepsTaken] = make(map[Pos][]Pos)
	}
	if ps, e := memo[stepsTaken][from]; e {
		return ps
	}

	if from == goal {
		return []Pos{goal}
	}

	if maxStepsTaken == stepsTaken {
		return nil
	}

	// if waiting and blizzard came, or we walked into a blizzard, this path ends
	if from.R >= 0 && from.R < len(m.rowStreams) {
		if m.HasBlizzard(stepsTaken, from.R, from.C) {
			memo[stepsTaken][from] = nil
			return nil
		}
	}

	nextPossible := make([]Pos, 0, 5)

	if from.R < 0 { // if at start, can move down only
		nextPossible = append(nextPossible, Pos{R: from.R + 1, C: from.C})
	} else if from.R == len(m.rowStreams) { // if at goal, can move up only
		nextPossible = append(nextPossible, Pos{R: from.R - 1, C: from.C})
	} else {
		// vertical movements
		if from.R == 0 {
			if from.C == m.start.C {
				nextPossible = append(nextPossible, m.start)
			}
		}
		if from.R == len(m.rowStreams)-1 {
			if from.C == m.goal.C {
				nextPossible = append(nextPossible, m.goal)
			}
		}
		if from.R > 0 {
			nextPossible = append(nextPossible, Pos{R: from.R - 1, C: from.C})
		}
		if from.R < len(m.rowStreams)-1 {
			nextPossible = append(nextPossible, Pos{R: from.R + 1, C: from.C})
		}

		// horizontal movements
		if from.C > 0 {
			nextPossible = append(nextPossible, Pos{R: from.R, C: from.C - 1})
		}
		if from.C < len(m.colStreams)-1 {
			nextPossible = append(nextPossible, Pos{R: from.R, C: from.C + 1})
		}
	}

	// can always wait
	nextPossible = append(nextPossible, from)

	var best []Pos
	resultFound := false

	for _, p := range nextPossible {
		result := m.StepsToGoal(p, goal, stepsTaken+1, maxStepsTaken)
		if result != nil {
			if resultFound {
				if len(result) < len(best) {
					best = result
				}
			} else {
				best = result
			}
			resultFound = true
		}
	}

	if resultFound {
		best = append([]Pos{from}, best...)
		memo[stepsTaken][from] = best
		return best
	}

	memo[stepsTaken][from] = nil
	return nil
}

type JSDir int

const (
	Forward JSDir = iota
	Backward
)

type JetStream struct {
	stream []byte
	length int
	dir    JSDir
}

func (js *JetStream) Place(b byte, p int) {
	js.stream[p] = b
	js.stream[p+js.length] = b
}

func (js *JetStream) GetBlizzardCount(stepsTaken, p int) (char byte, count int) {
	var pointer int
	if js.dir == Forward {
		pointer = 0 + (stepsTaken % js.length) + p
	} else if js.dir == Backward {
		pointer = js.length - (stepsTaken % js.length) + p
	}
	if js.stream[pointer] == ' ' {
		return ' ', 0
	}
	return js.stream[pointer], 1
}

func NewJetStream(length int, dir JSDir) *JetStream {
	js := new(JetStream)
	js.length = length
	js.dir = dir
	l := length + length
	js.stream = make([]byte, l, l)
	for i := 0; i < l; i++ {
		js.stream[i] = byte(' ')
	}
	return js
}

type ColumnStream struct {
	northward *JetStream
	southward *JetStream
}

func (cs *ColumnStream) PlaceNorthward(b byte, p int) {
	cs.northward.Place(b, p)
}

func (cs *ColumnStream) PlaceSouthward(b byte, p int) {
	cs.southward.Place(b, p)
}

type RowStream struct {
	westward *JetStream
	eastward *JetStream
}

func (rs *RowStream) PlaceWestward(b byte, p int) {
	rs.westward.Place(b, p)
}

func (rs *RowStream) PlaceEastward(b byte, p int) {
	rs.eastward.Place(b, p)
}

func NewColumnStream(length int) *ColumnStream {
	cs := new(ColumnStream)
	cs.northward = NewJetStream(length, Forward)
	cs.southward = NewJetStream(length, Backward)
	return cs
}

func NewRowStream(length int) *RowStream {
	rs := new(RowStream)
	rs.westward = NewJetStream(length, Forward)
	rs.eastward = NewJetStream(length, Backward)
	return rs
}
