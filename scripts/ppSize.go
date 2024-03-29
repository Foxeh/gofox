package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	"time"
)

type PpNumbers myTypes.PpNumbers

func PpRanking(dm *discordgo.Message, score int) {

	var ppNumbers PpNumbers

	user := dm.Author.String()
	timeStamp := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&ppNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&PpNumbers{User: user})
	}
	sqldb.DB.First(&ppNumbers, "User = ?", user)

	// Update user
	updatePpScores(ppNumbers, timeStamp, score)

	return
}

func updatePpScores(ppNumbers PpNumbers, timeStamp time.Time, score int) {
	// Update averages
	tries := ppNumbers.NumberTries + 1
	average := ((ppNumbers.AverageScore * ppNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&ppNumbers).Update("date", timeStamp)
	sqldb.DB.Model(&ppNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&ppNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if ppNumbers.CurrentScore < score {
		sqldb.DB.Model(&ppNumbers).Update("current_score", score)
	}
}
