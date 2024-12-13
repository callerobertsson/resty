package cli

import (
	"os"

	"github.com/callerobertsson/resty/dothttp"
)

type CLI struct {
	config   Config
	httpFile string
	dotHTTP  *dothttp.DotHTTP
	current  int
}

type Config struct {
	CurlCommand string // Default "curl"
	Editor      string // Default $EDITOR
}

func New(f string, c Config) *CLI {
	return &CLI{config: c, httpFile: f}
}

// ConfigFromFile create a Config instance created from the file content. Default values
// will be set for curl command and editor.
func ConfigFromFile(f string) (Config, error) {
	// TODO: Read config file
	c := Config{}

	if c.CurlCommand == "" {
		c.CurlCommand = "curl"
	}
	if c.Editor == "" {
		c.Editor = os.Getenv("EDITOR")
	}

	return Config{}, nil
}

func (cli *CLI) Start() error {

	// fmt.Printf("Resty: %#v\n", cli)

	maybeDH, err := dothttp.NewFromFile(cli.httpFile)
	if err != nil {
		return err
	}

	cli.dotHTTP = maybeDH

	return cli.commandLoop()
}
