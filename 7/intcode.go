package main

import "fmt"

const (
	opcodeAdd          = 1
	opcodeMultiply     = 2
	opcodeInput        = 3
	opcodeOutput       = 4
	opcodeJmpIfT       = 5
	opcodeJmpIfF       = 6
	opcodeLT           = 7
	opcodeEql          = 8
	opcodeDie          = 99
	modePosition   int = 0
	modeImmediate  int = 1
)

type intcodeComputer struct {
	name   string
	mem    []int
	ip     int
	input  <-chan int
	output chan<- int
}

func (c *intcodeComputer) jumpImpl(modes []int, cmp func(p int) bool) error {
	params, err := c.modalParams(pad(modes, 2)...)
	if err != nil {
		return err
	}

	if cmp(params[0]) {
		c.ip = params[1]
	}

	return nil
}

func (c *intcodeComputer) inputImpl() {
	target := c.read()
	in := <-c.input
	c.mem[target] = in
}

func (c *intcodeComputer) outputImpl(modes []int) error {
	params, err := c.modalParams(pad(modes, 1)...)
	if err != nil {
		return err
	}
	c.output <- params[0]
	return nil
}

func (c *intcodeComputer) modalParams(mode ...int) ([]int, error) {
	rv := make([]int, len(mode))
	for i, t := range mode {
		p := c.read()
		switch t {
		case modePosition:
			rv[i] = c.mem[p]
		case modeImmediate:
			rv[i] = p
		default:
			return nil, fmt.Errorf("unknown mode %d", t)
		}
	}
	return rv, nil
}

func (c *intcodeComputer) arithmeticImpl(parsedModes []int, f func(a, b int) int) error {
	modes := pad(parsedModes, 2)
	params, err := c.modalParams(modes...)
	if err != nil {
		return err
	}
	dest := c.read()
	c.mem[dest] = f(params[0], params[1])
	return nil
}

func (c *intcodeComputer) cmpImpl(parsedModes []int, f func(a, b int) bool) error {
	modes := pad(parsedModes, 2)
	params, err := c.modalParams(modes...)
	if err != nil {
		return err
	}
	dest := c.read()
	if f(params[0], params[1]) {
		c.mem[dest] = 1
	} else {
		c.mem[dest] = 0
	}
	return nil
}

func (c *intcodeComputer) Run() error {
	for {
		cmdDesc := c.read()
		opcode, parsedModes := parseOpcode(cmdDesc)

		switch opcode {
		case opcodeAdd:
			if err := c.arithmeticImpl(parsedModes, addOp); err != nil {
				return err
			}
		case opcodeMultiply:
			if err := c.arithmeticImpl(parsedModes, multOp); err != nil {
				return err
			}
		case opcodeInput:
			c.inputImpl()
		case opcodeOutput:
			if err := c.outputImpl(parsedModes); err != nil {
				return err
			}
		case opcodeJmpIfT:
			if err := c.jumpImpl(parsedModes, trueCmp); err != nil {
				return err
			}
		case opcodeJmpIfF:
			if err := c.jumpImpl(parsedModes, falseCmp); err != nil {
				return err
			}
		case opcodeLT:
			if err := c.cmpImpl(parsedModes, ltCmp); err != nil {
				return err
			}
		case opcodeEql:
			if err := c.cmpImpl(parsedModes, eqCmp); err != nil {
				return err
			}
		case opcodeDie:
			return nil
		default:
			return fmt.Errorf("illegal opcode %d", opcode)
		}
	}
}

func (c *intcodeComputer) read() int {
	rv := c.mem[c.ip]
	c.ip++
	return rv
}

func addOp(a, b int) int      { return a + b }
func multOp(a, b int) int     { return a * b }
func trueCmp(a int) bool      { return a != 0 }
func falseCmp(a int) bool     { return a == 0 }
func ltCmp(a int, b int) bool { return a < b }
func eqCmp(a int, b int) bool { return a == b }

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
