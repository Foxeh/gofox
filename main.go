package main

import (
	"github.com/Foxeh/gofox/log"
	"github.com/bwmarrin/discordgo"
	"github.com/paked/configure"
	"os"
)

var (
	conf         = configure.New()
	Discord, err = discordgo.New()
	botKey       = conf.String("botKey", "", "Bot key value")
)

// Current GoFox version
const Version = "v0.4.0-alpha"

func init() {
	// Pull in configuration
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewEnvironment())
	if _, err := os.Stat("config.json"); err == nil {
		conf.Use(configure.NewJSONFromFile("config.json"))
	}
	conf.Parse()
}

func main() {
	// Set bot token
	Discord.Token = *botKey

	// Open a websocket connection to Discord
	err = Discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer Discord.Close()

	err = Discord.UpdateStatus(0, "Golang FoxBot")
	errCheck("Error attempting to set my status", err)

	servers := Discord.State.Guilds
	log.Info.Printf("GoFox is running version: %s", Version)
	log.Info.Printf("GoFox has started on %d server(s)", len(servers))

	<-make(chan struct{})
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Error.Printf("%s: %+v", msg, err)
	}
}
