package datastore

import (
	"datastream/logs"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/IBM/sarama"
	_ "github.com/go-sql-driver/mysql"
)

func (k *KafkaConnection) Connect() error {
	kafkaProducer := sarama.NewConfig()
	kafkaProducer.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{k.Config.GetDSN()}, kafkaProducer)
	if err != nil {
		return err
	}
	kafkaConsumer := sarama.NewConfig()
	kafkaConsumer.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, kafkaConsumer)
	if err != nil {
		return err
	}

	k.Producer = producer
	k.Consumer = consumer
	return nil
}
func (k *KafkaConnection) Close() {
	k.Close()
}
func (k *KafkaConnection) setTopic(topic string) { k.Topic = topic }

func (k *KafkaConnection) SendMessage(value []byte, topic string) {
	defer k.Producer.Close()
	k.setTopic(topic)
	message := &sarama.ProducerMessage{
		Topic: k.Topic,
		Value: sarama.ByteEncoder(value),
	}
	k.Producer.SendMessage(message)
}

func (k *KafkaConnection) RetrieveMessage(chanel chan Contacts, topic string) {
	defer k.Consumer.Close()
	k.setTopic(topic)
	partitionConsumer, err := k.Consumer.ConsumePartition(k.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		logs.FileLog.Error("Error in consuming : %v", err)
	}
	defer partitionConsumer.Close()

	inactivityTimer := time.NewTimer(30 * time.Second)
	inactivityTimer.Stop()

	for {
		select {
		case <-inactivityTimer.C:
			return
		case message := <-partitionConsumer.Messages():
			parts := strings.SplitN(string(message.Value), ",", 3)

			if len(parts) != 3 {
				logs.FileLog.Error("Invalid message format: %s\n", string(message.Value))
				continue
			}
			data := Contacts{
				Name:    parts[0],
				Email:   parts[1],
				Details: parts[2],
			}
			chanel <- data
			inactivityTimer.Reset(5 * time.Second)
		case err := <-partitionConsumer.Errors():
			logs.FileLog.Error("Error consuming messages: %v", err)
		}
	}
}
