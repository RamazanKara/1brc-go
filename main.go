package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StationData holds the temperature data for a station.
type StationData struct {
	min, max, sum, count float64
}

// Constants for the number of workers and shards.
const (
	numWorkers = 16
	numShards  = 32
)

// Shard contains a map of station data and a mutex for concurrent access.
type Shard struct {
	data map[string]*StationData
	lock sync.Mutex
}

// StationMap holds shards for concurrent access to station data.
type StationMap struct {
	shards [numShards]*Shard
}

// NewStationMap initializes a new StationMap with the specified number of shards.
func NewStationMap() *StationMap {
	sm := &StationMap{}
	for i := 0; i < numShards; i++ {
		sm.shards[i] = &Shard{data: make(map[string]*StationData)}
	}
	return sm
}

// GetShard returns the shard for a given station key.
func (sm *StationMap) GetShard(key string) *Shard {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return sm.shards[hash.Sum32()%numShards]
}

// main is the entry point of the program.
func main() {
	startTime := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Usage: brc <file_path>")
		os.Exit(1)
	}
	fileName := os.Args[1]

	stationMap := processFile(fileName)

	printResults(stationMap)

	duration := time.Since(startTime)
	fmt.Printf("Processing completed in %s\n", duration)
}

// processFile processes the file and returns a populated StationMap.
func processFile(fileName string) *StationMap {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(numWorkers)
	var wg sync.WaitGroup

	sMap := NewStationMap()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(chunkStart int64) {
			defer wg.Done()
			processChunk(fileName, chunkStart, chunkSize, sMap)
		}(int64(i) * chunkSize)
	}

	wg.Wait()
	return sMap
}

// processChunk processes a chunk of the file.
func processChunk(fileName string, offset, size int64, sMap *StationMap) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.Seek(offset, 0); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	localMap := make(map[string]*StationData)

	if offset != 0 {
		_, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
	}

	var bytesRead int64
	for {
		line, err := reader.ReadString('\n')
		bytesRead += int64(len(line))

		if err == io.EOF || (offset+bytesRead) >= (offset+size) {
			break
		}
		if err != nil {
			panic(err)
		}

		processLine(strings.TrimSpace(line), localMap)
	}

	mergeLocalMap(localMap, sMap)
}

// mergeLocalMap merges a local map of station data into the global StationMap.
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

// processLine processes a single line of input and updates the local map.
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

// printResults prints the aggregated results from the StationMap.
func printResults(sm *StationMap) {
	consolidatedData := make(map[string]*StationData)
	for _, shard := range sm.shards {
		shard.lock.Lock()
		for station, data := range shard.data {
			consolidatedData[station] = data
		}
		shard.lock.Unlock()
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
