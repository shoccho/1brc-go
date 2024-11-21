package main

import (
	"1brc-go/advanced"
)

const (
	dataFilePath = "../1brc/measurements.txt"
)

func main() {
	// simple.Simple(dataFilePath)
	// advanced.CustomMmap(dataFilePath)
	advanced.ParallelMmap(dataFilePath)
}
