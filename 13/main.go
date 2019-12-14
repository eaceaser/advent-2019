package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	MemSize = 1024 * 16

	TypeEmpty tileType = iota - 1
	TileWall
	TileBlock
	TilePaddle
	TileBall
)

type tileType int

func (t tileType) String() string {
	switch t {
	case TypeEmpty:
		return "."
	case TileWall:
		return "-"
	case TileBlock:
		return "█"
	case TilePaddle:
		return "="
	case TileBall:
		return "o"
	}
	panic(fmt.Sprintf("unknown tile type %d\n", t))
}

type tile struct {
	x int
	y int
	t tileType
}

type screen struct {
	c  *intcodeComputer
	s  []tile
	in <-chan int
}

func (s *screen) draw() int {
	rv := 0
	xcur := 0
	ycur := 0
	for _, t := range s.s {
		for y := ycur; y < t.y; y++ {
			fmt.Print("\n")
		}
		for x := xcur; x < t.x; x++ {
			fmt.Print(" ")
		}
		if t.t == TileBlock {
			rv++
		}
		fmt.Print(t.t.String())
		xcur = t.x
		ycur = t.y
	}
	fmt.Println()
	return rv
}

func (s *screen) run() {
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		for {
			x, ok := <-s.in
			if !ok {
				wait.Done()
				return
			}
			y, ok := <-s.in
			if !ok {
				wait.Done()
				return
			}
			t, ok := <-s.in
			if !ok {
				wait.Done()
				return
			}

			n := sort.Search(len(s.s), func(i int) bool {
				el := s.s[i]
				if y > el.y {
					return false
				} else if y == el.y {
					return x < el.x
				} else {
					return true
				}
			})

			tile := tile{
				x: x,
				y: y,
				t: tileType(t),
			}
			origS := s.s
			s.s = append(origS[0:n], tile)
			if n < len(origS) {
				s.s = append(s.s, origS[n:]...)
			}
		}
	}()

	if err := s.c.Run(); err != nil {
		panic(err)
	}

	wait.Wait()
}

func mkScreen(mem []int) *screen {
	comm := make(chan int)

	c := &intcodeComputer{
		name:   "screen",
		mem:    mem,
		output: comm,
	}

	return &screen{
		c:  c,
		in: comm,
	}
}

func main() {
	memS, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}
	memS1 := strings.Split(strings.TrimSpace(string(memS)), ",")
	mem := make([]int, MemSize)
	for i, s := range memS1 {
		if mem[i], err = strconv.Atoi(s); err != nil {
			panic(err)
		}
	}

	s := mkScreen(mem)
	s.run()
	ans := s.draw()
	fmt.Printf("drew %d blocks\n", ans)
}
