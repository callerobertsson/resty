// Package cli config functions
package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	CurlCommand string // Default "curl"
	Editor      string // Default $EDITOR
	ColorMode   bool   // Default false, no color
}

// ConfigFromFile create a Config instance created from the file content. Default values
// will be set for curl command and editor.
func ConfigFromFile(f string) (*Config, error) {
	// TODO: Read config file
	c := Config{}

	return &c, nil
}

func ConfigJson(c Config) string {

	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return fmt.Sprintf("%v\n", string(bs))
}
