package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/callerobertsson/resty/dothttp"
	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) runAllRequests() error {
	rs := cli.dotHTTP.Requests

	for _, r := range rs {
		err := cli.runRequest(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cli *CLI) runCurrentRequest() error {
	r := cli.dotHTTP.Requests[cli.current]
	return cli.runRequest(r)
}

func (cli *CLI) runRequest(r dothttp.Request) error {
	args := r.BuildCurlArgs(cli.config.InsecureSSL)

	fmt.Println(utils.TITLE + "=== CURL ====================================================================" + utils.NORM)
	fmt.Printf("%s %s\n", cli.config.CurlCommand, strings.Join(args, " "))

	if r.Verb != "GET" && !confirmMessage("Are you sure?\n") {
		return nil
	}

	fmt.Println(utils.SUBTITLE + "Calling API..." + utils.NORM)

	resp, err := executeCommand(cli.config.CurlCommand, args...)
	if err != nil {
		return err
	}

	// Get format command and mime type if possible
	fmtCmd, mimeType := cli.getFormatterAndMimeType(r)

	var formatterErr error

	if fmtCmd != "" {
		var maybeResp string
		maybeResp, formatterErr = cli.formatResponse(resp, fmtCmd)
		if formatterErr == nil {
			resp = maybeResp
		}
	}

	fmt.Println(utils.TITLE + "=== Response ================================================================" + utils.NORM)
	fmt.Printf("%s\n", resp)
	fmt.Println(utils.TITLE + "=============================================================================" + utils.NORM)

	info := ", raw"
	if fmtCmd != "" && mimeType != "" {
		info = fmt.Sprintf(", formatted %q using %q", mimeType, fmtCmd)
	}

	fmt.Printf("%d chars%s\n", len(resp), info)

	if formatterErr != nil {
		fmt.Printf("Formatter failed: %v\n", formatterErr)
	}

	b := getByteMessage(utils.SUBTITLE + "Continue [y>] " + utils.NORM)
	if b == byte('>') {
		fmt.Printf(utils.SELECTED + "\nEnter file path: " + utils.NORM)
		reader := bufio.NewReader(os.Stdin)
		utils.SetBufferedInput()
		filePath, _ := reader.ReadString('\n')
		utils.SetUnbufferedInput()

		err = os.WriteFile(strings.TrimSpace(filePath), []byte(resp), 0666)
		if err != nil {
			fmt.Printf(utils.NOTICE+"Failed to write file: %v\n"+utils.NORM, err)
			stopMessage("\n")
		}
	}

	return nil
}

func (cli *CLI) getFormatterAndMimeType(r dothttp.Request) (string, string) {
	mimeType := r.Headers["accept"]

	if mimeType == "" {
		mimeType = "*"
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
		_, _ = io.WriteString(stdin, data)
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
