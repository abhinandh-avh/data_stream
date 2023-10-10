package main

import (
	"datastream/api"
	"datastream/logs"
	"datastream/routes"
	"log"
	"net/http"
)

func main() {
	logs.FileLog.Info("Welcome")

	router := routes.NewRouter()
	router.AddRoute("GET", "/", api.HomeHandler)
	router.AddRoute("GET", "/result", api.AboutHandler)
	router.AddRoute("POST", "/upload", api.InsertIntoKafkaHandler)
	router.AddRoute("POST", "/result", api.GetFromClickHouseHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	if error := http.ListenAndServe(":8080", router); error != nil {
		log.Fatalf("Server failed to start: %v", error)
	}

}
