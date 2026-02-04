package main

import (
	"fmt"
	"math/rand/v2"
)

const iterations = 10

func main() {
	intCh := make(chan int)
	squareCh := make(chan int)
	go getRandomInts(intCh)
	go getSquares(intCh, squareCh)
	var res []int

	for v := range squareCh {
		res = append(res, v)
	}
	for _, v := range res {
		fmt.Printf("%d ", v)
	}
}

func getRandomInts(out chan<- int) {
	defer close(out)
	var res []int
	for range iterations {
		res = append(res, rand.IntN(100))
	}
	for _, n := range res {
		out <- n
	}
}

func getSquares(in <-chan int, out chan<- int) {
	defer close(out)
	for n := range in {
		out <- n * n
	}
}
