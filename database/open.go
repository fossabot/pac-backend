package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"  //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite database driver
	"github.com/milutindzunic/pac-backend/config"
	"github.com/milutindzunic/pac-backend/data"
	"log"
)

func OpenDB(cnf *config.Config) (*gorm.DB, error) {
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

	db, err := gorm.Open(cnf.DbDriver, dbUrl)
	if err != nil {
		return nil, err
	}

	db.SingularTable(true)
	db.LogMode(cnf.LogPersistence)

	// Keep the schema up to date
	db = autoMigrate(db)

	return db, nil
}

func autoMigrate(db *gorm.DB) *gorm.DB {

	db = db.AutoMigrate(&data.Location{})
	db = db.AutoMigrate(&data.Event{})
	db = db.AutoMigrate(&data.Organization{})
	db = db.AutoMigrate(&data.Person{})
	db = db.AutoMigrate(&data.Room{})
	db = db.AutoMigrate(&data.Topic{})
	db = db.AutoMigrate(&data.Talk{})
	db = db.AutoMigrate(&data.TalkDate{})

	return db
}
