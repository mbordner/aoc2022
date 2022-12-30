package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"strconv"
)

type Node struct {
	v int
	n *Node
	p *Node
}

type Numbers struct {
	initial []int
	nodes   []*Node
	valMap  map[int]*Node
}

func NewNumbers(vals []int) *Numbers {
	n := new(Numbers)
	n.initial = vals
	n.nodes = make([]*Node, 0, len(vals))
	n.valMap = make(map[int]*Node)

	var first, last *Node

	first = &Node{v: vals[0]}
	n.valMap[vals[0]] = first
	n.nodes = append(n.nodes, first)

	prev := first
	for _, v := range vals[1:] {
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
	vals := make([]int, len(nums.initial), len(nums.initial))

	n := nums.nodes[0]
	for i := 0; i < len(vals); i++ {
		vals[i] = n.v
		n = n.n
	}

	fmt.Println(vals)
}

func (from *Node) Next(count int) *Node {
	n := from
	for i := 0; i < count; i++ {
		n = n.n
	}
	return n
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (from *Node) Prev(count int) *Node {
	p := from
	for i := 0; i < count; i++ {
		p = p.p
	}
	return p
}

// Node returns node that represents last value in the list
func (nums *Numbers) Node(v int) *Node {
	if n, e := nums.valMap[v]; e {
		return n
	}
	return nil
}

func (nums *Numbers) DoMix() {
	for i := 0; i < len(nums.initial); i++ {
		v := nums.initial[i]
		n := nums.nodes[i]
		if v > 0 {
			start := n.n
			amt := v - 1
			nums.Detach(n)
			target := start.Next(amt)
			nums.InsertAfter(n, target)
		} else if v < 0 {
			start := n.p
			amt := v + 1
			nums.Detach(n)
			target := start.Prev(abs(amt))
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
		nums.DoMix()
		n := nums.Node(0)
		n1000 := n.Next(1000)
		n2000 := n1000.Next(1000)
		n3000 := n2000.Next(1000)
		fmt.Println(n1000.v + n2000.v + n3000.v)
	}
}

func getNumbers(path string) *Numbers {
	lines, _ := file.GetLines(path)
	nums := make([]int, len(lines), len(lines))
	for i, n := range lines {
		nums[i], _ = strconv.Atoi(n)
	}
	return NewNumbers(nums)
}
