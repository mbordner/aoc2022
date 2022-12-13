package main

import (
	"encoding/json"
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"sort"
)

type Packet []interface{}
type Pair []Packet

func (p Pair) makeList(value float64) Packet {
	return []interface{}{value}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (p Pair) less(i, j Packet) int {
	minl := min(len(i), len(j))
	for index := 0; index < minl; index++ {
		lInt, lIsInt := i[index].(float64)
		rInt, rIsInt := j[index].(float64)
		lList, lIsList := i[index].([]interface{})
		rList, rIsList := j[index].([]interface{})

		if lIsInt && rIsInt {
			if lInt < rInt {
				return -1
			} else if lInt > rInt {
				return 1
			}
			continue
		}

		if lIsInt && rIsList {
			lList = p.makeList(lInt)
		} else if rIsInt && lIsList {
			rList = p.makeList(rInt)
		}

		lessVal := p.less(lList, rList)
		if lessVal != 0 {
			return lessVal
		}
	}
	if len(i) < len(j) {
		return -1
	} else if len(i) > len(j) {
		return 1
	}
	return 0
}

func (p Pair) Len() int {
	return len(p)
}
func (p Pair) Less(i, j int) bool {
	return p.less(p[i], p[j]) < 0
}
func (p Pair) Swap(i, j int) {

}

func main() {
	pairs, _ := getPairs("../data.txt")
	if len(pairs) > 0 {
		indexSum := 0
		for i := range pairs {
			if sort.IsSorted(pairs[i]) {
				fmt.Println("found pair in right order at index: ", i+1)
				indexSum += i + 1
			}
		}
		fmt.Println(indexSum)
	}
}

func getPairs(path string) ([]Pair, error) {
	lines, _ := file.GetLines(path)

	pairs := make([]Pair, 0, len(lines)/3+1)

	curPair := make([]Packet, 2, 2)
	lineCount := 0
	for _, line := range lines {
		if len(line) == 0 {
			pairs = append(pairs, curPair)
			curPair = make([]Packet, 2, 2)
			lineCount = 0
			continue
		}
		err := json.Unmarshal([]byte(line), &(curPair[lineCount]))
		if err != nil {
			return nil, err
		}
		lineCount++
	}

	pairs = append(pairs, curPair)
	return pairs, nil
}
