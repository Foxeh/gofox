package main

import (
	"github.com/Foxeh/gofox/router"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = router.New()

func init() {
	// Register the router OnMessageCreate handler that listens for and processes
	// all messages received.
	Discord.AddHandler(Router.OnMessageCreate)

	// Commands to look for
	// Comment out commands to turn them off

	_, _ = Router.Route("ping", "pong!", Router.Ping)
	_, _ = Router.Route("help", "Display this message.", Router.Help)
	_, _ = Router.Route("pp", "Check pp size.", Router.Pp)
	_, _ = Router.Route("wolfram", "Query Wolfram Alpha", Router.Wolfram)
	_, _ = Router.Route("stankrate", "Check stank levels.", Router.Stankrate)

}
