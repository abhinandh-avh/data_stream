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

type Contacts struct {
	Name    string
	Email   string
	Details string
}
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
// 	return nil, nil
// }

func (m *MySQLConnection) InsertData(contact []string, activitystring []string) {
	batchSize := 1000
	tx, err := m.DB.Begin()
	if err != nil {
		logs.FileLog.Error(fmt.Print(err))
	}
	for i := 0; i < len(contact); i++ {
		sqlStatement := fmt.Sprintf("INSERT INTO Contacts (ID, Name, Email, Details, Status) VALUES %s;", contact[i])
		_, err = tx.Exec(sqlStatement)
		if err != nil {
			tx.Rollback()

			logs.FileLog.Error(fmt.Print(err))
		}

		sqlStatements := fmt.Sprintf("INSERT INTO ContactActivity (ContactsID, CampaignID, ActivityType, ActivityDate) VALUES %s;",
			activitystring[i])
		_, err = tx.Exec(sqlStatements)
		if err != nil {
			tx.Rollback()

			logs.FileLog.Error(fmt.Print(err))
		}
		if i%batchSize == 0 {
			logs.FileLog.Info("Batch Inserted...")
			if err := tx.Commit(); err != nil {
				tx.Rollback()

				logs.FileLog.Error(fmt.Print(err))
			}
			tx, err = m.DB.Begin()
			if err != nil {
				logs.FileLog.Error(fmt.Print(err))
			}

		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()

		logs.FileLog.Error("EOF ->")
	}
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

	return nil, nil
}

func (c *ClickHouseConnection) InsertData(query string) error {

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

func (k *KafkaConnection) RetrieveMessage(chanel chan Contacts) {
	partitionConsumer, err := k.Consumer.ConsumePartition(k.Topic, 0, sarama.OffsetOldest)
	k.Flag = "retrieve"
	if err != nil {
		logs.FileLog.Warning(fmt.Sprintf("Error in consuming : %v", err))
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			parts := strings.SplitN(string(message.Value), ",", 3)
			if parts[0] == "EOF->" {
				logs.FileLog.Warning("Kafka messages are finished.")
				close(chanel)
				return
			}
			if len(parts) != 3 {
				logs.FileLog.Warning(fmt.Printf("Invalid message format: %s\n", string(message.Value)))
				continue
			}
			data := Contacts{
				Name:    parts[0],
				Email:   parts[1],
				Details: parts[2],
			}
			chanel <- data
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
