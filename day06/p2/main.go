package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
)

func main() {
	lines, _ := file.GetLines("../data.txt")
	for _, line := range lines {
		buffer := line
		for i := 13; i < len(buffer); i++ {
			test := make(map[rune]int)
			for _, r := range buffer[i-13 : i+1] {
				if _, exists := test[r]; !exists {
					test[r] = 0
				}
				test[r]++
			}
			if len(test) == 14 {
				fmt.Println(i + 1)
				break
			}
		}
	}
}
