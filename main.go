package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"  //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite database driver
	"github.com/justinas/alice"
	"github.com/milutindzunic/pac-backend/config"
	"github.com/milutindzunic/pac-backend/data"
	"github.com/milutindzunic/pac-backend/handlers"
	"github.com/milutindzunic/pac-backend/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func healthHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Health check handler called!")
	rw.WriteHeader(http.StatusNoContent)
}

func main() {

	// load the configuration
	cnf, err := config.LoadConfig()
	if err != nil {
		log.Println("Failed to read config")
		panic(err)
	}
	log.Printf("Loaded config: %+v\n", cnf)

	// set up application logger instance
	logger := hclog.New(&hclog.LoggerOptions{
		Output:          os.Stdout,
		Level:           hclog.LevelFromString(cnf.LogLevel),
		IncludeLocation: true,
	})

	// connect to database, defer closing
	db, err := openDB(cnf)
	if err != nil {
		logger.Error("Failed to connect to database")
		panic(err)
	}
	db.LogMode(cnf.LogPersistence)
	defer db.Close()

	// Keep the schema up to date
	db = db.AutoMigrate(&data.Location{})

	// create stores
	var locationStore data.LocationStore = data.NewLocationDBStore(db, logger)

	// create handlers
	lh := handlers.NewLocationsHandler(locationStore, logger)

	sm := mux.NewRouter()

	jsonChain := alice.New(middleware.EnforceJsonContentType)
	sm.HandleFunc("/", healthHandler)
	sm.Handle("/locations", http.HandlerFunc(lh.GetLocations)).Methods("GET")
	sm.Handle("/locations/{id:[0-9]+}", http.HandlerFunc(lh.GetLocation)).Methods("GET")
	sm.Handle("/locations", jsonChain.Then(http.HandlerFunc(lh.CreateLocation))).Methods("POST")
	sm.Handle("/locations/{id:[0-9]+}", jsonChain.Then(http.HandlerFunc(lh.UpdateLocation))).Methods("PUT")
	sm.Handle("/locations/{id:[0-9]+}", http.HandlerFunc(lh.DeleteLocation)).Methods("DELETE")

	// Prometheus metrics handler
	sm.Handle("/metrics", promhttp.Handler())

	// create Server
	s := http.Server{
		Addr:         cnf.BindAddress,
		Handler:      sm,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
	}

	go func() {
		logger.Info("Starting server on " + s.Addr)

		err := s.ListenAndServe()
		if err != nil {
			logger.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received.
	sig := <-c
	logger.Info("Received signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	logger.Info("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func openDB(cnf *config.Config) (*gorm.DB, error) {
	var dbUrl string

	switch cnf.DbDriver {
	case "sqlite3":
		dbUrl = cnf.DbName
		log.Println("Connecting to embedded sqlite3 database... file name: " + dbUrl)
	case "mysql":
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", cnf.DbUser, cnf.DbPassword, cnf.DbHost, cnf.DbPort, cnf.DbName)
		log.Println("Connecting to mysql database... uri: " + dbUrl)
	default:
		return nil, fmt.Errorf("error! Database driver must be one of: [sqlite3, mysql], was %s", cnf.DbDriver)
	}

	return gorm.Open(cnf.DbDriver, dbUrl)
}
