package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"github.com/mbordner/aoc2022/common/geom"
	"regexp"
	"strconv"
)

type Area map[geom.Pos]rune

type AreaMap struct {
	area         Area
	sensorBeacon []geom.Pos
	sensor       []geom.Pos
	distance     []int
}

func (am *AreaMap) print() {
	fmt.Println(fmt.Sprintf("min: %d,%d   max: %d,%d", bb.MinX, bb.MinY, bb.MaxX, bb.MaxY))

	for y := bb.MinY; y <= bb.MaxY; y++ {
		l := bb.MaxX - bb.MinX + 1
		line := make([]rune, l, l)
		for x, l := bb.MinX, 0; x <= bb.MaxX; x, l = x+1, l+1 {
			if r, exists := am.area[geom.Pos{X: x, Y: y, Z: 0}]; exists {
				line[l] = r
			} else {
				line[l] = '.'
			}
		}
		fmt.Println(fmt.Sprintf("%5d", y), string(line))
	}
}

var (
	bb       = &geom.BoundingBox{}
	reSensor = regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
)

func main() {
	am := getAreaMap("../data.txt")
	//am.print()

	area := checkRow(2000000, am)

	count := 0
	for _, r := range area {
		if r == '#' {
			count++
		}
	}

	fmt.Println(count)
}

func checkRow(row int, am *AreaMap) Area {
	checkArea := make(Area)
	for i := range am.sensor {
		if am.sensor[i].Y == row {
			checkArea[am.sensor[i]] = 'S'
		}
		if am.sensorBeacon[i].Y == row {
			checkArea[am.sensorBeacon[i]] = 'B'
		}
		rowDistance := geom.Abs(row - am.sensor[i].Y)
		if am.distance[i] >= rowDistance {
			extend := am.distance[i] - rowDistance

			for e := 0; e <= extend; e++ {
				checkP := geom.Pos{X: am.sensor[i].X - e, Y: row, Z: 0}
				if _, exists := checkArea[checkP]; !exists {
					checkArea[checkP] = '#'
				}
				checkP = geom.Pos{X: am.sensor[i].X + e, Y: row, Z: 0}
				if _, exists := checkArea[checkP]; !exists {
					checkArea[checkP] = '#'
				}
			}
		}
	}
	return checkArea
}

func getAreaMap(path string) *AreaMap {

	lines, _ := file.GetLines(path)

	am := &AreaMap{
		area:         make(Area),
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
				am.area[beacon] = 'B'
				am.area[sensor] = 'S'

				dis := sensor.ManhattanDistance(beacon)
				am.distance[i] = dis

				/*
					positions := sensor.GetXYPositionsWithinManhattanDistance(dis)

					for _, p := range positions {
						if _, exists := am.area[p]; !exists {
							am.area[p] = '#'
							bb.Extend(p)
						}
					}
				*/
			}
		}
	}

	return am
}
