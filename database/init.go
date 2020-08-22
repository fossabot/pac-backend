package database

import (
	"github.com/jinzhu/gorm"
	"github.com/milutindzunic/pac-backend/data"
	"time"
)

func Init(db *gorm.DB, ls data.LocationStore, es data.EventStore, os data.OrganizationStore, ps data.PersonStore, rs data.RoomStore, ts data.TopicStore, tlks data.TalkStore, tlkds data.TalkDateStore) {

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
	locationBelexpo, _ := ls.AddLocation(&data.Location{
		ID:   1,
		Name: "Belexpo Centar",
	})
	locationHotelPlaza, _ := ls.AddLocation(&data.Location{
		ID:   2,
		Name: "Hotel Plaza",
	})
	locationBelgradeFair, _ := ls.AddLocation(&data.Location{
		ID:   3,
		Name: "Belgrade Fair Building One",
	})

	// Events
	eventBestJavaConference, _ := es.AddEvent(&data.Event{
		ID:         1,
		Name:       "Best Java Conference",
		BeginDate:  time.Date(2021, time.Month(2), 12, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		LocationID: locationBelexpo.ID,
	})
	eventProdynaJobFair, _ := es.AddEvent(&data.Event{
		ID:         2,
		Name:       "Prodyna Job Fair",
		BeginDate:  time.Date(2021, time.Month(5), 2, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(6), 5, 0, 0, 0, 0, time.UTC),
		LocationID: locationHotelPlaza.ID,
	})
	eventITConnect, _ := es.AddEvent(&data.Event{
		ID:         3,
		Name:       "IT Connect",
		BeginDate:  time.Date(2021, time.Month(5), 10, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(5), 12, 0, 0, 0, 0, time.UTC),
		LocationID: locationHotelPlaza.ID,
	})
	/*eventCloudnativeConference*/ _, _ = es.AddEvent(&data.Event{
		ID:         4,
		Name:       "Prodyna Job Fair",
		BeginDate:  time.Date(2021, time.Month(5), 22, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(5), 23, 0, 0, 0, 0, time.UTC),
		LocationID: locationBelgradeFair.ID,
	})
	eventGoogleIO, _ := es.AddEvent(&data.Event{
		ID:         5,
		Name:       "Google I/O",
		BeginDate:  time.Date(2021, time.Month(6), 2, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2021, time.Month(6), 5, 0, 0, 0, 0, time.UTC),
		LocationID: locationBelgradeFair.ID,
	})

	// Organizations
	organizationProdyna, _ := os.AddOrganization(&data.Organization{
		ID:   1,
		Name: "Prodyna",
	})
	// Organizations
	organizationGoogle, _ := os.AddOrganization(&data.Organization{
		ID:   2,
		Name: "Google",
	})

	// Rooms
	roomRed, _ := rs.AddRoom(&data.Room{
		ID:             1,
		Name:           "Red Room",
		OrganizationID: organizationProdyna.ID,
	})
	roomWhite, _ := rs.AddRoom(&data.Room{
		ID:             2,
		Name:           "White Room",
		OrganizationID: organizationProdyna.ID,
	})
	roomBlue, _ := rs.AddRoom(&data.Room{
		ID:             3,
		Name:           "Blue Room",
		OrganizationID: organizationProdyna.ID,
	})
	roomGoogle, _ := rs.AddRoom(&data.Room{
		ID:             4,
		Name:           "Google Room",
		OrganizationID: organizationGoogle.ID,
	})

	// Topics
	topicJava, _ := ts.AddTopic(&data.Topic{
		ID:       1,
		Name:     "Java",
		Children: nil,
	})
	topicHibernate, _ := ts.AddTopic(&data.Topic{
		ID:       2,
		Name:     "Hibernate",
		Children: []data.Topic{*topicJava},
	})
	topicSpring, _ := ts.AddTopic(&data.Topic{
		ID:       3,
		Name:     "Spring",
		Children: []data.Topic{*topicJava, *topicHibernate},
	})
	topicKubernetes, _ := ts.AddTopic(&data.Topic{
		ID:       4,
		Name:     "Kubernetes",
		Children: nil,
	})
	topicJavaScript, _ := ts.AddTopic(&data.Topic{
		ID:       5,
		Name:     "JavaScript",
		Children: nil,
	})
	topicJobMarket, _ := ts.AddTopic(&data.Topic{
		ID:       6,
		Name:     "Job Market",
		Children: nil,
	})

	// Persons
	speakerDKrizic, _ := ps.AddPerson(&data.Person{
		ID:             1,
		Name:           "Darko Krizic",
		OrganizationID: organizationProdyna.ID,
	})
	speakerGGrujic, _ := ps.AddPerson(&data.Person{
		ID:             2,
		Name:           "Goran Grujic",
		OrganizationID: organizationProdyna.ID,
	})
	speakerMNikolic, _ := ps.AddPerson(&data.Person{
		ID:             3,
		Name:           "Milos Nikolic",
		OrganizationID: organizationProdyna.ID,
	})
	speakerAKoblin, _ := ps.AddPerson(&data.Person{
		ID:             4,
		Name:           "Aaron Koblin",
		OrganizationID: organizationGoogle.ID,
	})

	// Talks
	talkJavaSpringAndYou, _ := tlks.AddTalk(&data.Talk{
		ID:                1,
		Title:             "Java, Spring, and You",
		DurationInMinutes: 90,
		Language:          "English",
		Level:             data.BeginnerLevel,
		Persons:           []data.Person{*speakerDKrizic, *speakerGGrujic},
		Topics:            []data.Topic{*topicJava, *topicSpring, *topicHibernate},
		TalkDates:         nil,
	})
	talkFullStackJavaScriptOnKubernetes, _ := tlks.AddTalk(&data.Talk{
		ID:                2,
		Title:             "Fullstack JavaScript on Kubernetes",
		DurationInMinutes: 60,
		Language:          "Serbian",
		Level:             data.ExpertLevel,
		Persons:           []data.Person{*speakerMNikolic},
		Topics:            []data.Topic{*topicJavaScript, *topicKubernetes},
		TalkDates:         nil,
	})
	talkJavaForBeginners, _ := tlks.AddTalk(&data.Talk{
		ID:                3,
		Title:             "Java for Beginners",
		DurationInMinutes: 60,
		Language:          "English",
		Level:             data.BeginnerLevel,
		Persons:           []data.Person{*speakerGGrujic},
		Topics:            []data.Topic{*topicJava},
		TalkDates:         nil,
	})
	talkITJobMarketToday, _ := tlks.AddTalk(&data.Talk{
		ID:                4,
		Title:             "The IT Job Market Today",
		DurationInMinutes: 60,
		Language:          "English",
		Level:             data.BeginnerLevel,
		Persons:           []data.Person{*speakerAKoblin},
		Topics:            []data.Topic{*topicJava, *topicJobMarket},
		TalkDates:         nil,
	})

	// TalkDates
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         1,
		BeginDate:  time.Date(2021, time.Month(5), 12, 14, 0, 0, 0, time.UTC),
		TalkID:     talkJavaSpringAndYou.ID,
		RoomID:     roomRed.ID,
		EventID:    eventBestJavaConference.ID,
		LocationID: locationBelexpo.ID,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         2,
		BeginDate:  time.Date(2021, time.Month(5), 2, 10, 0, 0, 0, time.UTC),
		TalkID:     talkFullStackJavaScriptOnKubernetes.ID,
		RoomID:     roomWhite.ID,
		EventID:    eventProdynaJobFair.ID,
		LocationID: locationHotelPlaza.ID,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         3,
		BeginDate:  time.Date(2021, time.Month(5), 2, 13, 0, 0, 0, time.UTC),
		TalkID:     talkJavaForBeginners.ID,
		RoomID:     roomWhite.ID,
		EventID:    eventProdynaJobFair.ID,
		LocationID: locationHotelPlaza.ID,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         4,
		BeginDate:  time.Date(2021, time.Month(5), 10, 14, 0, 0, 0, time.UTC),
		TalkID:     talkJavaForBeginners.ID,
		RoomID:     roomBlue.ID,
		EventID:    eventITConnect.ID,
		LocationID: locationHotelPlaza.ID,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         5,
		BeginDate:  time.Date(2021, time.Month(6), 2, 15, 0, 0, 0, time.UTC),
		TalkID:     talkITJobMarketToday.ID,
		RoomID:     roomGoogle.ID,
		EventID:    eventGoogleIO.ID,
		LocationID: locationBelgradeFair.ID,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:         6,
		BeginDate:  time.Date(2021, time.Month(5), 10, 12, 0, 0, 0, time.UTC),
		TalkID:     talkITJobMarketToday.ID,
		RoomID:     roomBlue.ID,
		EventID:    eventITConnect.ID,
		LocationID: locationHotelPlaza.ID,
	})
}
