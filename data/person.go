package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
)

type Person struct {
	// gorm.Model
	ID             uint         `json:"id" gorm:"primary_key;auto_increment"`
	Name           string       `json:"name"`
	OrganizationID uint         `json:"-"`
	Organization   Organization `json:"organization" gorm:"foreignkey:OrganizationID"`
}

type PersonStore interface {
	GetPersons() ([]*Person, error)
	GetPersonByID(id uint) (*Person, error)
	UpdatePerson(id uint, person *Person) (*Person, error)
	AddPerson(person *Person) (*Person, error)
	DeletePersonByID(id uint) error
}

type PersonDBStore struct {
	*gorm.DB
	log hclog.Logger
}

type PersonNotFoundError struct {
	Cause error
}

func (e PersonNotFoundError) Error() string { return "Person not found! Cause: " + e.Cause.Error() }
func (e PersonNotFoundError) Unwrap() error { return e.Cause }

func NewPersonDBStore(db *gorm.DB, log hclog.Logger) *PersonDBStore {
	return &PersonDBStore{db, log}
}

func (db *PersonDBStore) GetPersons() ([]*Person, error) {
	db.log.Debug("Getting all persons...")

	var persons []*Person
	if err := db.Preload("Organization").Find(&persons).Error; err != nil {
		db.log.Error("Error getting all persons", "err", err)
		return []*Person{}, err
	}

	db.log.Debug("Returning persons", "persons", spew.Sprintf("%+v", persons))
	return persons, nil
}

func (db *PersonDBStore) GetPersonByID(id uint) (*Person, error) {
	db.log.Debug("Getting person by id...", "id", id)

	var person Person
	if err := db.Preload("Organization").First(&person, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Person not found by id", "id", id)
			return nil, &PersonNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting person by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning person", "person", hclog.Fmt("%+v", person))
	return &person, nil
}

func (db *PersonDBStore) UpdatePerson(id uint, person *Person) (*Person, error) {
	db.log.Debug("Updating person...", "person", hclog.Fmt("%+v", person))

	if err := db.Model(&Person{}).Where("id = ?", id).Take(&Person{}).Update(person).First(&person, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Person to be updated not found", "person", hclog.Fmt("%+v", person))
			return nil, &PersonNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating person", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated person", "person", hclog.Fmt("%+v", person))
	return person, nil
}

func (db *PersonDBStore) AddPerson(person *Person) (*Person, error) {
	db.log.Debug("Adding person...", "person", hclog.Fmt("%+v", person))

	if err := db.Create(&person).Error; err != nil {
		db.log.Error("Unexpected error creating person", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added person", "person", hclog.Fmt("%+v", person))
	return person, nil
}

func (db *PersonDBStore) DeletePersonByID(id uint) error {
	db.log.Debug("Deleting person by id...", "id", id)

	if err := db.Model(&Person{}).Where("id = ?", id).Take(&Person{}).Delete(&Person{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Person not found by id", "id", id)
			return &PersonNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting person", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted person")
	return nil
}
