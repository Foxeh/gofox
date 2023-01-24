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

func (m *Router) Gayrate(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	rand.Seed(time.Now().UTC().UnixNano())
	min := 0
	max := 100
	score := min + rand.Intn(max-min)
	strscore := strconv.Itoa(score)

	var descrip string
	mentions := dm.Mentions
	if len(mentions) > 0 {
		stripMention := strings.Split(dm.Mentions[0].String(), "#")
		descrip = stripMention[0] + " is " + strscore + "% gay "
		dm.Author = dm.Mentions[0]
	} else {
		descrip = "You are " + strscore + "% gay "
	}

	if score > 50 {
		descrip = descrip + ":rainbow_flag:"
	} else {
		descrip = descrip + ":straight_ruler:"
	}

	embed := new(discordgo.MessageEmbed)
	(*embed).Title = "gay r8 machine"
	(*embed).Description = descrip

	_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, embed)
	log.ErrCheck("Failed to send message", err)

	scripts.GayRanking(dm, score)

	return
}
