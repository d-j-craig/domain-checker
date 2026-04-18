package crawler

import (
	"domain-checker/types"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// takes a filename, should be csv and returns a list of domains
func ReadDomainsCsv(filename string) []string {

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
	domains := []string{}
	for _, record := range records {

		// record[0] to choose the current record
		domain := strings.TrimSpace(record[0])
		domain = strings.TrimPrefix(domain, "\ufeff")

		// check if domain contains the proper prefixes
		if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
			domain = "https://" + domain
		}

		// append domain to domains list
		domains = append(domains, domain)

	}

	return domains

}

// Takes in the slice of domain names from the csv, goes to the site and extracts the data
// it returns a slice of structs call
// sends http GET request to extract status code and estimate domain response time.
func GetStatus(dl []string) types.StatusDB {

	// create a go channel and create a count variable.
	c := make(chan types.StatusList)
	count := len(dl)

	// for each domain in the range take make a request, extract the data, if no errors return the response
	for _, domain := range dl {
		go func(domain string) {

			startTime := time.Now()
			resp, err := http.Get(domain)
			elapsed := time.Since(startTime).Milliseconds()

			// Define status struct for each domain
			status := types.StatusList{
				Name:         domain,
				Status:       0,
				ResponseTime: 0,
				ResponseSize: 0,
				Err:          err,
			}

			// If there are no errors add data to statusList
			if err == nil {
				// open response body and read it, if there are no errors append all data to the status.
				body, readErr := io.ReadAll(resp.Body)
				defer resp.Body.Close()

				if readErr == nil {
					status.Name = strings.Trim(domain, "https://")
					status.Status = resp.StatusCode
					status.ResponseTime = float64(elapsed) // convert to float
					status.ResponseSize = float64(len(body)/1000) // convert to float then bytes to kilobytes
				} else {
					status.Err = readErr
				}
			}

			// Send to go channel
			c <- status
		}(domain)
	}

	// Collect results in a statusDB
	var results types.StatusDB
	for range count {
		status := <-c
		results = append(results, status)
	}

	return results
}


