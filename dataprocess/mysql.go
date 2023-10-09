package dataprocess

import (
	"datastream/datastore"
	"datastream/logs"
)

func InsertIntoMysql(contactChannelTOSQL chan string, activityChannelTOSQL chan string) {
	mysql := datastore.DatastoreInstance("mysql")
	err := mysql.Connect()
	defer mysql.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	mysql.(*datastore.MySQLConnection).InsertData(contactChannelTOSQL, activityChannelTOSQL)
}
