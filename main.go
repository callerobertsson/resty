// Package main implements command line handling and configuration before starting the Resty CLI.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/callerobertsson/resty/cli"
)

var (
	// Version is the application version.
	Version string
	// Build contains the build date.
	Build string

	// Command line flags parsed in parseCommandLine().
	configFile     = ""
	httpFile       = ""
	showVersion    = false
	generateConfig = false
)

// main loads config and starts the CLI.
func main() {
	parseCommandLine()

	// Get configuration
	config, err := cli.GetConfigOrDefault(configFile)
	if err != nil {
		fmt.Printf("Error creating configuration: %s\n", err)
		os.Exit(1)
	}

	// Create CLI
	cli := cli.New(httpFile, config)

	if err = cli.Start(); err != nil {
		fmt.Printf("Error running Resty: %v\n", err)
	}
}

// parseCommandLine handles command line flags parsing.
func parseCommandLine() {
	flag.StringVar(&configFile, "c", "", "config file")
	flag.BoolVar(&showVersion, "v", false, "version information")
	flag.BoolVar(&generateConfig, "g", false, "print default config file")

	flag.Usage = usage

	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %q\nBuild: %v\n", Version, Build)
		os.Exit(0)
	}

	if generateConfig {
		fmt.Printf("%v\n", cli.ConfigJSON(cli.Config{}))
		os.Exit(1)
	}

	if len(flag.Args()) > 0 {
		httpFile = flag.Args()[0]
	}

	if httpFile == "" {
		fmt.Printf("Error: Need a .httpFile as argument\n")
		fmt.Printf("Use -h for help\n")
		os.Exit(1)
	}
}

func usage() {
	// fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	out := flag.CommandLine.Output()

	fmt.Fprintf(out, `RESTY
    SYNOPSIS
          resty [flags] <.http-file>

    DESCRIPTON
          Resty opens the <.http-file> and displays list of the requests.
          The list can be navigated using vim-like bindings or arrow keys.

          Requests can be run and the .http-file edited.

          For more info press '?' to show the available commands.

    OPTIONS
`)

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(out, "%8s-%-4s%s (default: %q)\n", "", f.Name, f.Usage, f.DefValue)
	})
}
