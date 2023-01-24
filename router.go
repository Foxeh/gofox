package main

func init() {
	// Commands to look for
	// Comment out commands to turn them off

	_, _ = Router.Route("help", "Display this message.", Router.Help)
	_, _ = Router.Route("wolfram", "Query Wolfram Alpha", Router.Wolfram)
	_, _ = Router.Route("pp", "Check pp size.", Router.Pp)
	_, _ = Router.Route("simp", "Check simp levels.", Router.Simprate)
	_, _ = Router.Route("gayrate", "Check gay levels.", Router.Gayrate)
	_, _ = Router.Route("waifu", "Check waifu levels.", Router.Waifu)
	_, _ = Router.Route("stankrate", "Check stank levels.", Router.Stankrate)
	_, _ = Router.Route("epicgamer", "Check epicgamer levels.", Router.Epicgamer)

}
