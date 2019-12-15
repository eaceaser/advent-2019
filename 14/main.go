package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	Ore  = "ORE"
	Fuel = "FUEL"

	MaxOre = 1000000000000
)

var pattern = regexp.MustCompile("(\\d+) (\\w+)")

type reagent struct {
	quantity int
	name     string
}

type spell struct {
	inputs []reagent
	output reagent
}

func (s spell) Ore() int {
	if s.inputs[0].name == Ore {
		return s.inputs[0].quantity
	}
	return 0
}

type plan struct {
	inventory map[string]int
	ore       int
}

func (p *plan) Save() string {
	var pairs []string
	for k, v := range p.inventory {
		pairs = append(pairs, fmt.Sprintf("%s:%d", k, v))
	}

	ss := sort.StringSlice(pairs)
	ss.Sort()
	return strings.Join(pairs, " ")
}

func (p *plan) Load(saved string) {
	rv := make(map[string]int, len(p.inventory))
	reagents := strings.Split(saved, " ")
	for _, r := range reagents {
		p2 := strings.SplitN(r, ":", 2)
		name := p2[0]
		amnt, err := strconv.Atoi(p2[1])
		if err != nil {
			panic(err)
		}
		rv[name] = amnt
	}
	p.inventory = rv
}

func calc(spellbook []spell) int {
	type fuelMemoed struct {
		fuel int
		ore  int
	}

	spellGraph := map[string][]spell{}
	for _, s := range spellbook {
		s1 := s
		spellGraph[s.output.name] = append(spellGraph[s.output.name], s1)
	}

	var f func(spell, *plan)
	f = func(spell spell, plan *plan) {
		if o := spell.Ore(); o > 0 {
			plan.ore += o
			return
		}

		for _, ir := range spell.inputs {
			onHand := plan.inventory[ir.name]
			if onHand < ir.quantity {
				reagentSpell := spellGraph[ir.name][0]
				needed := ir.quantity - onHand
				for needed > 0 {
					f(reagentSpell, plan)
					plan.inventory[reagentSpell.output.name] += reagentSpell.output.quantity
					needed -= reagentSpell.output.quantity
				}
			}
			plan.inventory[ir.name] -= ir.quantity
		}
	}

	p := plan{
		inventory: map[string]int{},
		ore:       0,
	}

	fuelSpell := spellGraph[Fuel][0]
	fuel := 0
	fuelMemo := map[string]fuelMemoed{}
	for p.ore < MaxOre {
		str := p.Save()
		if o, ok := fuelMemo[str]; ok {
			diff := p.ore - o.ore
			if MaxOre-p.ore < diff {
				goto Fallback
			}
			p.ore += diff
			fd := fuel - o.fuel
			fuel += fd
			continue
		}
	Fallback:
		fuelMemo[str] = fuelMemoed{
			fuel: fuel,
			ore:  p.ore,
		}
		if fuel%5000 == 0 {
			fmt.Printf("f=%d o=%d p=%f\n", fuel, p.ore, 100*float64(p.ore)/MaxOre)

		}
		f(fuelSpell, &p)
		fuel++
	}
	return fuel - 1
}

func mustParseReagent(s string) reagent {
	matches := pattern.FindAllStringSubmatch(s, -1)
	name := matches[0][2]
	amount, err := strconv.Atoi(matches[0][1])
	if err != nil {
		panic(err)
	}
	return reagent{
		quantity: amount,
		name:     name,
	}
}

func main() {
	var spellbook []spell
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	var spellStrs []string
	for scan.Scan() {
		spellStrs = append(spellStrs, scan.Text())
	}
	if err := scan.Err(); err != nil {
		panic(err)
	}

	for _, spellStr := range spellStrs {
		parts := strings.Split(spellStr, " => ")
		inputStr := parts[0]
		results := parts[1]
		is := strings.Split(inputStr, ", ")
		var inputs []reagent
		for _, s := range is {
			i := mustParseReagent(s)
			inputs = append(inputs, i)
		}
		result := mustParseReagent(results)
		s := spell{
			inputs: inputs,
			output: result,
		}
		spellbook = append(spellbook, s)
	}
	cost := calc(spellbook)
	fmt.Printf("calculated %d fuel\n", cost)
}
