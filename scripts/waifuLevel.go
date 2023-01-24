package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	"time"
)

type WaifuNumbers myTypes.WaifuNumbers

func WaifuRanking(dm *discordgo.Message, score int) {

	var waifuNumbers WaifuNumbers

	user := dm.Author.String()
	timeStamp := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&waifuNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&WaifuNumbers{User: user})
	}
	sqldb.DB.First(&waifuNumbers, "User = ?", user)

	// Update user
	updateWaifuScores(waifuNumbers, timeStamp, score)

	return
}

func updateWaifuScores(waifuNumbers WaifuNumbers, timeStamp time.Time, score int) {
	// Update averages
	tries := waifuNumbers.NumberTries + 1
	average := ((waifuNumbers.AverageScore * waifuNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&waifuNumbers).Update("date", timeStamp)
	sqldb.DB.Model(&waifuNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&waifuNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if waifuNumbers.CurrentScore < score {
		sqldb.DB.Model(&waifuNumbers).Update("current_score", score)
	}
}
