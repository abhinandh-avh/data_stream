package api

import (
	"datastream/dataprocess"
	"datastream/logs"

	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file for the "/" route
	http.ServeFile(w, r, "templates/HomePage.html")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the about.html file for the "/about" route
	http.ServeFile(w, r, "templates/ResultPage.html")
}

func InsertIntoKafkaHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file for the "/insert" route
	// Get the uploaded file from the request
	logs.FileLog.Info("hello")
	file, _, err := r.FormFile("file")
	if err != nil {
		logs.FileLog.Info("ello")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Pass the file to a function in the data processing package to insert into Kafka
	err = dataprocess.InsertCSVIntoKafka(file, "topic")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs.FileLog.Info("CSV data inserted into Kafka successfully!")
	dataprocess.ExtractFromKafka(1)

}

//	func GetFromClickHouseHandler(w http.ResponseWriter, r *http.Request) {
//		// Serve the index.html file for the "/get" route
//		http.ServeFile(w, r, "templates/main.html")
//	}
