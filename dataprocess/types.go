package dataprocess

import (
	"datastream/database"
	"datastream/logs"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ContactStatus struct {
	database.Contacts
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

func processData(item database.Contacts, id string, wg *sync.WaitGroup, chan1 chan string, chan2 chan string) {
	defer wg.Done()
	activity, flag := GenerateActivity(id)
	var new ContactStatus
	new.Contacts = item
	new.Id = id
	new.Status = flag
	values := formatActivity(activity)
	contactString := formatContact(new)
	chan1 <- contactString
	chan2 <- values
}
func formatContact(data ContactStatus) string {
	var values string
	values = fmt.Sprintf("('%s','%s','%s','%s',%d)", data.Id, data.Name, data.Email, data.Details, data.Status)

	return values
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
	max := 10000
	randomNumber := rand.Intn(max-min+1) + min
	switch {
	case randomNumber < 9000:
		if randomNumber < 200 {
			numbers = append(numbers, 1)
		} else if randomNumber < 3000 {
			numbers = append(numbers, 1, 3)
		} else if randomNumber < 6000 {
			numbers = append(numbers, 1, 3, 4)
		} else {
			numbers = append(numbers, 1, 3, 4, 7)
		}
	case randomNumber < 9900:
		if randomNumber < 9150 {
			numbers = append(numbers, 1, 3, 3)
		} else if randomNumber < 9300 {
			numbers = append(numbers, 1, 3, 4, 3)
		} else if randomNumber < 9450 {
			numbers = append(numbers, 1, 3, 4, 3, 4)
		} else if randomNumber < 9600 {
			numbers = append(numbers, 1, 3, 4, 7, 3)
		} else if randomNumber < 9750 {
			numbers = append(numbers, 1, 3, 4, 7, 3, 4)
		} else {
			numbers = append(numbers, 1, 3, 4, 7, 3, 4, 7)
		}
	case randomNumber <= 10000:
		flag = 0
		if randomNumber < 9920 {
			numbers = append(numbers, 2)
		} else if randomNumber < 9940 {
			numbers = append(numbers, 1, 3, 4, 5)
		} else if randomNumber < 9960 {
			numbers = append(numbers, 1, 3, 4, 6)
		} else if randomNumber < 9980 {
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
