/*
    The MIT License (MIT)

    Copyright (c) 2013 Ondřej Kupka

    Permission is hereby granted, free of charge, to any person obtaining a copy of
    this software and associated documentation files (the "Software"), to deal in
    the Software without restriction, including without limitation the rights to
    use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
    the Software, and to permit persons to whom the Software is furnished to do so,
    subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
    FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
    COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
    IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
    CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package gocli

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	UsageLine string
	Short     string
	Long      string

	Flags flag.FlagSet
	Args  []string

	parent  *Command
	subcmds []*Command
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	if i := strings.Index(name, " "); i >= 0 {
		name = name[:i]
	}
	return name
}

func (cmd *Command) Usage(msg ...string) {
	// Print message if supplied
	if len(msg) == 1 {
		fmt.Fprintf(os.Stderr, "%s\n\n", msg[0])
	}

	// Print usage line
	var usage = cmd.UsageLine
	for p := cmd.parent; p != nil; p = p.parent {
		usage = p.UsageLine + " " + usage
	}
	fmt.Fprintf(os.Stderr, "Usage: %s\n\n", usage)

	// Print long
	fmt.Fprintln(os.Stderr, strings.TrimSpace(cmd.Long))

	// Print options
	var opts bytes.Buffer
	cmd.Flags.SetOutput(&opts)
	cmd.Flags.PrintDefaults()
	if len(opts.String()) != 0 {
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, opts)
	}

	// Print options for the supercommands as well
	for p := cmd.parent; p != nil; p = p.parent {
		fmt.Fprintf(os.Stderr, "Options (supercommand '%s'):", p.Name())
		p.Flags.PrintDefaults()
	}

	// Print subcommands
	if cmd.subcmds != nil {
		fmt.Fprintf(os.Stderr, "\nSubcommands:\n")
		for _, subcmd := range cmd.subcmds {
			fmt.Fprintf(os.Stderr, "  %s \t %s\n", subcmd.Name(), subcmd.Short)
		}
	}
}

func (cmd *Command) Invoke() {

}

func (cmd *Command) ParseAndInvoke(args []string) {
	err := cmd.Flags.Parse(args)
	if err != nil {
		cmd.Usage(err.Error())
	}

	fArgs := cmd.Flags.Args()

	if cmd.Args != nil {
		if len(cmd.Args) != len(fArgs) {
			cmd.Usage("Invalid number of arguments")
		}
		for i := range cmd.Args {
			cmd.Args[i] = fArgs[i]
		}

		cmd.Invoke()
		return
	}

	if len(fArgs) == 0 {
		cmd.Invoke()
		return
	}

	if cmd.subcmds == nil {
		cmd.Usage("No subcommand defined: " + fArgs[0])
	}

	name := fArgs[0]
	for _, subcmd := range cmd.subcmds {
		if subcmd.Name() == name {
			subcmd.parent = cmd
			subcmd.ParseAndInvoke(fArgs[1:])
			return
		}
	}

	cmd.Usage("No subcommand defined: " + fArgs[0])
}

func (cmd *Command) MustRegisterSubcommand(subcmd *Command) {
	if cmd.subcmds == nil {
		cmd.subcmds = []*Command{subcmd}
		return
	}

	for _, c := range cmd.subcmds {
		if c.Name() == subcmd.Name() {
			panic(fmt.Sprintf("Subcommand %s already defined", subcmd.Name()))
		}
	}

	cmd.subcmds = append(cmd.subcmds, subcmd)
}
