package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/milutindzunic/pac-backend/data"
	"net/http"
)

type TopicsHandler struct {
	log   hclog.Logger
	store data.TopicStore
}

func NewTopicsHandler(store data.TopicStore, log hclog.Logger) *TopicsHandler {
	return &TopicsHandler{log, store}
}

func (lh *TopicsHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {

	topics, err := lh.store.GetTopics()
	if err != nil {
		writeJSONErrorWithStatus("Error getting all entities", err.Error(), rw, http.StatusInternalServerError)
		return
	}

	err = writeJSONWithStatus(topics, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TopicsHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	topic, err := lh.store.GetTopicByID(id)
	if err != nil {
		switch err.(type) {
		case *data.TopicNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(topic, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TopicsHandler) CreateTopic(rw http.ResponseWriter, r *http.Request) {

	topic := &data.Topic{}
	err := readJSON(r.Body, topic)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	topic, err = lh.store.AddTopic(topic)
	if err != nil {
		writeJSONErrorWithStatus("Error creating entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	err = writeJSONWithStatus(topic, rw, http.StatusCreated)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TopicsHandler) UpdateTopic(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	topic := &data.Topic{}
	err := readJSON(r.Body, topic)
	if err != nil {
		lh.log.Error("Error deserializing entity", err)
		writeJSONErrorWithStatus("Error deserializing entity", err.Error(), rw, http.StatusBadRequest)
		return
	}

	// TODO validacija ako treba

	topic, err = lh.store.UpdateTopic(id, topic)
	if err != nil {
		switch err.(type) {
		case *data.TopicNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	err = writeJSONWithStatus(topic, rw, http.StatusOK)
	if err != nil {
		lh.log.Error("Error serializing entity", err)
		return
	}
}

func (lh *TopicsHandler) DeleteTopic(rw http.ResponseWriter, r *http.Request) {
	id := readId(r)

	err := lh.store.DeleteTopicByID(id)
	if err != nil {
		switch err.(type) {
		case *data.TopicNotFoundError:
			writeJSONErrorWithStatus("Entity not found", err.Error(), rw, http.StatusNotFound)
			return
		default:
			writeJSONErrorWithStatus("Unexpected error occurred", err.Error(), rw, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (lh *TopicsHandler) GetTopicsByEventID(rw http.ResponseWriter, r *http.Request) {
	eventID := readId(r)

	events, err := lh.store.GetTopicsByEventID(eventID)
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
