package data

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
)

type Topic struct {
	// gorm.Model
	ID       uint    `json:"id" gorm:"primary_key;auto_increment"`
	Name     string  `json:"name" validate:"required" gorm:"not null"`
	Children []Topic `json:"children" gorm:"many2many:is_child_of;association_jointable_foreignkey:child_topic_id"`
}

type TopicStore interface {
	GetTopics() ([]*Topic, error)
	GetTopicByID(id uint) (*Topic, error)
	UpdateTopic(id uint, topic *Topic) (*Topic, error)
	AddTopic(topic *Topic) (*Topic, error)
	DeleteTopicByID(id uint) error
	GetTopicsByEventID(eventID uint) ([]*Topic, error)
}

type TopicDBStore struct {
	*gorm.DB
	validate *validator.Validate
	log hclog.Logger
}

type TopicNotFoundError struct {
	Cause error
}

func (e TopicNotFoundError) Error() string { return "Topic not found! Cause: " + e.Cause.Error() }
func (e TopicNotFoundError) Unwrap() error { return e.Cause }

func NewTopicDBStore(db *gorm.DB, log hclog.Logger) *TopicDBStore {
	return &TopicDBStore{db, validator.New(), log}
}

func (db *TopicDBStore) GetTopics() ([]*Topic, error) {
	db.log.Debug("Getting all topics...")

	var topics []*Topic
	if err := db.Preload("Children").Find(&topics).Error; err != nil {
		db.log.Error("Error getting all topics", "err", err)
		return []*Topic{}, err
	}

	db.log.Debug("Returning topics", "topics", spew.Sprintf("%+v", topics))
	return topics, nil
}

func (db *TopicDBStore) GetTopicByID(id uint) (*Topic, error) {
	db.log.Debug("Getting topic by id...", "id", id)

	var topic Topic
	if err := db.Preload("Children").First(&topic, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Topic not found by id", "id", id)
			return nil, &TopicNotFoundError{err}
		} else {
			db.log.Error("Unexpected error getting topic by id", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Returning topic", "topic", hclog.Fmt("%+v", topic))
	return &topic, nil
}

func (db *TopicDBStore) UpdateTopic(id uint, topic *Topic) (*Topic, error) {
	db.log.Debug("Updating topic...", "topic", hclog.Fmt("%+v", topic))

	err := db.validate.Struct(topic)
	if err != nil {
		db.log.Error("Error validating topic", "err", err)
		return nil, err
	}

	if err := db.Model(&Topic{}).Where("id = ?", id).Take(&Topic{}).Update(topic).First(&topic, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Topic to be updated not found", "topic", hclog.Fmt("%+v", topic))
			return nil, &TopicNotFoundError{err}
		} else {
			db.log.Error("Unexpected error updating topic", "err", err)
			return nil, err
		}
	}

	db.log.Debug("Successfully updated topic", "topic", hclog.Fmt("%+v", topic))
	return topic, nil
}

func (db *TopicDBStore) AddTopic(topic *Topic) (*Topic, error) {
	db.log.Debug("Adding topic...", "topic", hclog.Fmt("%+v", topic))

	err := db.validate.Struct(topic)
	if err != nil {
		db.log.Error("Error validating topic", "err", err)
		return nil, err
	}

	if err := db.Create(&topic).Error; err != nil {
		db.log.Error("Unexpected error creating topic", "err", err)
		return nil, err
	}

	db.log.Debug("Successfully added topic", "topic", hclog.Fmt("%+v", topic))
	return topic, nil
}

func (db *TopicDBStore) DeleteTopicByID(id uint) error {
	db.log.Debug("Deleting topic by id...", "id", id)

	if err := db.Model(&Topic{}).Where("id = ?", id).Take(&Topic{}).Delete(&Topic{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.log.Error("Topic not found by id", "id", id)
			return &TopicNotFoundError{err}
		} else {
			db.log.Error("Unexpected error deleting topic", "err", err)
			return err
		}
	}

	db.log.Debug("Successfully deleted topic")
	return nil
}

func (db *TopicDBStore) GetTopicsByEventID(eventID uint) ([]*Topic, error) {
	db.log.Debug("Getting topics by event id...", "eventID", eventID)

	var topics []*Topic
	if err := db.
		Table("topic").
		Select("DISTINCT *").
		Preload("Children").
		Joins("JOIN talk_topic ON talk_topic.topic_id = topic.id").
		Joins("JOIN talk_date ON talk_date.talk_id = talk_topic.talk_id").
		Where("talk_date.event_id = ?", eventID).
		Find(&topics).Error; err != nil {
			db.log.Error("Error getting topics", "err", err)
			return []*Topic{}, err
	}

	db.log.Debug("Returning topics", "topics", spew.Sprintf("%+v", topics))
	return topics, nil
}
