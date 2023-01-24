package main

import (
	"github.com/Foxeh/gofox/log"
	"github.com/Foxeh/gofox/router"
	"github.com/Foxeh/gofox/sqldb"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paked/configure"
	"os"
)

// Version GoFox
const Version = "v0.6.0"

var (
	Router = router.New()
	conf   = configure.New()
	botKey = conf.String("botKey", "", "Bot key value")
	status = conf.String("status", "", "Discord status for bot")
)

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

	// Start DB
	sqldb.ConnectDB()

	// Set discord bot token
	Discord, err := discordgo.New("Bot " + *botKey)

	// Verify a Token was provided
	if Discord.Token == "" {
		log.Warning.Printf("You must provide a Discord authentication token.")
		return
	}

	// Open a websocket connection to Discord
	err = Discord.Open()
	log.ErrCheck("Error opening connection to Discord", err)
	defer func(Discord *discordgo.Session) {
		err := Discord.Close()
		if err != nil {

		}
	}(Discord)

	err = Discord.UpdateGameStatus(0, *status)
	log.ErrCheck("Error attempting to set my status", err)

	servers := Discord.State.Guilds

	log.Info.Printf("GoFox is running version: %s", Version)
	log.Info.Printf("GoFox has started on %d server(s)", len(servers))

	Discord.AddHandler(Router.OnMessageCreate)

	<-make(chan struct{})
}
