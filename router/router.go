// Package router provides a simple Discord message route multiplexer that
// parses messages and then executes a matching registered handler, if found.
// Router can be used with both Disgord and the DiscordGo library.
package router

import (
	"github.com/Foxeh/gofox/log"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

func init() {
	// Call logger
	log.Init(os.Stdout, os.Stdout, os.Stderr)
}

// Route holds information about a specific message route handler
type Route struct {
	Pattern     string      // match pattern that should trigger this route handler
	Description string      // short description of this route
	Help        string      // detailed help string for this route
	Run         HandlerFunc // route handler function to call
}

// Context holds a bit of extra data we pass along to route handlers
// This way processing some of this only needs to happen once.
type Context struct {
	Fields          []string
	Content         string
	IsDirected      bool
	IsPrivate       bool
	HasPrefix       bool
	HasMention      bool
	HasMentionFirst bool
}

// HandlerFunc is the function signature required for a message route handler.
type HandlerFunc func(*discordgo.Session, *discordgo.Message, *Context)

// Router is the main struct for all Router methods.
type Router struct {
	Routes  []*Route
	Prefix  string
	Pattern string
}

// New returns a new Discord message route Router
func New() *Router {
	m := &Router{}
	m.Prefix = "pls"
	return m
}

// Route allows you to register a route
func (m *Router) Route(pattern, desc string, cb HandlerFunc) (*Route, error) {

	r := Route{}
	r.Pattern = pattern
	r.Description = desc
	r.Run = cb
	m.Routes = append(m.Routes, &r)

	return &r, nil
}

// FuzzyMatch attempts to find the best route match for a given message.
func (m *Router) FuzzyMatch(msg string) (*Route, []string) {
	// TODO: Change to be a bit more strict on matching
	// Tokenize the msg string into a slice of words
	fields := strings.Fields(msg)

	// no point to continue if there's no fields
	if len(fields) == 0 {
		return nil, nil
	}

	// Search though the command list for a match
	var r *Route
	var rank int

	var fk int
	for fk, fv := range fields {

		for _, rv := range m.Routes {

			// If we find an exact match, return that immediately.
			if rv.Pattern == fv {
				return rv, fields[fk:]
			}

			// Some "Fuzzy" searching...
			if strings.HasPrefix(rv.Pattern, fv) {
				if len(fv) > rank {
					r = rv
					rank = len(fv)
				}
			}
		}
	}
	return r, fields[fk:]
}

// OnMessageCreate is a DiscordGo Event Handler function.  This must be
// registered using the DiscordGo.Session.AddHandler function.  This function
// will receive all Discord messages and parse them for matches to registered
// routes.
func (m *Router) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {

	var err error

	// Ignore all messages created by the Bot account itself
	if mc.Author.ID == ds.State.User.ID {
		return
	}

	// Create Context struct that we can put various infos into
	ctx := &Context{
		Content: strings.TrimSpace(mc.Content),
	}

	// Fetch the channel for this Message
	var c *discordgo.Channel
	c, err = ds.State.Channel(mc.ChannelID)
	if err != nil {
		// Try fetching via REST API
		c, err = ds.Channel(mc.ChannelID)
		if err != nil {
			log.Warning.Printf("Unable to fetch Channel for Message, %s", err)
		} else {
			// Attempt to add this channel into our State
			err = ds.State.ChannelAdd(c)
			if err != nil {
				log.Warning.Printf("Error updating State with Channel, %s", err)
			}
		}
	}
	// Add Channel info into Context (if we successfully got the channel)
	if c != nil {
		if c.Type == discordgo.ChannelTypeDM {
			ctx.IsPrivate, ctx.IsDirected = true, true
		}
	}

	// Detect prefix mention
	if len(m.Prefix) > 0 {
		// TODO : Option to change prefix to support a per-guild user defined prefix
		if strings.HasPrefix(strings.ToLower(ctx.Content), m.Prefix) {
			ctx.IsDirected, ctx.HasPrefix, ctx.HasMentionFirst = true, true, true
			ctx.Content = strings.TrimPrefix(ctx.Content, m.Prefix)
		}
	}

	// For now, if we're not specifically mentioned we do nothing.
	if !ctx.IsDirected {
		return
	}

	// Try to find the "best match" command out of the message.
	r, fl := m.FuzzyMatch(ctx.Content)
	if r != nil {
		// TODO: Change to do something different when mentioned vs using prefix.
		// TODO: Possibly add rate limit?
		log.Info.Printf("Time: %s || Author: %s || Message: %s\n", mc.Timestamp, mc.Author, mc.Content)
		ctx.Fields = fl
		r.Run(ds, mc.Message, ctx)
		return
	}
}
