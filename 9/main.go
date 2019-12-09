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
)

func main() {
	memS, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}

	memS1 := strings.Split(strings.TrimSpace(string(memS)), ",")
	mem := make([]int, MemSize)
	for i, s := range memS1 {
		mem[i], err = strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
	}

	input := make(chan int, 1)
	input <- 2
	output := make(chan int)
	c := intcodeComputer{
		mem:    mem,
		input:  input,
		output: output,
	}

	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		for {
			select {
			case o, ok := <-output:
				if !ok {
					wait.Done()
					return
				}
				fmt.Printf("[OUT] %d\n", o)
			}
		}
	}()

	if err := c.Run(); err != nil {
		panic(err)
	}
	close(output)
	wait.Wait()
}
