// Package utils function for turning on and off buffered input.
package utils

import "os/exec"

func SetUnbufferedInput() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

// UnbufferedOff turn on nor
func SetBufferedInput() {
	exec.Command("reset").Run()
}
