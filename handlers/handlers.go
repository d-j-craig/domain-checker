package handlers

import (
	"domain-checker/crawler"
	"domain-checker/processors"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Renders and serves the input form template
func InputHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/input.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// This handlers procesess the csv file and loads the dashboard table
func FileHandler(w http.ResponseWriter, req *http.Request) (string, string){

	file, filename := processors.ParseCSVFile(w, req)

	domains := processors.ScanAndProcessCSV(file)

	domainData := crawler.GetStatus(domains)

	datestamp := time.Now().Format("2006-01-02")
	downloadName :=  datestamp + "_output_" + filename
	outputFilePath := "data/output/" + downloadName
	
	
	err := domainData.SaveToOutputCSV(outputFilePath)
	if err != nil {
		fmt.Println("Error saving CSV:", err)
	}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	data := map[string]any{
		"Domains": domainData,
	}

	tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return downloadName, outputFilePath
	
}



