package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"golang.org/x/exp/mmap"
)

// StationData holds the temperature data for a specific station.
type StationData struct {
	min, max, sum, count float64
}

const (
	numWorkers = 16   // Number of concurrent workers
	numShards  = 2048 // Number of shards for distributing data
)

// Shard represents a concurrent-safe structure holding station data.
type Shard struct {
	data map[string]*StationData
	lock sync.RWMutex
}

// StationMap aggregates multiple shards for station data.
type StationMap struct {
	shards [numShards]*Shard
}

// NewStationMap initializes a StationMap with predefined shards.
func NewStationMap() *StationMap {
	sm := &StationMap{}
	for i := 0; i < numShards; i++ {
		sm.shards[i] = &Shard{data: make(map[string]*StationData)}
	}
	return sm
}

// GetShard returns a specific shard based on the station key.
func (sm *StationMap) GetShard(key string) *Shard {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return sm.shards[hash.Sum32()%numShards]
}

func main() {
	startTime := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Usage: <program_name> <file_path>")
		os.Exit(1)
	}
	fileName := os.Args[1]

	stationMap := processFile(fileName)
	printResults(stationMap)

	duration := time.Since(startTime)
	fmt.Printf("Processing completed in %s\n", duration)
}

// processFile handles the file processing and returns a StationMap.
func processFile(fileName string) *StationMap {
	reader, err := mmap.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	fileSize := int64(reader.Len())
	chunkSize := fileSize / int64(numWorkers)
	sMap := NewStationMap()
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(chunkStart int) {
			defer wg.Done()
			actualStart, actualEnd := determineChunkBounds(reader, chunkStart, int(chunkSize))
			processChunk(reader, actualStart, actualEnd, sMap)
		}(int(i) * int(chunkSize))
	}

	wg.Wait()
	return sMap
}

// determineChunkBounds calculates the actual boundaries of a file chunk.
func determineChunkBounds(reader *mmap.ReaderAt, chunkStart, chunkSize int) (int, int) {
	var actualStart, actualEnd int

	if chunkStart != 0 {
		actualStart = chunkStart
		for actualStart < chunkStart+chunkSize {
			if reader.At(actualStart) == '\n' {
				actualStart++
				break
			}
			actualStart++
		}
	}

	actualEnd = chunkStart + chunkSize
	for actualEnd < int(reader.Len()) {
		if reader.At(actualEnd) == '\n' {
			actualEnd++
			break
		}
		actualEnd++
	}

	return actualStart, actualEnd
}

// processChunk handles the processing of a specific file chunk.
func processChunk(reader *mmap.ReaderAt, start, end int, sMap *StationMap) {
	localMap := make(map[string]*StationData)

	var currentPos int = start
	for currentPos < end {
		line, err := readLine(reader, currentPos, end)
		if err != nil {
			panic(err)
		}
		currentPos += int(len(line) + 1)

		processLine(strings.TrimSpace(line), localMap)
	}

	mergeLocalMap(localMap, sMap)
}

// readLine reads a line from the memory-mapped file.
func readLine(reader *mmap.ReaderAt, start, end int) (string, error) {
	var line []byte
	for start < end {
		b := reader.At(start)
		if b == '\n' {
			break
		}
		line = append(line, b)
		start++
	}
	return string(line), nil
}

// mergeLocalMap merges local station data into the global StationMap.
func mergeLocalMap(localMap map[string]*StationData, sm *StationMap) {
	for station, data := range localMap {
		shard := sm.GetShard(station)
		shard.lock.Lock()
		if sd, exists := shard.data[station]; exists {
			sd.sum += data.sum
			sd.count += data.count
			sd.min = min(sd.min, data.min)
			sd.max = max(sd.max, data.max)
		} else {
			shard.data[station] = data
		}
		shard.lock.Unlock()
	}
}

// processLine processes a single line of the file.
func processLine(line string, localMap map[string]*StationData) {
	parts := strings.SplitN(line, ";", 2)
	if len(parts) != 2 {
		return
	}

	station, tempStr := parts[0], parts[1]
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return
	}

	sd, exists := localMap[station]
	if !exists {
		sd = &StationData{min: temp, max: temp, sum: temp, count: 1}
		localMap[station] = sd
	} else {
		sd.sum += temp
		sd.count++
		if temp < sd.min {
			sd.min = temp
		}
		if temp > sd.max {
			sd.max = temp
		}
	}
}

// printResults outputs the aggregated station data.
func printResults(sm *StationMap) {
	consolidatedData := make(map[string]*StationData)
	for _, shard := range sm.shards {
		shard.lock.RLock()
		for station, data := range shard.data {
			consolidatedData[station] = data
		}
		shard.lock.RUnlock()
	}

	var keys []string
	for k := range consolidatedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Print("{")
	for i, key := range keys {
		sd := consolidatedData[key]
		mean := sd.sum / sd.count
		fmt.Printf("%s=%.1f/%.1f/%.1f", key, sd.min, mean, sd.max)
		if i < len(keys)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")
}
