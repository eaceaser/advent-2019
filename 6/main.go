package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	COM = "COM"
	YOU = "YOU"
	SAN = "SAN"
)

type object struct {
	name     string
	orbits   []*object
	orbiters []*object
}

var objects = make(map[string]*object, 0)

func countTransfers(from *object, to *object) int {
	visited := map[string]struct{}{}

	var f func(*object, int) int
	f = func(obj *object, depth int) int {
		toOrbit := append(obj.orbits, obj.orbiters...)
		for _, o := range toOrbit {
			if _, ok := visited[o.name]; ok {
				continue
			}
			visited[o.name] = struct{}{}

			if o.name == to.name {
				return depth - 1
			}

			rv := f(o, depth+1)
			if rv != 0 {
				return rv
			}
		}
		return 0
	}

	return f(from, 0)
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

	you, ok := objects[YOU]
	if !ok {
		panic("could not find YOU")
	}

	san, ok := objects[SAN]
	if !ok {
		panic("could not find SAN")
	}

	cnt := countTransfers(you, san)
	fmt.Printf("found %d transfers between %s and %s\n", cnt, you.name, san.name)
}
