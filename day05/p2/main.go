package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/datastructure"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
)

var (
	reInstruction = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	reStackIds    = regexp.MustCompile(`(\d+)`)
	reItems       = regexp.MustCompile(`\[(\w+)\]`)
)

type instruction struct {
	count       int
	source      int
	destination int
}

func main() {
	stacks, instructions := getData("../data.txt")
	for i := range instructions {
		instr := instructions[i]

		items := stacks[instr.source].PopN(instr.count)
		for _, item := range items {
			stacks[instr.destination].Push(item)
		}
	}

	items := make([]byte, len(stacks), len(stacks))
	for i := range stacks {
		items[i] = (stacks[i].Peek()).(string)[0]
	}

	fmt.Println(string(items))
}

func getData(path string) ([]*datastructure.Stack, []*instruction) {
	lines, _ := file.GetLines(path)

	var stacks []*datastructure.Stack

	var i int
	for ; len(lines[i]) > 0; i++ {
	}

	matches := reStackIds.FindAllStringSubmatch(lines[i-1], -1)
	if len(matches) > 0 {
		stacks = make([]*datastructure.Stack, len(matches), len(matches))
		for j := range matches {
			stacks[j] = datastructure.NewStack(i)
		}

		for j := i - 2; j >= 0; j-- {
			matchIndexes := reItems.FindAllStringSubmatchIndex(lines[j], -1)
			if len(matchIndexes) > 0 {
				for _, m := range matchIndexes {
					item := lines[j][m[2]:m[3]]
					stacks[m[0]/4].Push(item)
				}
			}
		}
	}

	instructions := make([]*instruction, 0, 100)
	for j := i + 1; j < len(lines); j++ {
		line := lines[j]
		matches := reInstruction.FindStringSubmatch(line)
		if len(matches) == 4 {
			instr := &instruction{}
			instr.count, _ = strconv.Atoi(matches[1])
			instr.source, _ = strconv.Atoi(matches[2])
			instr.destination, _ = strconv.Atoi(matches[3])
			instr.source--
			instr.destination--
			instructions = append(instructions, instr)
		}
	}

	return stacks, instructions
}
