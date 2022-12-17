package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"github.com/mbordner/aoc2022/common/geom"
	"github.com/mbordner/aoc2022/common/ranges"
	"regexp"
	"strconv"
)

type AreaMap struct {
	sensorBeacon []geom.Pos
	sensor       []geom.Pos
	distance     []int
}

var (
	bb       = &geom.BoundingBox{}
	reSensor = regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
)

func main() {
	limit := 4000000
	am := getAreaMap("../data.txt")

outer:
	for row := 0; row <= limit; row++ {

		rc := &ranges.Collection[int]{}

		for i := range am.sensor {
			s := am.sensor[i]
			d := am.distance[i]

			yDistance := geom.Max(row, s.Y) - geom.Min(row, s.Y)

			if yDistance <= d {
				dx := d - yDistance
				if (s.X-dx >= 0 && s.X-dx <= limit) || (s.X+dx >= 0 && s.X+dx <= limit) {
					left := geom.Max(0, s.X-dx)
					right := geom.Min(limit, s.X+dx)
					_, _ = rc.Add(left, right)
				}
			}
		}

		l := rc.Len()
		if l < limit+1 {

			values := rc.ValuePairs()
			if values[0] == 1 {
				fmt.Println(0, row)
				fmt.Println(row)
			} else if values[len(values)-1] == limit-1 {
				fmt.Println(limit, row)
				fmt.Println(uint64(limit)*uint64(limit) + uint64(row))
			} else {
				for i := 2; i < len(values); i += 2 {
					if values[i]-2 == values[i-1] {
						fmt.Println(values[i]-1, row)
						fmt.Println(uint64(values[i]-1)*uint64(limit) + uint64(row))
						break outer
					}
				}
			}
		}
	}

}

func getAreaMap(path string) *AreaMap {

	lines, _ := file.GetLines(path)

	am := &AreaMap{
		sensorBeacon: make([]geom.Pos, len(lines), len(lines)),
		sensor:       make([]geom.Pos, len(lines), len(lines)),
		distance:     make([]int, len(lines), len(lines)),
	}

	for i, line := range lines {
		if reSensor.MatchString(line) {
			matches := reSensor.FindStringSubmatch(line)
			if len(matches) == 5 {
				sensor := geom.Pos{}
				sensor.X, _ = strconv.Atoi(matches[1])
				sensor.Y, _ = strconv.Atoi(matches[2])
				beacon := geom.Pos{}
				beacon.X, _ = strconv.Atoi(matches[3])
				beacon.Y, _ = strconv.Atoi(matches[4])
				am.sensor[i] = sensor
				am.sensorBeacon[i] = beacon
				bb.Extend(beacon)
				bb.Extend(sensor)

				dis := sensor.ManhattanDistance(beacon)
				am.distance[i] = dis
			}
		}
	}

	return am
}
