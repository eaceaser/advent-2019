package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	WorldSize = 50
	StatusRow = WorldSize + 1
	cmdNorth  = 1
	cmdSouth  = 2
	cmdWest   = 3
	cmdEast   = 4
	codeWall  = 0
	codeMove  = 1
	codeFound = 2

	charDroid  = '@'
	charWall   = '='
	charSpace  = '.'
	charTarget = '*'

	modeManual    = 0
	modeSearching = 1
	modePathfind  = 2
	modeOxygen    = 3
)

type droid struct {
	c      *intcodeComputer
	mode   int
	world  []int
	pos    int
	steps  int
	target int
}

func (d *droid) coord() (int, int) {
	return coord(d.pos)
}

func (d *droid) nextOrientation(orientation int, tried map[int]struct{}) int {
	x, y := d.coord()
	searchX, searchY := move(x, y, orientation)
	searchPos := pos(searchX, searchY)
	if d.world[searchPos] == 0 {
		return orientation
	} else if d.world[searchPos] == charWall {
		tried[orientation] = struct{}{}
		orientation = turnRight(orientation)
		return d.nextOrientation(orientation, tried)
	} else if d.world[searchPos] == charSpace || d.world[searchPos] == charTarget {
		if len(tried) >= 3 {
			return orientation
		}
		tried[orientation] = struct{}{}
		return d.nextOrientation(turnRight(orientation), tried)
	}

	return 0
}

func (d *droid) oxygen() int {
	d.world[d.pos] = charSpace
	total := 0
	for _, n := range d.world {
		if n == charSpace {
			total++
		}
	}
	steps := 0
	seen := make([]bool, len(d.world))
	filled := []int{d.target}
	seen[d.target] = true
	for len(filled) < total {
		for _, node := range filled {
			x, y := coord(node)
			n := pos(move(x, y, cmdNorth))
			s := pos(move(x, y, cmdSouth))
			e := pos(move(x, y, cmdEast))
			w := pos(move(x, y, cmdWest))
			for _, tgt := range []int{n, s, e, w} {
				if d.world[tgt] == charSpace {
					if seen[tgt] {
						continue
					}
					filled = append(filled, tgt)
					seen[tgt] = true
				}
			}
		}
		steps++
	}

	return steps
}

func (d *droid) pathfind() int {
	seen := make([]bool, len(d.world))
	type node struct {
		pos  int
		dist int
	}
	q := []node{{d.pos, 0}}
	for {
		nd := q[0]
		dist := nd.dist
		q = q[1:]
		x, y := coord(nd.pos)
		seen[nd.pos] = true

		n := pos(move(x, y, cmdNorth))
		s := pos(move(x, y, cmdSouth))
		e := pos(move(x, y, cmdEast))
		w := pos(move(x, y, cmdWest))
		for _, next := range []int{n, s, e, w} {
			if seen[next] {
				continue
			}
			if d.world[next] == charTarget {
				return dist + 1
			}
			if d.world[next] == charSpace {
				q = append(q, node{next, dist + 1})
			}
		}
	}
}

func (d *droid) run() error {
	orientation := cmdWest
	x, y := d.coord()
	write(x, y, string(charDroid))
	orig := d.pos
	for {
		x, y := d.coord()
		state, err := d.c.Run()
		if err != nil {
			return err
		}

		write(0, StatusRow, fmt.Sprintf("x=%d y=%d t=%d", x, y, d.steps))
		switch state {
		case stateInput:
			switch d.mode {
			case modeManual:
				d.c.in = acceptInput()
			case modeSearching:
				next := d.nextOrientation(orientation, map[int]struct{}{})
				if next == 0 || (d.pos == orig && d.steps > 0) {
					d.mode = modeOxygen
					d.c.in = cmdNorth
					continue
				}
				orientation = next
				d.c.in = orientation
			case modePathfind:
				dist := d.pathfind()
				write(0, StatusRow+1, fmt.Sprintf("PATH: %d\n", dist))
				return nil
			case modeOxygen:
				steps := d.oxygen()
				write(0, StatusRow+1, fmt.Sprintf("OXYGEN TIME: %d\n", steps))
				return nil
			}
		case stateOutput:
			switch d.c.out {
			case codeWall:
				x, y := d.coord()
				wallX, wallY := move(x, y, d.c.in)
				d.world[pos(wallX, wallY)] = charWall
				write(wallX, wallY, string(charWall))
			case codeMove:
				x, y := d.coord()
				newX, newY := move(x, y, d.c.in)
				newPos := pos(newX, newY)
				var char int
				if d.pos == d.target {
					char = charTarget
				} else {
					char = charSpace
				}
				d.world[d.pos] = char
				write(x, y, string(char))
				d.world[newPos] = charDroid
				write(newX, newY, string(charDroid))
				d.pos = newPos
				d.steps++
			case codeFound:
				x, y := d.coord()
				newX, newY := move(x, y, d.c.in)
				d.world[d.pos] = charSpace
				newPos := pos(newX, newY)
				write(x, y, string(charSpace))
				d.world[newPos] = charTarget
				d.target = newPos
				write(newX, newY, string(charTarget))
				d.pos = newPos
				d.steps++
			}
		}
	}
}

func acceptInput() int {
	for {
		var b [3]byte
		os.Stdin.Read(b[:])
		if b == [3]byte{104, 0, 0} {
			return cmdWest
		} else if b == [3]byte{106, 0, 0} {
			return cmdSouth
		} else if b == [3]byte{107, 0, 0} {
			return cmdNorth
		} else if b == [3]byte{108, 0, 0} {
			return cmdEast
		}
	}
}

func setupTTY() {
	fmt.Print("\033[H\033[2J")
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func mkDroid(mem []int) *droid {
	c := &intcodeComputer{
		name: "droid",
		mem:  mem,
	}

	return &droid{
		c:     c,
		world: make([]int, WorldSize*WorldSize),
		pos:   pos(WorldSize/2, WorldSize/2),
		mode:  modeSearching,
	}
}

func main() {
	setupTTY()

	inputS, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}
	cmds := strings.Split(strings.TrimSpace(string(inputS)), ",")
	mem := make([]int, len(cmds))
	for i, s := range cmds {
		mem[i], err = strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
	}

	droid := mkDroid(mem)

	if err := droid.run(); err != nil {
		panic(err)
	}
}

func moveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y+1, x+1)
}

func move(x int, y int, dir int) (int, int) {
	switch dir {
	case cmdNorth:
		return x, y - 1
	case cmdSouth:
		return x, y + 1
	case cmdWest:
		return x - 1, y
	case cmdEast:
		return x + 1, y
	}
	panic("unknown direction")
}

func turnRight(dir int) int {
	switch dir {
	case cmdNorth:
		return cmdEast
	case cmdEast:
		return cmdSouth
	case cmdSouth:
		return cmdWest
	case cmdWest:
		return cmdNorth
	}
	panic("unknown dir")
}

func write(x int, y int, s string) {
	moveCursor(x, y)
	fmt.Print(s)
}

func coord(pos int) (int, int) {
	return pos % WorldSize, pos / WorldSize
}

func pos(x int, y int) int {
	return x + y*WorldSize
}
