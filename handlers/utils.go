package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// ErrorResponse is a generic error message returned by the server
type ErrorResponse struct {
	Message string `json:"Message"`
	Cause   string `json:"Cause"`
}

func readJSON(rc io.Reader, dst interface{}) error {

	dec := json.NewDecoder(rc)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		// TODO better error handling
		return err
	}

	return nil
}

func writeJSONWithStatus(i interface{}, rw http.ResponseWriter, status int) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(i)
}

func writeJSONErrorWithStatus(message string, cause string, rw http.ResponseWriter, status int) {

	errorResponse := ErrorResponse{message, cause}
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(&errorResponse)
}

func readId(r *http.Request) uint {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		// should not happen
		panic(err)
	}

	return uint(id)
}
