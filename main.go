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

const numWorkers = 16 // Number of worker goroutines

func main() {
    startTime := time.Now()

	// Adjust this to the path of your data file
    fileName := "./data/measurements.txt"
    stationData := processFile(fileName)

    printResults(stationData)

    duration := time.Since(startTime)
    fmt.Printf("Processing completed in %s\n", duration)
}

func processFile(fileName string) map[string]*StationData {
    linesCh := make(chan string, 1000)

    var wg sync.WaitGroup
    wg.Add(numWorkers)

    stationData := make(map[string]*StationData)
    var mu sync.Mutex

    // Worker pool pattern
    for i := 0; i < numWorkers; i++ {
        go worker(&wg, linesCh, stationData, &mu)
    }

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
    wg.Wait()

    return stationData
}

func worker(wg *sync.WaitGroup, lines <-chan string, data map[string]*StationData, mu *sync.Mutex) {
    defer wg.Done()
    for line := range lines {
        processLine(line, data, mu)
    }
}

func processLine(line string, data map[string]*StationData, mu *sync.Mutex) {
    parts := strings.Split(line, ";")
    if len(parts) != 2 {
        return
    }

    station, tempStr := parts[0], parts[1]
    temp, err := strconv.ParseFloat(tempStr, 64)
    if err != nil {
        return
    }

    mu.Lock()
    defer mu.Unlock()

    if sd, exists := data[station]; exists {
        sd.sum += temp
        sd.count++
        if temp < sd.min {
            sd.min = temp
        }
        if temp > sd.max {
            sd.max = temp
        }
    } else {
        data[station] = &StationData{min: temp, max: temp, sum: temp, count: 1}
    }
}

func printResults(stationData map[string]*StationData) {
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
