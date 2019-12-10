package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const (
	Empty    = '.'
	Asteroid = '#'
)

type coord struct {
	x int
	y int
}

type angles map[float64][]coord

func angle(src coord, dest coord) float64 {
	dy := float64(dest.y - src.y)
	dx := float64(dest.x - src.x)
	return math.Atan2(dy, dx)
}

func dist(src coord, dest coord) float64 {
	return math.Sqrt(math.Pow(float64(dest.x-src.x), 2) + math.Pow(float64(dest.y-src.y), 2))
}

func visible(grid []coord, obj coord) (int, angles) {
	angles := map[float64][]coord{}
	rv := 0
	for _, o2 := range grid {
		if obj == o2 {
			continue
		}

		s := angle(obj, o2)

		if arr, ok := angles[s]; ok {
			angles[s] = append(arr, o2)
			continue
		} else {
			angles[s] = []coord{o2}
		}

		rv++
	}

	for _, coords := range angles {
		sort.Slice(coords, func(i, j int) bool {
			c1 := coords[i]
			c2 := coords[j]

			d1 := dist(c1, obj)
			d2 := dist(c2, obj)
			return d1 < d2
		})
	}

	return rv, angles
}

func obliterate(angles angles) []coord {
	var rv []coord
	var angleArr []float64
	count := 0
	for k, coords := range angles {
		angleArr = append(angleArr, k)
		count += len(coords)
	}

	sort.Slice(angleArr, func(i int, j int) bool {
		f1 := angleArr[i]
		f2 := angleArr[j]

		if f1 < 0 {
			f1 += 2 * math.Pi
		}
		if f2 < 0 {
			f2 += 2 * math.Pi
		}

		f1 += math.Pi / 2
		f2 += math.Pi / 2

		for f1 >= 2*math.Pi {
			f1 -= 2 * math.Pi
		}
		for f2 >= 2*math.Pi {
			f2 -= 2 * math.Pi
		}

		return f1 < f2
	})

	for count > 0 {
		for _, angle := range angleArr {
			coords := angles[angle]
			if len(coords) > 0 {
				rv = append(rv, coords[0])
				angles[angle] = coords[1:]
				count--
			}
		}
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
	var maxAngles angles
	for _, obj := range objs {
		v, angles := visible(objs, obj)
		if v > maxVisible {
			maxVisible = v
			maxCoord = obj
			maxAngles = angles
		}
	}

	fmt.Printf("found %+v with visible %d\n", maxCoord, maxVisible)

	destroyed := obliterate(maxAngles)
	for i, c := range destroyed {
		ans := 100*c.x + c.y
		fmt.Printf("%d: %+v => %d\n", i+1, c, ans)
	}

}
