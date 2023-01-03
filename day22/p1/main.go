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

func main() {
	dir := Right
	data, actions := getData("../data.txt")

	pos := getStart(data)
actionsLoop:
	for a := range actions {
		if actions[a] == "L" || actions[a] == "R" {
			dir = getNextDir(dir, actions[a])
		} else {
			count, _ := strconv.Atoi(actions[a])
			for c := 0; c < count; c++ {
				var np Pos
				dir, np = getNextPos(data, dir, pos)
				if np == pos || data[np.R][np.C] == '#' {
					continue actionsLoop
				} else {
					pos = np
				}
			}
		}
	}

	fmt.Println(1000*(pos.R+1) + 4*(pos.C+1) + int(dir))
}

func getNextPos(data [][]byte, d Direction, p Pos) (nd Direction, np Pos) {
	nd = d
	switch d {
	case Up:
		np = Pos{R: p.R - 1, C: p.C}
		if np.R < 0 || data[np.R][np.C] == ' ' {
			for np.R = p.R; np.R < len(data) && np.C < len(data[np.R]) && data[np.R][np.C] != ' '; np.R++ {
			}
			np.R--
		}
	case Right:
		np = Pos{R: p.R, C: p.C + 1}
		if np.C == len(data[np.R]) || data[np.R][np.C] == ' ' {
			for np.C = p.C; np.C >= 0 && data[np.R][np.C] != ' '; np.C-- {
			}
			np.C++
		}
	case Down:
		np = Pos{R: p.R + 1, C: p.C}
		if np.R == len(data) || np.C > len(data[np.R]) || data[np.R][np.C] == ' ' {
			for np.R = p.R; np.R >= 0 && data[np.R][np.C] != ' '; np.R-- {
			}
			np.R++
		}
	case Left:
		np = Pos{R: p.R, C: p.C - 1}
		if np.C < 0 || data[np.R][np.C] == ' ' {
			for np.C = p.C; np.C < len(data[np.R]) && data[np.R][np.C] != ' '; np.C++ {
			}
			np.C--
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
