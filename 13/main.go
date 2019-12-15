package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	MemSize   = 1024 * 16
	PauseTime = 250 * time.Millisecond

	TileEmpty  tileType = 0
	TileWall   tileType = 1
	TileBlock  tileType = 2
	TilePaddle tileType = 3
	TileBall   tileType = 4
)

type tileType int

func (t tileType) String() string {
	switch t {
	case TileEmpty:
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

type game struct {
	c     *intcodeComputer
	t     int
	state int
}

var saves []intcodeComputer

func (g *game) awaitOutput() (int, error) {
	s, err := g.c.Run()
	if err != nil {
		return 0, err
	}
	if s != stateOutput {
		return 0, fmt.Errorf("awaiting output state, got state %d", s)
	}
	return g.c.out, nil
}

func (g *game) run() error {
	ymax := 0

	for {
		cs, err := g.c.Run()
		if err != nil {
			return err
		}

		switch cs {
		case stateOutput:
			x := g.c.out
			y, err := g.awaitOutput()
			if err != nil {
				return err
			}
			t, err := g.awaitOutput()
			if err != nil {
				return err
			}

			tt := tileType(t)
			if y > ymax {
				ymax = y
			}
			write(x, y, tt.String())

		case stateInput:
			goto Running
		case stateHalted:
			return fmt.Errorf("intcode halted during setup")
		}
	}

Running:
	g.c.in = -1
	steps := 0
	for {
		write(0, ymax+1, fmt.Sprintf("t=%d\n", steps))
		cs, err := g.c.Run()
		if err != nil {
			return err
		}

		switch cs {
		case stateOutput:
			x := g.c.out
			y, err := g.awaitOutput()
			if err != nil {
				return err
			}
			t, err := g.awaitOutput()
			if err != nil {
				return err
			}

			if x == -1 && y == 0 && t != 0 {
				write(0, ymax+2, fmt.Sprintf("s=%d\n", t))
				continue
			}

			tt := tileType(t)
			write(x, y, tt.String())
		case stateInput:
			var b [3]byte
			for {
				os.Stdin.Read(b[:])
				//if true {
				if b == [3]byte{27, 91, 67} {
					g.c.in = 1
					break
				} else if b == [3]byte{27, 91, 68} {
					g.c.in = -1
					break
				} else if b == [3]byte{32, 0, 0} {
					g.c.in = 0
					break
				} else if b == [3]byte{98, 0, 0} {
					if len(saves) > 0 {
						save := saves[0]
						g.c = &save
						steps--
						break
					}
				} else {
					write(0, ymax+3, fmt.Sprint("I got the byte", b, "("+string(b[:])+")"))
				}
			}
			save := g.c.copy()
			saves = append([]intcodeComputer{*save}, saves...)
			steps++
		case stateHalted:
			goto Done
		}
	}
Done:
	save := saves[len(saves)-10]
	saves = saves[10:]
	g.c = &save
	steps -= 10
	goto Running
	//return nil
}

func mkGame(mem []int) *game {
	c := &intcodeComputer{
		name: "game",
		mem:  mem,
	}

	return &game{
		c: c,
	}
}

func moveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y+1, x+1)
}

func main() {
	fmt.Print("\033[H\033[2J")
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

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

	mem[0] = 2
	s := mkGame(mem)
	if err := s.run(); err != nil {
		panic(err)
	}
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func write(x int, y int, s string) {
	moveCursor(x, y)
	fmt.Print(s)
}
