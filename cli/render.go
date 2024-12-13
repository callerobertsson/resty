package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func (cli *CLI) render() {
	for i, r := range cli.dotHTTP.Requests {
		indicator := " "
		if i == cli.current {
			indicator = ">"
		}

		fmt.Printf(" %s %2d: %-6s %s (%s)\n", indicator, i+1, r.Verb, r.Name, r.URLFormat)
		if i == cli.current {
			fmt.Printf("       Parsed URL: %q\n", r.ParsedURL)
			fmt.Printf("       URL Format: %q\n", r.URLFormat)
			fmt.Printf("       >> headers <<\n")
			fmt.Printf("       >> more info for selected <<\n")
		}
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
