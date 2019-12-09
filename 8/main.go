package main

import (
	"fmt"
	"io/ioutil"
)

const (
	Width       = 25
	Height      = 6
	Black       = byte('0')
	White       = byte('1')
	Transparent = byte('2')
)

func coord(x, y int) int {
	return y*Width + x%Width
}

func printImage(image []byte) {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			pixel := image[coord(x, y)]
			var char rune
			switch pixel {
			case Black:
				char = '.'
			case White:
				char = 'X'
			case Transparent:
				char = ' '
			default:
				panic("unknown pixel found")
			}

			fmt.Printf(" %c ", char)
		}
		fmt.Print("\n")
	}
}

func main() {
	raw, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}

	var layers [][]byte

	pos := 0
	for pos < len(raw)-1 {
		layer := make([]byte, Width*Height)
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				pixel := raw[pos]
				layer[coord(x, y)] = pixel
				pos++
			}
		}
		layers = append(layers, layer)
	}

	image := make([]byte, Width*Height)
	for _, layer := range layers {
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				pos := coord(x, y)
				pixel := layer[pos]
				curr := image[pos]

				if curr == White || curr == Black {
					continue
				}
				image[pos] = pixel
			}
		}
	}

	printImage(image)

}
