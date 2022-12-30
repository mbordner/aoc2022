package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"strconv"
)

type Node struct {
	v int64
	n *Node
	p *Node
}

type Numbers struct {
	initial []int64
	nodes   []*Node
	valMap  map[int64]*Node
	min     int64
	max     int64
}

func NewNumbers(vals []int64) *Numbers {
	n := new(Numbers)
	n.initial = vals
	n.nodes = make([]*Node, 0, len(vals))
	n.valMap = make(map[int64]*Node)

	var first, last *Node

	first = &Node{v: vals[0]}
	n.valMap[vals[0]] = first
	n.nodes = append(n.nodes, first)

	n.min = vals[0]
	n.max = vals[0]

	prev := first
	for _, v := range vals[1:] {
		if v < n.min {
			n.min = v
		}
		if v > n.max {
			n.max = v
		}
		tmp := &Node{v: v}
		tmp.p = prev
		prev.n = tmp
		prev = tmp
		n.valMap[v] = tmp // there will be dupes, but we'll use this only to find 0 which shouldn't have a dupe
		n.nodes = append(n.nodes, tmp)
	}

	last = prev

	first.p = last
	last.n = first

	return n
}

func (nums *Numbers) Print() {
	vals := make([]int64, len(nums.initial), len(nums.initial))

	n := nums.nodes[0]
	for i := 0; i < len(vals); i++ {
		vals[i] = n.v
		n = n.n
	}

	fmt.Println(vals)
}

func (from *Node) Next(count int64) *Node {
	n := from
	for i := int64(0); i < count; i++ {
		n = n.n
	}
	return n
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func (from *Node) Prev(count int64) *Node {
	p := from
	for i := int64(0); i < count; i++ {
		p = p.p
	}
	return p
}

// Node returns node that represents last value in the list
func (nums *Numbers) Node(v int64) *Node {
	if n, e := nums.valMap[v]; e {
		return n
	}
	return nil
}

func (nums *Numbers) getAmt(v int64) int64 {
	sz := int64(len(nums.initial))
	wraps := v / sz
	amt := (v % sz) + (wraps - 1)
	if amt < sz {
		return amt
	}
	return amt % (sz - 1)
}

func (nums *Numbers) DoMix() {
	for i := 0; i < len(nums.initial); i++ {
		v := nums.initial[i]
		n := nums.nodes[i]
		amt := nums.getAmt(abs(v))
		if v > 0 {
			start := n.n
			nums.Detach(n)
			target := start.Next(amt)
			nums.InsertAfter(n, target)
		} else if v < 0 {
			start := n.p
			nums.Detach(n)
			target := start.Prev(amt)
			nums.InsertBefore(n, target)
		}
	}
}

func (nums *Numbers) Detach(n *Node) {
	prev := n.p
	next := n.n
	n.p, n.n = nil, nil
	prev.n = next
	next.p = prev
}

func (nums *Numbers) InsertAfter(n *Node, at *Node) {
	next := at.n
	at.n = n
	n.p = at
	n.n = next
	next.p = n
}

func (num *Numbers) InsertBefore(n *Node, at *Node) {
	prev := at.p
	at.p = n
	n.n = at
	n.p = prev
	prev.n = n
}

// too low 9459
func main() {
	nums := getNumbers("../data.txt")

	if nums != nil {
		for i := 0; i < 10; i++ {
			nums.DoMix()
			fmt.Println("mixed ", i, " times")
		}
		n := nums.Node(0)
		n1000 := n.Next(1000)
		n2000 := n1000.Next(1000)
		n3000 := n2000.Next(1000)
		fmt.Println(n1000.v + n2000.v + n3000.v)
	}
}

func getNumbers(path string) *Numbers {
	lines, _ := file.GetLines(path)
	nums := make([]int64, len(lines), len(lines))
	for i, n := range lines {
		nums[i], _ = strconv.ParseInt(n, 10, 64)
		nums[i] *= 811589153
	}
	return NewNumbers(nums)
}
