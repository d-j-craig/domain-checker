package main

import (
	"domain-checker/db"
	"domain-checker/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// Check the status of web domains
// Input: input a list of domains
// Output: domains with status code

func main() {
	
	// loads .env into environment variables
	godotenv.Load()
	
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	// closes dbconnection
	defer db.CloseDB()

	mux := http.NewServeMux()

	// Serves form input to upload csv files
	mux.HandleFunc("/", handlers.InputHandler)
	
	// Processes CSV files and sends to dashboard template and also gets the output file name
	mux.HandleFunc("POST /process", handlers.FileHandler)

	// Download data in CSV format from the app.
	mux.HandleFunc("GET /downloadcsv", func(w http.ResponseWriter, req *http.Request){
		downloadName := req.URL.Query().Get("file")
		outputFilePath := req.URL.Query().Get("path")

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, downloadName))
		http.ServeFile(w, req, outputFilePath)
	})

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", mux)

}
