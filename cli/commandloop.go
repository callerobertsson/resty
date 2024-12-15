package cli

import (
	"fmt"
	"os"

	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) commandLoop() error {

	// Unbuffered input
	utils.SetUnbufferedInput()

	buf := make([]byte, 1)

	for {
		clear()
		conf := ", conf: none"
		if cli.config.configFile != "" {
			conf = ", conf: " + cli.config.configFile
		}
		fmt.Printf("RESTY\nfile: %s%s\n\n", cli.httpFile, conf)

		cli.renderUI()
		fmt.Printf("\n[revcq?] > ")

		// Read input rune
		os.Stdin.Read(buf)
		r := rune(buf[0])

		switch {
		case r == 'q':
			// Quit
			fmt.Println("\nbye!")
			utils.SetBufferedInput()
			os.Exit(0)

		case r == 'j' || r == 66:
			// Go down
			cli.current++
			if cli.current >= len(cli.dotHTTP.Requests) {
				cli.current--
			}

		case r == 'k' || r == 65:
			// Go up
			cli.current--
			if cli.current < 0 {
				cli.current = 0
			}

		case r == 'g':
			// Go to top
			cli.current = 0

		case r == 'G':
			// Go to bottom
			cli.current = len(cli.dotHTTP.Requests) - 1

		case r == 'e':
			// Edit .http-file (until no errors)
			for {
				clear()
				_, err := utils.EditFile(cli.httpFile, cli.config.Editor)
				if err != nil {
					stopMessage("Error editing %v: %v\n", cli.httpFile, err)
				}

				// Reload .http-file
				err = cli.dotHTTP.LoadHTTPFile(cli.httpFile)
				if err != nil {
					stopMessage("Error loading %v: %v\n", cli.httpFile, err)
					continue
				}

				cli.current = 0 // if something is deleted

				break
			}

		case r == 'c':
			// Config - show config file
			clear()
			fmt.Printf("Config in %s:\n%s", cli.config.configFile, ConfigJson(*cli.config))
			stopMessage("\n")

		case r == 'v':
			// Variables - show variables active for current request
			clear()
			cli.renderVariables()
			stopMessage("\n")

		case r == 'r' || r == '\n':
			// Run current request
			clear()
			// cli.renderRequestInfo()
			err := cli.runCurrentRequest()
			if err != nil {
				stopMessage("Error: %v\n", err)
			}

		case r == '?':
			// Show help
			clear()
			stopMessage("%v\n", keyHelpText)
		}
	}
}

var keyHelpText = `=== RESTY Keymap ===

Navigation

  Arrow keys or Vim (jk) for navigating endpoint list.

  g - Jump to first request
  G - Jump to last request

Commands

  r - Run currently selected request, <enter>
  v - View variables set for current request
  e - Edit the input file using Editor config setting or $EDITOR environment variable
  c - Show config settings

  ? - Show this help
  q - Quit Resty
`
