package router

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/scripts"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (m *Router) Simprate(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	rand.Seed(time.Now().UTC().UnixNano())
	min := 0
	max := 100
	score := min + rand.Intn(max-min)
	strscore := strconv.Itoa(score)

	var descrip string
	mentions := dm.Mentions
	if len(mentions) > 0 {
		stripMention := strings.Split(dm.Mentions[0].String(), "#")
		descrip = stripMention[0] + " is " + strscore + "% a simp "
		dm.Author = dm.Mentions[0]
	} else {
		descrip = "You are " + strscore + "% a simp "
	}

	if score < 30 {
		descrip = descrip + ":sunglasses:"
	} else if 30 < score && score < 70 {
		descrip = descrip + ":face_with_raised_eyebrow:"
	} else {
		descrip = descrip + ":joy:"
	}

	embed := new(discordgo.MessageEmbed)
	(*embed).Title = "simp r8 machine"
	(*embed).Description = descrip

	_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, embed)
	log.ErrCheck("Failed to send message", err)

	scripts.SimpRanking(dm, score)

	return
}
