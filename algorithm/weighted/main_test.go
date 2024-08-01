package main

import (
	"math/rand"
	"testing"
)

var g_items []int

func init() {
	g_items = make([]int, 500000)
	for i := range g_items {
		g_items[i] = rand.Intn(10)
	}
}

// func BenchmarkSW(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		// Replace this with the code you want to benchmark
// 		ExampleSWBench(g_items, 100000)
// 	}
// }

func BenchmarkRRW(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Replace this with the code you want to benchmark
		ExampleRRWBench(g_items, 100000)
	}
}

func BenchmarkWeightedRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Replace this with the code you want to benchmark
		WeightedRandBench(g_items, 100000)
	}
}
