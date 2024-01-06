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

## Performance Enhancements

- **Concurrent Map Implementation:** Introduced a sharded concurrent map to reduce lock contention. This allows for more efficient updates to the data structure in a multi-threaded environment.
- **Hash-Based Sharding:** Implemented hash-based sharding for distributing data across multiple shards, further reducing the chance of lock conflicts.
- **Optimized String Processing:** Refined the string handling logic to minimize overhead during file parsing.
- **Buffer Size Adjustments:** Tuned the buffer sizes for channels to balance throughput and memory usage.
- **Efficient Data Aggregation:** Streamlined the data aggregation process for improved efficiency.

Processing Time 5m19s. Tested with a Ryzen 5800x3d

## v3.0.0

## Key Enhancements

- **Parallel File Processing:** Implemented an advanced parallel processing approach where the input file is divided into chunks and processed independently in parallel, drastically reducing I/O bottleneck.
- **Optimized Memory Management:** Refined memory usage by processing data in chunks and employing local maps for data aggregation to reduce memory overhead.
- **Improved Data Aggregation:** Enhanced the efficiency of data aggregation through the use of sharded data structures, minimizing lock contention.

Processing Time: 1m3s. Tested with a Ryzen 5800x3d and 32 gigs Ram

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
{unak=38.8/38.8/38.8, Yuncheng=35.0/35.0/35.0, Yuncos=40.1/40.1/40.1, ...}
Processing completed in 9m1.812065864s
`

## Customization

- You can modify the number of workers in the program to match your CPU's core count for optimal performance.
- Adjust the file path in the program to point to your specific data file location.

## Notes

- Performance may vary based on the hardware specifications and the size of the input file.
- Ensure that the input file format is strictly adhered to for accurate results.
