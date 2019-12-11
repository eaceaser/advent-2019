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
		p.out <- ColorBlack
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

	fmt.Println(len(painter.painted))
}
