package api

import (
	"datastream/logs"
	"html/template"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of the logger
	logger := logs.NewSimpleLogger()

	tmpl, err := template.ParseFiles("templates/HomePage.html")
	if err != nil {
		// Log the error using the logger
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		// Log the error using the logger
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResultPageHandler() {
}

func UploadToKafka() {
	// Implement the logic to upload data to Kafka here
}

func GetDataFromClickHouse() {
	// Implement the logic to get data from ClickHouse here
}
