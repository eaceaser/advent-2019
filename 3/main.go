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

func (p point) Dist() int {
	return abs(p.x) + abs(p.y)
}

type seg struct {
	a point
	b point
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
			a: point{x,y},
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

	var intersections []point
	for _, s1 := range wires[0] {
		for _, s2 := range wires[1] {
			sl1 := s1.Slope()
			sl2 := s2.Slope()

			if sl1 == sl2 {
				if sl1 == 0 && s1.a.y == s2.a.y {
					minx1 := min(s1.a.x, s1.b.x)
					minx2 := min(s2.a.x, s2.b.x)
					maxminx := max(minx1, minx2)
					maxx1 := max(s1.a.x, s1.b.x)
					maxx2 := max(s2.a.x, s2.b.x)
					minmaxx := min(maxx1, maxx2)

					for i := maxminx; i <= minmaxx; i++ {
						i := point{i, s1.a.y}
						intersections = append(intersections, i)
					}
				} else if sl1 == 1 && s1.a.x == s2.a.x {
					miny1 := min(s1.a.y, s1.b.y)
					miny2 := min(s2.a.y, s2.b.y)
					maxminy := max(miny1, miny2)
					maxy1 := max(s1.a.y, s1.b.y)
					maxy2 := max(s2.a.y, s2.b.y)
					minmaxy := min(maxy1, maxy2)

					for i := maxminy; i <= minmaxy; i++ {
						i := point{s1.a.x, i}
						intersections = append(intersections, i)
					}
				}

			} else if sl1 == 0 {
				s1n := s1.Normalized()
				s2n := s2.Normalized()

				if s2.a.x <= s1n.b.x && s2.a.x >= s1n.a.x &&
					s1.a.y <= s2n.b.y && s1.a.y >= s2n.a.y {
					intersections = append(intersections, point{s2.a.x, s1.a.y})
				}
			} else {
				s1n := s1.Normalized()
				s2n := s2.Normalized()

				if s2.a.y <= s1n.b.y && s2.a.y >= s1n.a.y &&
					s1.a.x <= s2n.b.x && s1.b.x >= s2n.a.x {
					intersections = append(intersections, point{s1.a.x, s2.a.y})
				}
			}
		}
	}

	minDist := math.MaxInt64
	var minPt point
	for _, intersect := range intersections {
		if intersect == (point{}) {
			continue
		}

		dist := intersect.Dist()
		if dist < minDist {
			minDist = dist
			minPt = intersect
		}
	}

	if minPt != (point{}) {
		fmt.Printf("found min intersection: %+v: %d\n", minPt, minPt.Dist())
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
