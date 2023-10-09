package datastore

import (
	"database/sql"
	"datastream/logs"
	"fmt"
	"sync"
)

func (m *MySQLConnection) Connect() error {
	dsn := m.Config.GetDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logs.FileLog.Error("sql connection issue : %v", err)
		return err
	}

	m.DB = db
	return nil
}

func (m *MySQLConnection) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *MySQLConnection) QueryExec(query string) ([]map[string]interface{}, error) {
	return nil, nil
}

func (m *MySQLConnection) InsertData(contactChannelTOSQL chan string, activityChannelTOSQL chan string) {
	var wg sync.WaitGroup
	batchSize := 10
	counter1 := 0
	counter2 := 0
	tx, err := m.DB.Begin()
	if err != nil {
		logs.FileLog.Error("DB Begin : %v", err)
	}

	go func() {
		defer wg.Done()
		wg.Add(1)
		for contact := range contactChannelTOSQL {
			counter1++
			sqlStatement := fmt.Sprintf("INSERT INTO Contacts (ID, Name, Email, Details, Status) VALUES %s;", contact)
			_, err = tx.Exec(sqlStatement)
			if err != nil {
				tx.Rollback()
				logs.FileLog.Error("%v", err)
			}
			if counter1%batchSize == 0 {
				logs.FileLog.Info("Batch Inserted...")
				if err := tx.Commit(); err != nil {
					tx.Rollback()
					logs.FileLog.Error("%v", err)
				}
				tx, err = m.DB.Begin()
				if err != nil {
					logs.FileLog.Error("%v", err)
				}
			}
		}
		if err := tx.Commit(); err != nil {
			logs.FileLog.Error("%v", err)
		}
	}()
	go func() {
		defer wg.Done()
		wg.Add(1)
		for activity := range activityChannelTOSQL {
			counter2++
			sqlStatements :=
				fmt.Sprintf("INSERT INTO ContactActivity (ContactsID, CampaignID, ActivityType, ActivityDate) VALUES %s;",
					activity)
			_, err = tx.Exec(sqlStatements)
			if err != nil {
				tx.Rollback()

				logs.FileLog.Error("%v", err)
			}
			if counter2%batchSize == 0 {
				logs.FileLog.Info("Batch Inserted...")
				if err := tx.Commit(); err != nil {
					tx.Rollback()
					logs.FileLog.Error("%v", err)
				}
				tx, err = m.DB.Begin()
				if err != nil {
					logs.FileLog.Error("%v", err)
				}
			}
		}
		if err := tx.Commit(); err != nil {
			logs.FileLog.Error("%v", err)
		}
	}()
	wg.Wait()
	logs.FileLog.Info("SQL INSERTION COMPLETED...")
}
