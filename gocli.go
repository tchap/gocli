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
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Command struct {
	UsageLine string
	Short     string
	Long      string

	Flags     *flag.FlagSet

	// All subcommands registered. Although it can be accessed directly,
	// it is better to use MustRegisterSubcommand to append new subcommands,
	// also because it automatically adds help subcommands.
	Subcmds []*Command

	// Remember the parent command during flag parsing so that we can generate
	// the help output from the whole chain of subcommands.
	parent  *Command
}

func NewCommand() *Command {
	var cmd Command
	cmd.Flags = flag.NewFlagSet("", flag.ExitOnError)
	cmd.Flags.Usage = cmd.Usage
	return &cmd
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	if i := strings.Index(name, " "); i >= 0 {
		name = name[:i]
	}
	return name
}

var usageTemplate = `COMMAND:
    {{.Name}} - {{.Short}}
USAGE:
    {{.UsageLine}}
{{with $opts := call .OptionsString}}
OPTIONS:
    {{print $opts}}
{{end}}
DESCRIPTION:
    {{.Long}}
{{with .Subcmds}}
SUBCOMMANDS:
    {{range .}}{{.Name}}{{end}}
{{end}}
`

func (cmd *Command) Usage() {
	t := template.Must(template.New("usage").Parse(usageTemplate))
	t.Execute(os.Stderr, cmd)
}

func (cmd *Command) invoke(args []string) {

}

func (cmd *Command) OptionsString() string {
	return "Options go here"
}

func (cmd *Command) Execute(args []string) {
	cmd.Flags.Parse(args)
	subArgs := cmd.Flags.Args()

	if len(subArgs) == 0 {
		cmd.invoke(subArgs)
	}

	name := subArgs[0]
	for _, subcmd := range cmd.Subcmds {
		if subcmd.Name() == name {
			subcmd.parent = cmd
			subcmd.Execute(subArgs[1:])
			return
		}
	}

	cmd.invoke(subArgs)
}

func (cmd *Command) MustRegisterSubcommand(subcmd *Command) {
	// Easy if there are no subcommands defined yet.
	if cmd.Subcmds == nil {
		cmd.Subcmds = []*Command{subcmd}
		return
	}

	// Check for subcommand name collisions.
	for _, c := range cmd.Subcmds {
		if c.Name() == subcmd.Name() {
			panic(fmt.Sprintf("Subcommand %s already defined", subcmd.Name()))
		}
	}

	// If everything went well, register the subcommand.
	cmd.Subcmds = append(cmd.Subcmds, subcmd)
}
