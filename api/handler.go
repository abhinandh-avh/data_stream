package api

import (
	"datastream/dataprocess"
	"datastream/logs"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"net/http"
)

var (
	UniqueFileName string
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file for the "/" route
	http.ServeFile(w, r, "templates/HomePage.html")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the about.html file for the "/about" route
	http.ServeFile(w, r, "templates/ResultPage.html")
}

// func GetFromClickHouseHandler(w http.ResponseWriter, r *http.Request) {
// 	// Serve the index.html file for the "/get" route
// 	http.ServeFile(w, r, "templates/main.html")
// }

func InsertIntoKafkaHandler(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().UnixNano()
	uploadDir := "./uploads/"
	file, header, err := r.FormFile("file")
	UniqueFileName = fmt.Sprintf("%s%d", header.Filename, timestamp)
	fileName := filepath.Join(uploadDir, UniqueFileName)
	outFile, err := os.Create(fileName)
	if err != nil {
		logs.FileLog.Error("Error in filename creation %v", err)
		return
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error copying file to server", http.StatusInternalServerError)
		logs.FileLog.Error("Error copying file: %v", err)
		return
	}

	topic := UniqueFileName

	go dataprocess.InsertCSVIntoKafka(fileName, topic)

	logs.FileLog.Error("CSV data inserted into Kafka successfully!  FILENAME :: %s", fileName)

	http.ServeFile(w, r, "templates/HomePage.html")
	return
}
