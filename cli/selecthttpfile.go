// Package cli file selection loop
package cli

import (
	"fmt"
	"os"

	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) selectHTTPFileLoop(d string) error {
	for {
		httpFiles := utils.GetHTTPFilePaths(d, true)
		if len(httpFiles) < 1 {
			fmt.Printf("No .http-files found in %v\n", d)
			os.Exit(1)
		}

		title := cli.headerString()

		prompt := "\nfuzzy"
		f, err := utils.FuzzyListPicker(title, prompt, httpFiles)
		if err != nil {
			return err
		}

		if f == "" {
			fmt.Println("bye")
			return nil
		}

		if err = cli.StartFile(f); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		utils.RenderClear()
	}
}
