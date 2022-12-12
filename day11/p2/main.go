package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/bigexpression"
	"github.com/mbordner/aoc2022/common/file"
	"math/big"
	"regexp"
	"sort"
	"strconv"
)

var (
	reMonkey           = regexp.MustCompile(`^Monkey (\d+):`)
	reStartingItems    = regexp.MustCompile(`  Starting items: `)
	reItemNumbers      = regexp.MustCompile(`(\d+)`)
	reOperation        = regexp.MustCompile(`  Operation: new = (.*)`)
	reTestDivisibility = regexp.MustCompile(`^  Test: divisible by (\d+)`)
	reIfTrue           = regexp.MustCompile(`^    If true: throw to monkey (\d+)`)
	reIfFalse          = regexp.MustCompile(`^    If false: throw to monkey (\d+)`)
)

var (
	monkeyMap = make(map[int]*Monkey)
)

type Monkey struct {
	id                          int
	items                       []*big.Int
	expr                        *bigexpression.Parser
	testDivisibility            int64
	testDivisibilityTrueTarget  int
	testDivisibilityFalseTarget int
	limitModulo                 *big.Int

	inspectCount int
}

func (m *Monkey) addItem(item *big.Int) {
	m.items = append(m.items, item)
}

func (m *Monkey) inspectItems() {
	items := m.items
	m.items = make([]*big.Int, 0, 20)
	for i := 0; i < len(items); i++ {
		m.inspectCount++
		opValue := m.expr.Eval(map[string]*big.Int{"old": items[i]})
		newValue := opValue.Mod(opValue, m.limitModulo) // / int64(3)
		if big.NewInt(0).Mod(newValue, big.NewInt(m.testDivisibility)).Int64() == 0 {
			if tm, exists := monkeyMap[m.testDivisibilityTrueTarget]; exists {
				tm.addItem(opValue)
			} else {
				panic("huh?")
			}
		} else {
			if tm, exists := monkeyMap[m.testDivisibilityFalseTarget]; exists {
				tm.addItem(opValue)
			} else {
				panic("huh??")
			}
		}
	}
}

func main() {
	monkeys := getMonkeys("../data.txt")
	if len(monkeys) > 0 {
		for i := range monkeys {
			monkeyMap[monkeys[i].id] = monkeys[i]
		}
	} else {
		panic("wha?")
	}

	rounds := 0
	for rounds < 10000 {

		for i := range monkeys {
			monkeys[i].inspectItems()
		}

		rounds++

		counts := make([]int, len(monkeys), len(monkeys))
		for i := range monkeys {
			counts[i] = monkeys[i].inspectCount
		}

	}

	counts := make([]int, len(monkeys), len(monkeys))
	for i := range monkeys {
		counts[i] = monkeys[i].inspectCount
	}

	fmt.Println(counts)

	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	fmt.Println(counts[0] * counts[1])
}

func getMonkeys(path string) []*Monkey {

	lines, _ := file.GetLines(path)

	monkeys := make([]*Monkey, 0, 20)

	var curMonkey *Monkey

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if reMonkey.MatchString(line) {
			if curMonkey != nil {
				monkeys = append(monkeys, curMonkey)
			}
			curMonkey = &Monkey{}
			matches := reMonkey.FindStringSubmatch(line)
			if len(matches) == 2 {
				curMonkey.id, _ = strconv.Atoi(matches[1])
			}
		} else if reStartingItems.MatchString(line) {
			matches := reItemNumbers.FindAllStringSubmatch(line, -1)
			if len(matches) > 0 {
				curMonkey.items = make([]*big.Int, len(matches), len(matches))
				for i := range matches {
					int64Val, _ := strconv.ParseInt(matches[i][0], 10, 64)
					curMonkey.items[i] = big.NewInt(int64Val)
				}
			}
		} else if reOperation.MatchString(line) {
			matches := reOperation.FindStringSubmatch(line)
			if len(matches) == 2 {
				curMonkey.expr, _ = bigexpression.NewParser(matches[1])
			}
		} else if reTestDivisibility.MatchString(line) {
			matches := reTestDivisibility.FindStringSubmatch(line)
			if len(matches) == 2 {
				curMonkey.testDivisibility, _ = strconv.ParseInt(matches[1], 10, 64)
			}
		} else if reIfTrue.MatchString(line) {
			matches := reIfTrue.FindStringSubmatch(line)
			if len(matches) == 2 {
				curMonkey.testDivisibilityTrueTarget, _ = strconv.Atoi(matches[1])
			}
		} else if reIfFalse.MatchString(line) {
			matches := reIfFalse.FindStringSubmatch(line)
			if len(matches) == 2 {
				curMonkey.testDivisibilityFalseTarget, _ = strconv.Atoi(matches[1])
			}
		}
	}

	monkeys = append(monkeys, curMonkey)

	limitModulo := big.NewInt(monkeys[0].testDivisibility)
	for _, m := range monkeys[1:] {
		limitModulo.Mul(limitModulo, big.NewInt(m.testDivisibility))
	}

	for i := range monkeys {
		monkeys[i].limitModulo = limitModulo
	}

	return monkeys
}
