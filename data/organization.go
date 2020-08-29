package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
)

type Organization struct {
	// gorm.Model
	ID      uint     `json:"id" gorm:"primary_key;auto_increment"`
	Name    string   `json:"name" gorm:"unique;not null;default:null"`
}

type OrganizationStore interface {
	GetOrganizations() ([]*Organization, error)
	GetOrganizationByID(id uint) (*Organization, error)
	UpdateOrganization(id uint, organization *Organization) (*Organization, error)
	AddOrganization(organization *Organization) (*Organization, error)
	DeleteOrganizationByID(id uint) error
}

type OrganizationDBStore struct {
	*gorm.DB
	validate *validator.Validate
	log      hclog.Logger
}

type OrganizationNotFoundError struct {
	Cause error
}

func (e OrganizationNotFoundError) Error() string {
	return "Organization not found! Cause: " + e.Cause.Error()
}
func (e OrganizationNotFoundError) Unwrap() error { return e.Cause }

func NewOrganizationDBStore(db *gorm.DB, log hclog.Logger) *OrganizationDBStore {
	return &OrganizationDBStore{db, validator.New(), log}
}

func (db *OrganizationDBStore) GetOrganizations() ([]*Organization, error) {
	db.log.Debug("Getting all organizations...")

	var organizations []*Organization
	if err := db.Find(&organizations).Error; err != nil {
		db.log.Error("Error getting all organizations", "err", err)
		return []*Organization{}, err
	}

	db.log.Debug("Returning organizations", "organizations", spew.Sprintf("%+v", organizations))
	return organizations, nil
}

func (db *OrganizationDBStore) GetOrganizationByID(id uint) (*Organization, error) {
	db.log.Debug("Getting organization by id...", "id", id)

	var organization Organization
	if err := db.First(&organization, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Organization not found by id", "id", id)
			return nil, &OrganizationNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting organization by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning organization", "organization", hclog.Fmt("%+v", organization))
	return &organization, nil
}

func (db *OrganizationDBStore) UpdateOrganization(id uint, organization *Organization) (*Organization, error) {
	db.log.Debug("Updating organization...", "organization", hclog.Fmt("%+v", organization))

	err := db.validate.Struct(organization)
	if err != nil {
		db.log.Error("Error validating organization", "err", err)
		return nil, err
	}

	if err := db.Model(&Organization{}).Where("id = ?", id).Take(&Organization{}).Update(organization).First(&organization, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Organization to be updated not found", "organization", hclog.Fmt("%+v", organization))
			return nil, &OrganizationNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating organization", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated organization", "organization", hclog.Fmt("%+v", organization))
	return organization, nil
}

func (db *OrganizationDBStore) AddOrganization(organization *Organization) (*Organization, error) {
	db.log.Debug("Adding organization...", "organization", hclog.Fmt("%+v", organization))

	err := db.validate.Struct(organization)
	if err != nil {
		db.log.Error("Error validating organization", "err", err)
		return nil, err
	}

	if err := db.Create(&organization).Error; err != nil {
		db.log.Error("Unexpected error creating organization", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added organization", "organization", hclog.Fmt("%+v", organization))
	return organization, nil
}

func (db *OrganizationDBStore) DeleteOrganizationByID(id uint) error {
	db.log.Debug("Deleting organization by id...", "id", id)

	if err := db.Model(&Organization{}).Where("id = ?", id).Take(&Organization{}).Delete(&Organization{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Organization not found by id", "id", id)
			return &OrganizationNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting organization", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted organization")
	return nil
}
