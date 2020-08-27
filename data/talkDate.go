package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"time"
)

type TalkDate struct {
	// gorm.Model
	ID         uint      `json:"id" gorm:"primary_key;auto_increment"`
	BeginDate  time.Time `json:"beginDate" validate:"required" gorm:"not null"`
	TalkID     uint      `json:"-"`
	Talk       *Talk     `json:"talk" validate:"required"`
	RoomID     uint      `json:"-"`
	Room       *Room     `json:"room" validate:"required"`
	EventID    uint      `json:"-"`
	Event      *Event    `json:"event" validate:"required"`
	LocationID uint      `json:"-"`
	Location   *Location `json:"location" validate:"required"`
}

type TalkDateStore interface {
	GetTalkDates() ([]*TalkDate, error)
	GetTalkDateByID(id uint) (*TalkDate, error)
	UpdateTalkDate(id uint, talkDate *TalkDate) (*TalkDate, error)
	AddTalkDate(talkDate *TalkDate) (*TalkDate, error)
	DeleteTalkDateByID(id uint) error
	GetTalkDatesByEventID(eventID uint) ([]*TalkDate, error)
}

type TalkDateDBStore struct {
	*gorm.DB
	validate *validator.Validate
	log      hclog.Logger
}

type TalkDateNotFoundError struct {
	Cause error
}

func (e TalkDateNotFoundError) Error() string { return "TalkDate not found! Cause: " + e.Cause.Error() }
func (e TalkDateNotFoundError) Unwrap() error { return e.Cause }

func NewTalkDateDBStore(db *gorm.DB, log hclog.Logger) *TalkDateDBStore {
	return &TalkDateDBStore{db, validator.New(), log}
}

func (db *TalkDateDBStore) GetTalkDates() ([]*TalkDate, error) {
	db.log.Debug("Getting all talkDates...")

	var talkDates []*TalkDate
	if err := db.
		Preload("Talk").
		Preload("Talk.Persons").
		Preload("Talk.Topics").
		Preload("Talk.Topics.Children").
		Preload("Room").
		Preload("Event").
		Preload("Location").
		Find(&talkDates).Error; err != nil {
		db.log.Error("Error getting all talkDates", "err", err)
		return []*TalkDate{}, err
	}

	db.log.Debug("Returning talkDates", "talkDates", spew.Sprintf("%+v", talkDates))
	return talkDates, nil
}

func (db *TalkDateDBStore) GetTalkDateByID(id uint) (*TalkDate, error) {
	db.log.Debug("Getting talkDate by id...", "id", id)

	var talkDate TalkDate
	if err := db.
		Preload("Talk").
		Preload("Talk.Persons").
		Preload("Talk.Topics").
		Preload("Talk.Topics.Children").
		Preload("Room").
		Preload("Event").
		Preload("Location").
		First(&talkDate, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("TalkDate not found by id", "id", id)
			return nil, &TalkDateNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting talkDate by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning talkDate", "talkDate", hclog.Fmt("%+v", talkDate))
	return &talkDate, nil
}

func (db *TalkDateDBStore) UpdateTalkDate(id uint, talkDate *TalkDate) (*TalkDate, error) {
	db.log.Debug("Updating talkDate...", "talkDate", hclog.Fmt("%+v", talkDate))

	err := db.validate.Struct(talkDate)
	if err != nil {
		db.log.Error("Error validating talkDate", "err", err)
		return nil, err
	}

	if err := db.Model(&TalkDate{}).Where("id = ?", id).Take(&TalkDate{}).Update(talkDate).First(&talkDate, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("TalkDate to be updated not found", "talkDate", hclog.Fmt("%+v", talkDate))
			return nil, &TalkDateNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating talkDate", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated talkDate", "talkDate", hclog.Fmt("%+v", talkDate))
	return talkDate, nil
}

func (db *TalkDateDBStore) AddTalkDate(talkDate *TalkDate) (*TalkDate, error) {
	db.log.Debug("Adding talkDate...", "talkDate", hclog.Fmt("%+v", talkDate))

	err := db.validate.Struct(talkDate)
	if err != nil {
		db.log.Error("Error validating talkDate", "err", err)
		return nil, err
	}

	if err := db.Create(&talkDate).Error; err != nil {
		db.log.Error("Unexpected error creating talkDate", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added talkDate", "talkDate", hclog.Fmt("%+v", talkDate))
	return talkDate, nil
}

func (db *TalkDateDBStore) DeleteTalkDateByID(id uint) error {
	db.log.Debug("Deleting talkDate by id...", "id", id)

	if err := db.Model(&TalkDate{}).Where("id = ?", id).Take(&TalkDate{}).Delete(&TalkDate{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("TalkDate not found by id", "id", id)
			return &TalkDateNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting talkDate", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted talkDate")
	return nil
}

func (db *TalkDateDBStore) GetTalkDatesByEventID(eventID uint) ([]*TalkDate, error) {
	db.log.Debug("Getting talkDates by id...", "eventID", eventID)

	var talkDates []*TalkDate
	if err := db.
		Preload("Talk").
		Preload("Talk.Persons").
		Preload("Talk.Topics").
		Preload("Talk.Topics.Children").
		Preload("Room").
		Preload("Event").
		Preload("Location").
		Where(TalkDate{EventID:  eventID}).
		Find(&talkDates).Error; err != nil {
		db.log.Error("Error getting talkDates", "err", err)
		return []*TalkDate{}, err
	}

	db.log.Debug("Returning talkDates", "talkDates", spew.Sprintf("%+v", talkDates))
	return talkDates, nil
}
