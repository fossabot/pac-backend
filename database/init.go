package database

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/milutindzunic/pac-backend/data"
	"time"
)

func Init(db *gorm.DB, ls data.LocationStore, es data.EventStore, os data.OrganizationStore, ps data.PersonStore, rs data.RoomStore, ts data.TopicStore, tlks data.TalkStore, tlkds data.TalkDateStore, logger hclog.Logger) {

	logger.Info("Dropping all Tables...")
	// Drop the Entity tables
	db.DropTableIfExists(data.Event{})
	db.DropTableIfExists(data.Location{})
	db.DropTableIfExists(data.Organization{})
	db.DropTableIfExists(data.Person{})
	db.DropTableIfExists(data.Room{})
	db.DropTableIfExists(data.Talk{})
	db.DropTableIfExists(data.Topic{})
	db.DropTableIfExists(data.TalkDate{})
	// Drop the junction tables
	db.DropTableIfExists("is_child_of")
	db.DropTableIfExists("talk_topic")
	db.DropTableIfExists("talks_at")

	logger.Info("Recreating all Tables...")
	autoMigrate(db)

	logger.Info("Initializing DB with initial data...")
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
		ID:        1,
		Name:      "Best Java Conference",
		BeginDate: time.Date(2021, time.Month(5), 12, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.Month(5), 14, 0, 0, 0, 0, time.UTC),
		Location:  locationBelexpo,
	})
	eventProdynaJobFair, _ := es.AddEvent(&data.Event{
		ID:        2,
		Name:      "Prodyna Job Fair",
		BeginDate: time.Date(2021, time.Month(5), 2, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.Month(5), 5, 0, 0, 0, 0, time.UTC),
		Location:  locationHotelPlaza,
	})
	eventITConnect, _ := es.AddEvent(&data.Event{
		ID:        3,
		Name:      "IT Connect",
		BeginDate: time.Date(2021, time.Month(5), 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.Month(5), 12, 0, 0, 0, 0, time.UTC),
		Location:  locationHotelPlaza,
	})
	/*eventCloudnativeConference*/ _, _ = es.AddEvent(&data.Event{
		ID:        4,
		Name:      "Cloud Native Conference",
		BeginDate: time.Date(2021, time.Month(5), 22, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.Month(5), 23, 0, 0, 0, 0, time.UTC),
		Location:  locationBelgradeFair,
	})
	eventGoogleIO, _ := es.AddEvent(&data.Event{
		ID:        5,
		Name:      "Google I/O",
		BeginDate: time.Date(2021, time.Month(6), 2, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.Month(6), 5, 0, 0, 0, 0, time.UTC),
		Location:  locationBelgradeFair,
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
		ID:           1,
		Name:         "Red Room",
		Organization: organizationProdyna,
	})
	roomWhite, _ := rs.AddRoom(&data.Room{
		ID:           2,
		Name:         "White Room",
		Organization: organizationProdyna,
	})
	roomBlue, _ := rs.AddRoom(&data.Room{
		ID:           3,
		Name:         "Blue Room",
		Organization: organizationProdyna,
	})
	roomGoogle, _ := rs.AddRoom(&data.Room{
		ID:           4,
		Name:         "Google Room",
		Organization: organizationGoogle,
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
		ID:           1,
		Name:         "Darko Krizic",
		Organization: organizationProdyna,
	})
	speakerGGrujic, _ := ps.AddPerson(&data.Person{
		ID:           2,
		Name:         "Goran Grujic",
		Organization: organizationProdyna,
	})
	speakerMNikolic, _ := ps.AddPerson(&data.Person{
		ID:           3,
		Name:         "Milos Nikolic",
		Organization: organizationProdyna,
	})
	speakerAKoblin, _ := ps.AddPerson(&data.Person{
		ID:           4,
		Name:         "Aaron Koblin",
		Organization: organizationGoogle,
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
		ID:        1,
		BeginDate: time.Date(2021, time.Month(5), 12, 14, 0, 0, 0, time.UTC),
		Talk:      talkJavaSpringAndYou,
		Room:      roomRed,
		Event:     eventBestJavaConference,
		Location:  locationBelexpo,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:        2,
		BeginDate: time.Date(2021, time.Month(5), 2, 10, 0, 0, 0, time.UTC),
		Talk:      talkFullStackJavaScriptOnKubernetes,
		Room:      roomWhite,
		Event:     eventProdynaJobFair,
		Location:  locationHotelPlaza,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:        3,
		BeginDate: time.Date(2021, time.Month(5), 2, 13, 0, 0, 0, time.UTC),
		Talk:      talkJavaForBeginners,
		Room:      roomWhite,
		Event:     eventProdynaJobFair,
		Location:  locationHotelPlaza,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:        4,
		BeginDate: time.Date(2021, time.Month(5), 10, 14, 0, 0, 0, time.UTC),
		Talk:      talkJavaForBeginners,
		Room:      roomBlue,
		Event:     eventITConnect,
		Location:  locationHotelPlaza,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:        5,
		BeginDate: time.Date(2021, time.Month(6), 2, 15, 0, 0, 0, time.UTC),
		Talk:      talkITJobMarketToday,
		Room:      roomGoogle,
		Event:     eventGoogleIO,
		Location:  locationBelgradeFair,
	})
	_, _ = tlkds.AddTalkDate(&data.TalkDate{
		ID:        6,
		BeginDate: time.Date(2021, time.Month(5), 10, 12, 0, 0, 0, time.UTC),
		Talk:      talkITJobMarketToday,
		Room:      roomBlue,
		Event:     eventITConnect,
		Location:  locationHotelPlaza,
	})
}
