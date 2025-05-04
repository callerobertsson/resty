package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/callerobertsson/resty/utils"
)

func (cli *CLI) renderHeader() {
	fmt.Printf("%s", cli.headerString())
}

func (cli *CLI) headerString() string {
	http := "no http"
	if cli.httpFile != "" {
		http = cli.httpFile
	}
	env := "no env"
	if cli.envFile != "" {
		env = cli.envFile
	}
	conf := "no conf"
	if cli.config.configFile != "" {
		conf = cli.config.configFile
	}

	return fmt.Sprintf(utils.TITLE+"RESTY "+utils.SUBTITLE+"%s - %s - %s\n\n"+utils.NORM, http, env, conf)
}

func (cli *CLI) renderPrompt() {
	fmt.Printf(utils.SUBTITLE + "\n[rRevcq?] > " + utils.NORM)
}

func (cli *CLI) renderUI() {
	for i, r := range cli.dotHTTP.Requests {
		indicator := utils.NORM + " "
		if i == cli.current {
			indicator = utils.NOTICE + ">" + utils.SELECTED
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
		fmt.Printf(utils.TITLE + "Variables:\n" + utils.NORM)
		renderStringMap(cli.dotHTTP.Requests[cli.current].Variables, "  ")
	} else {
		fmt.Printf("No variables\n")
	}
}

func (cli *CLI) renderConfig() {
	fmt.Printf(utils.TITLE+"Config in %s:\n"+utils.NORM, cli.config.configFile)
	fmt.Printf("%s", ConfigJSON(*cli.config))
}

func renderStringMap(vars map[string]string, indent string) {
	for k, v := range vars {
		fmt.Printf("%s%s = %q\n", indent, k, v)
	}
}

func confirmMessage(f string, a ...any) bool {
	fmt.Printf(f, a...)
	fmt.Printf(utils.SELECTED + "-- press 'y' to continue --\n" + utils.NORM)
	bs := make([]byte, 1)
	_, _ = os.Stdin.Read(bs)

	return len(bs) > 0 && strings.ToLower(string(bs[0])) == "y"
}

func stopMessage(f string, a ...any) {
	fmt.Printf(f, a...)
	fmt.Printf(utils.SELECTED + "-- press any key to continue --\n" + utils.NORM)
	_, _ = os.Stdin.Read(make([]byte, 1))
}

func getByteMessage(f string, a ...any) byte {
	fmt.Printf(f, a...)
	// fmt.Printf(utils.SELECTED + "-- press any key to continue --\n" + utils.NORM)
	bs := make([]byte, 1)
	_, _ = os.Stdin.Read(bs)

	return bs[0]
}
