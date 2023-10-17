package main

import "fmt"

// https://zhuanlan.zhihu.com/p/105212518

var (
	net = [][]int{
		{1, 1, 0, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{0, 0, 0, 0},
	}

	u = []int{0, 0, 0, 0}
	v = []int{0, 0, 0, 0}

	result = []int{-1, -1, -1, -1} // same as v
)

func match(x int) int {
	u[x] = 1
	for y := range v {
		if v[y] == 0 && net[x][y] == 1 {
			v[y] = 1
			if result[y] == -1 || match(result[y]) == 1 {
				result[y] = x
				return 1
			}
		}
	}

	return 0
}

func reset() {
	u = []int{0, 0, 0, 0}
	v = []int{0, 0, 0, 0}
}

func main() {
	for i := range u {
		reset()
		match(i)
	}
	match(0)
	for y, x := range result {
		if x < 0 {
			continue
		}
		fmt.Printf("x -> y: %d -> %d\n", x, y)
	}
}
