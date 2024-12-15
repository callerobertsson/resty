package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (cli *CLI) renderUI() {
	for i, r := range cli.dotHTTP.Requests {
		indicator := " "
		if i == cli.current {
			indicator = ">"
		}

		name := r.Name
		if name == "" {
			name = r.URL
		}

		fmt.Printf(" %s %2d: %-6s %s\n", indicator, i+1, r.Verb, r.Name)
		indent := "       "
		if i == cli.current {
			fmt.Printf("%surl:    %q\n", indent, r.URL)
			fmt.Printf("%sformat: %q\n", indent, r.URLFormat)
			if len(r.Headers) > 0 {
				fmt.Printf("%sheaders:\n", indent)
				renderStringMap(r.Headers, indent+"  ")
			}
			if r.Body != "" {
				body := strings.Join(strings.Split(r.Body, "\n"), "\n"+indent+"  ")
				fmt.Printf("%sbody:\n%v\n", indent, indent+"  "+body)
			}
		}
	}
}

func (cli *CLI) renderRequestInfo() {
	r := cli.dotHTTP.Requests[cli.current]

	fmt.Printf("\n%s %s\n", r.Verb, r.URL)
	cli.renderVariables()
	cli.renderHeaders()
	fmt.Printf("\nBody:\n%s\n", r.Body)
	fmt.Printf("------------------------------------------\n")

}

func (cli *CLI) renderVariables() {
	if len(cli.dotHTTP.Requests[cli.current].Variables) > 0 {
		fmt.Printf("Variables:\n")
		renderStringMap(cli.dotHTTP.Requests[cli.current].Variables, "  ")
	} else {
		fmt.Printf("No variables\n")
	}
}

func (cli *CLI) renderHeaders() {
	if len(cli.dotHTTP.Requests[cli.current].Headers) > 0 {
		fmt.Printf("Headers:\n")
		renderStringMap(cli.dotHTTP.Requests[cli.current].Headers, "  ")
	} else {
		fmt.Printf("No headers\n")
	}
}

func renderStringMap(vars map[string]string, indent string) {
	for k, v := range vars {
		fmt.Printf("%s%s = %q\n", indent, k, v)
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func stopMessage(f string, a ...any) {
	fmt.Printf(f, a...)
	fmt.Printf("-- press any key to continue --\n")
	os.Stdin.Read(make([]byte, 1))
}
