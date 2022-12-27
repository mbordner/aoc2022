package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"strconv"
	"strings"
)

type PointType int

const (
	Unknown PointType = iota
	Solid
)

func main() {
	pts := getPoints("../data.txt")

	pm := make(PointMap)

	for _, p := range pts {
		pm[p] = Solid
	}

	surfaceArea := 0

	for p := range pm {
		if !pm.isAdjacentSolid(&p, 1, 0, 0) {
			surfaceArea++
		}
		if !pm.isAdjacentSolid(&p, -1, 0, 0) {
			surfaceArea++
		}
		if !pm.isAdjacentSolid(&p, 0, 1, 0) {
			surfaceArea++
		}
		if !pm.isAdjacentSolid(&p, 0, -1, 0) {
			surfaceArea++
		}
		if !pm.isAdjacentSolid(&p, 0, 0, 1) {
			surfaceArea++
		}
		if !pm.isAdjacentSolid(&p, 0, 0, -1) {
			surfaceArea++
		}
	}

	fmt.Println(surfaceArea)
}

func (pm PointMap) isAdjacentSolid(v *Vector, x, y, z int) bool {
	c := pm.GetFrom(v, x, y, z)
	if c == Solid {
		return true
	}
	return false
}

type PointMap map[Vector]PointType

func (pm PointMap) GetVector(v Vector) PointType {
	if t, exists := pm[v]; exists {
		return t
	}
	return Unknown
}

func (pm PointMap) Get(x, y, z int) PointType {
	return pm.GetVector(Vector{X: x, Y: y, Z: z})
}

func (pm PointMap) GetFrom(v *Vector, x, y, z int) PointType {
	t := v.Transform(x, y, z)
	return pm.GetVector(t)
}

func getPoints(path string) []Vector {
	lines, _ := file.GetLines(path)
	vs := make([]Vector, 0, len(lines))
	for _, line := range lines {
		s := strings.Split(line, ",")
		v := Vector{}
		v.X, _ = strconv.Atoi(s[0])
		v.Y, _ = strconv.Atoi(s[1])
		v.Z, _ = strconv.Atoi(s[2])
		vs = append(vs, v)
	}
	return vs
}

type Vector struct {
	X int
	Y int
	Z int
}

func (v Vector) Transform(x, y, z int) Vector {
	return Vector{X: v.X + x, Y: v.Y + y, Z: v.Z + z}
}
