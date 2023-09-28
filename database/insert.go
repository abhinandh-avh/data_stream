package database

import "datastream/types"

type DataStore interface {
	Add(contact types.Contacts) error
	GenerateActivity(contact types.ContactActivity) (types.ContactActivity, error)
}
