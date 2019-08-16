package router

import (
	"encoding/json"
	"fmt"
	"github.com/Foxeh/gofox/log"
	"github.com/bwmarrin/discordgo"
	"github.com/paked/configure"
	"github.com/subosito/shorturl"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	conf   = configure.New()
	wolfID = conf.String("wolframID", "", "Wolfram ID")
)

type WolframResponse struct {
	Queryresult struct {
		Pods []struct {
			Title   string `json:"title"`
			Subpods []struct {
				Plaintext string `json:"plaintext"`
			} `json:"subpods"`
		} `json:"pods"`
	} `json:"queryresult"`
}

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
	// Format query to URL
	content = url.QueryEscape(content)

	// Create URl to be sent to wolfram
	wolfURL := fmt.Sprintf("https://api.wolframalpha.com/v2/query?input=%s&appid=%s&output=JSON", content, *wolfID)

	// Get query
	res, err := http.Get(wolfURL)
	errCheck("Error querying wolfram alpha", err)

	// Format http.Get body
	data, err := ioutil.ReadAll(res.Body)
	errCheck("Failed to get request body", err)

	// Create tinyurl of query
	shortURL := urlShorten(content)

	// Format response
	result := "\n**Result:**\n"
	result += "```\n"
	result += getResult(data)
	result += "```\n"
	result += "**Full Result:**\n"
	result += "<" + shortURL + ">"

	_, _ = ds.ChannelMessageSend(dm.ChannelID, result)

	return
}

func getResult(data []byte) string {
	var pods WolframResponse
	err := json.Unmarshal(data, &pods)
	errCheck("Failed during json.Unmarshal", err)

	for i := range pods.Queryresult.Pods {
		if pods.Queryresult.Pods[i].Title == "Result" {
			res := pods.Queryresult.Pods[i].Subpods[0].Plaintext
			return res
		}
	}
	return "Result not found, go to link."
}

func urlShorten(content string) string {
	provider := "gggg"
	shortURL, err := shorturl.Shorten("https://www.wolframalpha.com/input/?i="+content, provider)
	errCheck("Error shortening URL", err)
	return string(shortURL)
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Error.Printf("%s: %+v", msg, err)
	}
}
