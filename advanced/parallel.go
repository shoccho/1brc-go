package advanced

import (
	"os"
	"runtime"

	"github.com/edsrzf/mmap-go"
)

type Chunk struct {
	start int
	end   int
}

func splitMem(mem mmap.MMap, n int) []Chunk {
	total := len(mem)
	chunkSize := total / n
	chunks := make([]Chunk, n)
	chunks[0].start = 0
	for i := 1; i < n; i++ {
		for j := i * chunkSize; j < i*chunkSize+50; j++ {
			if mem[j] == '\n' {
				chunks[i-1].end = j
				chunks[i].start = j + 1
				break
			}
		}
	}
	chunks[n-1].end = total - 1
	return chunks
}
func processChunk(ch chan map[string]*Measurement, data mmap.MMap, start, end int) {
	city := ""
	temperature := 0
	prev := start
	measurements := make(map[string]*Measurement)
	for i := start; i < end; i++ {
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
	ch <- measurements
}

func ParallelMmap(dataFilePath string) {
	maxCpu := min(runtime.NumCPU(), runtime.GOMAXPROCS(0))
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

	chunks := splitMem(data, maxCpu)
	allMeasurements := make(map[string]*Measurement)
	measurementChanel := make(chan map[string]*Measurement)

	for i := 0; i < maxCpu; i++ {
		go processChunk(measurementChanel, data, chunks[i].start, chunks[i].end)
	}
	for i := 0; i < maxCpu; i++ {
		measurements := <-measurementChanel
		for city, measurement := range measurements {
			storedMeasurement := allMeasurements[city]
			if storedMeasurement == nil {
				allMeasurements[city] = measurement
			} else {
				storedMeasurement.Min = min(storedMeasurement.Min, measurement.Min)
				storedMeasurement.Max = max(storedMeasurement.Max, measurement.Max)
				storedMeasurement.Sum += measurement.Sum
				storedMeasurement.Count += measurement.Count
			}
		}
	}

	printResult(allMeasurements)

}
