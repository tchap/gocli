package main

import "github.com/tchap/gocli"

var subcmd = &gocli.Command{
	UsageLine: "subcmd [-v]",
	Short:     "some kind of subcommand, you name it",
	Long:      "Brb, too tired to write long descriptions.",
}

var verbose bool

func init() {
	subcmd.Flags.BoolVar(&verbose, "v", false, "print verbose output")
	app.MustRegisterSubcommand(subcmd)
}
