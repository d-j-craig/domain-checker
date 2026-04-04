package types

import (
	"encoding/csv"
	"os"
	"strconv"
)

// intermediate type
type StatusList struct {
	Name         string
	Status       int
	ResponseTime float64
	ResponseSize float64
	Err          error
}

// output type
type StatusDB []StatusList

// write statusDB struct to a csv file
func (sd StatusDB) SaveToOutputCSV(filename string) error {

	// create a new file to store the data and handle the errors
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// defer closes the file at the end of the main function
	defer file.Close()

	// creates new csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// define the headers for the final csv
	writer.Write([]string{"domain", "status", "response_time_ms", "response_size_kb", "error"})

	// loop through the data in the DB struct and extract the data as strings
	for _, status := range sd {
		statusStr := strconv.Itoa(status.Status)
		statusResponseTime := strconv.FormatFloat(status.ResponseTime, 'f', 4, 64)
		statusResponseSize := strconv.FormatFloat(status.ResponseSize, 'f', 4, 64)
		errStr := ""
		if status.Err != nil {
			errStr = status.Err.Error()
		}
		// writes the final file, each time it does it it adds a row to the csv
		writer.Write([]string{status.Name, statusStr, statusResponseTime, statusResponseSize, errStr})
	}

	return nil

}
