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
	Port     int
}

func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

type ClickHouseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (c *ClickHouseConfig) GetDSN() string {
	return fmt.Sprintf(
		"tcp://%s:%d?username=%s&password=%s&database=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
	)
}

// LoadConfig loads configuration values from the environment file.
func LoadConfig(con string) DatabaseConfig {
	if err := godotenv.Load(".env"); err != nil {
		logs.FileLog.Error(fmt.Sprintf("Occured in reading .env file : %v ", err))
	}
	switch con {
	case "mysql":
		mysqlConfig := MySQLConfig{
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     3306,
			Username: os.Getenv("MYSQL_USERNAME"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Database: os.Getenv("MYSQL_DATABASE"),
		}
		return &mysqlConfig
	case "clickhouse":
		clickHouseConfig := ClickHouseConfig{
			Host:     os.Getenv("CLICKHOUSE_HOST"),
			Port:     9000,
			Username: os.Getenv("CLICKHOUSE_USERNAME"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
			Database: os.Getenv("CLICKHOUSE_DATABASE"),
		}

		return &clickHouseConfig
	default:
		return nil
	}
}
