package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"sort"
	"strconv"
)

func sum(values []int) int {
	s := 0
	for _, v := range values {
		s += v
	}
	return s
}

func main() {
	values := getData()

	sums := make([]int, 0, len(values))

	for _, v := range values {
		sums = append(sums, sum(v))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sums)))

	fmt.Println(sums[0] + sums[1] + sums[2])
}

func getData() [][]int {
	allValues := make([][]int, 0, 20)

	lines, _ := file.GetLines("../data.txt")
	values := make([]int, 0, 100)
	for _, line := range lines {
		if line == "" {
			allValues = append(allValues, values)
			values = make([]int, 0, 100)
		} else {
			i, _ := strconv.Atoi(line)
			values = append(values, i)
		}
	}
	return allValues
}
