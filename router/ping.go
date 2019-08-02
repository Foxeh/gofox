package router

import "github.com/bwmarrin/discordgo"

func (m *Router) Ping(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {

	resp := "Pong!"

	_, _ = ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
