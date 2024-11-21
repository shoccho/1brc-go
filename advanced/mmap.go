package advanced

import (
	"fmt"
	"os"
	"slices"

	"github.com/edsrzf/mmap-go"
)

type Measurement struct {
	Min   int
	Max   int
	Sum   int64
	Count int
}

func CustomMmap(dataFilePath string) {
	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	data, err := mmap.Map(dataFile, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer data.Unmap()
	city := ""
	temperature := 0
	prev := 0
	total := len(data)
	measurements := make(map[string]*Measurement)
	for i := 0; i < total; i++ {
		if data[i] == ';' {
			city = string(data[prev:i])
			temperature = 0
			i += 1
			negative := false
			for data[i] != '\n' {
				ch := data[i]
				if ch == '.' {
					i++
					continue
				} else if ch == '-' {
					negative = true
				} else {
					ch -= '0'
					if ch > 9 {
						panic("invalid")
					}
					temperature = temperature*10 + int(ch)
				}
				i++
			}
			if negative {
				temperature = -temperature
			}
			measurement := measurements[city]
			if measurement == nil {
				measurements[city] = &Measurement{
					Min:   temperature,
					Max:   temperature,
					Sum:   int64(temperature),
					Count: 1,
				}
			} else {
				measurement.Min = min(measurement.Min, temperature)
				measurement.Max = max(measurement.Max, temperature)
				measurement.Sum += int64(temperature)
				measurement.Count++
			}
			prev = i + 1
			city = ""
			temperature = 0
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
		mean := float64(measurement.Sum/10) / float64(measurement.Count)
		max := float64(measurement.Max) / 10
		min := float64(measurement.Min) / 10
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, min, mean, max)
		if idx < len(names)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Print("}\n")
}
