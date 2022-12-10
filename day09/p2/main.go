package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"github.com/mbordner/aoc2022/common/geom"
	"strconv"
)

var (
	bb = &geom.BoundingBox{}
)

type Moves []geom.Pos
type VisitedPositions map[geom.Pos]int

func (m *Moves) add(p geom.Pos) {
	*m = append(*m, p)
}

func (m *Moves) cur() geom.Pos {
	return (*m)[len(*m)-1]
}

func (ps *VisitedPositions) add(p geom.Pos) {
	if count, exists := (*ps)[p]; exists {
		(*ps)[p] = count + 1
	} else {
		(*ps)[p] = 1
	}
	bb.Extend(p)
}

type Knot struct {
	char    rune
	moves   Moves
	visited VisitedPositions
}

func NewKnot(char rune, init geom.Pos) *Knot {
	k := new(Knot)
	k.char = char
	k.moves = make(Moves, 0, 100)
	k.visited = make(VisitedPositions)

	k.moveTo(init)
	return k
}

func (k *Knot) moveTo(p geom.Pos) {
	k.moves.add(p)
	k.visited.add(p)
}

func (k *Knot) curPos() geom.Pos {
	return k.moves[len(k.moves)-1]
}

func (k *Knot) lastPos() geom.Pos {
	if len(k.moves) >= 2 {
		return k.moves[len(k.moves)-2]
	}
	return k.moves[0]
}

func abs(v int) int {
	if v >= 0 {
		return v
	}
	return v * -1
}

func normalize(v int) int {
	if v < 0 {
		return -1
	}
	return 1
}

func main() {
	knots := make([]*Knot, 10, 10)
	for i, j := 0, 9; i < 9; i, j = i+1, j-1 {
		knots[i] = NewKnot(rune('0'+j), geom.Pos{})
	}
	knots[9] = NewKnot('H', geom.Pos{})

	//knots = make([]*Knot, 2, 2)
	//knots[0] = NewKnot('T', geom.Pos{})
	//knots[1] = NewKnot('H', geom.Pos{})

	lines, _ := file.GetLines("../data.txt")
	for _, line := range lines {
		dir := string(line[0])
		amt, _ := strconv.Atoi(line[2:])

		headIndex := len(knots) - 1

		for i := 0; i < amt; i++ {

			nextPositions := make(geom.Positions, len(knots), len(knots))

			switch dir {
			case "U":
				nextPositions[headIndex] = knots[headIndex].curPos().Transform(0, 1, 0)
			case "D":
				nextPositions[headIndex] = knots[headIndex].curPos().Transform(0, -1, 0)
			case "L":
				nextPositions[headIndex] = knots[headIndex].curPos().Transform(-1, 0, 0)
			case "R":
				nextPositions[headIndex] = knots[headIndex].curPos().Transform(1, 0, 0)
			}

			for ki := headIndex - 1; ki >= 0; ki-- {
				kiCur := knots[ki].curPos()
				kiHeadNext := nextPositions[ki+1]

				v := kiHeadNext.Diff(kiCur)
				if kiCur.X == kiHeadNext.X || kiCur.Y == kiHeadNext.Y {
					if abs(v.X) > 1 {
						v.X = normalize(v.X)
					} else {
						v.X = 0
					}
					if abs(v.Y) > 1 {
						v.Y = normalize(v.Y)
					} else {
						v.Y = 0
					}
				} else {
					if abs(v.X) > 1 || abs(v.Y) > 1 {
						v.X = normalize(v.X)
						v.Y = normalize(v.Y)
					} else {
						v = geom.Pos{}
					}
				}

				nextPositions[ki] = kiCur.Transform(v.X, v.Y, v.Z)
			}

			for ki := 0; ki < len(nextPositions); ki++ {
				knots[ki].moveTo(nextPositions[ki])
			}

			//fmt.Println("======================")
			//PrintKnots(knots)
		}
	}

	fmt.Println("======================")

	fmt.Println(len(knots[0].visited))
}

func PrintKnots(knots []*Knot) {
	labels := make([]rune, 0, len(knots)+1)
	labels = append(labels, 's')
	pss := make(geom.Positions, 0, len(knots)+1)
	pss = append(pss, geom.Pos{})
	for i := range knots {
		labels = append(labels, knots[i].char)
		pss = append(pss, knots[i].curPos())
	}
	lines := bb.GetPrintLines('.', labels, pss)
	for _, line := range lines {
		fmt.Println(line)
	}
}
