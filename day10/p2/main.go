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

	cycle := 0
	for j := 0; j < 6; j++ {
		line := make([]byte, 40, 40)
		for i := 0; i < 40; i++ {
			sprite, _ := c.GetValueFromState("X", cycle+1)
			if i >= sprite-1 && i <= sprite+1 {
				line[i] = byte('#')
			} else {
				line[i] = byte('.')
			}
			cycle++
		}
		fmt.Println(string(line))
	}

}
