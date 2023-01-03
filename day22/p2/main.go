package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
)

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

var (
	reAction = regexp.MustCompile(`(\d+|R|L)`)
)

type Pos struct {
	R int
	C int
}

func (p Pos) String() string {
	return fmt.Sprintf("{%d,%d}", p.R, p.C)
}

func main() {
	dir := Right
	data, actions := getData("../data.txt")

	pos := getStart(data)
	//fmt.Print(pos)
actionsLoop:
	for a := range actions {
		if actions[a] == "L" || actions[a] == "R" {
			dir = getNextDir(dir, actions[a])
		} else {
			count, _ := strconv.Atoi(actions[a])
			for c := 0; c < count; c++ {
				var np Pos
				var nd Direction
				nd, np = getNextPos(data, dir, pos)
				if np == pos || data[np.R][np.C] == '#' {
					continue actionsLoop
				} else {
					//fmt.Println(getSide(pos), dir, pos)
					pos = np
					dir = nd
				}
			}
		}
	}

	fmt.Println(1000*(pos.R+1) + 4*(pos.C+1) + int(dir))
}

/*
*

	[1][2]
	[3]

[5][4]
[6]
*/
func getSide(p Pos) int {
	if p.R >= 0 && p.R <= 49 { // 1 or 2
		if p.C >= 50 && p.C <= 99 {
			return 1
		} else if p.C >= 100 && p.C <= 149 {
			return 2
		}
	} else if p.R >= 50 && p.R <= 99 { // 3
		if p.C >= 50 && p.C <= 99 {
			return 3
		}
	} else if p.R >= 100 && p.R <= 149 { // 4 or 5
		if p.C >= 0 && p.C <= 49 {
			return 5
		} else if p.C >= 50 && p.C <= 99 {
			return 4
		}
	} else if p.R >= 150 && p.R <= 199 {
		if p.C >= 0 && p.C <= 49 {
			return 6
		}
	}
	return -1
}

func getNextPos(data [][]byte, d Direction, p Pos) (nd Direction, np Pos) {
	curSide := getSide(p)

	switch d {
	case Right:
		np = Pos{R: p.R, C: p.C + 1}
	case Down:
		np = Pos{R: p.R + 1, C: p.C}
	case Left:
		np = Pos{R: p.R, C: p.C - 1}
	case Up:
		np = Pos{R: p.R - 1, C: p.C}
	}

	nextSide := getSide(np)
	if nextSide == curSide {
		return d, np
	}

	switch d {
	case Right:
		switch curSide {
		case 1:
			nextSide = 2
			nd = d
			// np is right
		case 2:
			nextSide = 4
			nd = Left
			np = Pos{R: 149 - p.R, C: 99}
		case 3:
			nextSide = 2
			nd = Up
			np = Pos{R: 49, C: (p.R - 50) + 100}
		case 4:
			nextSide = 2
			nd = Left
			np = Pos{R: 49 - (p.R - 100), C: 149}
		case 5:
			nextSide = 4
			nd = d
			// np is right
		case 6:
			nextSide = 4
			nd = Up
			np = Pos{R: 149, C: (p.R - 150) + 50}
		}
	case Down:
		switch curSide {
		case 1:
			nextSide = 3
			nd = d
			// np is right
		case 2:
			nextSide = 3
			nd = Left
			np = Pos{R: (p.C - 100) + 50, C: 99}
		case 3:
			nextSide = 4
			nd = d
			// np is right
		case 4:
			nextSide = 6
			nd = Left
			np = Pos{R: (p.C - 50) + 150, C: 49}
		case 5:
			nextSide = 6
			nd = d
			// np is right
		case 6:
			nextSide = 2
			nd = d
			np = Pos{R: 0, C: p.C + 100}
		}
	case Left:
		switch curSide {
		case 1:
			nextSide = 5
			nd = Right
			np = Pos{R: 149 - p.R, C: 0}
		case 2:
			nextSide = 1
			nd = d
			// np is right
		case 3:
			nextSide = 5
			nd = Down
			np = Pos{R: 100, C: p.R - 50}
		case 4:
			nextSide = 5
			nd = d
			// np is right
		case 5:
			nextSide = 1
			nd = Right
			np = Pos{R: 49 - (p.R - 100), C: 50}
		case 6:
			nextSide = 1
			nd = Down
			np = Pos{R: 0, C: 50 + (p.R - 150)}
		}
	case Up:
		switch curSide {
		case 1:
			nextSide = 6
			nd = Right
			np = Pos{R: (p.C - 50) + 150, C: 0}
		case 2:
			nextSide = 6
			nd = d
			np = Pos{R: 199, C: p.C - 100}
		case 3:
			nextSide = 1
			nd = d
			// np is right
		case 4:
			nextSide = 3
			nd = d
			// np is right
		case 5:
			nextSide = 3
			nd = Right
			np = Pos{R: 50 + p.C, C: 50}
		case 6:
			nextSide = 5
			nd = d
			// np is right
		}
	}

	return
}

func getNextDir(d Direction, action string) Direction {
	switch action {
	case "L":
		switch d {
		case Right:
			return Up
		case Down:
			return Right
		case Left:
			return Down
		case Up:
			return Left
		}
	case "R":
		switch d {
		case Right:
			return Down
		case Down:
			return Left
		case Left:
			return Up
		case Up:
			return Right
		}
	}
	return d
}

func getStart(data [][]byte) Pos {
	p := Pos{R: 0, C: 0}
	for data[0][p.C] == ' ' {
		p.C++
	}
	return p
}

func getData(path string) (data [][]byte, actions []string) {
	var empty int
	lines, _ := file.GetLines(path)
	for i := range lines {
		if len(lines[i]) == 0 {
			empty = i
			break
		}
	}

	data = make([][]byte, empty, empty)

	for i := 0; i < empty; i++ {
		data[i] = []byte(lines[i])
	}

	matches := reAction.FindAllStringSubmatch(lines[empty+1], -1)
	if len(matches) > 0 {
		actions = make([]string, len(matches), len(matches))
		for i := range matches {
			actions[i] = matches[i][0]
		}
	}

	return
}
