package types

type Contacts struct {
	id    int
	name  string
	email string
}

type ContactActivity struct {
	id             int
	contactid      int
	campaignid     int
	activitytytype int
	activitydate   string
}
type ContactDetails struct {
	contactid int
	dob       string
	city      string
	state     string
}
type InsertContact struct {
	contact Contacts
	status  int
}
type QuryOutput struct {
	// ......
}
