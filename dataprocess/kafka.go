package dataprocess

import (
	"bufio"
	"datastream/database"
	"io"
	"sync"

	"github.com/google/uuid"
)

var wg sync.WaitGroup

func InsertCSVIntoKafka(file io.Reader, topic string) error {
	kafka := database.Connections("kafka")
	kafka.Connect()
	defer kafka.Close()
	kafka.(*database.KafkaConnection).AddTopic(topic)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value := []byte(line)
		kafka.(*database.KafkaConnection).SendMessage(value)
	}
	kafka.(*database.KafkaConnection).SendMessage([]byte("EOF->"))
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func ExtractFromKafka(topic string) {
	outputChan := make(chan database.Contacts, 1000)
	outputChan1 := make(chan string, 1000)
	outputChan2 := make(chan string, 1000)

	kafka := database.Connections("kafka")
	kafka.Connect()
	defer kafka.Close()

	kafka.(*database.KafkaConnection).AddTopic(topic)

	kafka.(*database.KafkaConnection).RetrieveMessage(outputChan)
	for result := range outputChan {
		wg.Add(1)
		uniqueID := uuid.New().String()
		go processData(result, uniqueID, &wg, outputChan1, outputChan2)
	}
	go InsertIntoMysql(outputChan1, outputChan2)
	wg.Wait()
	close(outputChan1)
	close(outputChan2)
}
