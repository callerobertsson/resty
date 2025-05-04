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
	envFile        = ""
	path           = ""
	showVersion    = false
	generateConfig = false
)

// main loads config and starts the CLI.
func main() {
	parseCommandLine()

	// Create Resty CLI from configFile (if any)
	cli, err := cli.NewFromConfigFile(configFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Run Resty CLI on directory or file in path using variables in envFile (if any)
	if err = cli.Run(path, envFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// parseCommandLine handles command line flags parsing.
func parseCommandLine() {
	flag.StringVar(&configFile, "c", "", "config file")
	flag.StringVar(&envFile, "e", "", "env file")
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
		path = flag.Args()[0]
	}
}

func usage() {
	// fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	out := flag.CommandLine.Output()

	fmt.Fprintf(out, `RESTY
    SYNOPSIS
          resty [options] [<.http-file>|<directory>]

    DESCRIPTON
          By default, resty will find all .env-files in current directory and let the user select
          one or none, the all .http-files will be listed. The user can then select which one to 
          open. 

          If there is an argument on the command line, resty will check if it is a directory and
          open it in directory mode. If the argument is a file, resty will open it in file mode.

          Resty opens <.http-file> and displays list of the requests. The list can be navigated
          using vim-like bindings or arrow keys.

          Requests can be run by pressing '<enter>' and the .http-file edited by pressing 'e'.

          A default config file can be printed by using the '-g' flag. Modify it and put it in
          $HOME/.resty.json.

          For more info press '?' to show the available commands.

    OPTIONS
`)

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(out, "%8s-%-4s%s (default: %q)\n", "", f.Name, f.Usage, f.DefValue)
	})
}
