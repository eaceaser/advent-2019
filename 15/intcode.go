package main

import (
	"fmt"
)

const (
	opcodeAdd      = 1
	opcodeMultiply = 2
	opcodeInput    = 3
	opcodeOutput   = 4
	opcodeJmpIfT   = 5
	opcodeJmpIfF   = 6
	opcodeLT       = 7
	opcodeEql      = 8
	opcodeRelAdj   = 9
	opcodeDie      = 99
	modePosition   = 0
	modeImmediate  = 1
	modeRelative   = 2

	stateRunning computerState = 0
	stateHalted  computerState = 1
	stateInput   computerState = 2
	stateOutput  computerState = 3
)

type computerState int

type intcodeComputer struct {
	name  string
	mem   []int
	tp    *int
	ip    int
	rel   int
	in    int
	out   int
	state computerState
}

func (c *intcodeComputer) jumpImpl(modes []int, cmp func(p int) bool) error {
	params, err := c.modalParams(pad(modes, 2)...)
	if err != nil {
		return err
	}

	if cmp(*params[0]) {
		c.ip = *params[1]
	}

	return nil
}

func (c *intcodeComputer) inputImpl(modes []int) error {
	modes = pad(modes, 1)
	dest, err := c.modalParams(modes[0])
	if err != nil {
		return err
	}
	c.tp = dest[0]
	return nil
}

func (c *intcodeComputer) outputImpl(modes []int) error {
	params, err := c.modalParams(pad(modes, 1)...)
	if err != nil {
		return err
	}
	c.out = *params[0]
	return nil
}

func (c *intcodeComputer) modalParams(mode ...int) ([]*int, error) {
	rv := make([]*int, len(mode))
	for i, t := range mode {
		p := c.read()
		switch t {
		case modePosition:
			rv[i] = &c.mem[p]
		case modeImmediate:
			rv[i] = &p
		case modeRelative:
			rv[i] = &c.mem[c.rel+p]
		default:
			return nil, fmt.Errorf("unknown mode %d", t)
		}
	}
	return rv, nil
}

func (c *intcodeComputer) arithmeticImpl(parsedModes []int, f func(a, b int) int) error {
	modes := pad(parsedModes, 3)
	params, err := c.modalParams(modes...)
	if err != nil {
		return err
	}
	*params[2] = f(*params[0], *params[1])
	return nil
}

func (c *intcodeComputer) cmpImpl(parsedModes []int, f func(a, b int) bool) error {
	modes := pad(parsedModes, 3)
	params, err := c.modalParams(modes...)
	if err != nil {
		return err
	}
	dest := params[2]
	if f(*params[0], *params[1]) {
		*dest = 1
	} else {
		*dest = 0
	}
	return nil
}

func (c *intcodeComputer) relImpl(parsedModes []int) error {
	modes := pad(parsedModes, 1)
	params, err := c.modalParams(modes...)
	if err != nil {
		return err
	}
	c.rel += *params[0]
	return nil
}

func (c *intcodeComputer) runLoop() (computerState, error) {
	c.state = stateRunning
	for {
		cmdDesc := c.read()
		opcode, parsedModes := parseOpcode(cmdDesc)

		switch opcode {
		case opcodeAdd:
			if err := c.arithmeticImpl(parsedModes, addOp); err != nil {
				return 0, err
			}
		case opcodeMultiply:
			if err := c.arithmeticImpl(parsedModes, multOp); err != nil {
				return 0, err
			}
		case opcodeInput:
			if err := c.inputImpl(parsedModes); err != nil {
				return 0, err
			}
			c.state = stateInput
			return stateInput, nil
		case opcodeOutput:
			if err := c.outputImpl(parsedModes); err != nil {
				return 0, err
			}
			c.state = stateOutput
			return stateOutput, nil
		case opcodeJmpIfT:
			if err := c.jumpImpl(parsedModes, trueCmp); err != nil {
				return 0, err
			}
		case opcodeJmpIfF:
			if err := c.jumpImpl(parsedModes, falseCmp); err != nil {
				return 0, err
			}
		case opcodeLT:
			if err := c.cmpImpl(parsedModes, ltCmp); err != nil {
				return 0, err
			}
		case opcodeEql:
			if err := c.cmpImpl(parsedModes, eqCmp); err != nil {
				return 0, err
			}
		case opcodeRelAdj:
			if err := c.relImpl(parsedModes); err != nil {
				return 0, err
			}
		case opcodeDie:
			c.state = stateHalted
			return stateHalted, nil
		default:
			return 0, fmt.Errorf("illegal opcode %d", opcode)
		}
	}
}

func (c *intcodeComputer) Run() (computerState, error) {
	switch c.state {
	case stateRunning:
		return c.runLoop()
	case stateInput:
		*c.tp = c.in
		return c.runLoop()
	case stateOutput:
		return c.runLoop()
	default:
		return c.state, fmt.Errorf("invalid state")
	}
}

func (c *intcodeComputer) read() int {
	rv := c.mem[c.ip]
	c.ip++
	return rv
}

func (c *intcodeComputer) copy() *intcodeComputer {
	newMem := make([]int, len(c.mem))
	copy(newMem, c.mem)
	rv := &intcodeComputer{
		name:  c.name,
		mem:   newMem,
		tp:    c.tp,
		ip:    c.ip,
		rel:   c.rel,
		in:    c.in,
		out:   c.out,
		state: c.state,
	}
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
