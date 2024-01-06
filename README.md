# One Billion Row Challenge Processor in Go

## Overview

This Go program is designed to efficiently process a large dataset of temperature readings for different weather stations, as part of the One Billion Row Challenge. The program reads a text file containing temperature measurements, calculates the minimum, mean, and maximum temperature for each station, and outputs the results to standard output (stdout). Additionally, it measures and displays the total processing time.

## Key Features

- **Concurrency:** Uses goroutines for parallel processing, enhancing performance on multi-core processors.
- **Efficient File Reading:** Employs buffered reading for handling large files effectively.
- **Data Aggregation:** Calculates min, mean, and max temperatures for each station.
- **Performance Measurement:** Reports the total time taken for processing.

## Requirements

- Go Binaries ofc

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
Processing completed in 286.686139ms
`

YES, this really took only 2.87 s in go. Tested with a Ryzen 5800x3d

## Customization

- You can modify the number of goroutines in the program to match your CPU's core count for optimal performance.
- Adjust the file path in the program to point to your specific data file location.

## Notes

- Performance may vary based on the hardware specifications and the size of the input file.
- Ensure that the input file format is strictly adhered to for accurate results.
