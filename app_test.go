package gocli

import (
	"fmt"
)

// Examples -------------------------------------------------------------------

func ExampleApp() {
	// Create the root App object.
	app := NewApp("myapp")
	app.Short = "my bloody gocli app"
	app.Version = "1.2.3"
	app.Long = `
  This is a long description of my super uber cool app.`

	// Register a subcommand.
	var verbose bool

	var subcmd = &Command{
		UsageLine: "subcmd [-v]",
		Short:     "some kind of subcommand, you name it",
		Long:      "Brb, too tired to write long descriptions.",
		Action:    func(cmd *Command, args []string) {
			fmt.Printf("verbose mode set to %t\n", verbose)
		},
	}

	// Set up flags. This can be as well called in init() or so.
	subcmd.Flags.BoolVar(&verbose, "v", false, "print verbose output")

	// Register with the parent command. Also suitable for init().
	app.MustRegisterSubcommand(subcmd)

	// Run the whole thing.
	app.Run([]string{"subcmd"})
	app.Run([]string{"subcmd", "-v"})
	// Output:
	// verbose mode set to false
	// verbose mode set to true
}
