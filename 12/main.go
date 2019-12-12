package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const (
	NumSteps = 1000
)

var pattern = regexp.MustCompile("<x=(-?[0-9]+), y=(-?[0-9]+), z=(-?[0-9]+)>")

type coord struct {
	x int
	y int
	z int
}

func (c coord) add(o coord) coord {
	return coord{
		x: c.x + o.x,
		y: c.y + o.y,
		z: c.z + o.z,
	}
}

type obj struct {
	pos      coord
	velocity coord
}

func (o obj) potential() int {
	return abs(o.pos.x) + abs(o.pos.y) + abs(o.pos.z)
}

func (o obj) kinetic() int {
	return abs(o.velocity.x) + abs(o.velocity.y) + abs(o.velocity.z)
}

type model struct {
	t    int
	objs []obj
}

func dv(a int, b int) (int, int) {
	if a < b {
		return 1, -1
	} else if a > b {
		return -1, 1
	} else {
		return 0, 0
	}
}

func v(p1 coord, p2 coord) (coord, coord) {
	v1 := coord{}
	v2 := coord{}

	v1.x, v2.x = dv(p1.x, p2.x)
	v1.y, v2.y = dv(p1.y, p2.y)
	v1.z, v2.z = dv(p1.z, p2.z)

	return v1, v2
}

func (m *model) step() {
	for i := 0; i < len(m.objs); i++ {
		for j := i + 1; j < len(m.objs); j++ {
			o1 := m.objs[i]
			o2 := m.objs[j]
			v1, v2 := v(o1.pos, o2.pos)

			o1.velocity = o1.velocity.add(v1)
			o2.velocity = o2.velocity.add(v2)

			m.objs[i] = o1
			m.objs[j] = o2
		}
	}

	for i, o := range m.objs {
		o.pos = o.pos.add(o.velocity)
		m.objs[i] = o
	}

	m.t++
}

func (o obj) String() string {
	return fmt.Sprintf("pos=%+v vel=%+v pot=%d kin=%d total=%d", o.pos, o.velocity, o.potential(), o.kinetic(), o.potential()*o.kinetic())
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)
	var objs []obj
	for scan.Scan() {
		line := scan.Text()
		matches := pattern.FindAllStringSubmatch(line, -1)

		x, err := strconv.Atoi(matches[0][1])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(matches[0][2])
		if err != nil {
			panic(err)
		}
		z, err := strconv.Atoi(matches[0][3])
		if err != nil {
			panic(err)
		}

		pos := coord{x: x, y: y, z: z}
		objs = append(objs, obj{
			pos: pos,
		})
	}

	if len(objs) < 2 {
		panic("not enough objs")
	}

	if err := scan.Err(); err != nil {
		panic(err)
	}

	model := model{
		objs: objs,
	}

	for i := 0; i <= NumSteps; i++ {
		fmt.Printf("Step %d\n", i)
		sum := 0
		for j, o := range model.objs {
			fmt.Printf("%d: %s\n", j, o)
			sum += o.kinetic() * o.potential()
		}
		fmt.Printf("total energy: %d\n", sum)
		model.step()
	}
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}
