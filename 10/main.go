package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const (
	Empty    = '.'
	Asteroid = '#'
)

type coord struct {
	x int
	y int
}

func angle(src coord, dest coord) float64 {
	dy := float64(dest.y - src.y)
	dx := float64(dest.x - src.x)
	return math.Atan2(dy, dx)
}

func visible(grid []coord, obj coord) int {
	slopes := map[float64]struct{}{}
	rv := 0
	for _, o2 := range grid {
		if obj == o2 {
			continue
		}

		s := angle(obj, o2)

		if _, ok := slopes[s]; ok {
			continue
		}
		slopes[s] = struct{}{}

		rv++
	}

	return rv
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	var objs []coord
	y := 0
	for scan.Scan() {
		line := scan.Text()
		for x, c := range line {
			if c != Empty {
				coord := coord{x, y}
				objs = append(objs, coord)
			}
		}
		y++
	}

	if err := scan.Err(); err != nil {
		panic(err)
	}

	maxVisible := 0
	maxCoord := coord{}
	for _, obj := range objs {
		v := visible(objs, obj)
		//fmt.Printf("%+v: %d\n", obj, v)
		if v > maxVisible {
			maxVisible = v
			maxCoord = obj
		}
	}

	fmt.Printf("found %+v with visible %d\n", maxCoord, maxVisible)
}
