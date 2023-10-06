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

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	topic := "new35"
	// Pass the file to a function in the data processing package to insert into Kafka
	err = dataprocess.InsertCSVIntoKafka(file, topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		logs.FileLog.Info("CSV data inserted into Kafka successfully!")
		http.ServeFile(w, r, "templates/ResultPage.html")
		dataprocess.ExtractFromKafka(topic)
	}

}

func GetFromClickHouseHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file for the "/get" route
	http.ServeFile(w, r, "templates/main.html")
}
