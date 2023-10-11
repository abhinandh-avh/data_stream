package dataprocess

import (
	"datastream/datastore"
	"datastream/logs"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

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

	csvReader := csv.NewReader(fileToSend)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logs.FileLog.Error("Reading CSV: %v", err)
			return
		}

		// Convert the CSV line to a string (you can customize this as needed)
		lineStr := strings.Join(line, ",")
		// Send the line to Kafka
		kafkaInstance.(*datastore.KafkaConnection).SendMessage([]byte(lineStr))
	}

	go ExtractFromKafka(topic)
}

func ExtractFromKafka(topic string) {
	channelClossigngTrigger := make(chan bool)
	dataFromKafkaConsumer := make(chan datastore.Contacts)
	contactChannelTOSQL := make(chan string)
	activityChannelTOSQL := make(chan string)
	kafkaInstance := basicKafkaConnection(topic)
	go func() {
		kafkaInstance.(*datastore.KafkaConnection).RetrieveMessage(dataFromKafkaConsumer, topic)
		channelClossigngTrigger <- true
	}()
	go func() {
		for result := range dataFromKafkaConsumer {
			uniqueID := uuid.New().String()
			go processData(result, uniqueID, contactChannelTOSQL, activityChannelTOSQL)
		}
	}()

	go InsertIntoMysql(contactChannelTOSQL, activityChannelTOSQL)
	go func() {
		bools := <-channelClossigngTrigger
		time.Sleep(15 * time.Second)
		if bools {
			logs.FileLog.Info("Channels are closed")
		}
		close(contactChannelTOSQL)
		close(activityChannelTOSQL)
	}()
}
func basicKafkaConnection(topic string) datastore.DatabaseConnection {
	kafka := datastore.DatastoreInstance("kafka")
	kafka.Connect()
	kafka.(*datastore.KafkaConnection).SetTopic(topic)
	return kafka
}
