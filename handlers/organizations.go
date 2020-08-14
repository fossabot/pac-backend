package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type OrganizationsHandler struct {
	log   hclog.Logger
	store data.OrganizationStore
}

func NewOrganizationsHandler(store data.OrganizationStore, log hclog.Logger) *OrganizationsHandler {
	return &OrganizationsHandler{log, store}
}

func (lh *OrganizationsHandler) GetOrganizations(rw http.ResponseWriter, r *http.Request) {

	organizations, err := lh.store.GetOrganizations()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(organizations, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *OrganizationsHandler) GetOrganization(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	organization, err := lh.store.GetOrganizationByID(id)
	if err != nil {
		switch err.(type) {
		case *data.OrganizationNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(organization, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *OrganizationsHandler) CreateOrganization(rw http.ResponseWriter, r *http.Request) {

	organization := &data.Organization{}
	err := readJSON(r.Body, organization)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	organization, err = lh.store.AddOrganization(organization)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(organization, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *OrganizationsHandler) UpdateOrganization(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	organization := &data.Organization{}
	err := readJSON(r.Body, organization)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	organization, err = lh.store.UpdateOrganization(id, organization)
	if err != nil {
		switch err.(type) {
		case *data.OrganizationNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(organization, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *OrganizationsHandler) DeleteOrganization(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteOrganizationByID(id)
	if err != nil {
		switch err.(type) {
		case *data.OrganizationNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}
