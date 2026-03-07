package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Check the status of web domains
// Input: input a list of domains
// Output: domains with status code

// input type
type domainList []string

// intermediate type
type statusList struct {
	name         string
	status       int
	responseTime float64
	err          error
}

// output type
type statusDB []statusList

func main() {
	//only used for performance
	appStartTime := time.Now()

	fmt.Println("Started: ", appStartTime)

	filename := "./data/input/" + "tranco_list_200.csv"

	domains := readDomainsCsv(filename)

	//only used for performance
	crawlStartTime := time.Now()

	domains_status := getStatus(domains)

	//only used for performance
	crawlElapsedTime := time.Since(crawlStartTime)

	err := domains_status.saveToOutputCSV("./data/output/" + "output_" + filename)
	if err != nil {
		fmt.Println("Error saving CSV:", err)
	}

	//only used for performance
	appElapsed := time.Since(appStartTime)

	fmt.Println("Total Time: ", appElapsed)
	fmt.Println("Crawl Time: ", crawlElapsedTime)

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

	// define the domains as a blank domainList, for each record in the range above append the record to the list
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

// sends http GET request to extract status code and estimate domain response time.
func getStatus(dl domainList) statusDB {

	// create a go channel and create a count variable.
	c := make(chan statusList)
	count := len(dl)

	// for each domain in the range take make a request, extract the data, if no errors return the response
	for _, domain := range dl {
		go func(domain string) {

			startTime := time.Now()

			resp, err := http.Get(domain)

			elapsed := time.Since(startTime).Seconds()

			// Define status
			status := statusList{
				name:         domain,
				status:       0,
				responseTime: 0,
				err:          err,
			}

			// If there are no errors add data to statusList
			if err == nil {
				status.status = resp.StatusCode
				status.responseTime = elapsed
				resp.Body.Close()
			}

			// Send to go channel
			c <- status
		}(domain)
	}

	// Collect results in a statusDB
	var results statusDB
	for range count {
		status := <-c
		results = append(results, status)
	}

	return results
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
	writer.Write([]string{"domain", "status", "response_time_s", "error"})

	// loop through the data in the DB struct and extract the data as strings
	for _, status := range sd {
		statusStr := strconv.Itoa(status.status)
		statusResponseTime := strconv.FormatFloat(status.responseTime, 'f', 4, 64)
		errStr := ""
		if status.err != nil {
			errStr = status.err.Error()
		}
		// writes the final file, each time it does it it adds a row to the csv
		writer.Write([]string{status.name, statusStr, statusResponseTime, errStr})
	}

	return nil

}

// add print method to statusDB type
func (sd statusDB) Print() {
	for _, status := range sd {
		fmt.Println(status.name, status.status, status.responseTime, status.err)
	}
}
