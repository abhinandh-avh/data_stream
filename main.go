package main

import (
	"datastream/api"
	"datastream/logs"
	"datastream/routes"
	"log"
	"net/http"
)

func main() {
	logs.LogInstance()
	defer logs.LogClose()
	logs.FileLog.Info("Log File Created")

	router := routes.NewRouter()
	// Define routes
	router.AddRoute("GET", "/", api.HomeHandler)
	router.AddRoute("GET", "/result", api.AboutHandler)
	router.AddRoute("POST", "/upload", api.InsertIntoKafkaHandler)
	// router.AddRoute("POST", "/result", api.GetFromClickHouseHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	if error := http.ListenAndServe(":8080", router); error != nil {
		log.Fatalf("Server failed to start: %v", error)
	}

}
