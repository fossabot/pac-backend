package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type PersonsHandler struct {
	log   hclog.Logger
	store data.PersonStore
}

func NewPersonsHandler(store data.PersonStore, log hclog.Logger) *PersonsHandler {
	return &PersonsHandler{log, store}
}

func (lh *PersonsHandler) GetPersons(rw http.ResponseWriter, r *http.Request) {

	persons, err := lh.store.GetPersons()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(persons, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *PersonsHandler) GetPerson(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	person, err := lh.store.GetPersonByID(id)
	if err != nil {
		switch err.(type) {
		case *data.PersonNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(person, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *PersonsHandler) CreatePerson(rw http.ResponseWriter, r *http.Request) {

	person := &data.Person{}
	err := readJSON(r.Body, person)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	person, err = lh.store.AddPerson(person)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(person, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *PersonsHandler) UpdatePerson(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	person := &data.Person{}
	err := readJSON(r.Body, person)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	person, err = lh.store.UpdatePerson(id, person)
	if err != nil {
		switch err.(type) {
		case *data.PersonNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(person, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *PersonsHandler) DeletePerson(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeletePersonByID(id)
	if err != nil {
		switch err.(type) {
		case *data.PersonNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}
