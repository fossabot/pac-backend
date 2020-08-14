package dbinit

import (
	"github.com/jinzhu/gorm"
	"github.com/milutindzunic/pac-backend/data"
	"time"
)

func DB(db *gorm.DB, ls data.LocationStore, es data.EventStore, os data.OrganizationStore, ps data.PersonStore, rs data.RoomStore, ts data.TopicStore, tlks data.TalkStore, tlkds data.TalkDateStore) {

	// Parameter objects do not have primary keys - gorm will delete all the records!
	db.Delete(data.Location{})
	db.Delete(data.Event{})
	db.Delete(data.Organization{})
	db.Delete(data.Person{})
	db.Delete(data.Room{})
	db.Delete(data.Topic{})
	db.Delete(data.Talk{})
	db.Delete(data.TalkDate{})

	// Locations
	belexpoLocation, _ := ls.AddLocation(&data.Location{
		ID:   1,
		Name: "Belexpo Centar",
	})
	hotelPlazaLocation, _ := ls.AddLocation(&data.Location{
		ID:   2,
		Name: "Hotel Plaza",
	})

	// Events
	bestJavaConf, _ := es.AddEvent(&data.Event{
		ID:         1,
		Name:       "Best Java Conference",
		BeginDate:  time.Date(2021, time.Month(2), 12, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		LocationID: belexpoLocation.ID,
	})
	prodynaConf, _ := es.AddEvent(&data.Event{
		ID:         2,
		Name:       "Prodyna Job Fair",
		BeginDate:  time.Date(2021, time.Month(5), 2, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(6), 5, 0, 0, 0, 0, time.UTC),
		LocationID: hotelPlazaLocation.ID,
	})

	// Organizations
	prodyna, _ := os.AddOrganization(&data.Organization{
		ID:   1,
		Name: "Prodyna",
	})

	// Rooms
	redRoom, _ := rs.AddRoom(&data.Room{
		ID:             1,
		Name:           "Red Room",
		OrganizationID: prodyna.ID,
	})
	whiteRoom, _ := rs.AddRoom(&data.Room{
		ID:             2,
		Name:           "White Room",
		OrganizationID: prodyna.ID,
	})

	// Topics
	javaTopic, _ := ts.AddTopic(&data.Topic{
		ID:       1,
		Name:     "Java",
		Children: nil,
	})
	hibernateTopic, _ := ts.AddTopic(&data.Topic{
		ID:       2,
		Name:     "Java",
		Children: []data.Topic{*javaTopic},
	})
	springTopic, _ := ts.AddTopic(&data.Topic{
		ID:       3,
		Name:     "Java",
		Children: []data.Topic{*javaTopic, *hibernateTopic},
	})
	kubernetesTopic, _ := ts.AddTopic(&data.Topic{
		ID:       4,
		Name:     "Kubernetes",
		Children: nil,
	})
	jsTopic, _ := ts.AddTopic(&data.Topic{
		ID:       5,
		Name:     "JavaScript",
		Children: nil,
	})

	// Persons
	dkrizic, _ := ps.AddPerson(&data.Person{
		ID:             1,
		Name:           "Darko Krizic",
		OrganizationID: prodyna.ID,
	})
	ggrujic, _ := ps.AddPerson(&data.Person{
		ID:             2,
		Name:           "Goran Grujic",
		OrganizationID: prodyna.ID,
	})
	mnikolic, _ := ps.AddPerson(&data.Person{
		ID:             3,
		Name:           "Milos Nikolic",
		OrganizationID: prodyna.ID,
	})

	// Talks
	javaSpringAndYouTalk, _ := tlks.AddTalk(&data.Talk{
		ID:                1,
		Title:             "Java, Spring, and You",
		DurationInMinutes: 90,
		Language:          "English",
		Level:             data.BeginnerLevel,
		Persons:           []data.Person{*dkrizic, *ggrujic},
		Topics:            []data.Topic{*javaTopic, *springTopic, *hibernateTopic},
		TalkDates:         nil,
	})
	fullStackJavaScriptOnKubernetesTalk, _ := tlks.AddTalk(&data.Talk{
		ID:                2,
		Title:             "Fullstack JavaScript on Kubernetes",
		DurationInMinutes: 60,
		Language:          "Serbian",
		Level:             data.ExpertLevel,
		Persons:           []data.Person{*mnikolic},
		Topics:            []data.Topic{*jsTopic, *kubernetesTopic},
		TalkDates:         nil,
	})

	// TalkDates
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         1,
		BeginDate:  time.Date(2021, time.Month(2), 12, 14, 0, 0, 0, time.UTC),
		TalkID:     javaSpringAndYouTalk.ID,
		RoomID:     redRoom.ID,
		EventID:    bestJavaConf.ID,
		LocationID: belexpoLocation.ID,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         2,
		BeginDate:  time.Date(2021, time.Month(5), 2, 10, 0, 0, 0, time.UTC),
		TalkID:     fullStackJavaScriptOnKubernetesTalk.ID,
		RoomID:     whiteRoom.ID,
		EventID:    prodynaConf.ID,
		LocationID: hotelPlazaLocation.ID,
	})
}
