package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
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

	max := 0
	for _, v := range values {
		s := sum(v)
		if s > max {
			max = s
		}
	}

	fmt.Println(max)
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
