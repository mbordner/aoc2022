package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
)

// AX 1 ROCK, BY 2 PAPER, CZ 3 SCISSORS
func getScore(game string) int {
	you := int(game[2]-'X') + 1
	opp := int(game[0]-'A') + 1
	score := (you - opp + 4) % 3
	return score*3 + you
}

func main() {
	lines, _ := file.GetLines("../test2.txt")

	score := 0
	for _, line := range lines {
		s := getScore(line)
		fmt.Println(s)
		score += getScore(line)
	}

	fmt.Println(score)

}
