package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const COM = "COM"

type object struct {
	name     string
	orbits   []*object
	orbiters []*object
}

var objects = make(map[string]*object, 0)

// le lazy
func triangle(num int) int {
	rv := 0
	for i := num; i > 0; i-- {
		rv += i
	}
	return rv
}

func countOrbits(head *object) int {
	visited := map[string]struct{}{}

	var f func(*object, int) int
	f = func(obj *object, depth int) int {
		rv := depth
		for _, o := range obj.orbiters {
			if _, ok := visited[o.name]; ok {
				continue
			}
			rv += f(o, depth+1)
		}
		return rv
	}

	return f(head, 0)
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := scanner.Text()
		objectNames := strings.Split(l, ")")
		if len(objectNames) != 2 {
			panic(fmt.Sprintf("malformed line %s", l))
		}
		orbiteeName := objectNames[0]
		orbiterName := objectNames[1]

		orbitee, ok := objects[orbiteeName]
		if !ok {
			orbitee = &object{
				name: orbiteeName,
			}
			objects[orbiteeName] = orbitee
		}

		orbiter, ok := objects[orbiterName]
		if !ok {
			orbiter = &object{
				name:   orbiterName,
				orbits: nil,
			}
			objects[orbiterName] = orbiter
		}

		orbitee.orbiters = append(orbitee.orbiters, orbiter)
		orbiter.orbits = append(orbiter.orbits, orbitee)
	}

	com, ok := objects[COM]
	if !ok {
		panic("could not find COM")
	}

	orbits := countOrbits(com)
	fmt.Printf("found %d orbits\n", orbits)
}
