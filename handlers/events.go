package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type EventsHandler struct {
	log   hclog.Logger
	store data.EventStore
}

func NewEventsHandler(store data.EventStore, log hclog.Logger) *EventsHandler {
	return &EventsHandler{log, store}
}

func (lh *EventsHandler) GetEvents(rw http.ResponseWriter, r *http.Request) {

	events, err := lh.store.GetEvents()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(events, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *EventsHandler) GetEvent(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	event, err := lh.store.GetEventByID(id)
	if err != nil {
		switch err.(type) {
		case *data.EventNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(event, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *EventsHandler) CreateEvent(rw http.ResponseWriter, r *http.Request) {

	event := &data.Event{}
	err := readJSON(r.Body, event)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	event, err = lh.store.AddEvent(event)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(event, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *EventsHandler) UpdateEvent(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	event := &data.Event{}
	err := readJSON(r.Body, event)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	event, err = lh.store.UpdateEvent(id, event)
	if err != nil {
		switch err.(type) {
		case *data.EventNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(event, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *EventsHandler) DeleteEvent(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteEventByID(id)
	if err != nil {
		switch err.(type) {
		case *data.EventNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (lh *EventsHandler) GetEventsByTalkID(rw http.ResponseWriter, r *http.Request) {
	talkID := readId(r)

	events, err := lh.store.GetEventsByTalkID(talkID)
	if err != nil {
		writeJSONErrorWithStatus("Error getting entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(events, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}