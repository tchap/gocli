package main

import (
	"fmt"
	"os"

	"github.com/tchap/gocli"
)

var app = gocli.NewApp("myapp")
func init() {
	app.Short = "my bloody gocli app"
	app.Version = "1.2.3"
	app.Long = `
  This is a long description of my super uber cool app.`
}

func main() {
	app.Run(os.Args[1:])
	if verbose {
		fmt.Println("subcmd verbose mode ON")
	}
}
