package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	"time"
)

type EpicgamerNumbers myTypes.EpicgamerNumbers

func EpicRanking(dm *discordgo.Message, score int) {

	var epicgamerNumbers EpicgamerNumbers

	user := dm.Author.String()
	timeStamp := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&epicgamerNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&EpicgamerNumbers{User: user})
	}
	sqldb.DB.First(&epicgamerNumbers, "User = ?", user)

	// Update user
	updateEpicScores(epicgamerNumbers, timeStamp, score)

	return
}

func updateEpicScores(epicgamerNumbers EpicgamerNumbers, timeStamp time.Time, score int) {
	// Update averages
	tries := epicgamerNumbers.NumberTries + 1
	average := ((epicgamerNumbers.AverageScore * epicgamerNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&epicgamerNumbers).Update("date", timeStamp)
	sqldb.DB.Model(&epicgamerNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&epicgamerNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if epicgamerNumbers.CurrentScore < score {
		sqldb.DB.Model(&epicgamerNumbers).Update("current_score", score)
	}
}
