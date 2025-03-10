# One Billion Row Challenge Processor in Go

## Overview

This Go program is designed to efficiently process a large dataset of temperature readings for different weather stations, as part of the [One Billion Row Challenge](https://github.com/gunnarmorling/1brc). The program reads a text file containing temperature measurements, calculates the minimum, mean, and maximum temperature for each station, and outputs the results to stdout. Additionally, it measures and displays the total processing time.

## Key Features (v1.0.0)

- **Concurrency:** Uses goroutines for parallel processing 2 enhance performance on multi-core processors.
- **Efficient File Reading:** Employs buffered reading for handling the 12 gb of dataset more effectively.
- **Data Aggregation:** Calculates min, mean, and max temperatures for each station.
- **Performance Measurement:** Reports the total time taken for processing.

Processing Time: 9m21s. Tested with a Ryzen 5800x3d

## v1.1.0

The program has undergone several optimizations to improve its processing time:

- **Concurrency Model Improved:** Implemented a worker pool pattern for dynamic goroutine management and balanced workload distribution.
- **Buffered Channels:** Increased channel buffer sizes to reduce blocking and increase throughput.
- **Batch Processing:** Process multiple lines of data in a single goroutine to reduce overhead.
- **I/O Enhancements:** Adjusted file reading for larger chunks to reduce I/O bottlenecks.

Processing Time: 6m53s. Tested with a Ryzen 5800x3d

## v2.0.0

Version 2.0 of the One Billion Row Challenge Processor introduces significant optimizations, leading to a substantial reduction in processing time. This release focuses on enhancing concurrency handling and reducing contention, along with other performance improvements.

- **Concurrent Map Implementation:** Introduced a sharded concurrent map to reduce lock contention. This allows for more efficient updates to the data structure in a multi-threaded environment.
- **Hash-Based Sharding:** Implemented hash-based sharding for distributing data across multiple shards, further reducing the chance of lock conflicts.
- **Optimized String Processing:** Refined the string handling logic to minimize overhead during file parsing.
- **Buffer Size Adjustments:** Tuned the buffer sizes for channels to balance throughput and memory usage.
- **Efficient Data Aggregation:** Streamlined the data aggregation process for improved efficiency.

Processing Time 5m19s. Tested with a Ryzen 5800x3d

## v3.0.0

- **Parallel File Processing:** Implemented an advanced parallel processing approach where the input file is divided into chunks and processed independently in parallel, drastically reducing I/O bottleneck.
- **Optimized Memory Management:** Refined memory usage by processing data in chunks and employing local maps for data aggregation to reduce memory overhead.
- **Improved Data Aggregation:** Enhanced the efficiency of data aggregation through the use of sharded data structures, minimizing lock contention.

Processing Time: 1m3s. Tested with a Ryzen 5800x3d and 32 gigs Ram

## v3.1.0

- **Reduced Calls to global map**
- **Implemented a function to determine chunk bounds**

Processing Time: 59.6s. Tested with a Ryzen 5800x3d and 32 gigs Ram

I got this down to 59 Seconds and achieved my goal of getting it to under 1 minute. I am pretty happy with that for a single day session of coding. Further improvements could be made, and if I would continue working on it I would probably directly use a syscall with mmap and use the 8-byte hash of id as a key for an unsafe maphash. And maybe write some tests.

## v3.2.0

- **Memory-Mapped File Processing:** The program now uses `mmap` for file I/O operations, allowing the operating system to handle file paging and reducing overhead.
- **Read-Write Locks:** Replaced `sync.Mutex` with `sync.RWMutex` in the `Shard` struct to reduce lock contention and allow multiple readers concurrently.
- **Optimized Chunk Boundary Determination:** Improved the `determineChunkBounds` function to minimize the number of file seeks and scans, using a more efficient method to find the actual start and end of chunks.

Processing Time: 55.2s. Tested with a Ryzen 5800x3d and 32 gigs Ram

## Requirements

- Go Runtime ofc (1.21)
- Having the Dataset Up and Ready, see here for further instructions: [One Billion Row Challenge](https://github.com/gunnarmorling/1brc)

## How to Run the Program

1. **Prepare the Data File:**
   - Ensure your data file is in the correct format: `<station_name>;<temperature_measurement>` per line.

2. **Compile the Program:**
   - Navigate to the directory containing the program.
   - Run `go build -o brc`.

3. **Execute the Program:**
   - Run the compiled binary: `./brc <file-path>`.
   - The program will read the data file, process the information, and output the results to stdout.
   - The total processing time will be displayed at the end.

## Sample Output

`
{Tampa=-26.5/22.9/80.2, Tashkent=-35.5/14.8/67.7, Tauranga=-32.8/14.8/65.2, ...}
Processing completed in 55.2s
`

## Customization

- You can modify the number of workers in the program to match your CPU's core count for optimal performance.

## Notes

- Performance may vary based on the hardware specifications

## Optimizations

### Memory-Mapped File Processing

The program now uses `mmap` for file I/O operations. This allows the operating system to handle file paging, which can significantly reduce the overhead of file I/O operations, especially for large files.

### Read-Write Locks

The `Shard` struct now uses `sync.RWMutex` instead of `sync.Mutex`. This change allows multiple readers to access the data concurrently, reducing lock contention and improving performance when multiple goroutines are reading from the same shard.

### Optimized Chunk Boundary Determination

The `determineChunkBounds` function has been improved to minimize the number of file seeks and scans. This optimization uses a more efficient method to find the actual start and end of chunks, reducing the time spent on determining chunk boundaries.
