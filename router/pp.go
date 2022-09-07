package router

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/scripts"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
	"time"
)

func (m *Router) Pp(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	rand.Seed(time.Now().UTC().UnixNano())
	min := 0
	max := 30
	score := min + rand.Intn(max-min)

	var descrip string
	mentions := dm.Mentions
	if len(mentions) > 0 {
		stripMention := strings.Split(dm.Mentions[0].String(), "#")
		descrip = stripMention[0] + "'s penis\n"
		dm.Author = dm.Mentions[0]
	} else {
		stripAuthor := strings.Split(dm.Author.String(), "#")
		descrip = stripAuthor[0] + "'s penis\n"
	}

	pp := ""
	for i := 0; i < score; i++ {
		pp = pp + "="
	}

	descrip = descrip + "8" + pp + "D"

	embed := new(discordgo.MessageEmbed)
	(*embed).Title = "peepee size machine"
	(*embed).Description = descrip

	_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, embed)
	log.ErrCheck("Failed to send message", err)

	scripts.PpRanking(dm, score)

	return
}
