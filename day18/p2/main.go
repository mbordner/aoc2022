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
	OutSideAir
	InsideAir
)

func main() {
	pts := getPoints("../data.txt")

	pm := make(PointMap)
	outside := make(PointMap)
	solids := make(PointMap)

	minX, minY, minZ, maxX, maxY, maxZ := 0, 0, 0, 0, 0, 0

	for _, p := range pts {
		pm[p] = Solid
		solids[p] = Solid
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Z < minZ {
			minZ = p.Z
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.Z > maxZ {
			maxZ = p.Z
		}
	}

	minX -= 2
	minY -= 2
	minZ -= 2
	maxX += 2
	maxY += 2
	maxZ += 2

	for z := minZ; z <= maxZ; z++ {
		for y := minY; y <= maxY; y++ {

			x := minX
			for pm.Get(x, y, z) == Unknown && x <= maxX {
				v := Vector{X: x, Y: y, Z: z}
				pm[v] = OutSideAir
				outside[v] = OutSideAir
				x++
			}

			x = maxX
			for pm.Get(x, y, z) == Unknown && x >= 0 {
				v := Vector{X: x, Y: y, Z: z}
				pm[v] = OutSideAir
				outside[v] = OutSideAir
				x--
			}

		}
	}

	var fillAir func(v Vector)
	fillAir = func(v Vector) {
		pm[v] = OutSideAir
		outside[v] = OutSideAir
		if pm.isAdjacentUnknown(&v, 1, 0, 0) {
			fillAir(v.Transform(1, 0, 0))
		}
		if pm.isAdjacentUnknown(&v, -1, 0, 0) {
			fillAir(v.Transform(-1, 0, 0))
		}
		if pm.isAdjacentUnknown(&v, 0, 1, 0) {
			fillAir(v.Transform(0, 1, 0))
		}
		if pm.isAdjacentUnknown(&v, 0, -1, 0) {
			fillAir(v.Transform(0, -1, 0))
		}
		if pm.isAdjacentUnknown(&v, 0, 0, 1) {
			fillAir(v.Transform(0, 0, 1))
		}
		if pm.isAdjacentUnknown(&v, 0, 0, -1) {
			fillAir(v.Transform(0, 0, -1))
		}
	}

	for z := minZ + 1; z < maxZ; z++ {
		for y := minY + 1; y < maxY; y++ {
			for x := minX + 1; x < maxX; x++ {
				t := pm.Get(x, y, z)
				if t == OutSideAir {
					fillAir(Vector{X: x, Y: y, Z: z})
				}
			}
		}
	}

	surfaceArea := 0

	for p := range solids {
		if t := pm.GetFrom(&p, 1, 0, 0); t == OutSideAir {
			surfaceArea++
		}
		if t := pm.GetFrom(&p, -1, 0, 0); t == OutSideAir {
			surfaceArea++
		}
		if t := pm.GetFrom(&p, 0, 1, 0); t == OutSideAir {
			surfaceArea++
		}
		if t := pm.GetFrom(&p, 0, -1, 0); t == OutSideAir {
			surfaceArea++
		}
		if t := pm.GetFrom(&p, 0, 0, 1); t == OutSideAir {
			surfaceArea++
		}
		if t := pm.GetFrom(&p, 0, 0, -1); t == OutSideAir {
			surfaceArea++
		}
	}

	fmt.Println(surfaceArea)
}

func (pm PointMap) isAdjacentType(t PointType, v *Vector, x, y, z int) bool {
	c := pm.GetFrom(v, x, y, z)
	if c == t {
		return true
	}
	return false
}

func (pm PointMap) isAdjacentSolid(v *Vector, x, y, z int) bool {
	return pm.isAdjacentType(Solid, v, x, y, z)
}

func (pm PointMap) isAdjacentUnknown(v *Vector, x, y, z int) bool {
	return pm.isAdjacentType(Unknown, v, x, y, z)
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
