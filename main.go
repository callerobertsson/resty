// Package main implements command line handling and configuration before starting the Resty CLI.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/callerobertsson/resty/cli"
)

var (
	// Version is the application version
	Version string
	// Build contains the build date
	Build string

	// Command line flags parsed in init()
	configFile     = ""
	httpFile       = ""
	showVersion    = false
	generateConfig = false

	// Constants
	defaultConfigFileName = ".resty.json"
)

// init handles command line flags parsing
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

	if httpFile == "" {
		fmt.Printf("Need a .httpFile as argument\n")
		os.Exit(1)
	}
}

// main loads config and starts the CLI.
func main() {

	// Get configuration
	config, err := getConfigOrDefault(configFile)
	if err != nil {
		fmt.Printf("Error creating configuration: %s\n", err)
		os.Exit(1)
	}

	// Create CLI
	cli := cli.New(httpFile, config)

	if err := cli.Start(); err != nil {
		fmt.Printf("Error running Resty: %v\n", err)
	}
}

func getConfigOrDefault(f string) (*cli.Config, error) {

	// Default config
	config := cli.Config{}
	config.CurlCommand = "curl"
	config.Editor = os.Getenv("EDITOR")

	// If no config file
	if f == "" {
		df, err := resolveDefaultConfigFilePath()
		if err != nil {
			return &config, nil // ignore error
		}

		// Return default config, if no config file exists
		if !fileExists(f) {
			return &config, nil
		}

		// Use default config file
		f = df
	}

	// Let cli create config from file
	return cli.ConfigFromFile(f)
}

func resolveDefaultConfigFilePath() (string, error) {

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v%c%v", u.HomeDir, os.PathSeparator, defaultConfigFileName), nil
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}
