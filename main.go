package main

import (
    "bufio"
    "fmt"
    "hash/fnv"
    "os"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

type StationData struct {
    min, max, sum, count float64
}

const numWorkers = 16
const numShards = 32 // Number of shards in the concurrent map

type ConcurrentMap struct {
    shards [numShards]map[string]*StationData
    locks  [numShards]*sync.Mutex
}

func NewConcurrentMap() *ConcurrentMap {
    cMap := &ConcurrentMap{}
    for i := 0; i < numShards; i++ {
        cMap.shards[i] = make(map[string]*StationData)
        cMap.locks[i] = &sync.Mutex{}
    }
    return cMap
}

func (cMap *ConcurrentMap) GetShard(key string) (shard map[string]*StationData, lock *sync.Mutex) {
    hash := fnv.New32()
    hash.Write([]byte(key))
    shardIndex := hash.Sum32() % numShards
    return cMap.shards[shardIndex], cMap.locks[shardIndex]
}

func main() {
    startTime := time.Now()

    if len(os.Args) < 2 {
        fmt.Println("Usage: brc <file_path>")
        os.Exit(1)
    }
    fileName := os.Args[1]

    stationData := processFile(fileName)

    printResults(stationData)

    duration := time.Since(startTime)
    fmt.Printf("Processing completed in %s\n", duration)
}

func processFile(fileName string) *ConcurrentMap {
    linesCh := make(chan string, 1000)
    var wg sync.WaitGroup
    wg.Add(numWorkers)

    cMap := NewConcurrentMap()

    for i := 0; i < numWorkers; i++ {
        go worker(&wg, linesCh, cMap)
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

    return cMap
}

func worker(wg *sync.WaitGroup, lines <-chan string, cMap *ConcurrentMap) {
    defer wg.Done()
    for line := range lines {
        processLine(line, cMap)
    }
}

func processLine(line string, cMap *ConcurrentMap) {
    parts := strings.Split(line, ";")
    if len(parts) != 2 {
        return
    }

    station, tempStr := parts[0], parts[1]
    temp, err := strconv.ParseFloat(tempStr, 64)
    if err != nil {
        return
    }

    shard, lock := cMap.GetShard(station)
    lock.Lock()
    data, exists := shard[station]
    if !exists {
        data = &StationData{min: temp, max: temp, sum: temp, count: 1}
        shard[station] = data
    } else {
        data.sum += temp
        data.count++
        if temp < data.min {
            data.min = temp
        }
        if temp > data.max {
            data.max = temp
        }
    }
    lock.Unlock()
}

func printResults(cMap *ConcurrentMap) {
    // Consolidate data from shards
    consolidatedData := make(map[string]*StationData)
    for _, shard := range cMap.shards {
        for station, data := range shard {
            consolidatedData[station] = data
        }
    }

    // Sort the station names
    keys := make([]string, 0, len(consolidatedData))
    for station := range consolidatedData {
        keys = append(keys, station)
    }
    sort.Strings(keys)

    // Print sorted results
    fmt.Print("{")
    for i, key := range keys {
        data := consolidatedData[key]
        mean := data.sum / data.count
        fmt.Printf("%s=%.1f/%.1f/%.1f", key, data.min, mean, data.max)
        if i < len(keys)-1 {
            fmt.Print(", ")
        }
    }
    fmt.Println("}")
}
