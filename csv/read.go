package csv

import (
	"domain-checker/types"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// takes a filename, should be csv and returns a list of domains
func ReadDomainsCsv(filename string) types.DomainList {

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
	domains := types.DomainList{}
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
