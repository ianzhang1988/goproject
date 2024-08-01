package main

import (
	"fmt"

	"github.com/mroth/weightedrand/v2"
	"github.com/smallnest/weighted"
)

// too slow
// func ExampleSW_Next() {
// 	w := &weighted.SW{}
// 	w.Add("a", 5)
// 	w.Add("b", 2)
// 	w.Add("c", 3)

// 	for i := 0; i < 9; i++ {
// 		fmt.Printf("%s ", w.Next())
// 	}

// 	fmt.Println()

// 	mapCnt := map[string]int{}

// 	for i := 0; i < 10000; i++ {
// 		c := w.Next().(string)
// 		if _, ok := mapCnt[c]; !ok {
// 			mapCnt[c] = 0
// 		}
// 		mapCnt[c] = 1 + mapCnt[c]
// 	}

// 	fmt.Println(mapCnt)
// }

// func ExampleSWBench(items []int, queryNum int) {

// 	w := &weighted.SW{}
// 	for i := range items {
// 		w.Add(i, items[i])
// 	}

// 	for i := 0; i < queryNum; i++ {
// 		_ = w.Next()
// 	}

// }

func ExampleRRW_Next() {

	w := &weighted.RRW{}
	w.Add("a", 5)
	w.Add("b", 2)
	w.Add("c", 3)

	for i := 0; i < 9; i++ {
		fmt.Printf("%s ", w.Next())
	}

	fmt.Println()

	w = &weighted.RRW{}
	for i := 1; i <= 100; i++ {
		w.Add(i, i)
	}

	mapCnt := map[int]int{}

	for i := 0; i < 101000; i++ {
		c := w.Next().(int)
		if _, ok := mapCnt[c]; !ok {
			mapCnt[c] = 0
		}
		mapCnt[c] = 1 + mapCnt[c]
	}

	fmt.Println(mapCnt)

}

func ExampleRRWBench(items []int, queryNum int) {

	w := &weighted.RRW{}
	for i := range items {
		w.Add(i, items[i])
	}

	counter := 0
	for i := 0; i < queryNum; i++ {
		c := w.Next()
		counter += c.(int)
	}

	// fmt.Println(counter)
}

func WeightedRand() {
	chooser, _ := weightedrand.NewChooser(
		weightedrand.NewChoice('🍒', 0),
		weightedrand.NewChoice('🍋', 1),
		weightedrand.NewChoice('🍊', 1),
		weightedrand.NewChoice('🍉', 3),
		weightedrand.NewChoice('🥑', 5),
	)
	// The following will print 🍋 and 🍊 with 0.1 probability, 🍉 with 0.3
	// probability, and 🥑 with 0.5 probability. 🍒 will never be printed. (Note
	// the weights don't have to add up to 10, that was just done here to make
	// the example easier to read.)
	for i := 0; i < 20; i++ {
		result := chooser.Pick()
		fmt.Printf("%s ", string(result))
	}
	fmt.Println()

	choices := []weightedrand.Choice[int, int]{}
	for i := 1; i <= 100; i++ {
		choices = append(choices, weightedrand.NewChoice(i, i))
	}

	chooser2, _ := weightedrand.NewChooser(
		choices...,
	)

	mapCnt := map[int]int{}

	for i := 0; i < 101000; i++ {
		c := chooser2.Pick()
		if _, ok := mapCnt[c]; !ok {
			mapCnt[c] = 0
		}
		mapCnt[c] = 1 + mapCnt[c]
	}

	fmt.Println(mapCnt)
}

func WeightedRandBench(items []int, queryNum int) {

	choices := []weightedrand.Choice[int, int]{}
	for i := range items {
		choices = append(choices, weightedrand.Choice[int, int]{Item: i, Weight: items[i]})
	}

	chooser2, _ := weightedrand.NewChooser(choices...)

	counter := 0
	for i := 0; i < queryNum; i++ {
		c := chooser2.Pick()
		counter += c
	}
	// fmt.Println(counter)
}

func main() {
	// ExampleSW_Next()
	// fmt.Println("---------------------------")
	ExampleRRW_Next()
	fmt.Println("---------------------------")
	WeightedRand()
}
