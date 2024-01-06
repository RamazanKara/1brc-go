package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StationData stores min, mean, max, and total sum of temperatures and count of readings for a station
type StationData struct {
	min, max, sum, count float64
}

func main() {
	startTime := time.Now()
	// Adjust this to the path of your data file
	fileName := "./data/weather_stations.csv"

	// Read and process file concurrently
	stationData := processFileConcurrently(fileName)

	// Prepare output
	outputResults(stationData)

	duration := time.Since(startTime)
	fmt.Printf("Processing completed in %s\n", duration)
}

func processFileConcurrently(fileName string) map[string]*StationData {
	// Number of goroutines to use (can be tuned based on CPU cores)
	const numGoroutines = 16

	// Channel for passing lines to processing goroutines
	linesCh := make(chan string, numGoroutines)

	// WaitGroup to wait for all processing goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Mutex for synchronizing access to the map
	var mu sync.Mutex

	// Map to store the aggregated data
	stationData := make(map[string]*StationData)

	// Start processing goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for line := range linesCh {
				// Process line and update data
				parts := strings.Split(line, ";")
				if len(parts) != 2 {
					continue // Skip malformed lines
				}
				station, tempStr := parts[0], parts[1]
				temp, err := strconv.ParseFloat(tempStr, 64)
				if err != nil {
					continue // Skip lines with invalid temperature
				}

				mu.Lock()
				data, exists := stationData[station]
				if !exists {
					data = &StationData{min: temp, max: temp}
					stationData[station] = data
				}
				data.sum += temp
				data.count++
				if temp < data.min {
					data.min = temp
				}
				if temp > data.max {
					data.max = temp
				}
				mu.Unlock()
			}
		}()
	}

	// Open file and buffer reading
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linesCh <- scanner.Text()
	}
	close(linesCh)

	// Wait for all processing to be done
	wg.Wait()

	return stationData
}

func outputResults(stationData map[string]*StationData) {
	// Extract keys and sort them
	keys := make([]string, 0, len(stationData))
	for key := range stationData {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Generate output
	fmt.Print("{")
	for i, key := range keys {
		data := stationData[key]
		mean := data.sum / data.count
		fmt.Printf("%s=%.1f/%.1f/%.1f", key, data.min, mean, data.max)
		if i < len(keys)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")
}
