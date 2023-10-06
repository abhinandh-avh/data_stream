package dataprocess

import (
	"datastream/database"
	"datastream/logs"
)

func InsertIntoMysql(chan1 chan string, chan2 chan string) {
	mysql := database.Connections("mysql")
	err := mysql.Connect()
	defer mysql.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	mysql.(*database.MySQLConnection).InsertData(chan1, chan2)
	if err != nil {
		logs.FileLog.Error("Batch insertion doesn't works properly")
	} else {
		logs.FileLog.Info("SQL insertion completed")
	}
}
