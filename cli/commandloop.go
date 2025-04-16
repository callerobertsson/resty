package cli

import (
	"fmt"
	"os"

	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) commandLoop() error {
	// Unbuffered input
	utils.SetUnbufferedInput()
	defer utils.SetBufferedInput()

	// Single byte buffer for keypresses
	buf := make([]byte, 1)

	// Command loop
	for {
		utils.RenderClear()

		cli.renderHeader()
		cli.renderUI()
		cli.renderPrompt()

		// Read input rune
		_, _ = os.Stdin.Read(buf)
		r := rune(buf[0])

		switch {
		case r == 'q' || r == 27:
			handleQuit()
			return nil
		case r == 'g':
			cli.handleFirst()
		case r == 'k' || r == 65:
			cli.handleUp()
		case r == 'j' || r == 66:
			cli.handleDown()
		case r == 'G':
			cli.handleLast()
		case r == 'e':
			cli.handleEdit()
		case r == 'c':
			cli.handleConfig()
		case r == 'v':
			cli.handleVariables()
		case r == 'r' || r == '\n':
			cli.handleRun()
		case r == 'R' || r == '\n':
			cli.handleRunAll()
		case r == '?':
			handleHelp()
		}
	}
}

func handleQuit() {
	fmt.Println("\nbye!")
	utils.SetBufferedInput()
	// os.Exit(0)
}

func (cli *CLI) handleFirst() {
	cli.current = 0
}

func (cli *CLI) handleUp() {
	cli.current--
	if cli.current < 0 {
		cli.current = 0
	}
}

func (cli *CLI) handleDown() {
	cli.current++
	if cli.current >= len(cli.dotHTTP.Requests) {
		cli.current--
	}
}

func (cli *CLI) handleLast() {
	cli.current = len(cli.dotHTTP.Requests) - 1
}

func (cli *CLI) handleEdit() {
	// Edit .http-file (until no errors)
	for {
		utils.RenderClear()
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
}

func (cli *CLI) handleConfig() {
	// Config - show config file
	utils.RenderClear()
	cli.renderConfig()
	stopMessage("\n")
}

func (cli *CLI) handleVariables() {
	// Variables - show variables active for current request
	utils.RenderClear()
	cli.renderVariables()
	stopMessage("\n")
}

func (cli *CLI) handleRun() {
	// Run current request
	utils.RenderClear()
	// cli.renderRequestInfo()
	err := cli.runCurrentRequest()
	if err != nil {
		stopMessage("Error: %v\n", err)
	}
}

func (cli *CLI) handleRunAll() {
	// Run current request
	utils.RenderClear()
	// cli.renderRequestInfo()
	err := cli.runAllRequests()
	if err != nil {
		stopMessage("Error: %v\n", err)
	}
}

func handleHelp() {
	// Show help
	utils.RenderClear()
	stopMessage("%v\n", keyHelpText)
}

var keyHelpText = `=== RESTY Keymap ===

Navigation

  Arrow keys or Vim (jk) for navigating endpoint list.

  g - Jump to first request
  G - Jump to last request

Commands

  r - Run currently selected request, <enter>
  R - Run all requests in current .http-file
  v - View variables set for current request
  e - Edit the input file using Editor config setting or $EDITOR environment variable
  c - Show config settings

  ? - Show this help
  q - Quit Resty
`
