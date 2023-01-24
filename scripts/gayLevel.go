package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	"time"
)

type GayNumbers myTypes.GayNumbers

func GayRanking(dm *discordgo.Message, score int) {

	var gayNumbers GayNumbers

	user := dm.Author.String()
	timeStamp := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&gayNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&GayNumbers{User: user})
	}
	sqldb.DB.First(&gayNumbers, "User = ?", user)

	// Update user
	updateGayScores(gayNumbers, timeStamp, score)

	return
}

func updateGayScores(gayNumbers GayNumbers, timeStamp time.Time, score int) {
	// Update averages
	tries := gayNumbers.NumberTries + 1
	average := ((gayNumbers.AverageScore * gayNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&gayNumbers).Update("date", timeStamp)
	sqldb.DB.Model(&gayNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&gayNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if gayNumbers.CurrentScore < score {
		sqldb.DB.Model(&gayNumbers).Update("current_score", score)
	}
}
