package dataprocess

import (
	"bufio"
	"datastream/database"
	"datastream/logs"
	"io"
	"sync"

	"github.com/google/uuid"
)

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
	count := 0
	var wg sync.WaitGroup
	outputChan := make(chan database.Contacts, 1000)

	kafka := database.Connections("kafka")
	kafka.Connect()
	defer kafka.Close()

	kafka.(*database.KafkaConnection).AddTopic(topic)

	go kafka.(*database.KafkaConnection).RetrieveMessage(outputChan)
	for result := range outputChan {
		count++
		wg.Add(1)
		uniqueID := uuid.New().String()
		go processData(result, uniqueID, &wg)
		if count%1000 == 0 {
			wg.Wait()
		}
	}
	logs.FileLog.Info("Wating...")
	wg.Wait()
	InsertIntoMysql()
}
