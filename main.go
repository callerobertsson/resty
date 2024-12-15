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

	configFile     = ""
	httpFile       = ""
	showVersion    = false
	generateConfig = false
)

func init() {

	flag.StringVar(&configFile, "c", "", "config file")
	flag.BoolVar(&showVersion, "v", false, "version information")
	flag.BoolVar(&generateConfig, "g", false, "print default config file")

	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %q\nBuild: %v\n", Version, Build)
		os.Exit(0)
	}

	if generateConfig {
		fmt.Printf("%v\n", cli.ConfigJson(cli.Config{}))
		os.Exit(1)
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

	config := &cli.Config{}
	if configFile != "" {
		maybeConfig, err := cli.ConfigFromFile(configFile)
		if err != nil {
			fmt.Printf("Error reading config file %q: %v\n", configFile, err)
			os.Exit(1)
		}
		config = maybeConfig
	}

	if config.CurlCommand == "" {
		config.CurlCommand = "curl"
	}
	if config.Editor == "" {
		config.Editor = os.Getenv("EDITOR")
	}

	cli := cli.New(httpFile, config)

	if err := cli.Start(); err != nil {
		fmt.Printf("Error running Resty: %v\n", err)
	}
}
