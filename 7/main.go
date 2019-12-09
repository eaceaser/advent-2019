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

var inputs = [numAmplifiers]int{5, 6, 7, 8, 9}

type relay struct {
	in   <-chan int
	out  chan<- int
	halt <-chan struct{}
	last int
}

func (r *relay) run() {
	for {
		select {
		case m := <-r.in:
			r.last = m
			r.out <- m
		case <-r.halt:
			return
		}
	}
}

type amplifier struct {
	computer *intcodeComputer
	in       chan<- int
	out      <-chan int
	recv     <-chan int
}

func (a *amplifier) run() error {
	err := a.computer.Run()
	return err
}

func mkAmplifier(name string, mem []int) *amplifier {
	input := make(chan int, 1)
	output := make(chan int, 1)

	c := &intcodeComputer{
		name:   name,
		mem:    mem,
		input:  input,
		output: output,
	}

	return &amplifier{
		computer: c,
		in:       input,
		out:      output,
	}
}

func tune(phases []int, amps []*amplifier) {
	for i, phase := range phases {
		amps[i].in <- phase
	}
}

func run(mem []int, phases []int) int {
	amps := make([]*amplifier, numAmplifiers)

	for i := 0; i < numAmplifiers; i++ {
		cmem := make([]int, len(mem))
		copy(cmem, mem)
		name := fmt.Sprintf("amp:%d", i)
		amp := mkAmplifier(name, cmem)
		amps[i] = amp
	}

	ampWait := sync.WaitGroup{}
	ampWait.Add(len(amps))
	for _, a := range amps {
		go func(a1 *amplifier) {
			if err := a1.run(); err != nil {
				panic(err)
			}
			ampWait.Done()
		}(a)
	}

	tune(phases, amps)
	amps[0].in <- 0

	var relayHalts []chan<- struct{}
	for i := 1; i < len(amps); i++ {
		halt := make(chan struct{})
		relay := relay{
			in:   amps[i-1].out,
			out:  amps[i].in,
			halt: halt,
		}

		go func() {
			relay.run()
		}()

		relayHalts = append(relayHalts, halt)
	}

	feedbackHalt := make(chan struct{})
	feedbackRelay := relay{
		in:   amps[len(amps)-1].out,
		out:  amps[0].in,
		halt: feedbackHalt,
	}
	go func() {
		feedbackRelay.run()
	}()
	relayHalts = append(relayHalts, feedbackHalt)

	ampWait.Wait()

	for _, halt := range relayHalts {
		halt <- struct{}{}
	}

	return feedbackRelay.last
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
