package gotopo30

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LatRange struct {
	min   float64
	max   float64
	label string
}

type LonRange struct {
	min   float64
	max   float64
	label string
}

var latRanges = []LatRange{
	{40, 90, "N90"},
	{-10, 40, "N40"},
	{-60, -10, "S10"},
	{-90, -60, "S60"},
}

var lonRanges = []LonRange{
	{-180, -140, "W180"},
	{-140, -100, "W140"},
	{-100, -60, "W100"},
	{-60, -20, "W060"},
	{-20, 20, "W020"},
	{20, 60, "E020"},
	{60, 100, "E060"},
	{100, 140, "E100"},
	{140, 180, "E140"},
}

func getFileName(lat float64, lon float64) (string, error) {
	var latLabel string
	for _, lr := range latRanges {
		if lat >= lr.min && lat <= lr.max {
			latLabel = lr.label
			break
		}
	}

	var lonLabel string
	for _, lr := range lonRanges {
		if lon >= lr.min && lon <= lr.max {
			lonLabel = lr.label
			break
		}
	}

	if latLabel == "S60" {
		switch {
		case lon >= -180 && lon < -120:
			return "W180S60", nil
		case lon >= -120 && lon < -60:
			return "W120S60", nil
		case lon >= -60 && lon < 0:
			return "W060S60", nil
		case lon >= 0 && lon < 60:
			return "W000S60", nil
		case lon >= 60 && lon < 120:
			return "E060S60", nil
		case lon >= 120 && lon <= 180:
			return "E120S60", nil
		}
	}

	if latLabel != "" && lonLabel != "" {
		return lonLabel + latLabel, nil
	}
	return "", errors.New("could not get filename for given lat/lon")
}

type GTOPOHeaderInfo struct {
	byteOrder     rune
	layout        string
	nrows         int64
	ncols         int64
	nbands        int32
	nbits         int64
	bandRowBytes  int64
	totalRowBytes int64
	bandGapBytes  int32
	nodata        int64
	ulxmap        float64
	ulymap        float64
	xdim          float64
	ydim          float64
}

func readHeaderFile(filePath string) (*GTOPOHeaderInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	headerInfo := &GTOPOHeaderInfo{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		key, value := fields[0], fields[1]

		switch key {
		case "BYTEORDER":
			if value == "M" {
				headerInfo.byteOrder = 'M'
			} else {
				headerInfo.byteOrder = 'L'
			}
		case "LAYOUT":
			headerInfo.layout = value
		case "NROWS":
			headerInfo.nrows, err = strconv.ParseInt(value, 10, 64)
		case "NCOLS":
			headerInfo.ncols, err = strconv.ParseInt(value, 10, 64)
		case "NBANDS":
			bands, _ := strconv.ParseInt(value, 10, 32)
			headerInfo.nbands = int32(bands)
		case "NBITS":
			headerInfo.nbits, err = strconv.ParseInt(value, 10, 64)
		case "BANDROWBYTES":
			headerInfo.bandRowBytes, err = strconv.ParseInt(value, 10, 64)
		case "TOTALROWBYTES":
			headerInfo.totalRowBytes, err = strconv.ParseInt(value, 10, 64)
		case "BANDGAPBYTES":
			bandGapBytes, _ := strconv.ParseInt(value, 10, 32)
			headerInfo.bandGapBytes = int32(bandGapBytes)
		case "NODATA":
			headerInfo.nodata, err = strconv.ParseInt(value, 10, 64)
		case "ULXMAP":
			headerInfo.ulxmap, err = strconv.ParseFloat(value, 64)
		case "ULYMAP":
			headerInfo.ulymap, err = strconv.ParseFloat(value, 64)
		case "XDIM":
			headerInfo.xdim, err = strconv.ParseFloat(value, 64)
		case "YDIM":
			headerInfo.ydim, err = strconv.ParseFloat(value, 64)
		}

		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headerInfo, nil
}

func getElevation(filename string, nrows int, ncols int, targetLat float64, targetLon float64, ulxmap float64, ulymap float64, xdim float64, ydim float64) (int16, error) {

	file, err := os.Open(filename)
	if err != nil {
		return -9999, fmt.Errorf("failed to open file :%v", err)
	}
	defer file.Close()

	data := make([]int16, nrows*ncols)

	err = binary.Read(file, binary.BigEndian, data)
	if err != nil {
		return 0, fmt.Errorf("failed to read binary data: %v", err)
	}

	j := int((targetLon - ulxmap) / xdim)
	i := int((ulymap - targetLat) / ydim)

	if i >= 0 && i < nrows && j >= 0 && j < ncols {
		return data[i*ncols+j], nil
	}

	return -9999, fmt.Errorf("coordinates out of bounds")
}

func GetGTOPOElevation(lat float64, lon float64, baseFilePath string) (int16, error) {
	fileName, err := getFileName(lat, lon)
	if err != nil {
		fmt.Println(err)
		return -9999, fmt.Errorf("failed to get file name for lat/lon: %v", err)
	}

	headerInfo, err := readHeaderFile(baseFilePath + "/" + fileName + ".HDR")
	if err != nil {
		return -9999, fmt.Errorf("failed to get header file information: %v", err)
	}

	DEMFilepath := baseFilePath + "/" + fileName + ".DEM"
	elevation, err := getElevation(DEMFilepath, int(headerInfo.nrows), int(headerInfo.ncols), lat, lon, headerInfo.ulxmap, headerInfo.ulymap, headerInfo.xdim, headerInfo.ydim)
	if err != nil {
		fmt.Println(err)
		return -9999, fmt.Errorf("failed to get elevation: %v", err)
	}

	return elevation, nil
}
