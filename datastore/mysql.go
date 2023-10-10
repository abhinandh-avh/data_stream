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
	trigerContactCommit := make(chan bool)
	trigerActivityCommit := make(chan bool)
	var activitytriger bool
	var contacttriger bool
	returnFromCommit := make(chan bool)
	var wg sync.WaitGroup
	var waitg sync.WaitGroup
	batchSize := 1000
	counter1 := 0
	counter2 := 0
	tx, err := m.DB.Begin()
	if err != nil {
		logs.FileLog.Error("DB Begin : %v", err)
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for contact := range contactChannelTOSQL {
			counter1++
			sqlStatement := fmt.Sprintf("INSERT INTO Contacts (ID, Name, Email, Details, Status) VALUES %s;", contact)
			_, err = tx.Exec(sqlStatement)
			if err != nil {
				tx.Rollback()
				logs.FileLog.Error("Excecuting : %v", err)
			}
			if counter1%batchSize == 0 {
				trigerContactCommit <- true
				status := <-returnFromCommit
				if status {
					status = false
				}
			}
		}
		if err := tx.Commit(); err != nil {
			logs.FileLog.Error("Finishing...")
		}
	}()
	go func() {
		defer wg.Done()
		for activity := range activityChannelTOSQL {
			counter2++
			sqlStatements :=
				fmt.Sprintf("INSERT INTO ContactActivity (ContactsID, CampaignID, ActivityType, ActivityDate) VALUES %s;",
					activity)
			_, err = tx.Exec(sqlStatements)
			if err != nil {
				tx.Rollback()
				logs.FileLog.Error("Excecuting : %v", err)
			}
			if counter2%batchSize == 0 {
				trigerActivityCommit <- true
				status := <-returnFromCommit
				if status {
					status = false
				}
			}
		}
		if err := tx.Commit(); err != nil {
			logs.FileLog.Error("Finishing...")
		}
	}()

	go func() {
		for {
			logs.FileLog.Info("____________")
			waitg.Add(2)
			go func() {
				defer waitg.Done()
				activitytriger = <-trigerActivityCommit
			}()
			go func() {
				defer waitg.Done()
				contacttriger = <-trigerContactCommit
			}()
			waitg.Wait()
			if contacttriger == true && activitytriger == true {
				if err := tx.Commit(); err == nil {
					logs.FileLog.Info("Batch Inserted...")
				} else {
					tx.Rollback()
					logs.FileLog.Error("%v", err)
				}
				tx, err = m.DB.Begin()
				if err != nil {
					logs.FileLog.Error("%v", err)
				}
			}
			contacttriger = false
			activitytriger = false
			returnFromCommit <- true
			returnFromCommit <- true
		}
	}()
	wg.Wait()
	logs.FileLog.Info("SQL INSERTION COMPLETED...")
}
