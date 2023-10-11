package dataprocess

import (
	"datastream/datastore"
	"datastream/logs"
)

func QueryFromClickhouse() []map[string]interface{} {
	clickhouse := datastore.DatastoreInstance("clickhouse")
	err := clickhouse.Connect()
	defer clickhouse.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	query := "select ContactsID, totalcount from contact_activity_summary_result order by totalcount desc limit 5;"
	queryResult, err := clickhouse.(*datastore.ClickHouseConnection).QueryExec(query)
	if err != nil {
		logs.FileLog.Error("Query doesn't works properly")
	}
	return queryResult

}
