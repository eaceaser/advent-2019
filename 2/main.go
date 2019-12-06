package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	OpcodeAdd      = 1
	OpcodeMultiply = 2
	OpcodeDie      = 99
	Target         = 19690720
)

type IntcodeComputer struct {
	mem []int
	ip  int
}

func (c *IntcodeComputer) add() {
	a := c.read()
	b := c.read()
	x := c.read()

	c.mem[x] = c.mem[a] + c.mem[b]
}

func (c *IntcodeComputer) multiply() {
	a := c.read()
	b := c.read()
	x := c.read()

	c.mem[x] = c.mem[a] * c.mem[b]
}

func (c *IntcodeComputer) Run() (int, error) {
	for {
		opcode := c.read()
		switch opcode {
		case OpcodeAdd:
			c.add()
		case OpcodeMultiply:
			c.multiply()
		case OpcodeDie:
			return c.mem[0], nil
		default:
			return 0, fmt.Errorf("illegal opcode %d", opcode)
		}
	}
}

func (c *IntcodeComputer) read() int {
	rv := c.mem[c.ip]
	c.ip++
	return rv
}

func prep(mem []int, a, b int) {
	mem[1] = a
	mem[2] = b
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	memS, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	memS1 := strings.Split(string(memS), ",")
	mem := make([]int, len(memS1))
	for i, m := range memS1 {
		mi, err := strconv.Atoi(strings.TrimSpace(m))
		if err != nil {
			panic(err)
		}
		mem[i] = mi
	}

	prep(mem, 12, 2)

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			mem2 := make([]int, len(mem))
			copy(mem2, mem)
			prep(mem2, i, j)
			c := IntcodeComputer{
				mem: mem2,
			}
			res, err := c.Run()
			if err != nil {
				panic(err)
			}

			if res == Target {
				fmt.Printf("noun=%d verb=%d answer=%d\n", i, j, 100*i+j)
				return
			}
		}
	}
}
