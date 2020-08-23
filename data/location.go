package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
)

type Location struct {
	// gorm.Model
	ID   uint   `json:"id" gorm:"primary_key;auto_increment"`
	Name string `json:"name" validate:"required" gorm:"not null"`
}

type LocationStore interface {
	GetLocations() ([]*Location, error)
	GetLocationByID(id uint) (*Location, error)
	UpdateLocation(id uint, loc *Location) (*Location, error)
	AddLocation(loc *Location) (*Location, error)
	DeleteLocationByID(id uint) error
}

type LocationDBStore struct {
	*gorm.DB
	validate *validator.Validate
	log hclog.Logger
}

type LocationNotFoundError struct {
	Cause error
}

func (e LocationNotFoundError) Error() string { return "Location not found! Cause: " + e.Cause.Error() }
func (e LocationNotFoundError) Unwrap() error { return e.Cause }

func NewLocationDBStore(db *gorm.DB, log hclog.Logger) *LocationDBStore {
	return &LocationDBStore{db, validator.New(), log}
}

func (db *LocationDBStore) GetLocations() ([]*Location, error) {
	db.log.Debug("Getting all locations...")

	var locations []*Location
	if err := db.Find(&locations).Error; err != nil {
		db.log.Error("Error getting all locations", "err", err)
		return []*Location{}, err
	}

	db.log.Debug("Returning locations", "locations", spew.Sprintf("%+v", locations))
	return locations, nil
}

func (db *LocationDBStore) GetLocationByID(id uint) (*Location, error) {
	db.log.Debug("Getting location by id...", "id", id)

	var location Location
	if err := db.First(&location, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Location not found by id", "id", id)
			return nil, &LocationNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting location by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning location", "location", hclog.Fmt("%+v", location))
	return &location, nil
}

func (db *LocationDBStore) UpdateLocation(id uint, location *Location) (*Location, error) {
	db.log.Debug("Updating location...", "location", hclog.Fmt("%+v", location))

	err := db.validate.Struct(location)
	if err != nil {
		db.log.Error("Error validating location", "err", err)
		return nil, err
	}

	if err := db.Model(&Location{}).Where("id = ?", id).Take(&Location{}).Update(location).First(&location, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Location to be updated not found", "location", hclog.Fmt("%+v", location))
			return nil, &LocationNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating location", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated location", "location", hclog.Fmt("%+v", location))
	return location, nil
}

func (db *LocationDBStore) AddLocation(location *Location) (*Location, error) {
	db.log.Debug("Adding location...", "location", hclog.Fmt("%+v", location))

	err := db.validate.Struct(location)
	if err != nil {
		db.log.Error("Error validating location", "err", err)
		return nil, err
	}

	if err := db.Create(&location).Error; err != nil {
		db.log.Error("Unexpected error creating location", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added location", "location", hclog.Fmt("%+v", location))
	return location, nil
}

func (db *LocationDBStore) DeleteLocationByID(id uint) error {
	db.log.Debug("Deleting location by id...", "id", id)

	if err := db.Model(&Location{}).Where("id = ?", id).Take(&Location{}).Delete(&Location{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Location not found by id", "id", id)
			return &LocationNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting location", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted location")
	return nil
}
