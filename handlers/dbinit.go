package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/milutindzunic/pac-backend/data"
	"github.com/milutindzunic/pac-backend/database"
	"net/http"
)

type DBInitHandler struct {
	db                *gorm.DB
	logger            hclog.Logger
	eventStore        data.EventStore
	locationStore     data.LocationStore
	organizationStore data.OrganizationStore
	personStore       data.PersonStore
	roomStore         data.RoomStore
	topicStore        data.TopicStore
	talkStore         data.TalkStore
	talkDateStore     data.TalkDateStore
}

func NewDBInitHandler(db *gorm.DB, ls data.LocationStore, es data.EventStore, os data.OrganizationStore, ps data.PersonStore, rs data.RoomStore, ts data.TopicStore, tlks data.TalkStore, tlkds data.TalkDateStore, logger hclog.Logger) *DBInitHandler {
	return &DBInitHandler{
		db:                db,
		logger:            logger,
		eventStore:        es,
		locationStore:     ls,
		organizationStore: os,
		personStore:       ps,
		roomStore:         rs,
		topicStore:        ts,
		talkStore:         tlks,
		talkDateStore:     tlkds,
	}
}

func (ih *DBInitHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	ih.logger.Debug("Init database endpoint called...")

	database.Init(ih.db, ih.locationStore, ih.eventStore, ih.organizationStore, ih.personStore, ih.roomStore, ih.topicStore, ih.talkStore, ih.talkDateStore, ih.logger)

	rw.WriteHeader(http.StatusNoContent)
}
