package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/callerobertsson/resty/cli"
)

var (
	Version string
	Build   string

	configFile = ""
	httpFile   = ""
	showVersion = false
)

func init() {

	flag.BoolVar(&showVersion, "v", false, "version information")
	flag.StringVar(&configFile, "c", "", "config file")

	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %q\nBuild: %v\n", Version, Build)
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		httpFile = flag.Args()[0]
	}
}

func main() {

	if httpFile == "" {
		fmt.Printf("Need a .httpFile as argument\n")
		os.Exit(1)
	}

	// TODO: Implement support of config files

	config := cli.Config{}
	if configFile != "" {
		maybeConfig, err := cli.ConfigFromFile(configFile)
		if err != nil {
			fmt.Printf("Error reading config file %q: %v\n", configFile, err)
			os.Exit(1)
		}
		config = maybeConfig
	}

	cli := cli.New(httpFile, config)

	cli.Start()
}
