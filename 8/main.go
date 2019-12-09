package main

import (
	"fmt"
	"io/ioutil"
	"math"
)

const (
	Width  = 25
	Height = 6
)

func main() {
	raw, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}

	var layers []map[int]int

	pos := 0
	for pos < len(raw)-1 {
		layer := map[int]int{}
		for y := 0; y < Height; y++ {
			for j := 0; j < Width; j++ {
				digitS := raw[pos]
				digit := int(digitS - '0')
				layer[digit] += 1
				pos++
			}
		}
		layers = append(layers, layer)
	}

	min0 := math.MaxInt64
	var minL int
	var ans int
	for i, layer := range layers {
		if layer[0] < min0 {
			min0 = layer[0]
			minL = i
			ans = layer[1] * layer[2]
		}
	}

	fmt.Printf("layer %d has minimal zeros (%d). ans=%d\n", minL, min0, ans)
}
