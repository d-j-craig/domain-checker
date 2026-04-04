package processors

import (
	"bufio"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

func ParseCSVFile(w http.ResponseWriter, req *http.Request)(file multipart.File, filename string){

	// Parse form of max size 10MB
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("ParseMultipartForm error:", err)
		fmt.Println("Content-Type:", req.Header.Get("Content-Type"))
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}
	// extracts csv from form gets filename for later use
	file, handler, err := req.FormFile("csv-file")
	filename = handler.Filename
	if err != nil {
		http.Error(w, "File could not be retrieved", http.StatusBadRequest)
		return
	}
	defer file.Close()
	// Returns multipart file and filename
	return file, filename

}

// Scans CSV, processes the file and returns the list of domains to be crawled
func ScanAndProcessCSV(file multipart.File)(domains []string){	

	// scans content of file 
	scanner := bufio.NewScanner(file)

	// cleans up white space extracts domain names and makes a slice
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		domains = append(domains, domain)
	}
	return domains

	
	

	

	

	

}