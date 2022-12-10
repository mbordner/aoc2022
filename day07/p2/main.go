package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/file"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	rootDir = NewDir("/", nil)
	dirs    = make(map[string]*fsNode)
	reEntry = regexp.MustCompile(`(dir|\d+) (.*)`)
)

type fsNode struct {
	name    string
	entries map[string]*fsNode
	size    int
	parent  *fsNode
}

func (fsn *fsNode) isDir() bool {
	return fsn.entries != nil
}

func NewDir(name string, parent *fsNode) *fsNode {
	d := &fsNode{}
	d.entries = make(map[string]*fsNode)
	d.parent = parent
	d.name = name
	return d
}

func (fsn *fsNode) addEntry(e *fsNode) {
	if fsn.entries != nil {
		fsn.entries[e.name] = e
	}
}
func (fsn *fsNode) getSize() int {
	if !fsn.isDir() || fsn.size > 0 {
		return fsn.size
	}
	size := 0
	for _, e := range fsn.entries {
		size += e.getSize()
	}
	fsn.size = size
	return size
}
func (fsn *fsNode) getAbsPath() string {
	if fsn.name == "/" {
		return fsn.name
	}
	return filepath.Join(fsn.parent.getAbsPath(), fsn.name)
}

func NewFile(name string, size int, parent *fsNode) *fsNode {
	f := &fsNode{}
	f.entries = nil
	f.parent = parent
	f.name = name
	f.size = size
	return f
}

type fsNodeList []*fsNode

func (fl fsNodeList) Len() int           { return len(fl) }
func (fl fsNodeList) Less(i, j int) bool { return fl[i].getSize() < fl[j].getSize() }
func (fl fsNodeList) Swap(i, j int)      { fl[i], fl[j] = fl[j], fl[i] }

func main() {

	cwd := rootDir
	dirs[rootDir.getAbsPath()] = rootDir

	lines, _ := file.GetLines("../data.txt")
	for _, line := range lines {
		if len(line) > 0 {
			if strings.HasPrefix(line, "$ ") {
				cmd := line[2:]
				if strings.HasPrefix(cmd, "cd ") {
					targetPath := cmd[3:]
					if targetPath == "/" {
						cwd = rootDir
					} else {
						if d, exists := cwd.entries[targetPath]; exists {
							cwd = d
						} else if targetPath == ".." {
							if cwd.parent != nil {
								cwd = cwd.parent
							} else {
								panic("root has no parent")
							}
						} else {
							panic("it should exist")
						}
					}
				} // ignore others
			} else {
				// dir or file entry
				if !reEntry.MatchString(line) {
					panic("what is this")
				} else {
					matches := reEntry.FindStringSubmatch(line)
					if len(matches) == 3 {
						if matches[1] == "dir" {
							d := NewDir(matches[2], cwd)
							dirs[d.getAbsPath()] = d
							cwd.addEntry(d)
						} else {
							s, _ := strconv.Atoi(matches[1])
							f := NewFile(matches[2], s, cwd)
							cwd.addEntry(f)
						}
					}
				}
			}
		}
	}

	sum := 0

	fsnl := make(fsNodeList, 0, len(dirs))

	for _, d := range dirs {
		fsnl = append(fsnl, d)
		dsz := d.getSize()
		if dsz <= 100000 {
			sum += dsz
		}
	}

	curUnused := 70000000 - dirs["/"].getSize()
	needToFree := 30000000 - curUnused
	sort.Sort(fsnl)

	for i := range fsnl {
		fmt.Println(fsnl[i].getAbsPath(), fsnl[i].getSize())
	}

	for i := range fsnl {
		if fsnl[i].getSize() >= needToFree {
			fmt.Println("we need to free: ", needToFree)
			fmt.Println("smallest we can delete: ", fsnl[i].getAbsPath(), fsnl[i].getSize())
			break
		}
	}

}
