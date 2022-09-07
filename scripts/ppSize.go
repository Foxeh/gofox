package scripts

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/sqldb"
	myTypes "github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
)

type PpNumbers myTypes.PpNumbers

func PpRanking(dm *discordgo.Message, score int) {

	var ppNumbers PpNumbers

	//Only need this if taking scores from another bot
	//re := regexp.MustCompile(`[-]?\d[\d,]*[]?[\d{2}]*`)
	//strScore := re.FindString(dm.Embeds[0].Description)
	//score, err := strconv.Atoi(strScore)
	//log.ErrCheck("Failed to convert score to int", err)

	user := dm.Author.String()
	time := dm.Timestamp

	// Check if user exists
	if err := sqldb.DB.Where("User = ?", user).First(&ppNumbers).Error; err != nil {
		log.Info.Printf("User: " + user + " does not exists in db, creating profile")
		sqldb.DB.Create(&PpNumbers{User: user})
	}
	sqldb.DB.First(&ppNumbers, "User = ?", user)

	// Update user
	updatePpScores(ppNumbers, time, score)

	return
}

func updatePpScores(ppNumbers PpNumbers, time discordgo.Timestamp, score int) {
	// Update averages
	tries := ppNumbers.NumberTries + 1
	average := ((ppNumbers.AverageScore * ppNumbers.NumberTries) + score) / tries
	sqldb.DB.Model(&ppNumbers).Update("date", time)
	sqldb.DB.Model(&ppNumbers).Update("number_tries", tries)
	sqldb.DB.Model(&ppNumbers).Update("average_score", average)

	// Check if current score is better then previous
	if ppNumbers.CurrentScore < score {
		sqldb.DB.Model(&ppNumbers).Update("current_score", score)
	}
}
