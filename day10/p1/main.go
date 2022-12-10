package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"github.com/mbordner/aoc2022/day10/cpu"
)

func main() {

	c := cpu.NewCPU()

	lines, _ := file.GetLines("../data.txt")

	for _, line := range lines {
		c.Process(line)
	}

	cycles := []int{20, 60, 100, 140, 180, 220}
	signalStrengths := make([]int, len(cycles), len(cycles))
	sum := 0

	for i := range cycles {
		xVal, _ := c.GetValueFromState("X", cycles[i])
		fmt.Println("x value at ", cycles[i], " is ", xVal)
		signalStrengths[i] = xVal * cycles[i]
		sum += signalStrengths[i]
	}

	fmt.Println(c)
	fmt.Println(signalStrengths)

	fmt.Println(sum)

}
