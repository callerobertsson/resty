// Package cli config functions
package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/callerobertsson/resty/utils"
)

// DefaultConfigFileName is the default configuration file name.
const DefaultConfigFileName = ".resty.json"

// Config holds the settings.
type Config struct {
	configFile  string // Config file path, set by application
	CurlCommand string // Default "curl"
	Editor      string // Default $EDITOR
	ColorMode   bool   // TODO: implement - Default false, no color
	InsecureSSL bool   // Default false

	// TODO: Add config settings
	// - add formatter per header accept types

}

// ConfigFromReader constructs a Config from JSON data read from the Reader.
func ConfigFromJSONReader(r io.Reader) (*Config, error) {
	c := &Config{}

	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, c)
	if err != nil {
		return c, err
	}

	c.configFile = "<reader>"

	return c, nil
}

// ConfigFromFile uses ConfigFromReader to create a Config from the JSON file content.
func ConfigFromJSONFile(f string) (*Config, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(file)

	c, err := ConfigFromJSONReader(r)
	if err != nil {
		return nil, err
	}

	c.configFile = f

	return c, nil
}

func ConfigJSON(c Config) string {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return fmt.Sprintf("%v\n", string(bs))
}

// GetConfigOrDefault returns a Config constructed from the JSON file. If f is
// the empty string it tries to read the default config file. If that fails an
// empty Config, with default values, are returned.
func GetConfigOrDefault(f string) (*Config, error) {
	// Default config
	config := Config{
		CurlCommand: "curl",
		Editor:      os.Getenv("EDITOR"),
	}

	// If no config file
	if f == "" {
		df, err := resolveDefaultConfigFilePath()
		if err != nil {
			return &config, nil // ignore error
		}

		// Return default config, if no config file exists
		if !utils.FileExists(df) {
			return &config, nil
		}

		// Use default config file
		f = df
	}

	// Create config from file
	return ConfigFromJSONFile(f)
}

// resolveDefaultConfigFilePath returns a path to the default config file in
// the user home directory.
func resolveDefaultConfigFilePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v%c%v", u.HomeDir, os.PathSeparator, DefaultConfigFileName), nil
}
