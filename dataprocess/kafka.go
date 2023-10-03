package dataprocess

import (
	"bufio"
	"datastream/database"
	"io"
)

func InsertCSVIntoKafka(file io.Reader, topic string) error {
	kafka := database.Connections("kafka")
	kafka.Connect()
	defer kafka.Close()
	kafka.(*database.KafkaConnection).AddTopic("topic")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		value := []byte(line)
		kafka.(*database.KafkaConnection).SendMessage(value)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func ExtractFromKafka(partition int32) {
	kafka := database.Connections("kafka")
	kafka.Connect()
	defer kafka.Close()
	kafka.(*database.KafkaConnection).AddTopic("topic")
	kafka.(*database.KafkaConnection).RetrieveMessage(partition)
}
