package main

import simple "1brc-go/simple"

const (
	dataFilePath = "../1brc/measurements.txt"
)

func main() {
	simple.Simple(dataFilePath)
}
