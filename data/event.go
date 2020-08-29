package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"time"
)

type Event struct {
	// gorm.Model
	ID         uint      `json:"id" gorm:"primary_key;auto_increment"`
	Name       string    `json:"name" gorm:"not null"`
	BeginDate  time.Time `json:"beginDate" gorm:"not null"`
	EndDate    time.Time `json:"endDate" gorm:"not null"`
	LocationID uint      `json:"-"`
	Location   *Location `json:"location,omitempty" gorm:"association_autoupdate:false"`
}

type EventStore interface {
	GetEvents() ([]*Event, error)
	GetEventByID(id uint) (*Event, error)
	UpdateEvent(id uint, event *Event) (*Event, error)
	AddEvent(event *Event) (*Event, error)
	DeleteEventByID(id uint) error
	GetEventsByTalkID(talkID uint) ([]*Event, error)
}

type EventDBStore struct {
	*gorm.DB
	validate *validator.Validate
	log      hclog.Logger
}

type EventNotFoundError struct {
	Cause error
}

func (e EventNotFoundError) Error() string { return "Event not found! Cause: " + e.Cause.Error() }
func (e EventNotFoundError) Unwrap() error { return e.Cause }

func NewEventDBStore(db *gorm.DB, log hclog.Logger) *EventDBStore {
	return &EventDBStore{db, validator.New(), log}
}

func (db *EventDBStore) GetEvents() ([]*Event, error) {
	db.log.Debug("Getting all events...")

	var events []*Event
	if err := db.Preload("Location").Find(&events).Error; err != nil {
		db.log.Error("Error getting all events", "err", err)
		return []*Event{}, err
	}

	db.log.Debug("Returning events", "events", spew.Sprintf("%+v", events))
	return events, nil
}

func (db *EventDBStore) GetEventByID(id uint) (*Event, error) {
	db.log.Debug("Getting event by id...", "id", id)

	var event Event
	if err := db.Preload("Location").First(&event, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Event not found by id", "id", id)
			return nil, &EventNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting event by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning event", "event", hclog.Fmt("%+v", event))
	return &event, nil
}

func (db *EventDBStore) UpdateEvent(id uint, event *Event) (*Event, error) {
	db.log.Debug("Updating event...", "event", hclog.Fmt("%+v", event))

	err := db.validate.Struct(event)
	if err != nil {
		db.log.Error("Error validating event", "err", err)
		return nil, err
	}

	if err := db.Model(&Event{}).Where("id = ?", id).Take(&Event{}).Update(event).First(&event, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Event to be updated not found", "event", hclog.Fmt("%+v", event))
			return nil, &EventNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating event", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated event", "event", hclog.Fmt("%+v", event))
	return db.GetEventByID(event.ID)
}

func (db *EventDBStore) AddEvent(event *Event) (*Event, error) {
	db.log.Debug("Adding event...", "event", hclog.Fmt("%+v", event))

	err := db.validate.Struct(event)
	if err != nil {
		db.log.Error("Error validating event", "err", err)
		return nil, err
	}

	if err := db.Create(&event).Error; err != nil {
		db.log.Error("Unexpected error creating event", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added event", "event", hclog.Fmt("%+v", event))
	return db.GetEventByID(event.ID)
}

func (db *EventDBStore) DeleteEventByID(id uint) error {
	db.log.Debug("Deleting event by id...", "id", id)

	if err := db.Model(&Event{}).Where("id = ?", id).Take(&Event{}).Delete(&Event{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Event not found by id", "id", id)
			return &EventNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting event", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted event")
	return nil
}

func (db *EventDBStore) GetEventsByTalkID(talkID uint) ([]*Event, error) {
	db.log.Debug("Getting event by talk id...", "talkID", talkID)

	var events []*Event
	if err := db.
		Preload("Location").
		Where("id IN ?", db.Table("talk_date").Select("event_id").Where("talk_id = ?", talkID).SubQuery()).
		Find(&events).Error; err != nil {
		db.log.Error("Error getting events", "err", err)
		return []*Event{}, err
	}

	db.log.Debug("Returning events", "events", spew.Sprintf("%+v", events))
	return events, nil
}
