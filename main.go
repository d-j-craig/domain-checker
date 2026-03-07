package main

import (
	"domain-checker/crawler"
	"domain-checker/csv"
	"flag"
	"fmt"
	"time"
)

// Check the status of web domains
// Input: input a list of domains
// Output: domains with status code

func main() {
	//input filename
	filename := flag.String("file", "data/input/domains.csv", "CSV file with domains")
	outputFile := flag.String("output", "data/output/output_domains.csv", "Output CSV File")
	flag.Parse()

	fmt.Println("Started: ", time.Now())

	domains := csv.ReadDomainsCsv(*filename)

	//only used for performance
	crawlStartTime := time.Now()

	domainsStatus := crawler.GetStatus(domains)

	//only used for performance
	crawlElapsedTime := time.Since(crawlStartTime)

	err := domainsStatus.SaveToOutputCSV(*outputFile)
	if err != nil {
		fmt.Println("Error saving CSV:", err)
	}

	fmt.Println("Crawl Time: ", crawlElapsedTime)

}
