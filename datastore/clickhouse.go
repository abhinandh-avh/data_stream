package datastore

import "database/sql"

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
