package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type TalksHandler struct {
	log   hclog.Logger
	store data.TalkStore
}

func NewTalksHandler(store data.TalkStore, log hclog.Logger) *TalksHandler {
	return &TalksHandler{log, store}
}

func (lh *TalksHandler) GetTalks(rw http.ResponseWriter, r *http.Request) {

	talks, err := lh.store.GetTalks()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(talks, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalksHandler) GetTalk(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	talk, err := lh.store.GetTalkByID(id)
	if err != nil {
		switch err.(type) {
		case *data.TalkNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(talk, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalksHandler) CreateTalk(rw http.ResponseWriter, r *http.Request) {

	talk := &data.Talk{}
	err := readJSON(r.Body, talk)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	talk, err = lh.store.AddTalk(talk)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(talk, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalksHandler) UpdateTalk(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	talk := &data.Talk{}
	err := readJSON(r.Body, talk)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	talk, err = lh.store.UpdateTalk(id, talk)
	if err != nil {
		switch err.(type) {
		case *data.TalkNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(talk, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TalksHandler) DeleteTalk(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteTalkByID(id)
	if err != nil {
		switch err.(type) {
		case *data.TalkNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}

