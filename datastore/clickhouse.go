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

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to store the query results
	var queryResult []map[string]interface{}

	// Iterate through the rows and populate the queryResult slice
	for rows.Next() {
		var column1 string
		var column2 int
		if err := rows.Scan(&column1, &column2); err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"Column1": column1,
			"Column2": column2,
		}

		queryResult = append(queryResult, result)
	}

	return queryResult, nil
}

func (c *ClickHouseConnection) InsertData(query string) error {

	return nil
}
