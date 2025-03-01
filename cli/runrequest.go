package cli

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/callerobertsson/resty/dothttp"
	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) runCurrentRequest() error {
	r := cli.dotHTTP.Requests[cli.current]
	args := r.BuildCurlArgs(cli.config.InsecureSSL)

	fmt.Println(utils.TITLE + "=== CURL ======================================================================" + utils.NORM)
	fmt.Printf("%s %s\n", cli.config.CurlCommand, strings.Join(args, " "))

	if r.Verb != "GET" && !confirmMessage("Are you sure?\n") {
		return nil
	}

	bs, err := executeCommand(cli.config.CurlCommand, args...)
	if err != nil {
		return err
	}

	// Get format command and mime type if possible
	fmtCmd, mimeType := cli.getFormatterAndMimeType(r)

	resp := string(bs)
	var formatterErr error

	if fmtCmd != "" {
		maybeResp, err := cli.formatResponse(resp, fmtCmd)
		if err != nil {
			formatterErr = err
		} else {
			resp = maybeResp
		}
	}

	fmt.Println(utils.TITLE + "=== Response ==================================================================" + utils.NORM)
	fmt.Printf("%s\n", resp)
	fmt.Println(utils.TITLE + "===============================================================================" + utils.NORM)

	info := ", raw"
	if fmtCmd != "" && mimeType != "" {
		info = fmt.Sprintf(", formatted %q using %q", mimeType, fmtCmd)
	}

	fmt.Printf("%d chars%s\n", len(resp), info)

	if formatterErr != nil {
		fmt.Printf("Formatter failed: %v\n", formatterErr)
	}

	stopMessage("\n") // TODO: maybe ('>' for) saving to file?

	return nil
}

func (cli *CLI) getFormatterAndMimeType(r dothttp.Request) (string, string) {
	mimeType, ok := r.Headers["accept"]
	if !ok {
		return "", ""
	}

	formatter, ok := cli.config.Formatters[mimeType]
	if !ok {
		return "", mimeType
	}

	return formatter, mimeType
}

func (cli *CLI) formatResponse(s string, fmtCmd string) (string, error) {

	fs, err := executeStdinCommand(s, fmtCmd)
	if err != nil {
		return "", err
	}

	return fs, nil
}

func executeStdinCommand(data, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, data)
	}()

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func executeCommand(cmd string, args ...string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
