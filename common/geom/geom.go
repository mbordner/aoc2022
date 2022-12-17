package geom

import (
	"fmt"
	"math"
	"strings"
)

type Direction int

const (
	Unknown Direction = 0
	North   Direction = 1
	South   Direction = 2
	West    Direction = 4
	East    Direction = 8
)

func (d Direction) Is(dirs []Direction) bool {
	for _, dir := range dirs {
		if int(d)&int(dir) == int(dir) {
			return true
		}
	}
	return false
}

func (d Direction) Opposite() Direction {
	o := 0
	if int(d)&int(North) == int(North) {
		o |= int(South)
	}
	if int(d)&int(West) == int(West) {
		o |= int(East)
	}
	if int(d)&int(South) == int(South) {
		o |= int(North)
	}
	if int(d)&int(East) == int(East) {
		o |= int(West)
	}
	return Direction(o)
}

type BoundingBox struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
	MinZ int
	MaxZ int
}

func (bb *BoundingBox) SetExtents(x1, y1, z1, x2, y2, z2 int) {
	bb.MinX = x1
	bb.MinY = y1
	bb.MinZ = z1
	bb.MaxX = x2
	bb.MaxY = y2
	bb.MaxZ = z2
}

func (bb BoundingBox) XMin() int {
	return bb.MinX
}

func (bb BoundingBox) XMax() int {
	return bb.MaxX
}

func (bb BoundingBox) YMin() int {
	return bb.MinY
}

func (bb BoundingBox) YMax() int {
	return bb.MaxY
}

func (bb BoundingBox) ZMin() int {
	return bb.MinZ
}

func (bb BoundingBox) ZMax() int {
	return bb.MaxZ
}

func (bb BoundingBox) String() string {
	p1 := Pos{X: bb.MinX, Y: bb.MinY}
	p2 := Pos{X: bb.MaxX, Y: bb.MaxY}
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

func (bb *BoundingBox) Extend(p Pos) {
	if p.X < bb.MinX {
		bb.MinX = p.X
	}
	if p.X > bb.MaxX {
		bb.MaxX = p.X
	}
	if p.Y > bb.MaxY {
		bb.MaxY = p.Y
	}
	if p.Y < bb.MinY {
		bb.MinY = p.Y
	}
	if p.Z < bb.MinZ {
		bb.MinZ = p.Z
	}
	if p.Z > bb.MaxZ {
		bb.MaxZ = p.Z
	}
}

func (bb *BoundingBox) Contains(p Pos) bool {
	if p.X < bb.MinX || p.X > bb.MaxX {
		return false
	}
	if p.Y < bb.MinY || p.Y > bb.MaxY {
		return false
	}
	if p.Z < bb.MinZ || p.Z > bb.MaxZ {
		return false
	}
	return true
}

func (bb *BoundingBox) Surrounds(obb *BoundingBox) bool {
	if obb.MinX < bb.MinX {
		return false
	}
	if obb.MaxX > bb.MaxX {
		return false
	}
	if obb.MinY < bb.MinY {
		return false
	}
	if obb.MaxY > bb.MaxY {
		return false
	}
	if obb.MinZ < bb.MinZ {
		return false
	}
	if obb.MaxZ > bb.MaxZ {
		return false
	}
	return true
}

func (bb *BoundingBox) GetDirection(p Pos) Direction {
	dir := 0
	if p.X < bb.MinX {
		dir |= int(West)
	}
	if p.X > bb.MaxX {
		dir |= int(East)
	}
	if p.Y > bb.MaxY {
		dir |= int(North)
	}
	if p.Y < bb.MinY {
		dir |= int(South)
	}
	return Direction(dir)
}

func (bb *BoundingBox) Intersects(p1, p2 Pos) bool {
	return false
}

func (bb *BoundingBox) GetPrintLines(defaultChar rune, chars []rune, pss Positions) []string {
	lines := make([]string, 0, bb.MaxY-bb.MinY+1)

	charMap := make(map[Pos]rune)
	if chars != nil && pss != nil && len(chars) == len(pss) {
		for i := range chars {
			pss[i].Z = 0
			charMap[pss[i]] = chars[i]
		}
	}
	for i := bb.MaxY; i >= bb.MinY; i-- {
		l := bb.MaxX - bb.MinX + 1
		lineRunes := make([]rune, l, l)
		for j, c := bb.MinX, 0; j <= bb.MaxX; j, c = j+1, c+1 {
			lineRunes[c] = defaultChar
			if char, exists := charMap[Pos{X: j, Y: i, Z: 0}]; exists {
				lineRunes[c] = char
			}
		}
		lines = append(lines, string(lineRunes))
	}

	return lines
}

func (bb *BoundingBox) DistanceFromEdge(p Pos) int {
	d := math.MaxInt64

	t := bb.MaxX - p.X
	if t < d {
		d = t
	}

	t = p.X - bb.MinX
	if t < d {
		d = t
	}

	t = p.Y - bb.MinY
	if t < d {
		d = t
	}

	t = bb.MaxY - p.Y
	if t < d {
		d = t
	}

	return d
}

type Positions []Pos

type Pos struct {
	X int
	Y int
	Z int
}

func (p Pos) Transform(x, y, z int) Pos {
	return Pos{X: p.X + x, Y: p.Y + y, Z: p.Z + z}
}

func (p Pos) Diff(o Pos) Pos {
	v := Pos{
		X: p.X - o.X,
		Y: p.Y - o.Y,
		Z: p.Z - o.Z,
	}
	return v
}

func (p Pos) ManhattanDistance(o Pos) int {
	return Max(p.X, o.X) - Min(p.X, o.X) +
		Max(p.Y, o.Y) - Min(p.Y, o.Y) +
		Max(p.Z, o.Z) - Min(p.Z, o.Z)
}

func (p Pos) GetXYPositionsAtManhattanDistance(d int) Positions {
	ps := make(Positions, 0, ((d-1)*4)+4)

	left := Pos{X: p.X - d, Y: p.Y}
	right := Pos{X: p.X + d, Y: p.Y}
	top := Pos{X: p.X, Y: p.Y - d}
	bottom := Pos{X: p.X, Y: p.Y + d}

	ps = append(ps, Positions{top, bottom, left, right}...)

	for j := 1; j < d; j++ {
		dx := d - j
		topY, bottomY := p.Y-j, p.Y+j
		topLeft := Pos{X: p.X - dx, Y: topY}
		topRight := Pos{X: p.X + dx, Y: topY}
		bottomLeft := Pos{X: p.X - dx, Y: bottomY}
		bottomRight := Pos{X: p.X + dx, Y: bottomY}

		ps = append(ps, Positions{topLeft, topRight, bottomLeft, bottomRight}...)
	}

	return ps
}

// GetXYPositionsWithinManhattanDistance returns positions within distance in x,y plane
func (p Pos) GetXYPositionsWithinManhattanDistance(d int) Positions {
	pm := make(map[Pos]bool)

	for j := 0; j <= d; j++ {
		for i := 0; i <= d-j; i++ {
			topY, bottomY := p.Y-j, p.Y+j
			topLeft := Pos{X: p.X - i, Y: topY}
			topRight := Pos{X: p.X + i, Y: topY}
			bottomLeft := Pos{X: p.X - i, Y: bottomY}
			bottomRight := Pos{X: p.X + i, Y: bottomY}
			pm[topLeft] = true
			pm[topRight] = true
			pm[bottomLeft] = true
			pm[bottomRight] = true
		}
	}

	ps := make(Positions, 0, len(pm))
	for p := range pm {
		ps = append(ps, p)
	}
	return ps
}

func (p Pos) Clone() Pos {
	return Pos{X: p.X, Y: p.Y, Z: p.Z}
}

func (ps Positions) String() string {
	strs := make([]string, 0, len(ps))
	for _, p := range ps {
		strs = append(strs, p.String())
	}
	return strings.Join(strs, ",")
}

func (p Pos) String() string {
	return fmt.Sprintf("{x:%d, y:%d, z:%d}", p.X, p.Y, p.Z)
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (bb *BoundingBox) GetPositionsSize() uint64 {
	xs := Abs(bb.MaxX-bb.MinX) + 1
	ys := Abs(bb.MaxY-bb.MinY) + 1
	zs := Abs(bb.MaxZ-bb.MinZ) + 1
	return uint64(xs) * uint64(ys) * uint64(zs)
}

func (bb *BoundingBox) GetPositions() Positions {
	poss := make(Positions, 0, ((bb.MaxX-bb.MinX)+1)*((bb.MaxY-bb.MinY)+1*((bb.MaxZ-bb.MinZ)+1)))
	for z := bb.MinZ; z <= bb.MaxZ; z++ {
		for y := bb.MinY; y <= bb.MaxY; y++ {
			for x := bb.MinX; x <= bb.MaxX; x++ {
				poss = append(poss, Pos{Z: z, Y: y, X: x})
			}
		}
	}
	return poss
}

func (ps *Positions) Transform(x, y, z int) Positions {
	for i := 0; i < len(*ps); i++ {
		(*ps)[i].X += x
		(*ps)[i].Y += y
		(*ps)[i].Z += z
	}
	return *ps
}
