package config

import (
	"datastream/logs"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig interface {
	GetDSN() string
}

type MySQLConfig struct {
	Host     string
	Username string
	Password string
	Database string
	Port     string
}

func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

type ClickHouseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (c *ClickHouseConfig) GetDSN() string {
	return fmt.Sprintf(
		"tcp://%s:%s?username=%s&password=%s&database=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
	)
}

type KafkaConfig struct {
	Host string
	Port string
}

func (c *KafkaConfig) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s",
		c.Host,
		c.Port,
	)
}

// LoadConfig loads configuration values from the environment file.
func LoadConfig(con string) DatabaseConfig {
	if err := godotenv.Load(".env"); err != nil {
		logs.FileLog.Error("Occured in reading .env file : %v", err)
	}
	switch con {
	case "mysql":
		mysqlConfig := MySQLConfig{
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     os.Getenv("MYSQL_PORT"),
			Username: os.Getenv("MYSQL_USERNAME"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Database: os.Getenv("MYSQL_DATABASE"),
		}
		return &mysqlConfig
	case "clickhouse":
		clickHouseConfig := ClickHouseConfig{
			Host:     os.Getenv("CLICKHOUSE_HOST"),
			Port:     os.Getenv("CLICKHOUSE_PORT"),
			Username: os.Getenv("CLICKHOUSE_USERNAME"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
			Database: os.Getenv("CLICKHOUSE_DATABASE"),
		}

		return &clickHouseConfig
	case "kafka":
		KafkaConfig := KafkaConfig{
			Host: os.Getenv("CLICKHOUSE_HOST"),
			Port: os.Getenv("KAFKA_PORT"),
		}

		return &KafkaConfig
	default:
		return nil
	}
}
