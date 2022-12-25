package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"sort"
	"strings"
)

const (
	gridSize = 50
)

/*
....#..
....#..
..###..
...#...
..###..
...#...
..####.
..###..
..###..
#.##.#.
#.#####
####.#.
#.####.
height: 13
*/

type memStats struct {
	r int64
	h int64
}

func main() {
	moves := getMoves("../data.txt")
	g := NewGrid()

	if len(moves) > 0 {

		numRocks := int64(1000000000000)
		nextShapeIndex := int64(0)
		nextMoveIndex := 0

		skippedHeightToAdd := int64(0)

		patternMap := make(map[string]memStats)

	nextRock:
		for r := int64(0); r < numRocks; r++ {
			var s *Shape
			s = getNextShape(nextShapeIndex)
			nextShapeIndex++
			g.MoveShapeToStart(s)

			for {
				var d Dir
				moveIndex := nextMoveIndex
				nextMoveIndex, d = getNextMove(moveIndex, moves)

				// jet movement, may not happen
				if g.CanMove(s, d) {
					switch d {
					case Right:
						s.Move(1, 0)
					case Left:
						s.Move(-1, 0)
					}
				}

				if g.CanMove(s, Down) {
					s.Move(0, -1)
				} else {
					if skippedHeightToAdd == 0 {
						patternDescription := g.getPatternDescription(s, nextMoveIndex)
						heightNow := g.GetShapesHeight()
						g.AddShape(s)
						if prevStats, exists := patternMap[patternDescription]; exists {
							numShapesInRepeatingGroup := r - prevStats.r
							repeatingShapesHeight := heightNow - prevStats.h

							skipCount := (numRocks - r) / numShapesInRepeatingGroup

							skippedHeightToAdd = skipCount * repeatingShapesHeight

							r += skipCount * numShapesInRepeatingGroup

							fmt.Println(prevStats.h, heightNow)
						} else {
							patternMap[patternDescription] = memStats{r: r, h: heightNow}
						}
					} else {
						g.AddShape(s) // fix this shape to the grid permanently
					}

					continue nextRock
				}
			}

		}

		g.Print()
		fmt.Println("height:", g.GetShapesHeight()+skippedHeightToAdd)

	}

}

func getNextMove(moveIndex int, moves []rune) (int, Dir) {
	var d Dir
	m := moves[moveIndex]
	nextMoveIndex := moveIndex + 1
	if nextMoveIndex == len(moves) {
		nextMoveIndex = 0
	}
	switch m {
	case '<':
		d = Left
	case '>':
		d = Right
	}
	return nextMoveIndex, d
}

func getNextShape(shapeIndex int64) *Shape {
	switch shapeIndex % 5 {
	case 0:
		return NewShape1()
	case 1:
		return NewShape2()
	case 2:
		return NewShape3()
	case 3:
		return NewShape4()
	case 4:
		return NewShape5()
	}
	return NewShape1()
}

func getMoves(path string) []rune {
	lines, _ := file.GetLines(path)
	return []rune(lines[0])
}

func max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

type Dir int

const (
	Left Dir = iota
	Right
	Down
	Up
)

type topStats struct {
	minX int64
	maxX int64
	tops map[int64]int64
}

func (ts *topStats) track(p *Point) {
	if p.X < ts.minX {
		ts.minX = p.X
	}
	if p.X > ts.maxX {
		ts.maxX = p.X
	}
	if v, e := ts.tops[p.X]; e {
		if p.Y > v {
			ts.tops[p.X] = p.Y
		}
	} else {
		ts.tops[p.X] = p.Y
	}
}

func (ts *topStats) getTop(x int64) int64 {
	if v, e := ts.tops[x]; e {
		return v
	}
	return 0
}

func (ts *topStats) getTopPoints() Points {
	count := abs(ts.maxX-ts.minX) + 1
	pss := make(Points, count, count)
	for x, i := ts.minX, int64(0); i < count; x, i = x+1, i+1 {
		pss[i] = &Point{X: x, Y: ts.getTop(x)}
	}
	sort.Sort(pss)
	return pss
}

func (ts *topStats) topDiff(o *Point) Points {
	pss := ts.getTopPoints()
	for i, p := range pss {
		pss[i] = o.Diff(p)
	}
	return pss
}

func NewTopStats() *topStats {
	ts := new(topStats)
	ts.tops = make(map[int64]int64)
	return ts
}

type GridShapeCollection []*Shape

func (gsc *GridShapeCollection) GetMaxShapeY() int64 {
	if len(*gsc) > 0 {
		return (*gsc)[0].Bounds().Max.Y
	}
	return -1
}

func (gsc *GridShapeCollection) Add(s *Shape) {
	*gsc = append(*gsc, s)
	sort.Sort(gsc)
}

func (gsc *GridShapeCollection) Remove(s *Shape) {
	for i := range *gsc {
		if (*gsc)[i] == s {
			*gsc = append((*gsc)[0:i], (*gsc)[i+1:]...)
			break
		}
	}
}

func (gsc *GridShapeCollection) Len() int {
	return len(*gsc)
}

// Less ensures that the shapes with highest Y extents are first
func (gsc *GridShapeCollection) Less(a, b int) bool {
	aBounds := (*gsc)[a].Bounds()
	bBounds := (*gsc)[b].Bounds()
	if bBounds.Max.Y > aBounds.Max.Y {
		return false
	}
	return true
}

func (gsc *GridShapeCollection) Swap(a, b int) {
	(*gsc)[a], (*gsc)[b] = (*gsc)[b], (*gsc)[a]
}

type BoundsCollection []Bounds

func (bc *BoundsCollection) Add(b Bounds) {
	*bc = append(*bc, b)
	sort.Sort(bc)
}

func (bc *BoundsCollection) Len() int {
	return len(*bc)
}

func (bc *BoundsCollection) Swap(a, b int) {
	(*bc)[a], (*bc)[b] = (*bc)[b], (*bc)[a]
}

// Less ensures that the bounds with the highest Y extends are first
func (bc *BoundsCollection) Less(a, b int) bool {
	if (*bc)[b].Max.Y > (*bc)[a].Max.Y {
		return false
	}
	return true
}

type Grid struct {
	tops      *topStats
	bounds    BoundsCollection
	boundsMap map[Bounds]*GridShapeCollection
}

func NewGrid() *Grid {
	g := new(Grid)
	g.tops = NewTopStats()
	g.bounds = make(BoundsCollection, 0, gridSize)
	g.boundsMap = make(map[Bounds]*GridShapeCollection)
	return g
}

func (g *Grid) getPatternDescription(s *Shape, nextMoveIndex int) string {
	o := s.Origin()
	return fmt.Sprintf("%s|%d|%d|%s", string(s.char), nextMoveIndex, o.X, g.tops.topDiff(&o).String())
}

func (g *Grid) Print() {
	bounds := Bounds{}
	bounds.Min.X = 2147483647
	bounds.Min.Y = 2147483647
	points := make(map[Point]rune)

	for _, b := range g.bounds {
		c := g.boundsMap[b]
		for _, s := range *c {
			for _, p := range s.points {
				point := *p
				points[point] = s.char
				if p.X < bounds.Min.X {
					bounds.Min.X = p.X
				}
				if p.X > bounds.Max.X {
					bounds.Max.X = p.X
				}
				if p.Y < bounds.Min.Y {
					bounds.Min.Y = p.Y
				}
				if p.Y > bounds.Max.Y {
					bounds.Max.Y = p.Y
				}
			}
		}
		if len(*c) > 0 {
			break
		}
	}

	for y := bounds.Max.Y; y >= bounds.Min.Y; y-- {
		l := abs(bounds.Max.X-bounds.Min.X) + 1
		line := make([]rune, l, l)
		for x, l := bounds.Min.X, 0; x <= bounds.Max.X; x, l = x+1, l+1 {
			if r, exists := points[Point{X: x, Y: y}]; exists {
				line[l] = r
			} else {
				line[l] = '.'
			}
		}
		fmt.Println(string(line))
	}
}

func (g *Grid) CanMove(s *Shape, d Dir) bool {
	dx, dy := int64(0), int64(0)
	switch d {
	case Left:
		dx = -1
	case Right:
		dx = 1
	case Down:
		dy = -1
	case Up:
		dy = 1
	}
	ts := s.Transform(dx, dy)

	tsBounds := ts.Bounds()
	if tsBounds.Min.X <= -3 { // would hit left wall
		return false
	}

	if tsBounds.Max.X >= 5 { // would hit right wall
		return false
	}

	if tsBounds.Min.Y < 0 { // would hit floor
		return false
	}

	// find bounds to check.
	// the origin of our shape is grabbed first
	// but if the shape extends up, right or up and right past this bounds
	// grab the other 3 to check as well.
	// once we have the bounds to check, we can get the collections and scan
	// for collisions
	checkGridBounds := make([]Bounds, 0, 4)
	tsOriginBounds := g.GetGridBounds(ts)
	checkGridBounds = append(checkGridBounds, tsOriginBounds)

	if tsBounds.Max.Y > tsOriginBounds.Max.Y {
		checkGridBounds = append(checkGridBounds,
			g.GetGridBoundsForPoint(Point{X: tsBounds.Min.X, Y: tsBounds.Max.Y}))
	}

	if tsBounds.Max.X > tsOriginBounds.Max.X {
		checkGridBounds = append(checkGridBounds,
			g.GetGridBoundsForPoint(Point{X: tsBounds.Max.X, Y: tsBounds.Min.Y}))
	}

	if tsBounds.Max.Y > tsOriginBounds.Max.Y && tsBounds.Max.X > tsOriginBounds.Max.X {
		checkGridBounds = append(checkGridBounds,
			g.GetGridBoundsForPoint(Point{X: tsBounds.Max.X, Y: tsBounds.Max.Y}))
	}

	for cgbIndex := range checkGridBounds {
		gridBounds := checkGridBounds[cgbIndex]
		for _, o := range *g.boundsMap[gridBounds] {
			if ts.Collides(o) {
				return false
			}
		}
	}

	return true
}

func (g *Grid) GetMaxShapeY() int64 {
	boundsYMax := int64(-1)
	for _, b := range g.bounds {
		if len(*g.boundsMap[b]) > 0 {
			boundsYMax = b.Max.Y
			break
		}
	}
	shapeYMax := int64(-1)
	for i := range g.bounds {
		if g.bounds[i].Max.Y < boundsYMax {
			break
		}
		shapeCol := g.boundsMap[g.bounds[i]]
		shapeYMax = max(shapeYMax, shapeCol.GetMaxShapeY())
	}
	return shapeYMax
}

func (g *Grid) GetShapesHeight() int64 {
	return g.GetMaxShapeY() + 1
}

// MoveShapeToStart controls the logic where new shapes should spawn, so it will move new shapes where they belong
func (g *Grid) MoveShapeToStart(s *Shape) {
	if len(g.bounds) == 0 {
		b := Bounds{Min: Point{X: 0, Y: 0}, Max: Point{X: gridSize, Y: gridSize}}
		g.TrackNewBounds(b)
	}

	shapeYMax := g.GetMaxShapeY()

	// There is a wall at X == -2
	// Floor is Y == 0 (-1, shapes stop at 0)
	// There is a wall at X == 4
	// width of space == 7 between walls
	// origin Y start is at 3 above maxY (could be floor, so would be 3)
	s.Move(0, shapeYMax+int64(4))
}

func (g *Grid) AddShape(s *Shape) {
	boundsToAdd := make([]Bounds, 0, 4)
	shapeBounds := s.Bounds()
	shapeOriginBounds := g.GetGridBounds(s)
	boundsToAdd = append(boundsToAdd, shapeOriginBounds)

	if shapeBounds.Max.Y > shapeOriginBounds.Max.Y {
		boundsToAdd = append(boundsToAdd,
			g.GetGridBoundsForPoint(Point{X: shapeBounds.Min.X, Y: shapeBounds.Max.Y}))
	}

	if shapeBounds.Max.X > shapeOriginBounds.Max.X {
		boundsToAdd = append(boundsToAdd,
			g.GetGridBoundsForPoint(Point{X: shapeBounds.Max.X, Y: shapeBounds.Min.Y}))
	}

	if shapeBounds.Max.Y > shapeOriginBounds.Max.Y && shapeBounds.Max.X > shapeOriginBounds.Max.X {
		boundsToAdd = append(boundsToAdd,
			g.GetGridBoundsForPoint(Point{X: shapeBounds.Max.X, Y: shapeBounds.Max.Y}))
	}

	for _, b := range boundsToAdd {
		g.boundsMap[b].Add(s)
	}

	for _, p := range s.points {
		g.tops.track(p)
	}

	topPoints := g.tops.getTopPoints()
	minY := shapeBounds.Min.Y
	for _, p := range topPoints {
		if p.Y < minY {
			minY = p.Y
		}
	}
	minGridBounds := g.GetGridBoundsForPoint(Point{X: shapeBounds.Min.X, Y: minY})
	var indexToDrop int
	for indexToDrop < len(g.bounds) && g.bounds[indexToDrop].Min.Y >= minGridBounds.Min.Y {
		indexToDrop++
	}
	if indexToDrop < len(g.bounds) {
		toDrop := g.bounds[indexToDrop:]
		g.bounds = g.bounds[0:indexToDrop]
		for _, b := range toDrop {
			delete(g.boundsMap, b)
		}
	}
}

func (g *Grid) GetGridBoundsForPoint(p Point) Bounds {
	bXMin := p.X / gridSize * gridSize
	bYMin := p.Y / gridSize * gridSize
	bXMax := bXMin + gridSize - 1
	bYMax := bYMin + gridSize - 1

	b := Bounds{Min: Point{X: bXMin, Y: bYMin}, Max: Point{X: bXMax, Y: bYMax}}
	if _, e := g.boundsMap[b]; !e {
		g.TrackNewBounds(b)
	}

	return b
}

func (g *Grid) GetGridBounds(s *Shape) Bounds {
	return g.GetGridBoundsForPoint(s.Origin())
}

func (g *Grid) TrackNewBounds(b Bounds) {
	g.bounds.Add(b)
	col := make(GridShapeCollection, 0, gridSize)
	g.boundsMap[b] = &col
}

type Bounds struct {
	Min Point
	Max Point
}

func (b Bounds) Contains(p *Point) bool {
	if p.X >= b.Min.X && p.X <= b.Max.X && p.Y >= b.Min.Y && p.Y <= b.Max.Y {
		return true
	}
	return false
}

type Points []*Point

func (ps Points) Len() int {
	return len(ps)
}

func (ps Points) Less(a, b int) bool {
	if ps[a].X < ps[b].X {
		return true
	}
	return false
}

func (ps Points) Swap(a, b int) {
	ps[a], ps[b] = ps[b], ps[a]
}

func (ps Points) Transform(x, y int64) Points {
	tPS := make([]*Point, len(ps), len(ps))
	for i := range ps {
		tPS[i] = ps[i].Transform(x, y)
	}
	return tPS
}

func (ps Points) String() string {
	ss := make([]string, len(ps), len(ps))
	for i := range ps {
		ss[i] = ps[i].String()
	}
	return strings.Join(ss, ";")
}

type Point struct {
	X int64
	Y int64
}

func (p *Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p *Point) Transform(x, y int64) *Point {
	return &Point{X: p.X + x, Y: p.Y + y}
}

func (p *Point) Diff(o *Point) *Point {
	d := Point{X: p.X - o.X, Y: p.Y - o.Y}
	return &d
}

type Shape struct {
	boundsOffset Bounds
	points       Points
	char         rune
}

func (s *Shape) Bounds() Bounds {
	return Bounds{Min: *(s.points[0]).Transform(s.boundsOffset.Min.X, s.boundsOffset.Min.Y),
		Max: *(s.points[len(s.points)-1]).Transform(s.boundsOffset.Max.X, s.boundsOffset.Max.Y)}
}

func (s *Shape) BoundsOverlaps(ob *Bounds) bool {
	b := s.Bounds()
	if b.Contains(&ob.Min) || b.Contains(&ob.Max) ||
		b.Contains(&Point{X: ob.Min.X, Y: ob.Max.Y}) ||
		b.Contains(&Point{X: ob.Max.X, Y: ob.Min.Y}) {
		return true
	}
	return false
}

func (s *Shape) ShapeBoundsOverlaps(o *Shape) bool {
	ob := o.Bounds()
	return s.BoundsOverlaps(&ob)
}

func (s *Shape) Collides(o *Shape) bool {
	pMap := make(map[Point]bool)
	for _, sp := range s.points {
		pMap[*sp] = true
	}
	for _, op := range o.points {
		if _, exists := pMap[*op]; exists {
			return true
		}
	}
	return false
}

func (s *Shape) MoveTo(ps Points) {
	s.points = ps
}

func (s *Shape) Move(x, y int64) {
	s.MoveTo(s.points.Transform(x, y))
}

func (s *Shape) Origin() Point {
	return s.Bounds().Min
}

func (s *Shape) Transform(x, y int64) *Shape {
	o := new(Shape)
	o.boundsOffset = s.boundsOffset
	o.points = s.points.Transform(x, y)
	return o
}

func (s *Shape) MoveToTransformedShape(o *Shape) {
	s.MoveTo(o.points)
}

/*
*
####
*/
func NewShape1() *Shape {
	s := new(Shape)
	s.char = '1'
	s.boundsOffset = Bounds{Min: Point{X: 0, Y: 0}, Max: Point{X: 0, Y: 0}}

	s.points = make([]*Point, 4, 4)
	s.points[0] = &Point{X: 0, Y: 0}
	s.points[1] = &Point{X: 1, Y: 0}
	s.points[2] = &Point{X: 2, Y: 0}
	s.points[3] = &Point{X: 3, Y: 0}

	return s
}

/*
.#.
###
.#.
*/
func NewShape2() *Shape {
	s := new(Shape)
	s.char = '2'
	s.boundsOffset = Bounds{Min: Point{X: -1, Y: 0}, Max: Point{X: 1, Y: 0}}

	s.points = make([]*Point, 5, 5)
	s.points[0] = &Point{X: 1, Y: 0}
	s.points[1] = &Point{X: 0, Y: 1}
	s.points[2] = &Point{X: 1, Y: 1}
	s.points[3] = &Point{X: 2, Y: 1}
	s.points[4] = &Point{X: 1, Y: 2}

	return s
}

/*
..#
..#
###
*/
func NewShape3() *Shape {
	s := new(Shape)
	s.char = '3'
	s.boundsOffset = Bounds{Min: Point{X: 0, Y: 0}, Max: Point{X: 0, Y: 0}}

	s.points = make([]*Point, 5, 5)
	s.points[0] = &Point{X: 0, Y: 0}
	s.points[1] = &Point{X: 1, Y: 0}
	s.points[2] = &Point{X: 2, Y: 0}
	s.points[3] = &Point{X: 2, Y: 1}
	s.points[4] = &Point{X: 2, Y: 2}

	return s
}

/*
#
#
#
#
*/
func NewShape4() *Shape {
	s := new(Shape)
	s.char = '4'
	s.boundsOffset = Bounds{Min: Point{X: 0, Y: 0}, Max: Point{X: 0, Y: 0}}

	s.points = make([]*Point, 4, 4)
	s.points[0] = &Point{X: 0, Y: 0}
	s.points[1] = &Point{X: 0, Y: 1}
	s.points[2] = &Point{X: 0, Y: 2}
	s.points[3] = &Point{X: 0, Y: 3}

	return s
}

/*
##
##
*/
func NewShape5() *Shape {
	s := new(Shape)
	s.char = '5'
	s.boundsOffset = Bounds{Min: Point{X: 0, Y: 0}, Max: Point{X: 0, Y: 0}}

	s.points = make([]*Point, 4, 4)
	s.points[0] = &Point{X: 0, Y: 0}
	s.points[1] = &Point{X: 1, Y: 0}
	s.points[2] = &Point{X: 0, Y: 1}
	s.points[3] = &Point{X: 1, Y: 1}

	return s
}
