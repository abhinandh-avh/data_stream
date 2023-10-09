package dataprocess

import (
	"datastream/datastore"
	"sync"

	"github.com/google/uuid"
)

func InsertCSVIntoKafka(file []byte, topic string) {
	kafkaInstance := basicKafkaConnection()
	kafkaInstance.(*datastore.KafkaConnection).SendMessage(file, topic)
}

func ExtractFromKafka(topic string) {
	var wg sync.WaitGroup
	dataFromKafkaConsumer := make(chan datastore.Contacts)
	contactChannelTOSQL := make(chan string)
	activityChannelTOSQL := make(chan string)

	kafkaInstance := basicKafkaConnection()
	go kafkaInstance.(*datastore.KafkaConnection).RetrieveMessage(dataFromKafkaConsumer, topic)
	for result := range dataFromKafkaConsumer {
		wg.Add(1)
		uniqueID := uuid.New().String()
		go processData(result, uniqueID, &wg, contactChannelTOSQL, activityChannelTOSQL)
	}
	go InsertIntoMysql(contactChannelTOSQL, activityChannelTOSQL)
	wg.Wait()
	close(contactChannelTOSQL)
	close(activityChannelTOSQL)
}
func basicKafkaConnection() datastore.DatabaseConnection {
	kafka := datastore.DatastoreInstance("kafka")
	kafka.Connect()
	return kafka
}
