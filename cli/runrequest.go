package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (cli *CLI) runCurrentRequest() error {
	args := cli.dotHTTP.Requests[cli.current].BuildCurlArgs()

	fmt.Printf(TITLE + "=== CURL ======================================================================\n" + NORM)
	fmt.Printf("%s %s\n", cli.config.CurlCommand, strings.Join(args, " "))
	fmt.Printf(TITLE + "=== Response ==================================================================\n" + NORM)

	err := runProcess(cli.config.CurlCommand, args...)
	if err != nil {
		return err
	}

	stopMessage("\n")

	return nil
}

func runProcess(cmd string, args ...string) error {
	c, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	as := append([]string{c}, args...)

	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	p, err := os.StartProcess(c, as, &attr)
	if err != nil {
		return err
	}
	_, _ = p.Wait()

	return nil
}
