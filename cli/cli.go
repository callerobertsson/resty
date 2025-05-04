package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/callerobertsson/resty/dothttp"
	"github.com/callerobertsson/resty/utils"
)

type CLI struct {
	config   *Config
	httpFile string
	envFile  string
	env      map[string]string
	dotHTTP  *dothttp.DotHTTP
	current  int
}

func NewFromConfigFile(configFile string) (*CLI, error) {
	cli := &CLI{}

	// Get configuration
	config, err := getConfigOrDefault(configFile)
	if err != nil {
		return nil, fmt.Errorf("error creating configuration: %w", err)
	}

	cli.config = config

	return cli, nil
}

// Run starts the Resty CLI.
// Parameter path can be a directory or a file path or be emtpy. Defaults to current dir.
// Parameter envFile can contain a path to an environment file or be empty.
func (cli *CLI) Run(path, envFile string) error {
	if path == "" {
		path = "."
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %q does not exist", path)
	}

	cli.envFile = envFile

	switch {
	case fi.IsDir():
		if err = cli.startDirectory(path); err != nil {
			return err
		}
	default:
		if err = cli.startFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (cli *CLI) startFile(f string) error {
	cli.httpFile = f

	if cli.envFile != "" {
		env, err := readEnvFile(cli.envFile)
		if err != nil {
			return fmt.Errorf("could not read env file %q: %w", cli.envFile, err)
		}
		cli.env = env
	}

	maybeDH, err := dothttp.NewFromFile(cli.httpFile, cli.env)
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

func (cli *CLI) startDirectory(d string) error {
	utils.ColorOff()
	if cli.config.ColorMode {
		utils.ColorOn()
	}

	if cli.envFile == "" {
		envFile, err := cli.selectEnvFile(d)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		fmt.Println("got an env:", envFile)

		cli.envFile = envFile

		if envFile != "" {
			cli.envFile = envFile
			cli.env, err = readEnvFile(cli.envFile)
			if err != nil {
				return err
			}
		}
	}

	return cli.selectHTTPFileLoop(d)
}

func readEnvFile(f string) (map[string]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return map[string]string{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	env := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		kv := strings.Split(line, "=")

		if len(kv) < 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(strings.Join(kv[1:], "="))

		env[key] = val
	}

	return env, nil
}
