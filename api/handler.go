package api

import (
	"datastream/logs"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Create an instance of the logger
	logger := logs.NewSimpleLogger()

	tmpl, err := template.ParseFiles("templates/main.html")
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
