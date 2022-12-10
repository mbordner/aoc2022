package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
)

func priority(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r-'a') + 1
	}
	return int(r-'A') + 27
}

type compartment map[rune]int

func getCompartment(s string) compartment {
	c := make(compartment)
	for _, r := range s {
		if count, exists := c[r]; exists {
			c[r] = count + 1
		} else {
			c[r] = 1
		}
	}
	return c
}

func main() {
	sum := 0
	lines, _ := file.GetLines("../data.txt")

	for _, line := range lines {
		l := len(line) / 2
		s1 := string([]byte(line)[0:l])
		s2 := string([]byte(line)[l:])

		c1 := getCompartment(s1)
		c2 := getCompartment(s2)

		for r := range c1 {
			if _, exists := c2[r]; exists {
				sum += priority(r)
				break
			}
		}
	}

	fmt.Println(sum)

}
