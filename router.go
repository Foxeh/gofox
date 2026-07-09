package main

func init() {
	// Commands to look for
	// Comment out commands to turn them off

	Router.Route("help", "Display this message.", Router.Help)
	Router.Route("pp", "Check pp size.", Router.Pp)
	Router.Route("simp", "Check simp levels.", Router.Simprate)
	Router.Route("gayrate", "Check gay levels.", Router.Gayrate)
	Router.Route("waifu", "Check waifu levels.", Router.Waifu)
	Router.Route("stankrate", "Check stank levels.", Router.Stankrate)
	Router.Route("epicgamer", "Check epicgamer levels.", Router.Epicgamer)

}
