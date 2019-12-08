package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	numAmplifiers = 5
)

var inputs = [numAmplifiers]int{0, 1, 2, 3, 4}

type amplifier struct {
	phase    int
	computer *intcodeComputer
	in       chan<- int
	out      <-chan int
	recv     <-chan int
}

func (a *amplifier) run() {
	halt := make(chan struct{})
	go func() {
		a.in <- a.phase
		for {
			select {
			case r := <-a.recv:
				a.in <- r
			case <-halt:
				return
			}
		}
	}()

	err := a.computer.Run()
	halt <- struct{}{}
	if err != nil {
		panic(err)
	}
}

func mkAmplifier(name string, mem []int, phase int) *amplifier {
	input := make(chan int)
	output := make(chan int)

	c := &intcodeComputer{
		name:   name,
		mem:    mem,
		input:  input,
		output: output,
	}

	return &amplifier{
		phase:    phase,
		computer: c,
		in:       input,
		out:      output,
	}
}

func run(mem []int, phases []int) int {
	var rv int

	input := make(chan int)
	finalHalt := make(chan struct{})
	amps := [numAmplifiers]*amplifier{}

	for i := 0; i < numAmplifiers; i++ {
		cmem := make([]int, len(mem))
		copy(cmem, mem)
		name := fmt.Sprintf("amp:%d", i)
		amp := mkAmplifier(name, cmem, phases[i])
		amps[i] = amp
	}

	first := amps[0]
	first.recv = input
	for i := 1; i < numAmplifiers; i++ {
		amps[i].recv = amps[i-1].out
	}
	last := amps[len(amps)-1]
	output := last.out

	ampWait := sync.WaitGroup{}
	ampWait.Add(1)
	go func() {
		input <- 0
		ampWait.Done()
	}()

	finalWait := sync.WaitGroup{}
	finalWait.Add(1)

	go func() {
		for {
			select {
			case o := <-output:
				rv = o
			case <-finalHalt:
				finalWait.Done()
				return
			}
		}
	}()

	for _, a := range amps {
		ampWait.Add(1)
		go func(a1 *amplifier) {
			a1.run()
			ampWait.Done()
		}(a)
	}

	ampWait.Wait()
	finalHalt <- struct{}{}
	finalWait.Wait()

	return rv
}

func permutations(in []int) [][]int {
	var rv [][]int
	var f func([]int, int)

	f = func(in []int, n int) {
		if n == 1 {
			tmp := make([]int, len(in))
			copy(tmp, in)
			rv = append(rv, tmp)
		} else {
			for i := 0; i < n; i++ {
				f(in, n-1)
				if n%2 == 1 {
					tmp := in[i]
					in[i] = in[n-1]
					in[n-1] = tmp
				} else {
					tmp := in[0]
					in[0] = in[n-1]
					in[n-1] = tmp
				}
			}
		}
	}

	f(in, len(in))
	return rv
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

	tests := permutations(inputs[:])
	max := 0
	var maxPhases [numAmplifiers]int
	for _, test := range tests {
		cmem := make([]int, len(mem))
		copy(cmem, mem)
		val := run(cmem, test)
		if val > max {
			max = val
			copy(maxPhases[:], test[:])
		}
	}

	fmt.Printf("max signal: %d for phases %+v\n", max, maxPhases)
}
