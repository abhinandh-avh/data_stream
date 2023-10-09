package dataprocess

import (
	"datastream/datastore"
	"datastream/logs"
	"io"
	"os"
	"sync"

	"github.com/google/uuid"
)

func InsertCSVIntoKafka(fileName string, topic string) {
	kafkaInstance := basicKafkaConnection(topic)
	fileToSend, err := os.Open(fileName)
	if err != nil {
		logs.FileLog.Error("Opening fileToSend: %v", err)
		return
	}
	defer fileToSend.Close()

	buffer := make([]byte, 1024*1024)
	for {
		numeberOflines, err := fileToSend.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logs.FileLog.Error("Reading file: %v", err)
			return
		}
		kafkaInstance.(*datastore.KafkaConnection).SendMessage(buffer[:numeberOflines])
	}

	go ExtractFromKafka(topic)
}

func ExtractFromKafka(topic string) {
	var wg sync.WaitGroup
	dataFromKafkaConsumer := make(chan datastore.Contacts)
	contactChannelTOSQL := make(chan string)
	activityChannelTOSQL := make(chan string)

	kafkaInstance := basicKafkaConnection(topic)
	go kafkaInstance.(*datastore.KafkaConnection).RetrieveMessage(dataFromKafkaConsumer, topic)
	go func() {
		for result := range dataFromKafkaConsumer {
			wg.Add(1)
			uniqueID := uuid.New().String()
			go processData(result, uniqueID, &wg, contactChannelTOSQL, activityChannelTOSQL)
		}
	}()
	go InsertIntoMysql(contactChannelTOSQL, activityChannelTOSQL)
	wg.Wait()
	close(contactChannelTOSQL)
	close(activityChannelTOSQL)
	kafkaInstance.Close()
}
func basicKafkaConnection(topic string) datastore.DatabaseConnection {
	kafka := datastore.DatastoreInstance("kafka")
	kafka.Connect()
	kafka.(*datastore.KafkaConnection).SetTopic(topic)
	return kafka
}
