// Package cli config functions
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	"github.com/callerobertsson/resty/utils"
)

// Constants
const DefaultConfigFileName = ".resty.json"

// Config holds the settings.
type Config struct {
	CurlCommand string // Default "curl"
	Editor      string // Default $EDITOR
	ColorMode   bool   // Default false, no color

	// TODO:
	// - add formatter per header accept types

	configFile string // Config file path, set by application
}

// ConfigFromFile create a Config instance created from the file content. Default values
// will be set for curl command and editor.
func ConfigFromFile(f string) (*Config, error) {

	c := &Config{configFile: f}

	bs, err := os.ReadFile(f)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(bs, c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func ConfigJson(c Config) string {

	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return fmt.Sprintf("%v\n", string(bs))
}

func GetConfigOrDefault(f string) (*Config, error) {

	// Default config
	config := Config{}
	config.CurlCommand = "curl"
	config.Editor = os.Getenv("EDITOR")

	// If no config file
	if f == "" {
		df, err := resolveDefaultConfigFilePath()
		if err != nil {
			return &config, nil // ignore error
		}

		// Return default config, if no config file exists
		if !utils.FileExists(f) {
			return &config, nil
		}

		// Use default config file
		f = df
	}

	// Create config from file
	return ConfigFromFile(f)
}

func resolveDefaultConfigFilePath() (string, error) {

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v%c%v", u.HomeDir, os.PathSeparator, DefaultConfigFileName), nil
}
