# One Billion Row Challenge Processor in Go

## Overview

This Go program is designed to efficiently process a large dataset of temperature readings for different weather stations, as part of the [One Billion Row Challenge](https://github.com/gunnarmorling/1brc). The program reads a text file containing temperature measurements, calculates the minimum, mean, and maximum temperature for each station, and outputs the results to stdout. Additionally, it measures and displays the total processing time.

## Key Features (v1.0.0)

- **Concurrency:** Uses goroutines for parallel processing 2 enhance performance on multi-core processors.
- **Efficient File Reading:** Employs buffered reading for handling the 12 gb of dataset more effectively.
- **Data Aggregation:** Calculates min, mean, and max temperatures for each station.
- **Performance Measurement:** Reports the total time taken for processing.

Processing Time: 9m21s. Tested with a Ryzen 5800x3d

## Recent Optimizations (v1.1.0)

The program has undergone several optimizations to improve its processing time:

- **Concurrency Model Improved:** Implemented a worker pool pattern for dynamic goroutine management and balanced workload distribution.
- **Buffered Channels:** Increased channel buffer sizes to reduce blocking and increase throughput.
- **Batch Processing:** Process multiple lines of data in a single goroutine to reduce overhead.
- **I/O Enhancements:** Adjusted file reading for larger chunks to reduce I/O bottlenecks.

Processing Time: 6m53s. Tested with a Ryzen 5800x3d

## Requirements

- Go Runtime ofc (1.21)
- Having the Dataset Up and Ready, see here for further instructions: [One Billion Row Challenge](https://github.com/gunnarmorling/1brc)

## How to Run the Program

1. **Prepare the Data File:**
   - Ensure your data file is in the correct format: `<station_name>;<temperature_measurement>` per line.
   - Place the data file in `./data/weather_stations.csv`, or update the file path in the program accordingly.

2. **Compile the Program:**
   - Navigate to the directory containing the program.
   - Run `go build -o brc`.

3. **Execute the Program:**
   - Run the compiled binary: `./brc`.
   - The program will read the data file, process the information, and output the results to stdout.
   - The total processing time will be displayed at the end.

## Sample Output

`
{unak=38.8/38.8/38.8, Yuncheng=35.0/35.0/35.0, Yuncos=40.1/40.1/40.1, ...}
Processing completed in 9m 21
`

## Customization

- You can modify the number of workers in the program to match your CPU's core count for optimal performance.
- Adjust the file path in the program to point to your specific data file location.

## Notes

- Performance may vary based on the hardware specifications and the size of the input file.
- Ensure that the input file format is strictly adhered to for accurate results.
