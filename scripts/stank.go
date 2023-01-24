package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	"time"
)

type StankNumbers myTypes.StankNumbers

func StankRanking(dm *discordgo.Message, score int) {

	var stankNumbers StankNumbers

	user := dm.Author.String()
	timeStamp := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&stankNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&StankNumbers{User: user})
	}
	sqldb.DB.First(&stankNumbers, "User = ?", user)

	// Update user
	updateStankScores(stankNumbers, timeStamp, score)

	return
}

func updateStankScores(stankNumbers StankNumbers, timeStamp time.Time, score int) {
	// Update averages
	tries := stankNumbers.NumberTries + 1
	average := ((stankNumbers.AverageScore * stankNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&stankNumbers).Update("date", timeStamp)
	sqldb.DB.Model(&stankNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&stankNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if stankNumbers.CurrentScore < score {
		sqldb.DB.Model(&stankNumbers).Update("current_score", score)
	}
}
