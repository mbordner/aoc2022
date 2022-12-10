package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
)

var (
	rePairs = regexp.MustCompile(`(\d+)-(\d+),(\d+)-(\d+)`)
)

func main() {
	count := 0
	pairs := getPairs("../data.txt")
	if len(pairs) > 0 {
		for _, p := range pairs {
			// start1 contained in set2     or end1 contained in set2
			if (p[0] >= p[2] && p[0] <= p[3]) || (p[1] <= p[3] && p[1] >= p[2]) ||
				// start2 contained in set1     or end2 contained in set1
				(p[2] >= p[0] && p[2] <= p[1]) || (p[3] <= p[1] && p[3] >= p[0]) {
				count += 1
			}
		}
	}
	fmt.Println(count)
}

func getVal(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func getPairs(path string) [][4]int {

	lines, _ := file.GetLines(path)
	pairs := make([][4]int, 0, 20)

	for _, line := range lines {
		matches := rePairs.FindStringSubmatch(line)
		if len(matches) > 0 {
			var p [4]int
			p[0] = getVal(matches[1])
			p[1] = getVal(matches[2])
			p[2] = getVal(matches[3])
			p[3] = getVal(matches[4])
			pairs = append(pairs, p)
		}
	}

	return pairs
}
