package router

import (
	"github.com/Foxeh/gofox/log"
	"github.com/bwmarrin/discordgo"
	"github.com/paked/configure"
	"net/url"
	"os"
	"strings"
)

var (
	conf   = configure.New()
	wolfID = conf.String("wolframID", "", "Wolfram ID")
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

func (m *Router) Wolfram(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {

	// Get content to be queried by wolfram alpha, removing prefix/command
	content := strings.Trim(dm.Content, m.Prefix+m.Pattern)
	// Format query for URL
	content = url.QueryEscape(content)

	// Create URl to be sent to wolfram

	// wolfURL := fmt.Sprintf("https://api.wolframalpha.com/v2/query?input=%s&appid=%s", content, *wolfID)

	// Get query

	// res, err := http.Get(wolfURL)
	// errCheck("Error querying wolfram alpha", err)

	// TODO: Do something with res

	_, _ = ds.ChannelMessageSend(dm.ChannelID, content)

	return
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Error.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
