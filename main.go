package main

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/justinas/alice"
	"github.com/milutindzunic/pac-backend/auth"
	"github.com/milutindzunic/pac-backend/config"
	"github.com/milutindzunic/pac-backend/data"
	"github.com/milutindzunic/pac-backend/database"
	"github.com/milutindzunic/pac-backend/handlers"
	"github.com/milutindzunic/pac-backend/middleware"
	"github.com/milutindzunic/pac-backend/middleware/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// load the configuration
	cnf, err := config.LoadConfig()
	if err != nil {
		log.Println("Failed to load config")
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
	db, err := database.OpenDB(cnf)
	if err != nil {
		logger.Error("Failed to connect to database", "err", err)
		panic(err)
	}
	defer db.Close()

	// create stores
	var locationStore data.LocationStore = data.NewLocationDBStore(db, logger)
	var eventStore data.EventStore = data.NewEventDBStore(db, logger)
	var organizationStore data.OrganizationStore = data.NewOrganizationDBStore(db, logger)
	var personStore data.PersonStore = data.NewPersonDBStore(db, logger)
	var roomStore data.RoomStore = data.NewRoomDBStore(db, logger)
	var topicStore data.TopicStore = data.NewTopicDBStore(db, logger)
	var talkStore data.TalkStore = data.NewTalkDBStore(db, logger)
	var talkDateStore data.TalkDateStore = data.NewTalkDateDBStore(db, logger)

	// create handlers
	hh := handlers.NewHealthHandler(db, logger)
	lh := handlers.NewLocationsHandler(locationStore, logger)
	eh := handlers.NewEventsHandler(eventStore, logger)
	oh := handlers.NewOrganizationsHandler(organizationStore, logger)
	ph := handlers.NewPersonsHandler(personStore, logger)
	rh := handlers.NewRoomsHandler(roomStore, logger)
	th := handlers.NewTopicsHandler(topicStore, logger)
	tkh := handlers.NewTalksHandler(talkStore, logger)
	tdh := handlers.NewTalkDatesHandler(talkDateStore, logger)
	ih := handlers.NewDBInitHandler(db, locationStore, eventStore, organizationStore, personStore, roomStore, topicStore, talkStore, talkDateStore, logger)

	// TODO ovde radi testiranja
	database.Init(db, locationStore, eventStore, organizationStore, personStore, roomStore, topicStore, talkStore, talkDateStore, logger)

	// Authentication
	oauth, err := auth.NewProvider(auth.OauthConfig{
		Enabled:      cnf.OAuthEnable,
		Issuer:       cnf.OAuthIssuer,
		ClientID:     cnf.OAuthClientId,
		ClientSecret: cnf.OAuthClientSecret,
		RedirectURL:  cnf.OAuthRedirectUrl,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}, logger)
	if err != nil {
		logger.Error("Failed to create Oauth2 configuration", "err", err)
		panic(err)
	}

	// Handler chains
	defaultChain := alice.New(metrics.Prometheus)
	secureChain := defaultChain.Append(middleware.EnforceJsonContentType)
	secureJsonChain := secureChain
	if cnf.OAuthEnable {
		secureJsonChain = secureChain.Append(oauth.Middleware)
	}

	sm := mux.NewRouter()
	sm.Use(middleware.AllowCORS)

	// Register handlers
	// Health
	sm.HandleFunc("/", hh.Handle)
	// Locations
	sm.Handle("/locations", defaultChain.Then(http.HandlerFunc(lh.GetLocations))).Methods("GET")
	sm.Handle("/locations/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(lh.GetLocation))).Methods("GET")
	sm.Handle("/locations", secureJsonChain.Then(http.HandlerFunc(lh.CreateLocation))).Methods("POST", "OPTIONS")
	sm.Handle("/locations/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(lh.UpdateLocation))).Methods("PUT", "OPTIONS")
	sm.Handle("/locations/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(lh.DeleteLocation))).Methods("DELETE", "OPTIONS")
	// Events
	sm.Handle("/events", defaultChain.Then(http.HandlerFunc(eh.GetEvents))).Methods("GET")
	sm.Handle("/events/talk/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(eh.GetEventsByTalkID))).Methods("GET")
	sm.Handle("/events/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(eh.GetEvent))).Methods("GET")
	sm.Handle("/events", secureJsonChain.Then(http.HandlerFunc(eh.CreateEvent))).Methods("POST", "OPTIONS")
	sm.Handle("/events/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(eh.UpdateEvent))).Methods("PUT", "OPTIONS")
	sm.Handle("/events/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(eh.DeleteEvent))).Methods("DELETE", "OPTIONS")
	// Organizations
	sm.Handle("/organizations", defaultChain.Then(http.HandlerFunc(oh.GetOrganizations))).Methods("GET")
	sm.Handle("/organizations/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(oh.GetOrganization))).Methods("GET")
	sm.Handle("/organizations", secureJsonChain.Then(http.HandlerFunc(oh.CreateOrganization))).Methods("POST", "OPTIONS")
	sm.Handle("/organizations/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(oh.UpdateOrganization))).Methods("PUT", "OPTIONS")
	sm.Handle("/organizations/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(oh.DeleteOrganization))).Methods("DELETE", "OPTIONS")
	// Persons
	sm.Handle("/persons", defaultChain.Then(http.HandlerFunc(ph.GetPersons))).Methods("GET")
	sm.Handle("/persons/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(ph.GetPerson))).Methods("GET")
	sm.Handle("/persons", secureJsonChain.Then(http.HandlerFunc(ph.CreatePerson))).Methods("POST", "OPTIONS")
	sm.Handle("/persons/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(ph.UpdatePerson))).Methods("PUT", "OPTIONS")
	sm.Handle("/persons/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(ph.DeletePerson))).Methods("DELETE", "OPTIONS")
	// Rooms
	sm.Handle("/rooms", defaultChain.Then(http.HandlerFunc(rh.GetRooms))).Methods("GET")
	sm.Handle("/rooms/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(rh.GetRoom))).Methods("GET")
	sm.Handle("/rooms", secureJsonChain.Then(http.HandlerFunc(rh.CreateRoom))).Methods("POST", "OPTIONS")
	sm.Handle("/rooms/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(rh.UpdateRoom))).Methods("PUT", "OPTIONS")
	sm.Handle("/rooms/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(rh.DeleteRoom))).Methods("DELETE", "OPTIONS")
	// Topics
	sm.Handle("/topics", defaultChain.Then(http.HandlerFunc(th.GetTopics))).Methods("GET")
	sm.Handle("/topics/event/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(th.GetTopicsByEventID))).Methods("GET")
	sm.Handle("/topics/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(th.GetTopic))).Methods("GET")
	sm.Handle("/topics", secureJsonChain.Then(http.HandlerFunc(th.CreateTopic))).Methods("POST", "OPTIONS")
	sm.Handle("/topics/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(th.UpdateTopic))).Methods("PUT", "OPTIONS")
	sm.Handle("/topics/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(th.DeleteTopic))).Methods("DELETE", "OPTIONS")
	// Talks
	sm.Handle("/talks", defaultChain.Then(http.HandlerFunc(tkh.GetTalks))).Methods("GET")
	sm.Handle("/talks/event/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(tkh.GetTalksByEventID))).Methods("GET")
	sm.Handle("/talks/person/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(tkh.GetTalksByPersonID))).Methods("GET")
	sm.Handle("/talks/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(tkh.GetTalk))).Methods("GET")
	sm.Handle("/talks", secureJsonChain.Then(http.HandlerFunc(tkh.CreateTalk))).Methods("POST", "OPTIONS")
	sm.Handle("/talks/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(tkh.UpdateTalk))).Methods("PUT", "OPTIONS")
	sm.Handle("/talks/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(tkh.DeleteTalk))).Methods("DELETE", "OPTIONS")
	// Talk Dates
	sm.Handle("/talkDates", defaultChain.Then(http.HandlerFunc(tdh.GetTalkDates))).Methods("GET")
	sm.Handle("/talkDates/event/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(tdh.GetTalkDatesByEventID))).Methods("GET")
	sm.Handle("/talkDates/{id:[0-9]+}", defaultChain.Then(http.HandlerFunc(tdh.GetTalkDate))).Methods("GET")
	sm.Handle("/talkDates", secureJsonChain.Then(http.HandlerFunc(tdh.CreateTalkDate))).Methods("POST", "OPTIONS")
	sm.Handle("/talkDates/{id:[0-9]+}", secureJsonChain.Then(http.HandlerFunc(tdh.UpdateTalkDate))).Methods("PUT", "OPTIONS")
	sm.Handle("/talkDates/{id:[0-9]+}", secureChain.Then(http.HandlerFunc(tdh.DeleteTalkDate))).Methods("DELETE", "OPTIONS")

	// OAuth2 callback
	sm.Handle("/oauth2/callback", oauth.CallbackHandler())

	// Prometheus metrics handler
	sm.Handle("/metrics", promhttp.Handler())

	// Database init handler
	sm.Handle("/initDB", http.HandlerFunc(ih.Handle))

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
