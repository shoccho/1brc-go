package Simple

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Measurement struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func Simple(dataFilePath string) {
	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()
	measurements := make(map[string]*Measurement)

	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		rawString := fileScanner.Text()
		city, tempStr, found := strings.Cut(rawString, ";")
		if !found {
			// panic("Invalid line ")
			continue
		}
		temperature, err := strconv.ParseFloat(tempStr, 32)
		if err != nil {
			panic(rawString)
		}
		measurement := measurements[city]
		if measurement == nil {
			measurements[city] = &Measurement{
				Min:   temperature,
				Max:   temperature,
				Sum:   temperature,
				Count: 1,
			}
		} else {
			measurement.Max = max(measurement.Max, temperature)
			measurement.Min = min(measurement.Min, temperature)
			measurement.Sum += temperature
			measurement.Count++
		}
	}
	printResult(measurements)
}
func printResult(results map[string]*Measurement) {
	names := make([]string, 0, len(results))
	for name := range results {
		names = append(names, name)
	}
	slices.Sort(names)
	fmt.Print("{")
	for idx, name := range names {
		measurement := results[name]
		mean := measurement.Sum / float64(measurement.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, measurement.Min, mean, measurement.Max)
		if idx < len(names)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Print("}\n")
}
