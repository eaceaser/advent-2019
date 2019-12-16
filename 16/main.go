package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var base = [4]int{0, 1, 0, -1}

type pattern struct {
	n   int
	idx int
	pos int
}

func (p *pattern) next() int {
	rv := base[p.pos]
	p.idx = (p.idx + 1) % p.n
	if p.idx == 0 {
		p.pos = (p.pos + 1) % len(base)
	}
	return rv
}

func mkPattern(n int) *pattern {
	if n == 1 {
		return &pattern{
			n:   n,
			idx: 0,
			pos: 1,
		}
	}
	return &pattern{
		n:   n,
		idx: 1,
		pos: 0,
	}
}

func fft(seq []int, phases int) []int {
	rv := make([]int, len(seq))
	copy(rv, seq)
	for p := 0; p < phases; p++ {
		next := make([]int, len(seq))
		for i := 0; i < len(rv); i++ {
			pattern := mkPattern(i + 1)
			res := 0
			for _, s := range rv {
				p := pattern.next()
				res += s * p
			}
			res = abs(res) % 10
			next[i] = res
		}
		copy(rv, next)
	}

	return rv
}

func main() {
	in, err := ioutil.ReadFile("input")
	if err != nil {
		panic(err)
	}
	str := strings.TrimSpace(string(in))

	seq := make([]int, len(str))
	for i, c := range str {
		seq[i] = int(c - '0')
	}

	res := fft(seq, 100)
	fmt.Println(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res[0:8])), ""), "[]"))
}

func abs(i int) int {
	if i < 0 {
		return -1 * i
	}
	return i
}
