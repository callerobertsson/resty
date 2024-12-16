package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Terminal color constants.
const (
	NORM     = "\033[0m"
	TITLE    = "\033[1;33m" // title - yellow bold
	SUBTITLE = "\033[37m"   // greyish
	SELECTED = "\033[1;32m" // selected - dark green
	NOTICE   = "\033[1;31m" // red, bold
)

func (cli *CLI) renderHeader() {
	conf := ", conf: none"
	if cli.config.configFile != "" {
		conf = " - " + cli.config.configFile
	}
	fmt.Printf(TITLE+"RESTY "+SUBTITLE+"%s%s\n\n"+NORM, cli.httpFile, conf)
}

func (cli *CLI) renderPrompt() {
	fmt.Printf(SUBTITLE + "\n[revcq?] > " + NORM)
}

func (cli *CLI) renderUI() {
	for i, r := range cli.dotHTTP.Requests {
		indicator := NORM + " "
		if i == cli.current {
			indicator = NOTICE + ">" + SELECTED
		}

		if r.Name == "" {
			r.Name = r.URL
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

func (cli *CLI) renderVariables() {
	if len(cli.dotHTTP.Requests[cli.current].Variables) > 0 {
		fmt.Printf(TITLE + "Variables:\n" + NORM)
		renderStringMap(cli.dotHTTP.Requests[cli.current].Variables, "  ")
	} else {
		fmt.Printf("No variables\n")
	}
}

func (cli *CLI) renderConfig() {
	fmt.Printf(TITLE+"Config in %s:\n"+NORM, cli.config.configFile)
	fmt.Printf("%s", ConfigJSON(*cli.config))
}

func renderStringMap(vars map[string]string, indent string) {
	for k, v := range vars {
		fmt.Printf("%s%s = %q\n", indent, k, v)
	}
}

func renderClear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		return
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func confirmMessage(f string, a ...any) bool {
	fmt.Printf(f, a...)
	fmt.Printf(SELECTED + "-- press 'y' to continue --\n" + NORM)
	bs := make([]byte, 1)
	_, _ = os.Stdin.Read(bs)

	return len(bs) > 0 && strings.ToLower(string(bs[0])) == "y"
}

func stopMessage(f string, a ...any) {
	fmt.Printf(f, a...)
	fmt.Printf(SELECTED + "-- press any key to continue --\n" + NORM)
	_, _ = os.Stdin.Read(make([]byte, 1))
}
