package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// want to check the status of web domains
// input: input a list of domains
// output: domains with status code

// input type
type domainList []string

// intermediate type
type statusList struct {
	name   string
	status int
	err    error
}

// output type
type statusDB []statusList

func main() {

	domains := readDomainsCsv("domains.csv")

	domains_status := getStatus(domains)

	err := domains_status.saveToOutputCSV("domains_output.csv")
	if err != nil {
		fmt.Println("Error saving CSV:", err)
	}

}

// takes a filename, should be csv and returns a list of domains
func readDomainsCsv(filename string) domainList {

	// read in the data from the csv, if there is an error flag it.
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()

	//create a new file reader and read in all the values and return an error if necessary
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading csv:", err)
	}

	// define the domains as a blank domainList, for each record in the range above append the record to the lsit
	domains := domainList{}
	for _, record := range records {

		// record[0] to choose the current record
		domain := strings.TrimSpace(record[0])
		domain = strings.TrimPrefix(domain, "\ufeff")

		// check if domain contains the proper prefixes
		if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
			domain = "http://" + domain
		}

		// append domain to domains list
		domains = append(domains, domain)

	}

	return domains

}

// gets status code of website and prints a slice of structs
func getStatus(dl domainList) statusDB {

	statusDB := statusDB{}

	for _, domain := range dl {
		// get response from domain with error
		resp, err := http.Get(domain)
		// if there is an error status code is 0 and return error message, append this to the statusDB object.
		if err != nil {

			status := statusList{
				name:   domain,
				status: 0,
				err:    err,
			}

			statusDB = append(statusDB, status)
			continue

		}
		// if no error then extract the details from the domain response and populate the status struct
		status := statusList{
			name:   domain,
			status: resp.StatusCode,
			err:    err,
		}
		// append the status list struct to the statusdb to be returned, close the response body afterwards.
		statusDB = append(statusDB, status)

		resp.Body.Close()

	}

	return statusDB

}

// write statusDB struct to a csv file
func (sd statusDB) saveToOutputCSV(filename string) error {
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
	writer.Write([]string{"domain", "status", "error"})
	// loop through the data in the DB struct and extract the data as strings

	for _, status := range sd {
		statusStr := strconv.Itoa(status.status)
		errStr := ""
		if status.err != nil {
			errStr = status.err.Error()
		}
		// writes the final file, each time it does it it adds a row to the csv
		writer.Write([]string{status.name, statusStr, errStr})
	}

	return nil

}

// add print method to statusDB type
func (sd statusDB) Print() {
	for _, status := range sd {
		fmt.Println(status.name, status.status, status.err)
	}
}
