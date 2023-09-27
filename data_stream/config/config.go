package config

import "database/sql"

type DBConnector interface {
	Connect() (*sql.DB, error)
	Close() error
}

// MySQLConfig holds the MySQL connection configuration.
type MySQLConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	DBName   string
}

// MySQLConnector is a struct that implements the DBConnector interface for MySQL.
type MySQLConnector struct {
	config MySQLConfig
}

// KafkaConfig holds the Kafka connection configuration.
type KafkaConfig struct {
	Broker string
}

// KafkaConnector is a struct that implements the DBConnector interface for Kafka.
type KafkaConnector struct {
	config KafkaConfig
}

// ClickHouseConfig holds the ClickHouse connection configuration.
type ClickHouseConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	DBName   string
}

// ClickHouseConnector is a struct that implements the DBConnector interface for ClickHouse.
type ClickHouseConnector struct {
	config ClickHouseConfig
}
