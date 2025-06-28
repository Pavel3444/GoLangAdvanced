package main

import (
	"fmt"
	"math/rand"
)

func main() {
	sliceLength := 10
	inputChan := make(chan int, sliceLength)
	outputChan := make(chan int, sliceLength)

	go initialSlice(sliceLength, 0, 100, inputChan)
	go operationSliceElements(inputChan, outputChan)

	for n := range outputChan {
		fmt.Printf("%d ", n)
	}
}

func initialSlice(length int, minValue int, maxValue int, inputChan chan int) {
	for i := 0; i < length; i++ {
		n := rand.Intn(maxValue-minValue+1) + minValue
		inputChan <- n
	}
	close(inputChan)
}

func operationSliceElements(out chan int, outputChan chan int) {
	for n := range out {
		outputChan <- n * n
	}
	close(outputChan)
}
