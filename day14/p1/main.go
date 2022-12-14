package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"strconv"
	"strings"
)

var (
	cave  = make(Cave)
	drain = Pos{x: 500, y: 0}
)

type Cave map[Pos]rune

type Pos struct {
	x int
	y int
}

func (c Cave) extents() (minX int, minY int, maxX int, maxY int) {
	minX, maxX = drain.x, drain.x
	minY, maxY = drain.y, drain.y

	for p := range c {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	return
}

func (c Cave) print() {
	minX, minY, maxX, maxY := c.extents()
	fmt.Println(fmt.Sprintf("min: %d,%d   max: %d,%d", minX, minY, maxX, maxY))

	for y := minY; y <= maxY; y++ {
		l := maxX - minX + 1
		line := make([]rune, l, l)
		for x, l := minX, 0; x <= maxX; x, l = x+1, l+1 {
			if r, exists := c[Pos{x: x, y: y}]; exists {
				line[l] = r
			} else {
				line[l] = '.'
			}
		}
		fmt.Println(string(line))
	}

}

type Sand struct {
	cur Pos
}

func (s Sand) next() *Pos {
	below := Pos{x: s.cur.x, y: s.cur.y + 1}
	if _, exists := cave[below]; !exists {
		return &below
	}
	belowLeft := Pos{x: s.cur.x - 1, y: s.cur.y + 1}
	if _, exists := cave[belowLeft]; !exists {
		return &belowLeft
	}
	belowRight := Pos{x: s.cur.x + 1, y: s.cur.y + 1}
	if _, exists := cave[belowRight]; !exists {
		return &belowRight
	}
	return nil
}

func sortInts(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func main() {
	cave[drain] = '+'

	drawRocks("../data.txt")

	minX, _, maxX, maxY := cave.extents()

	sandCount := 0

outer:
	for {
		sand := &Sand{cur: Pos{x: drain.x, y: drain.y}}
		for p := sand.next(); p != nil; p = sand.next() {
			if p.x < minX || p.x > maxX || p.y > maxY {
				break outer
			}
			sand.cur = *p
		}
		cave[sand.cur] = 'o'
		sandCount++
	}

	cave.print()

	fmt.Println("sand count: ", sandCount)
}

func drawRocks(path string) {

	lines, _ := file.GetLines(path)

	for _, line := range lines {
		positionStrings := strings.Split(line, " -> ")
		positions := make([]Pos, len(positionStrings), len(positionStrings))

		for i := range positions {
			tokens := strings.Split(positionStrings[i], ",")
			p := Pos{}
			p.x, _ = strconv.Atoi(tokens[0])
			p.y, _ = strconv.Atoi(tokens[1])
			positions[i] = p
		}

		for i := 1; i < len(positions); i++ {
			fromX, toX := sortInts(positions[i-1].x, positions[i].x)
			fromY, toY := sortInts(positions[i-1].y, positions[i].y)
			for x := fromX; x <= toX; x++ {
				cave[Pos{x: x, y: toY}] = '#'
			}
			for y := fromY; y <= toY; y++ {
				cave[Pos{x: toX, y: y}] = '#'
			}
		}
	}

}
