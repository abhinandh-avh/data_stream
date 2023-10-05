package dataprocess

import (
	"datastream/database"
	"datastream/logs"
)

func InsertIntoMysql() {
	mysql := database.Connections("mysql")
	err := mysql.Connect()
	defer mysql.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	mysql.(*database.MySQLConnection).InsertData(ContactsSlice, ActivitySlice)
	if err != nil {
		logs.FileLog.Error("Batch insertion doesn't works properly")
	} else {
		logs.FileLog.Info("SQL insertion completed")
	}
}
