package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

const (
	MemSize = 64 * 1024

	ColorBlack = 0
	ColorWhite = 1

	TurnLeft  = 0
	TurnRight = 1

	HeadingNorth heading = 0
	HeadingEast  heading = 1
	HeadingSouth heading = 2
	HeadingWest  heading = 3
)

type heading byte

func (h heading) left() heading {
	rv := h - 1
	if rv > 3 {
		return 3
	}
	return rv
}

func (h heading) right() heading {
	return (h + 1) % 4
}

type coord struct {
	x int
	y int
}

func (c coord) step(h heading) coord {
	switch h {
	case HeadingNorth:
		return coord{x: c.x, y: c.y + 1}
	case HeadingEast:
		return coord{x: c.x + 1, y: c.y}
	case HeadingSouth:
		return coord{x: c.x, y: c.y - 1}
	case HeadingWest:
		return coord{x: c.x - 1, y: c.y}
	}

	panic("unknown heading")
}

type painter struct {
	c       *intcodeComputer
	pos     coord
	heading heading
	painted map[coord]bool
	in      <-chan int
	out     chan<- int
	max     coord
	min     coord
}

func mkPainter(mem []int) *painter {
	in := make(chan int, 1)
	out := make(chan int)

	c := &intcodeComputer{
		mem:    mem,
		input:  in,
		output: out,
	}

	return &painter{
		c:       c,
		pos:     coord{},
		heading: 0,
		painted: map[coord]bool{},
		in:      out,
		out:     in,
	}
}

func (p *painter) run() error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		p.out <- ColorWhite
		for {
			i, ok := <-p.in
			if !ok {
				wg.Done()
				return
			}
			p.paint(i)
			m := <-p.in
			p.move(m)
			color := ColorBlack
			white, ok := p.painted[p.pos]
			if ok && white {
				color = ColorWhite
			}
			p.out <- color
		}
	}()

	err := p.c.Run()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (p *painter) paint(color int) {
	white := false
	if color == ColorWhite {
		white = true
	}

	if p.pos.x > p.max.x {
		p.max.x = p.pos.x
	}
	if p.pos.y > p.max.y {
		p.max.y = p.pos.y
	}
	if p.pos.x < p.min.x {
		p.min.x = p.pos.x
	}
	if p.pos.y < p.min.y {
		p.min.y = p.pos.y
	}

	p.painted[p.pos] = white
}

func (p *painter) move(dir int) {
	switch dir {
	case TurnLeft:
		p.heading = p.heading.left()
	case TurnRight:
		p.heading = p.heading.right()
	}
	p.pos = p.pos.step(p.heading)
}

func (p *painter) print() {
	for y := p.max.y; y >= p.min.y; y-- {
		for x := p.min.x; x <= p.max.x; x++ {
			cur := coord{x: x, y: y}
			white, ok := p.painted[cur]
			if !ok {
				fmt.Print(" ")
				continue
			}

			if white {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
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

	painter := mkPainter(mem)
	if err := painter.run(); err != nil {
		panic(err)
	}

	painter.print()
}
