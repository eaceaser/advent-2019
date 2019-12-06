package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func checkPw(pw [6]int) bool {
	adjFound := false
	for i := 1; i < 6; i++ {
		a := pw[i-1]
		b := pw[i]
		if a == b {
			adjFound = true
		}
		if b < a {
			return false
		}
	}
	return adjFound
}

func pwToArr(pw int) [6]int {
	var rv [6]int
	for i := 0; i < 6; i++ {
		digit := (pw / int(math.Pow10(5-i))) % 10
		rv[i] = digit
	}
	return rv
}

func main() {
	minS := os.Args[1]
	maxS := os.Args[2]

	min, err := strconv.Atoi(minS)
	if err != nil {
		panic(err)
	}

	max, err := strconv.Atoi(maxS)
	if err != nil {
		panic(err)
	}

	found := 0
	for i := min; i <= max; i++ {
		arr := pwToArr(i)
		if checkPw(arr) {
			fmt.Printf("pw found: %d\n", i)
			found += 1
		}
	}

	fmt.Printf("found %d passwords\n", found)
}
