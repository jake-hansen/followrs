package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/jake-hansen/followrs/config"
	"github.com/jake-hansen/followrs/server"
)

// main parses any given flags and starts the server.
func main() {
	var environment *string = flag.String("e", "dev", "environment to run in")
	flag.Usage = func() {
		fmt.Println("Usage: serve -e {environment}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)

	title := "FOLLOWRS"
	printTitle(title)
	fmt.Println("running in environment " + *environment)

	server.Init(*environment, time.Now())
}

// printTitle prints the given title to the screen in ASCII art format.
func printTitle(title string) {
	figure.NewFigure(title, "larry3d", false).Print()
}
