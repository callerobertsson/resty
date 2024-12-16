// Package utils function for turning on and off buffered input.
package utils

import (
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

/*
   // disable input buffering
   exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
   // do not display entered characters on the screen
   exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
   // restore the echoing state when exiting
   defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
*/

var oldTermState *term.State

func SetUnbufferedInput() {
	if runtime.GOOS == "windows" {
		state, _ := term.MakeRaw(int(os.Stdin.Fd()))
		oldTermState = state
		return
	}

	// Disable input buffering
	_ = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// No character echo
	_ = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

// UnbufferedOff turns on normal buffering and character echo.
func SetBufferedInput() {
	if runtime.GOOS == "windows" {
		_ = term.Restore(int(os.Stdin.Fd()), oldTermState)
	} else {
		_ = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	}
}
