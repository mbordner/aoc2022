package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"math"
	"strconv"
)

func main() {
	lines, _ := file.GetLines("../data.txt")
	sum := 0
	for _, line := range lines {
		sum += snafuToDec(line)
	}
	fmt.Println(decToSnafu(sum))
}

var (
	snafuValues = map[byte]int{'0': 0, '1': 1, '2': 2, '-': -1, '=': -2}
)

func pow(a, e int) int {
	return int(math.Pow(float64(a), float64(e)))
}

func snafuToDec(s string) int {
	val := 0
	for p, i := len(s)-1, 0; p >= 0; p, i = p-1, i+1 {
		val += pow(5, i) * snafuValues[s[p]]
	}
	return val
}

// 37 base 5 = 122
// 13 base 5 = 23
// 23442 base 5 = 1747

//	    2   stays
//	  1-2   4 becomes -1, and carries 1
//	 10-2   4 becomes -1, and carries 1
//	1-0-2   3 becomes -2, and carries 1
//
// 1=-0-2   2 stays, adds to 1 and becomes 3, switches to = carries 1
func decToSnafu(v int) string {
	t := strconv.FormatInt(int64(v), 5)
	s := ""
	carry := 0
	for i := len(t) - 1; i >= 0; i-- {
		prevCarry := carry
		b := int(t[i] - '0')
		if b > 2 {
			carry = 1
			if b == 4 {
				b = -1
			} else if b == 3 {
				b = -2
			}
		} else {
			carry = 0
		}
		b += prevCarry
		if b >= 0 {
			if b <= 2 {
				s = string(byte(b)+'0') + s
			} else {
				carry = 1
				if b == 3 {
					s = "=" + s
				} else if b == 4 {
					s = "-" + s
				}
			}
		} else {
			if b == -2 {
				s = "=" + s
			} else {
				s = "-" + s
			}
		}
	}
	if carry > 0 {
		s = string(byte(carry)+'0') + s
	}
	return s
}
