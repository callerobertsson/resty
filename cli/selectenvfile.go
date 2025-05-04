// Package cli file selection loop
package cli

import (
	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) selectEnvFile(d string) (string, error) {
	noFile := "No Environment file"

	envFiles := utils.GetEnvFilePaths(d, true)
	if len(envFiles) < 1 {
		return "", nil
	}

	title := cli.headerString()
	title += utils.TITLE + "Select Environment File:\n\n" + utils.NORM

	envFiles = append([]string{noFile}, envFiles...)

	prompt := "\nfuzzy"
	f, err := utils.FuzzyListPicker(title, prompt, envFiles)
	if err != nil {
		return "", err
	}

	if f == noFile {
		return "", nil
	}

	return f, nil
}
