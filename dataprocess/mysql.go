package dataprocess

import (
	"datastream/datastore"
	"datastream/logs"
)

func InsertIntoMysql(chan1 chan string, chan2 chan string) {
	mysql := datastore.DatastoreInstance("mysql")
	err := mysql.Connect()
	defer mysql.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	mysql.(*datastore.MySQLConnection).InsertData(chan1, chan2)
}
