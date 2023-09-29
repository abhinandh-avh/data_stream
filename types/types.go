package types

type Contacts struct {
	id      int
	name    string
	email   string
	details string
}

type ContactActivity struct {
	id             int
	contactid      int
	campaignid     int
	activitytytype int
	activitydate   string
}

type ContactStatus struct {
	contact Contacts
	status  int
}
type QueryOutput struct {
}
