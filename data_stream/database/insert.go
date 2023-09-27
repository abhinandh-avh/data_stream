package database

import "datastream/types"

type DataStore interface {
	InsertContact(contact types.InsertContacts) error
	GenerateActivity(contact types.ContactActivity) (types.ContactActivity, error)
}
