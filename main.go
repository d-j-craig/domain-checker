package main

import (
	"domain-checker/crawler"
	"domain-checker/csv"
	"fmt"
	"time"
)

// Check the status of web domains
// Input: input a list of domains
// Output: domains with status code

func main() {
	//only used for performance
	appStartTime := time.Now()

	fmt.Println("Started: ", appStartTime)

	filename := "data/input/" + "tranco_list_250.csv"

	domains := csv.ReadDomainsCsv(filename)

	//only used for performance
	crawlStartTime := time.Now()

	domains_status := crawler.GetStatus(domains)

	//only used for performance
	crawlElapsedTime := time.Since(crawlStartTime)

	err := domains_status.SaveToOutputCSV("./data/output/" + "output_" + filename)
	if err != nil {
		fmt.Println("Error saving CSV:", err)
	}

	//only used for performance
	appElapsed := time.Since(appStartTime)

	fmt.Println("Total Time: ", appElapsed)
	fmt.Println("Crawl Time: ", crawlElapsedTime)

}
