package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"strconv"
)

type Direction int

const (
	N Direction = iota
	S
	W
	E
)

type treeData struct {
	i          int
	j          int
	visibility [4]bool
	treeCount  [4]int
}

func (td treeData) visible() bool {
	for i := 0; i < len(td.visibility); i++ {
		if td.visibility[i] {
			return true
		}
	}
	return false
}

func (td treeData) score() int {
	score := td.treeCount[0] * td.treeCount[1]
	for i := 2; i < 4; i++ {
		score *= td.treeCount[i]
	}
	return score
}

func visible(trees [][]int, i, j int) treeData {
	td := treeData{i: i, j: j}

	if j == 0 {
		td.visibility[W] = true
	} else {
		visible := true
		for p := j - 1; p >= 0; p-- {
			td.treeCount[W]++
			if trees[i][p] >= trees[i][j] {
				visible = false
				break
			}
		}
		td.visibility[W] = visible
	}

	if i == 0 {
		td.visibility[N] = true
	} else {
		visible := true
		for p := i - 1; p >= 0; p-- {
			td.treeCount[N]++
			if trees[p][j] >= trees[i][j] {
				visible = false
				break
			}
		}
		td.visibility[N] = visible
	}

	if j == len(trees[i])-1 {
		td.visibility[E] = true
	} else {
		visible := true
		for p := j + 1; p < len(trees[i]); p++ {
			td.treeCount[E]++
			if trees[i][p] >= trees[i][j] {
				visible = false
				break
			}
		}
		td.visibility[E] = visible
	}

	if i == len(trees)-1 {
		td.visibility[S] = true
	} else {
		visible := true
		for p := i + 1; p < len(trees); p++ {
			td.treeCount[S]++
			if trees[p][j] >= trees[i][j] {
				visible = false
				break
			}
		}
		td.visibility[S] = visible
	}

	return td
}

func main() {
	trees := getData("../data.txt")
	if len(trees) > 0 {
		visibleCount := 0
		maxScore := 0
		for i := 0; i < len(trees); i++ {
			for j := 0; j < len(trees[i]); j++ {
				td := visible(trees, i, j)
				score := td.score()
				if score > maxScore {
					maxScore = score
				}
				if td.visible() {
					visibleCount++
				}
			}
		}
		fmt.Println(visibleCount, maxScore)
	}
}

func getData(path string) [][]int {

	lines, _ := file.GetLines(path)

	trees := make([][]int, len(lines), len(lines))

	for i, line := range lines {
		trees[i] = make([]int, len(line), len(line))
		for j, r := range line {
			val, _ := strconv.Atoi(string(r))
			trees[i][j] = val
		}
	}

	return trees
}
