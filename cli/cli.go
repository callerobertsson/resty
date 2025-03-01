package cli

import (
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

	utils.ColorOff()
	if cli.config.ColorMode {
		utils.ColorOn()
	}

	return cli.commandLoop()
}

func (cli *CLI) StartDirectory(d string) error {

	utils.ColorOff()
	if cli.config.ColorMode {
		utils.ColorOn()
	}

	return cli.selectLoop(d)

}
