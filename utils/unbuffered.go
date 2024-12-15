// Package utils function for turning on and off buffered input.
package utils

import "os/exec"

func SetUnbufferedInput() {
	_ = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

// UnbufferedOff turn on normal buffering.
func SetBufferedInput() {
	_ = exec.Command("reset").Run()
}
