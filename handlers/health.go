package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"net/http"
)

type HealthHandler struct {
	log hclog.Logger
	db  *gorm.DB
}

func NewHealthHandler(db *gorm.DB, log hclog.Logger) *HealthHandler {
	return &HealthHandler{log, db}
}

func (hh *HealthHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	hh.log.Debug("Health endpoint called...")

	if err := hh.db.DB().Ping(); err != nil {
		hh.log.Error("Error! Cannot ping database!", "err", err)
		writeUnhealthy(rw)
		return
	}

	hh.log.Debug("Returning healthy!")
	writeHealthy(rw)
}

func writeHealthy(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("."))
}

func writeUnhealthy(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusInternalServerError)
	_, _ = rw.Write([]byte("."))
}
