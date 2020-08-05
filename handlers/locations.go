package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type LocationsHandler struct {
	log   hclog.Logger
	store data.LocationStore
}

func NewLocationsHandler(store data.LocationStore, log hclog.Logger) *LocationsHandler {
	return &LocationsHandler{log, store}
}

func (lh *LocationsHandler) GetLocations(rw http.ResponseWriter, r *http.Request) {

	locations, err := lh.store.GetLocations()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(locations, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *LocationsHandler) GetLocation(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	location, err := lh.store.GetLocationByID(id)
	if err != nil {
		switch err.(type) {
		case *data.LocationNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(location, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *LocationsHandler) CreateLocation(rw http.ResponseWriter, r *http.Request) {

	location := &data.Location{}
	err := readJSON(r.Body, location)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	location, err = lh.store.AddLocation(location)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(location, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *LocationsHandler) UpdateLocation(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	location := &data.Location{}
	err := readJSON(r.Body, location)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	location, err = lh.store.UpdateLocation(id, location)
	if err != nil {
		switch err.(type) {
		case *data.LocationNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(location, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *LocationsHandler) DeleteLocation(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteLocationByID(id)
	if err != nil {
		switch err.(type) {
		case *data.LocationNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}
