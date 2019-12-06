package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type intersection struct {
	point
	dist int
}

func (p point) Dist() int {
	return abs(p.x) + abs(p.y)
}

type seg struct {
	a   point
	b   point
	len int
}

func (s seg) Slope() int {
	if s.a.x != s.b.x {
		return 0
	} else {
		return 1
	}
}

func (s seg) Normalized() seg {
	rv := seg{}
	rv.a.x = min(s.a.x, s.b.x)
	rv.b.x = max(s.a.x, s.b.x)
	rv.a.y = min(s.a.y, s.b.y)
	rv.b.y = max(s.a.y, s.b.y)
	return rv
}

func pathToSeg(path []string) ([]seg, error) {
	var rv []seg
	x := 0
	y := 0
	for _, p := range path {
		dir := p[0]
		ls := p[1:]
		l, err := strconv.Atoi(ls)
		if err != nil {
			return nil, err
		}
		s := seg{
			a:   point{x, y},
			len: l,
		}

		switch dir {
		case 'R':
			s.b.x = x + l
			s.b.y = y
		case 'L':
			s.b.x = x - l
			s.b.y = y
		case 'U':
			s.b.x = x
			s.b.y = y + l
		case 'D':
			s.b.x = x
			s.b.y = y - l
		default:
			return nil, fmt.Errorf("unknown direction: %b", dir)
		}

		x = s.b.x
		y = s.b.y
		rv = append(rv, s)
	}

	return rv, nil
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var wires [][]seg
	for scanner.Scan() {
		wireS := scanner.Text()
		wire := strings.Split(wireS, ",")
		segs, err := pathToSeg(wire)
		if err != nil {
			panic(err)
		}

		wires = append(wires, segs)
	}

	if len(wires) != 2 {
		panic("didn't find 2 wires")
	}

	if wires == nil {
		panic("no wires")
	}

	var intersections []intersection
	w1d := 0
	for _, s1 := range wires[0] {
		w2d := 0
		for _, s2 := range wires[1] {
			baseLen := w1d + w2d
			sl1 := s1.Slope()
			sl2 := s2.Slope()

			if sl1 == sl2 {
				if sl1 == 0 && s1.a.y == s2.a.y {
					s1n := s1.Normalized()
					s2n := s2.Normalized()
					maxminx := max(s1n.a.x, s2n.a.x)
					minmaxx := min(s1n.b.x, s2n.b.x)

					for i := maxminx; i <= minmaxx; i++ {
						length := baseLen + (i - maxminx) + (s1.a.y - s2n.a.y)
						i := intersection{point{i, s1.a.y}, length}
						intersections = append(intersections, i)
					}
				} else if sl1 == 1 && s1.a.x == s2.a.x {
					s1n := s1.Normalized()
					s2n := s2.Normalized()
					maxminy := max(s1n.a.y, s2n.a.y)
					minmaxy := min(s1n.b.y, s2n.b.y)

					for i := maxminy; i <= minmaxy; i++ {
						length := baseLen + (i - maxminy) + (s1.a.x - s2n.a.x)
						i := intersection{point{s1.a.x, i}, length}
						intersections = append(intersections, i)
					}
				}

			} else if sl1 == 0 {
				s1n := s1.Normalized()
				s2n := s2.Normalized()

				if s2.a.x <= s1n.b.x && s2.a.x >= s1n.a.x &&
					s1.a.y <= s2n.b.y && s1.a.y >= s2n.a.y {
					s1d := abs(s1.a.x - s2.a.x)
					s2d := abs(s1.a.y - s2.a.y)
					length := baseLen + s1d + s2d
					intersect := intersection{point{s2.a.x, s1.a.y}, length}
					intersections = append(intersections, intersect)
				}
			} else {
				s1n := s1.Normalized()
				s2n := s2.Normalized()

				if s2.a.y <= s1n.b.y && s2.a.y >= s1n.a.y &&
					s1.a.x <= s2n.b.x && s1.b.x >= s2n.a.x {
					s1d := abs(s1.a.x - s2.a.x)
					s2d := abs(s2.a.y - s1.a.y)
					length := baseLen + s1d + s2d
					intersect := intersection{point{s1.a.x, s2.a.y}, length}
					intersections = append(intersections, intersect)
				}
			}
			w2d += s2.len
		}
		w1d += s1.len
	}

	minLength := math.MaxInt64
	var minPt point
	for _, intersect := range intersections {
		if intersect.point == (point{}) {
			continue
		}

		length := intersect.dist
		if length < minLength {
			minLength = length
			minPt = intersect.point
		}
	}

	if minPt != (point{}) {
		fmt.Printf("found min intersection: %+v: %d\n", minPt, minLength)
	} else {
		fmt.Println("could not find min intersection")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -1 * a
	}
	return a
}
