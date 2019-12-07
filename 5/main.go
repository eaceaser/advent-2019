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
	OpcodeInput    = 3
	OpcodeOutput   = 4
	OpcodeDie      = 99

	ModePosition  int = 0
	ModeImmediate int = 1
)

type IntcodeComputer struct {
	mem    []int
	ip     int
	input  <-chan int
	output chan<- int
	halt chan<- struct{}
}

func (c *IntcodeComputer) inputImpl() {
	target := c.read()
	in := <-c.input
	c.mem[target] = in
}

func (c *IntcodeComputer) outputImpl(modes []int) error {
	params, err := c.modalParams(pad(modes, 1)...)
	if err != nil {
		return err
	}
	c.output <- params[0]
	return nil
}

func (c *IntcodeComputer) modalParams(mode ...int) ([]int, error) {
	rv := make([]int, len(mode))
	for i, t := range mode {
		p := c.read()
		switch t {
		case ModePosition:
			rv[i] = c.mem[p]
		case ModeImmediate:
			rv[i] = p
		default:
			return nil, fmt.Errorf("unknown mode %d", t)

		}
	}
	return rv, nil
}

func (c *IntcodeComputer) arithmeticImpl(parsedModes []int, f func(a, b int) int) error {
	modes := pad(parsedModes, 2)
	params, err := c.modalParams(modes...)
	if err != nil {
		return err
	}
	dest := c.read()
	c.mem[dest] = f(params[0], params[1])
	return nil
}

func (c *IntcodeComputer) Run() error {
	for {
		cmdDesc := c.read()
		opcode, parsedModes := parseOpcode(cmdDesc)

		switch opcode {
		case OpcodeAdd:
			if err := c.arithmeticImpl(parsedModes, addOp); err != nil {
				return err
			}
		case OpcodeMultiply:
			if err := c.arithmeticImpl(parsedModes, multOp); err != nil {
				return err
			}
		case OpcodeInput:
			c.inputImpl()
		case OpcodeOutput:
			if err := c.outputImpl(parsedModes); err != nil {
				return err
			}
		case OpcodeDie:
			close(c.halt)
			return nil
		default:
			return fmt.Errorf("illegal opcode %d", opcode)
		}
	}
}

func (c *IntcodeComputer) read() int {
	rv := c.mem[c.ip]
	c.ip++
	return rv
}

func addOp(a, b int) int  { return a + b }
func multOp(a, b int) int { return a * b }

func parseOpcode(code int) (opcode int, modes []int) {
	opcode = code % 100
	rest := code / 100

	modes = make([]int, 0)
	for rest > 0 {
		mode := rest % 10
		modes = append(modes, mode)
		rest /= 10
	}

	return opcode, modes
}

func pad(m []int, sz int) []int {
	return append(m, make([]int, sz-len(m))...)
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

	input := make(chan int)
	output := make(chan int)
	halt := make(chan struct{})

	go func() {
		input <- 1
	}()

	go func() {
		for {
			select {
			case o := <-output:
				fmt.Printf("[OUT] %d\n", o)
			case <-halt:
				fmt.Println("[DONE]")
				return
			}

		}
	}()

	c := IntcodeComputer{
		mem:    mem,
		ip:     0,
		input:  input,
		output: output,
		halt: halt,
	}

	err = c.Run()
	if err != nil {
		panic(err)
	}
}
