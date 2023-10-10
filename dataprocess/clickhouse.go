package dataprocess

import (
	"datastream/datastore"
	"datastream/logs"
)

func QueryFromClickhouse() {
	clickhouse := datastore.DatastoreInstance("clickhouse")
	err := clickhouse.Connect()
	defer clickhouse.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	clickhouse.(*datastore.MySQLConnection).QueryExec("shdgf")
}
