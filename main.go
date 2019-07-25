package main

import (
	"github.com/Foxeh/gofox/log"
	"github.com/bwmarrin/discordgo"
	"github.com/paked/configure"
	"os"
)

var (
	botID  string
	conf   = configure.New()
	botKey = conf.String("botKey", "", "Bot key value")
)

// Current GoFox version
const Version = "v0.1.0-alpha"

func main() {
	// Call logger
	log.Init(os.Stdout, os.Stdout, os.Stderr)

	// Pull in configuration
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewEnvironment())
	if _, err := os.Stat("config.json"); err == nil {
		conf.Use(configure.NewJSONFromFile("config.json"))
	}
	conf.Parse()

	// Create new bot instance
	discord, err := discordgo.New(*botKey)
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	// Set BotID and command handler
	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Golang FoxBot")
		errCheck("Error attempting to set my status", err)

		servers := discord.State.Guilds
		log.Info.Printf("GoFox is running version: %s", Version)
		log.Info.Printf("GoFox has started on %d server(s)", len(servers))
	})

	// Open a websocket connection to Discord
	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	<-make(chan struct{})
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Error.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	content := message.Content

	if content == "!Test" {
		_, err := discord.ChannelMessageSend(message.ChannelID, "Testing...")
		errCheck("Error sending message.", err)
		log.Warning.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
	}
}
