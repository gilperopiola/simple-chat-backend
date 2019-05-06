package main

import (
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DatabaseActions interface {
	Setup()
	Migrate()
	Seed()
	Purge()
}

type MyDatabase struct {
	*gorm.DB
}

func (database *MyDatabase) Setup() {
	var err error
	database.DB, err = gorm.Open(config.DATABASE.TYPE, config.DATABASE.USERNAME+":"+config.DATABASE.PASSWORD+"@tcp("+config.DATABASE.HOSTNAME+":"+
		config.DATABASE.PORT+")/"+config.DATABASE.SCHEMA+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	database.Migrate()
}

func (database *MyDatabase) Migrate() {
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Room{})
	database.AutoMigrate(&Message{})
}

func (database *MyDatabase) Seed() {
	for i := 0; i < 10; i++ {
		database.Create(&User{Username: "test_username_" + strconv.Itoa(i), Password: hash("test_username_"+strconv.Itoa(i), "password")})
		database.Create(&Room{})
	}
}

func (database *MyDatabase) Purge() {
	database.Delete(User{})
	database.Delete(Room{})
	database.Delete(Message{})
}

func (database *MyDatabase) GetTestingData() (User, Room) {
	var users []User
	database.Find(&users)

	var rooms []Room
	database.Find(&rooms)

	return users[0], rooms[0]
}

func (database *MyDatabase) BeautifyError(err error) string {
	s := err.Error()

	if strings.Contains(s, "Duplicate entry") {
		duplicateField := strings.Split(s, "'")[3]
		return duplicateField + " already in use"
	}

	return s
}
