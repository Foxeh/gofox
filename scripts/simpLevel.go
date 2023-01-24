package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	"time"
)

type SimpNumbers myTypes.SimpNumbers

func SimpRanking(dm *discordgo.Message, score int) {

	var simpNumbers SimpNumbers

	user := dm.Author.String()
	timeStamp := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&simpNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&SimpNumbers{User: user})
	}
	sqldb.DB.First(&simpNumbers, "User = ?", user)

	// Update user
	updateSimpScores(simpNumbers, timeStamp, score)

	return
}

func updateSimpScores(simpNumbers SimpNumbers, timeStamp time.Time, score int) {
	// Update averages
	tries := simpNumbers.NumberTries + 1
	average := ((simpNumbers.AverageScore * simpNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&simpNumbers).Update("date", timeStamp)
	sqldb.DB.Model(&simpNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&simpNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if simpNumbers.CurrentScore < score {
		sqldb.DB.Model(&simpNumbers).Update("current_score", score)
	}
}
