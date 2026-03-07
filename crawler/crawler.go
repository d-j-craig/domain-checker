package crawler

import (
	"domain-checker/types"
	"net/http"
	"time"
)

// sends http GET request to extract status code and estimate domain response time.
func GetStatus(dl types.DomainList) types.StatusDB {

	// create a go channel and create a count variable.
	c := make(chan types.StatusList)
	count := len(dl)

	// for each domain in the range take make a request, extract the data, if no errors return the response
	for _, domain := range dl {
		go func(domain string) {

			startTime := time.Now()

			resp, err := http.Get(domain)

			elapsed := time.Since(startTime).Seconds()

			// Define status
			status := types.StatusList{
				Name:         domain,
				Status:       0,
				ResponseTime: 0,
				Err:          err,
			}

			// If there are no errors add data to statusList
			if err == nil {
				status.Status = resp.StatusCode
				status.ResponseTime = elapsed
				resp.Body.Close()
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
