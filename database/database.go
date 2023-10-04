package database

import (
	"database/sql"
	"datastream/config"
	"datastream/logs"
	"fmt"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/IBM/sarama"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConnection interface {
	Connect() error
	Close()
}

type MySQLConnection struct {
	Config config.DatabaseConfig
	DB     *sql.DB
}

func (m *MySQLConnection) Connect() error {
	dsn := m.Config.GetDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logs.FileLog.Warning("sql connection issue")
		return err
	}

	m.DB = db
	return nil
}

func (m *MySQLConnection) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

// func (m *MySQLConnection) QueryExec(query string) ([]map[string]interface{}, error) {
// 	// Implement the MySQL query logic here
// 	return nil, nil
// }

func (m *MySQLConnection) InsertData(contacts ContactStatus, activitystring string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		logs.FileLog.Error(fmt.Print(err))
	}

	_, err = tx.Exec("INSERT INTO Users (id, name, email, details, status) VALUES (?, ?, ?, ?, ?)",
		contacts.Id, contacts.Name, contacts.Email, contacts.Details, contacts.Status)
	if err != nil {
		tx.Rollback()

		logs.FileLog.Error(fmt.Print(err))
	}

	sqlStatement := fmt.Sprintf("INSERT INTO Activity (contact_id, campaign_id, activity_type, date) VALUES %s;",
		activitystring)
	_, err = tx.Exec(sqlStatement)
	if err != nil {
		tx.Rollback()

		logs.FileLog.Error(fmt.Print(err))
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()

		logs.FileLog.Error(fmt.Print(err))
	}
	return nil
}

type ClickHouseConnection struct {
	Config config.DatabaseConfig
	DB     *sql.DB
}

func (c *ClickHouseConnection) Connect() error {
	dsn := c.Config.GetDSN()
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return err
	}
	c.DB = db
	return nil
}

func (c *ClickHouseConnection) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

func (c *ClickHouseConnection) QueryExec(query string) ([]map[string]interface{}, error) {
	// Implement the ClickHouse query logic here
	// ...
	return nil, nil
}

func (c *ClickHouseConnection) InsertData(query string) error {
	// Implement the ClickHouse execution logic here
	// ...
	return nil
}

type KafkaConnection struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	Topic    string
	Flag     string
}

func (k *KafkaConnection) Connect() error {
	kafkaProducer := sarama.NewConfig()
	kafkaProducer.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, kafkaProducer)
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
	switch k.Flag {
	case "send":
		k.Producer.Close()
	case "retrieve":
		k.Consumer.Close()
	}
}
func (k *KafkaConnection) AddTopic(topic string) { k.Topic = topic }

func (k *KafkaConnection) SendMessage(value []byte) {
	message := &sarama.ProducerMessage{
		Topic: k.Topic,
		Value: sarama.ByteEncoder(value),
	}
	k.Flag = "send"
	k.Producer.SendMessage(message)
}

func (k *KafkaConnection) RetrieveMessage() {
	inputChan := make(chan Contacts, 10)
	partitionConsumer, err := k.Consumer.ConsumePartition(k.Topic, 0, sarama.OffsetOldest)
	k.Flag = "retrieve"
	if err != nil {
		logs.FileLog.Warning(fmt.Sprintf("Error in consuming : %v", err))
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():

			var data Contacts
			parts := strings.SplitN(string(message.Value), ",", 3)
			if parts[0] == "EOF->" {
				logs.FileLog.Warning("Kafka messages are finished.")
				return
			}
			if len(parts) != 3 {
				logs.FileLog.Warning(fmt.Printf("Invalid message format: %s\n", string(message.Value)))
				continue
			}
			data = Contacts{
				Name:    parts[0],
				Email:   parts[1],
				Details: parts[2],
			}
			go ProcessData(inputChan)
			inputChan <- data
		case err := <-partitionConsumer.Errors():
			logs.FileLog.Warning(fmt.Sprintf("Error consuming messages: %v", err))
		}
	}
}

func Connections(conn string) DatabaseConnection {
	connector := config.LoadConfig(conn)
	switch conn {
	case "mysql":
		return &MySQLConnection{Config: connector}
	case "clickhouse":
		return &ClickHouseConnection{Config: connector}
	case "kafka":
		return &KafkaConnection{}
	default:
		logs.FileLog.Error("Error connecting to database...")
	}
	return nil
}
