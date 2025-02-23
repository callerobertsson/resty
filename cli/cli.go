package cli

import (
	"fmt"
	"os"

	"github.com/callerobertsson/resty/dothttp"
	"github.com/callerobertsson/resty/utils"
)

type CLI struct {
	config   *Config
	httpFile string
	rootDir  string
	dotHTTP  *dothttp.DotHTTP
	current  int
}

func New(c *Config) *CLI {
	return &CLI{config: c}
}

func (cli *CLI) StartFile(f string) error {
	cli.httpFile = f

	maybeDH, err := dothttp.NewFromFile(cli.httpFile)
	if err != nil {
		return err
	}

	cli.dotHTTP = maybeDH

	colorOff()
	if cli.config.ColorMode {
		colorOn()
	}

	return cli.commandLoop()
}
func (cli *CLI) StartDirectory(d string) error {
	// Loop until quit
	for {
		// Find all .http-files recursively
		httpFiles := utils.GetHttpFilePaths(d, true)
		if len(httpFiles) < 1 {
			fmt.Printf("No .http-files found in %v\n", d)
			os.Exit(1)
		}

		title := "Resty File Selector\n(up, down to navigate, ctrl-c to exit)\n"
		prompt := "fuzzy find file"
		f, err := utils.FuzzyListPicker(title, prompt, httpFiles)
		if err != nil {
			return err
		}

		if f == "" {
			fmt.Println("bye")
			return nil
		}

		if err := cli.StartFile(f); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		utils.RenderClear()
	}
}
