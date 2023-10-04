package database

import (
	"datastream/logs"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Contacts struct {
	Name    string
	Email   string
	Details string
}
type ContactStatus struct {
	Contacts
	Id     string
	Status int
}
type ContactActivity struct {
	Activitydate []string
	Activitytype []int
	Contactid    string
	Campaignid   int
}

type QueryOutput struct {
}

func InsertInMysql(contacts ContactStatus, activitystring string) {
	mysql := Connections("mysql")
	err := mysql.Connect()
	defer mysql.Close()
	if err != nil {
		logs.FileLog.Error("Connection doesn't works properly")
	}
	err = mysql.(*MySQLConnection).InsertData(contacts, activitystring)
	if err != nil {
		logs.FileLog.Error("Batch insertion doesn't works properly")
	}
}

func ProcessData(input <-chan Contacts) {
	for item := range input {
		uniqueID := uuid.New().String()
		activity, flag := GenerateActivity(uniqueID)
		var new ContactStatus
		new.Contacts = item
		new.Id = uniqueID
		new.Status = flag
		values := formatActivity(activity)
		// logs.FileLog.Info(fmt.Sprintf("...loading : %s , %s", new.Email, values))
		InsertInMysql(new, values)
	}
}

func formatActivity(activity []ContactActivity) string {
	var values string
	for _, ins := range activity {
		for i := 0; i < len(ins.Activitytype); i++ {
			values += fmt.Sprintf("('%s',%d,%d,'%s'),", ins.Contactid, ins.Campaignid, ins.Activitytype[i], ins.Activitydate[i])
		}
	}
	activityString := strings.TrimRight(values, ",")
	return activityString
}

func GenerateActivity(id string) ([]ContactActivity, int) {
	var activity []ContactActivity
	var flag int
	startDate := "2023-01-01"
	for count := 1; count <= 100; count++ {
		var ins ContactActivity
		var types []int

		if count%15 == 0 {
			newDate := addMonth(startDate)
			startDate = newDate
		}

		ins.Contactid = id
		ins.Campaignid = count
		types, flag = activitype()
		ins.Activitytype = types
		dates := generateDateSlice(startDate, len(types))
		ins.Activitydate = dates
		if flag == 0 {
			activity = append(activity, ins)
			return activity, flag
		}
		activity = append(activity, ins)
	}
	return activity, flag
}

func activitype() ([]int, int) {
	var numbers []int
	flag := 1
	min := 1
	max := 1000
	randomNumber := rand.Intn(max-min+1) + min
	switch {
	case randomNumber < 900:
		if randomNumber < 50 {
			numbers = append(numbers, 1)
		} else if randomNumber < 300 {
			numbers = append(numbers, 1, 3)
		} else if randomNumber < 600 {
			numbers = append(numbers, 1, 3, 4)
		} else {
			numbers = append(numbers, 1, 3, 4, 7)
		}
	case randomNumber < 960:
		if randomNumber < 910 {
			numbers = append(numbers, 1, 3, 3)
		} else if randomNumber < 920 {
			numbers = append(numbers, 1, 3, 4, 3)
		} else if randomNumber < 930 {
			numbers = append(numbers, 1, 3, 4, 3, 4)
		} else if randomNumber < 940 {
			numbers = append(numbers, 1, 3, 4, 7, 3)
		} else if randomNumber < 650 {
			numbers = append(numbers, 1, 3, 4, 7, 3, 4)
		} else {
			numbers = append(numbers, 1, 3, 4, 7, 3, 4, 7)
		}
	case randomNumber < 1000:
		flag = 0
		if randomNumber < 970 {
			numbers = append(numbers, 2)
		} else if randomNumber < 975 {
			numbers = append(numbers, 1, 3, 4, 5)
		} else if randomNumber < 980 {
			numbers = append(numbers, 1, 3, 4, 6)
		} else if randomNumber < 990 {
			numbers = append(numbers, 1, 3, 5)
		} else {
			numbers = append(numbers, 1, 3, 6)
		}
		break
	}
	return numbers, flag
}

func generateDateSlice(startDate string, length int) []string {
	parsedStartDate, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		logs.FileLog.Error("Error parsing date")
	}
	dateSlice := make([]string, length)
	min := 0
	max := 4
	randomNumber := rand.Intn(max-min+1) + min
	for i := 0; i < length; i++ {
		dateStr := parsedStartDate.Format("2006-01-02")
		dateSlice[i] = dateStr
		parsedStartDate = parsedStartDate.AddDate(0, 0, randomNumber)
	}

	return dateSlice
}
func addMonth(startDate string) string {
	parsedStartDate, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		logs.FileLog.Error("Error parsing date")
	}
	newDate := parsedStartDate.AddDate(0, 1, 0).Format("2006-01-02")
	return newDate
}
