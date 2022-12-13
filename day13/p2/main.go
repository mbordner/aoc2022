package main

import (
	"encoding/json"
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"sort"
)

type Packet []interface{}
type Packets []Packet

func (p Packet) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

func (p Packets) makeList(value float64) Packet {
	return []interface{}{value}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (p Packets) less(i, j Packet) int {
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

func (p Packets) Len() int {
	return len(p)
}
func (p Packets) Less(i, j int) bool {
	return p.less(p[i], p[j]) < 0
}
func (p Packets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	pairs, _ := getPairs("../data.txt")
	packets := make(Packets, 0, len(pairs)*2+2)

	if len(pairs) > 0 {
		for _, pair := range pairs {
			packets = append(packets, pair...)
		}
	}

	dividerPackets := []string{`[[2]]`, `[[6]]`}

	pair := make(Packets, 2, 2)
	_ = json.Unmarshal([]byte(dividerPackets[0]), &(pair[0]))
	_ = json.Unmarshal([]byte(dividerPackets[1]), &(pair[1]))

	packets = append(packets, pair...)

	sort.Sort(packets)

	indexes := make(map[string]int)

	for i, p := range packets {
		indexes[p.String()] = i
		if _, exists := indexes[dividerPackets[0]]; exists {
			if _, exists = indexes[dividerPackets[1]]; exists {
				break
			}
		}
	}

	fmt.Println((indexes[dividerPackets[0]] + 1) * (indexes[dividerPackets[1]] + 1))
}

func getPairs(path string) ([]Packets, error) {
	lines, _ := file.GetLines(path)

	pairs := make([]Packets, 0, len(lines)/3+1)

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
