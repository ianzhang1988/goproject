package main

import "fmt"

type KMAlgo struct {
	weight  *[][]int
	visx    []bool
	visy    []bool
	lx      []int
	ly      []int
	results []int
	d       int
}

func (km *KMAlgo) init() {
	// make weight n*n

	n := len(*km.weight)

	km.visx = make([]bool, n)
	km.visy = make([]bool, n)
	km.lx = make([]int, n)
	km.ly = make([]int, n)
	km.results = make([]int, n)
	km.d = 99999999 // max(weight) ?

	for i := 0; i < n; i++ {
		km.results[i] = -1
		km.lx[i] = -99999999
		for j := 0; j < n; j++ {
			if (*km.weight)[i][j] > km.lx[i] {
				km.lx[i] = (*km.weight)[i][j]
			}
		}
	}
}

func (km *KMAlgo) cleanvis() {
	n := len(*km.weight)
	km.visx = make([]bool, n)
	km.visy = make([]bool, n)
}

func (km *KMAlgo) match(x int) bool {
	km.visx[x] = true // not really used
	for y := range km.visy {
		if !km.visy[y] {
			d := km.lx[x] + km.ly[y] - (*km.weight)[x][y]
			if d == 0 {
				km.visy[y] = true
				if km.results[y] == -1 || km.match(km.results[y]) {
					km.results[y] = x
					return true
				}
			} else {
				if km.d > d {
					km.d = d
				}
			}
		}
	}

	return false
}

func (km *KMAlgo) PrintResult() {
	for y, x := range km.results {
		if x < 0 {
			continue
		}
		fmt.Printf("x -> y: %d -> %d\n", x, y)
	}
}

func (km *KMAlgo) KuhMunkras(weight_input *[][]int) {
	km.weight = weight_input
	km.init()

	for i := range km.visx {
		// fmt.Println(i)
		// km.PrintResult()

		km.d = 99999999
		for {
			km.cleanvis()
			if km.match(i) {
				break
			}

			for i := 0; i < len(*km.weight); i++ {
				if km.visx[i] {
					km.lx[i] -= km.d
				}
				if km.visy[i] {
					km.ly[i] += km.d
				}
			}
		}
	}
}

func main() {
	km := KMAlgo{}
	weight := &[][]int{
		{2, 0, 0, 4},
		{5, 10, 0, 0},
		{0, 3, 4, 0},
		{0, 0, 0, 0},
	}
	km.KuhMunkras(weight)
	km.PrintResult()
}
