package crawler

import (
	"domain-checker/types"
	"io"
	"net/http"
	"time"
)

// Takes in the slice of domain names from the form, goes to the site and extracts the data
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
			resp, err := http.Get("https://" + domain)
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
					status.Status = resp.StatusCode
					status.ResponseTime = float64(elapsed) // convert to float
					status.ResponseSize = float64(len(body)/1000) // convert to flaot then bytes to kilobytes
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
