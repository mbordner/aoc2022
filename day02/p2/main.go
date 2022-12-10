package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
)

// A 1 ROCK, B 2 PAPER, C 3 SCISSORS
// X lose, Y draw, Z win
func getScore(game string) int {
	opp := int(game[0]-'A') + 1
	score := int(game[2] - 'X')
	you := ((score + opp + 4) % 3) + 1
	return score*3 + you
}

func main() {
	lines, _ := file.GetLines("../data.txt")

	score := 0
	for _, line := range lines {
		s := getScore(line)
		fmt.Println(s)
		score += getScore(line)
	}

	fmt.Println(score)

}
