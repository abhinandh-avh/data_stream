package datastore

import (
	"database/sql"
	"datastream/config"
	"datastream/logs"

	"github.com/IBM/sarama"
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

type ClickHouseConnection struct {
	Config config.DatabaseConfig
	DB     *sql.DB
}

type KafkaConnection struct {
	Config   config.DatabaseConfig
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	Topic    string
}

func DatastoreInstance(dataStore string) DatabaseConnection {
	connector := config.LoadConfig(dataStore)
	switch dataStore {
	case "mysql":
		return &MySQLConnection{Config: connector}
	case "clickhouse":
		return &ClickHouseConnection{Config: connector}
	case "kafka":
		return &KafkaConnection{Config: connector}
	default:
		logs.FileLog.Error("Error connecting to database...  %s", dataStore)
	}
	return nil
}
