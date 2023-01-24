package sqldb

import (
	"github.com/Foxeh/gofox/log"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type StankNumbers struct {
	gorm.Model
	Date         time.Time
	User         string
	CurrentScore int
	NumberOfWins int
	LastWinDate  time.Time
	NumberTries  int
	AverageScore int
}

type PpNumbers struct {
	gorm.Model
	Date         time.Time
	User         string
	CurrentScore int
	NumberOfWins int
	LastWinDate  time.Time
	NumberTries  int
	AverageScore int
}

type GayNumbers struct {
	gorm.Model
	Date         time.Time
	User         string
	CurrentScore int
	NumberOfWins int
	LastWinDate  time.Time
	NumberTries  int
	AverageScore int
}

type EpicgamerNumbers struct {
	gorm.Model
	Date         time.Time
	User         string
	CurrentScore int
	NumberOfWins int
	LastWinDate  time.Time
	NumberTries  int
	AverageScore int
}

type SimpNumbers struct {
	gorm.Model
	Date         time.Time
	User         string
	CurrentScore int
	NumberOfWins int
	LastWinDate  time.Time
	NumberTries  int
	AverageScore int
}

type WaifuNumbers struct {
	gorm.Model
	Date         time.Time
	User         string
	CurrentScore int
	NumberOfWins int
	LastWinDate  time.Time
	NumberTries  int
	AverageScore int
}

var DB *gorm.DB

func ConnectDB() {

	// Setup DB
	db, err := gorm.Open(sqlite.Open("gofox.db"), &gorm.Config{})
	log.ErrCheck("Failed to start DB.", err)

	// Create Tables
	err = db.AutoMigrate(&StankNumbers{})
	log.ErrCheck("Failed to create tables", err)

	err = db.AutoMigrate(&PpNumbers{})
	log.ErrCheck("Failed to create tables", err)

	log.Info.Printf("Tables created.")

	DB = db

}
