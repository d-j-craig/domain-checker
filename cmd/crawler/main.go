package main

import (
	"domain-checker/crawler"
	"domain-checker/db"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)


func main() {

	// loads .env into environment variables
	godotenv.Load()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	defer db.CloseDB()

	filename := "domains.csv"

	domains := crawler.ReadDomainsCsv("./data/input/" + filename)

	domainsStatus := crawler.GetStatus(domains)

	db.InsertDomainData(domainsStatus)
	
	fmt.Println("Data written to DB.")

}