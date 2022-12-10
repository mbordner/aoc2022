package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/array/strings"
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

func exists(c compartment, r rune) bool {
	if _, exists := c[r]; exists {
		return true
	}
	return false
}

func findCommon(cs []compartment) rune {
	var common rune
outer:
	for common, _ = range cs[0] {
		for _, oc := range cs[1:] {
			if !exists(oc, common) {
				continue outer
			}
		}
		break outer
	}
	return common
}

func main() {
	sum := 0
	lines, _ := file.GetLines("../data.txt")

	groups := strings.Group(lines, 3)
	for _, group := range groups {
		compartments := make([]compartment, 0, len(group))
		for _, s := range group {
			compartments = append(compartments, getCompartment(s))
		}
		common := findCommon(compartments)
		sum += priority(common)
	}

	fmt.Println(sum)

}
