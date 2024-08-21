# GoTOPO30

GoTOPO30 is a Golang package designed to read GTOPO30 files and extract elevation data from them. GTOPO30 is a global digital elevation model (DEM) that provides elevation data with a horizontal grid spacing of 30 arc seconds (approximately 1 kilometer). I have included some additional information at the bottom of the README below the usage.

## Features

- **Read GTOPO30 DEM Files**: Parse and read elevation data from GTOPO30 files in Golang.
- **Retrieve Elevation Data**: Get the elevation for any latitude and longitude within the bounds of the GTOPO30 dataset.
- **Supports Large Datasets**: Efficiently handles large binary raster files using Golang's robust file handling capabilities.

## Installation

To install GoTOPO30, simply run:

```bash
go get github.com/1dylan1/gotopo30
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/gotopo30"
)

func main() {
    lat := 37.7749   // Latitude for San Francisco, CA
    lon := -122.4194 // Longitude for San Francisco, CA
    baseFilePath := "/path/to/gtopo30/files"

    elevation, err := gotopo30.GetGTOPOElevation(lat, lon, baseFilePath)
    if err != nil {
        log.Fatalf("Failed to get elevation: %v", err)
    }

    fmt.Printf("Elevation at (%.4f, %.4f): %d meters\n", lat, lon, elevation)
}
```

## `GetGTOPOElevation` Function
The `GetGTOPOElevation` function retrieves the elevation for a given latitude and longitude using the GTOPO30 data files.

### Function Signature
```go
func GetGTOPOElevation(lat float64, lon float64, baseFilePath string) (int16, error)
```
### Parameters
- `lat` (float64): The latitude of the point for which you want to get the elevation.
- `lon` (float64): The longitude of the point.
- `baseFilePath` (string): The base directory path where the GTOPO30 `.DEM` and `.HDR` files are stored.

### Returns
- `int16`: The elevation at the given latitude and longitude in meters.
- `error`: An error message if the elevation could not be retrieved.


## GTOPO Overview 
GTOPO30 is a digital elevation model developed by the United States Geological Survey (USGS). It provides elevation data globally with a 30 arc-second resolution, covering the Earth's land surface between 90 degrees north and 90 degrees south latitude. GTOPO is spread up into 27 tiles, with 6 tiles covering Antarctica, Each non-Antarctic tile will cover 50 degrees of latitude and 40 degrees of longitude, whereas the Antarctic tiles will cover 30 degrees of latitude and 60 degrees of longitude. The tiles' names refer to the longitude and latitude of the upper-left (northwest) corner of the tile. For example, the coordinates of the upper-left corner of the tile E020N40 are 20 degrees east longitude and 40 degrees north latitude. There is one additional tile that covers all of Antarctica with data in a polar stereographic projection, but I do **not** handle this. 

The horizontal coordinate system is
decimal degrees of latitude and longitude referenced to WGS84. The vertical units represent
elevation in meters above mean sea level. The elevation values range from -407 to 8,752 meters.
In the DEM, ocean areas have been masked as "no data" and have been assigned a value of
9999 (hence why I return -9999 for error handling purposes). Lowland coastal areas have an elevation of at least 1 meter, so in the event that a user
reassigns the ocean value from -9999 to 0 the land boundary portrayal will be maintained. Due to
the nature of the raster structure of the DEM, small islands in the ocean less than approximately 1
square kilometer will not be represented.

## Contributing
I'm always open to contributions as it's possible I've missed an edge case or fundamentally misunderstood something. Feel free to open an issue or PR.