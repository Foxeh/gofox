package main

// This file adds the Disgord message route multiplexer, aka "command router".
// to the Disgord bot. This is an optional addition however it is included
// by default to demonstrate how to extend the Disgord bot.

import (
	"github.com/Foxeh/gofox/router"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = router.New()

func init() {
	// Register the mux OnMessageCreate handler that listens for and processes
	// all messages received.
	Discord.AddHandler(Router.OnMessageCreate)

	// Register the build-in help command.
	_, _ = Router.Route("help", "Display this message.", Router.Help)
	_, _ = Router.Route("ping", "pong!", Router.Ping)
}
