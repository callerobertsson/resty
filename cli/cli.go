package cli

import (
	"github.com/callerobertsson/resty/dothttp"
)

type CLI struct {
	config   *Config
	httpFile string
	dotHTTP  *dothttp.DotHTTP
	current  int
}

func New(f string, c *Config) *CLI {
	return &CLI{config: c, httpFile: f}
}

func (cli *CLI) Start() error {
	maybeDH, err := dothttp.NewFromFile(cli.httpFile)
	if err != nil {
		return err
	}

	cli.dotHTTP = maybeDH

	return cli.commandLoop()
}
