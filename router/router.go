// Package router provides a simple Discord message route multiplexer for the
// DiscordGo library that parses messages and then executes a matching
// registered handler, if found.
package router

import (
	"fmt"
	"github.com/Foxeh/gofox/log"
	"github.com/bwmarrin/discordgo"
	"strings"
)

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
	Fields     []string
	Content    string
	IsDirected bool
	IsPrivate  bool
	HasPrefix  bool
}

// HandlerFunc is the function signature required for a message route handler.
type HandlerFunc func(*discordgo.Session, *discordgo.Message, *Context)

// Router is the main struct for all Router methods.
type Router struct {
	Routes []*Route
	Prefix string
}

// New returns a new Discord message route Router
func New() *Router {
	m := &Router{}
	m.Prefix = "pls"
	return m
}

// Route registers a command pattern with the handler to run for it.
func (m *Router) Route(pattern, desc string, cb HandlerFunc) {
	m.Routes = append(m.Routes, &Route{
		Pattern:     pattern,
		Description: desc,
		Run:         cb,
	})
}

// FuzzyMatch attempts to find the best route match for a given message.
// Matching is case-insensitive and runs in tiers of decreasing strictness
// against the first word of the message:
//
//  1. an exact match of a command name
//  2. a prefix that unambiguously identifies one command ("stank" -> "stankrate")
//  3. the unique closest command within a small edit distance ("hlep" -> "help")
//
// If the first word matches nothing, the remaining words are scanned for an
// exact command name, so commands can still be mixed into a sentence.
// The returned fields are the words following the matched command.
func (m *Router) FuzzyMatch(msg string) (*Route, []string) {
	return m.match(strings.Fields(strings.ToLower(msg)))
}

// match implements FuzzyMatch on an already-tokenized, lowercased message.
func (m *Router) match(fields []string) (*Route, []string) {
	if len(fields) == 0 {
		return nil, nil
	}

	cmd, args := fields[0], fields[1:]

	if r := m.exactMatch(cmd); r != nil {
		return r, args
	}

	if pm := m.prefixMatches(cmd); len(pm) == 1 {
		return pm[0], args
	}

	if r := m.closestMatch(cmd); r != nil {
		return r, args
	}

	// Fallback: a command name written exactly somewhere later in the message.
	for i, fv := range args {
		if r := m.exactMatch(fv); r != nil {
			return r, args[i+1:]
		}
	}

	return nil, nil
}

// Suggest returns command patterns the given word plausibly meant, for use in
// an "unknown command" reply: commands the word abbreviates, or failing that,
// commands within a small edit distance of it.
func (m *Router) Suggest(word string) []string {
	var out []string
	for _, r := range m.prefixMatches(word) {
		out = append(out, r.Pattern)
	}
	if len(out) > 0 {
		return out
	}
	for _, r := range m.Routes {
		if d := editDistance(word, r.Pattern); d <= 2 && d < len(word) {
			out = append(out, r.Pattern)
		}
	}
	return out
}

func (m *Router) exactMatch(word string) *Route {
	for _, r := range m.Routes {
		if r.Pattern == word {
			return r
		}
	}
	return nil
}

func (m *Router) prefixMatches(word string) []*Route {
	var matches []*Route
	for _, r := range m.Routes {
		if strings.HasPrefix(r.Pattern, word) {
			matches = append(matches, r)
		}
	}
	return matches
}

// closestMatch returns the route whose pattern is closest to word by edit
// distance, when that match is both close enough to be a likely typo and
// unambiguous. Words shorter than three characters are skipped, since nearly
// anything is within one edit of them.
func (m *Router) closestMatch(word string) *Route {
	if len(word) < 3 {
		return nil
	}
	maxDist := 1
	if len(word) > 5 {
		maxDist = 2
	}

	var best *Route
	bestDist := maxDist + 1
	unique := true
	for _, r := range m.Routes {
		switch d := editDistance(word, r.Pattern); {
		case d < bestDist:
			best, bestDist, unique = r, d, true
		case d == bestDist:
			unique = false
		}
	}
	if best == nil || !unique {
		return nil
	}
	return best
}

// editDistance returns the optimal string alignment distance between a and b:
// the minimum number of single-character insertions, deletions, substitutions
// and adjacent transpositions needed to turn one into the other. Command
// names are plain ASCII, so this operates on bytes.
func editDistance(a, b string) int {
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
		d[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		d[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+1)
			}
		}
	}
	return d[len(a)][len(b)]
}

// OnMessageCreate is a DiscordGo Event Handler function.  This must be
// registered using the DiscordGo.Session.AddHandler function.  This function
// will receive all Discord messages and parse them for matches to registered
// routes.
func (m *Router) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {

	// Ignore our own messages and those of any other bot or webhook, so two
	// bots can never trigger each other back and forth.
	if mc.Author.ID == ds.State.User.ID || mc.Author.Bot {
		return
	}

	// Create Context struct that we can put various infos into
	ctx := &Context{
		Content: strings.TrimSpace(mc.Content),
	}

	// A message with no guild is a direct message, no channel lookup needed.
	if mc.GuildID == "" {
		ctx.IsPrivate, ctx.IsDirected = true, true
	}

	// Detect prefix mention
	if len(m.Prefix) > 0 {
		// TODO : Option to change prefix to support a per-guild user defined prefix
		if strings.HasPrefix(strings.ToLower(ctx.Content), m.Prefix) {
			ctx.IsDirected, ctx.HasPrefix = true, true
			// Slice rather than TrimPrefix so a mixed-case prefix ("Pls") is
			// stripped too.
			ctx.Content = strings.TrimSpace(ctx.Content[len(m.Prefix):])
		}
	}

	// For now, if we're not specifically mentioned we do nothing.
	if !ctx.IsDirected {
		return
	}

	// Try to find the "best match" command out of the message.
	fields := strings.Fields(strings.ToLower(ctx.Content))
	r, fl := m.match(fields)
	if r != nil {
		// TODO: Change to do something different when mentioned vs using prefix.
		// TODO: Possibly add rate limit?
		log.Info.Printf("Time: %s || Author: %s || Message: %s\n", mc.Timestamp, mc.Author, mc.Content)
		ctx.Fields = fl
		r.Run(ds, mc.Message, ctx)
		return
	}

	// The message was directed at the bot but matched nothing; if the first
	// word was close to a known command, suggest it. Otherwise stay silent
	// so casual prefix use doesn't draw replies.
	if len(fields) == 0 {
		return
	}
	sugg := m.Suggest(fields[0])
	if len(sugg) == 0 {
		return
	}
	reply := fmt.Sprintf("Unknown command %q. Did you mean: %s? Try `%s help`.",
		fields[0], strings.Join(sugg, ", "), m.Prefix)
	_, _ = ds.ChannelMessageSend(mc.ChannelID, reply)
}
