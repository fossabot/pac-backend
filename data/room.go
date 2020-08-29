package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
)

type Room struct {
	// gorm.Model
	ID             uint          `json:"id" gorm:"primary_key;auto_increment"`
	Name           string        `json:"name" gorm:"not null;default:''"`
	OrganizationID uint          `json:"-" gorm:"not null"`
	Organization   *Organization `json:"organization,omitempty" gorm:"association_autoupdate:false"`
}

type RoomStore interface {
	GetRooms() ([]*Room, error)
	GetRoomByID(id uint) (*Room, error)
	UpdateRoom(id uint, room *Room) (*Room, error)
	AddRoom(room *Room) (*Room, error)
	DeleteRoomByID(id uint) error
}

type RoomDBStore struct {
	*gorm.DB
	validate *validator.Validate
	log      hclog.Logger
}

type RoomNotFoundError struct {
	Cause error
}

func (e RoomNotFoundError) Error() string { return "Room not found! Cause: " + e.Cause.Error() }
func (e RoomNotFoundError) Unwrap() error { return e.Cause }

func NewRoomDBStore(db *gorm.DB, log hclog.Logger) *RoomDBStore {
	return &RoomDBStore{db, validator.New(), log}
}

func (db *RoomDBStore) GetRooms() ([]*Room, error) {
	db.log.Debug("Getting all rooms...")

	var rooms []*Room
	if err := db.Preload("Organization").Find(&rooms).Error; err != nil {
		db.log.Error("Error getting all rooms", "err", err)
		return []*Room{}, err
	}

	db.log.Debug("Returning rooms", "rooms", spew.Sprintf("%+v", rooms))
	return rooms, nil
}

func (db *RoomDBStore) GetRoomByID(id uint) (*Room, error) {
	db.log.Debug("Getting room by id...", "id", id)

	var room Room
	if err := db.Preload("Organization").First(&room, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Room not found by id", "id", id)
			return nil, &RoomNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting room by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning room", "room", hclog.Fmt("%+v", room))
	return &room, nil
}

func (db *RoomDBStore) UpdateRoom(id uint, room *Room) (*Room, error) {
	db.log.Debug("Updating room...", "room", hclog.Fmt("%+v", room))

	err := db.validate.Struct(room)
	if err != nil {
		db.log.Error("Error validating room", "err", err)
		return nil, err
	}

	if err := db.Model(&Room{}).Where("id = ?", id).Take(&Room{}).Update(room).First(&room, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Room to be updated not found", "room", hclog.Fmt("%+v", room))
			return nil, &RoomNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating room", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated room", "room", hclog.Fmt("%+v", room))
	return db.GetRoomByID(room.ID)
}

func (db *RoomDBStore) AddRoom(room *Room) (*Room, error) {
	db.log.Debug("Adding room...", "room", hclog.Fmt("%+v", room))

	err := db.validate.Struct(room)
	if err != nil {
		db.log.Error("Error validating room", "err", err)
		return nil, err
	}

	if err := db.Create(&room).Error; err != nil {
		db.log.Error("Unexpected error creating room", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added room", "room", hclog.Fmt("%+v", room))
	return db.GetRoomByID(room.ID)
}

func (db *RoomDBStore) DeleteRoomByID(id uint) error {
	db.log.Debug("Deleting room by id...", "id", id)

	if err := db.Model(&Room{}).Where("id = ?", id).Take(&Room{}).Delete(&Room{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Room not found by id", "id", id)
			return &RoomNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting room", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted room")
	return nil
}
