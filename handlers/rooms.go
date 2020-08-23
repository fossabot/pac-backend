package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type RoomsHandler struct {
	log   hclog.Logger
	store data.RoomStore
}

func NewRoomsHandler(store data.RoomStore, log hclog.Logger) *RoomsHandler {
	return &RoomsHandler{log, store}
}

func (lh *RoomsHandler) GetRooms(rw http.ResponseWriter, r *http.Request) {

	rooms, err := lh.store.GetRooms()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(rooms, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *RoomsHandler) GetRoom(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	room, err := lh.store.GetRoomByID(id)
	if err != nil {
		switch err.(type) {
		case *data.RoomNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(room, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *RoomsHandler) CreateRoom(rw http.ResponseWriter, r *http.Request) {

	room := &data.Room{}
	err := readJSON(r.Body, room)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	room, err = lh.store.AddRoom(room)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(room, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *RoomsHandler) UpdateRoom(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	room := &data.Room{}
	err := readJSON(r.Body, room)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	room, err = lh.store.UpdateRoom(id, room)
	if err != nil {
		switch err.(type) {
		case *data.RoomNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(room, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *RoomsHandler) DeleteRoom(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteRoomByID(id)
	if err != nil {
		switch err.(type) {
		case *data.RoomNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}
