package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func fuel(mass int) int {
	return (mass / 3) - 2
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	sum := 0
	for s.Scan() {
		l := s.Text()
		mass, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}

		f := 0
		for mass > 0 {
			mass = fuel(mass)
			if mass > 0 {
				f += mass
			}
		}
		sum += f
	}

	fmt.Println(sum)
}
