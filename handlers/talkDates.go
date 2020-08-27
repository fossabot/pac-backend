package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type TalkDatesHandler struct {
	log   hclog.Logger
	store data.TalkDateStore
}

func NewTalkDatesHandler(store data.TalkDateStore, log hclog.Logger) *TalkDatesHandler {
	return &TalkDatesHandler{log, store}
}

func (lh *TalkDatesHandler) GetTalkDates(rw http.ResponseWriter, r *http.Request) {

	talkDates, err := lh.store.GetTalkDates()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(talkDates, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalkDatesHandler) GetTalkDate(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	talkDate, err := lh.store.GetTalkDateByID(id)
	if err != nil {
		switch err.(type) {
		case *data.TalkDateNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(talkDate, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalkDatesHandler) CreateTalkDate(rw http.ResponseWriter, r *http.Request) {

	talkDate := &data.TalkDate{}
	err := readJSON(r.Body, talkDate)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	talkDate, err = lh.store.AddTalkDate(talkDate)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(talkDate, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalkDatesHandler) UpdateTalkDate(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	talkDate := &data.TalkDate{}
	err := readJSON(r.Body, talkDate)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	talkDate, err = lh.store.UpdateTalkDate(id, talkDate)
	if err != nil {
		switch err.(type) {
		case *data.TalkDateNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(talkDate, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalkDatesHandler) DeleteTalkDate(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteTalkDateByID(id)
	if err != nil {
		switch err.(type) {
		case *data.TalkDateNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (lh *TalkDatesHandler) GetTalkDatesByEventID(rw http.ResponseWriter, r *http.Request) {
	eventID := readId(r)

	talkDates, err := lh.store.GetTalkDatesByEventID(eventID)
	if err != nil {
		writeJSONErrorWithStatus("Error getting entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(talkDates, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}