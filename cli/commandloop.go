package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func (cli *CLI) commandLoop() error {

	// Unbuffered input
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	buf := make([]byte, 1)

	for {
		clear()
		fmt.Println("RESTY:")
		fmt.Println("? for help, q to quit")

		cli.render()

		// Read input rune
		os.Stdin.Read(buf)
		r := rune(buf[0])

		switch {
		case r == 'q':
			fmt.Println("\nbye!")
			exec.Command("reset").Run()
			os.Exit(0)

		case r == 'j':
			cli.current++
			if cli.current >= len(cli.dotHTTP.Requests) {
				cli.current--
			}

		case r == 'k':
			cli.current--
			if cli.current < 0 {
				cli.current = 0
			}

		case r == 'r':
			// Run current
			err := cli.runCurrentRequest()
			if err != nil {
				clear()
				fmt.Printf("Error: %v\n-- press any key to continue --\n", err)
				os.Stdin.Read(buf)
			}

		}
	}
}
