package advanced

import (
	"bytes"
	"os"
	"runtime"

	"github.com/edsrzf/mmap-go"
)

type item struct {
	key         []byte
	measurement *Measurement
}

const (
	numBuckets = 1 << 17

	// FNV-1 64-bit constants from hash/fnv.
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

func processCustomChunk(ch chan map[string]*Measurement, data mmap.MMap, start, end int) {
	temperature := 0
	prev := start
	hash := uint64(offset64)
	// measurements := make(map[string]*Measurement)

	items := make([]item, numBuckets)
	size := 0

	for i := start; i < end; i++ {
		hash ^= uint64(data[i])
		hash *= prime64
		if data[i] == ';' {
			city_bytes := data[prev:i]
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
			// find index of hash
			hashIndex := int(hash & uint64(numBuckets-1))

			for {
				if items[hashIndex].key == nil {
					// key := make([]byte, len(city))
					// copy(key, []byte(city))
					items[hashIndex] = item{
						key: city_bytes,
						measurement: &Measurement{
							Min:   temperature,
							Max:   temperature,
							Sum:   int64(temperature),
							Count: 1,
						},
					}
					size++
					break
				}
				if bytes.Equal(items[hashIndex].key, city_bytes) {
					measurement := items[hashIndex].measurement
					measurement.Min = min(measurement.Min, temperature)
					measurement.Max = max(measurement.Max, temperature)
					measurement.Sum += int64(temperature)
					measurement.Count++
					break
				}
				hashIndex++
				if hashIndex >= numBuckets {
					// panic("collision")
					hashIndex = 0
				}

			}
			prev = i + 1
			temperature = 0
			hash = uint64(offset64)
		}
	}
	measurements := make(map[string]*Measurement)
	for _, item := range items {
		if item.key == nil {
			continue
		}
		measurements[string(item.key)] = item.measurement
	}
	ch <- measurements
}

func CustomHashMap(dataFilePath string) {

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
		go processCustomChunk(measurementChanel, data, chunks[i].start, chunks[i].end)
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
