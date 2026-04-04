package main

import (
	"domain-checker/handlers"
	"fmt"
	"net/http"
)

// Check the status of web domains
// Input: input a list of domains
// Output: domains with status code

func main() {

	

	// Serves form input to upload csv files
	http.HandleFunc("/", handlers.InputHandler)
	
	// Processes CSV files and sends to dashboard template and also gets the output file name
	var downloadName, outputFilePath string

	http.HandleFunc("/process", func(w http.ResponseWriter, req *http.Request) {
		downloadName, outputFilePath = handlers.FileHandler(w, req)
	})

	// Download data in CSV format from the app.
	http.HandleFunc("/downloadcsv", func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, downloadName))
		http.ServeFile(w, req, outputFilePath)
	})
	
	http.ListenAndServe(":8080", nil)

}
