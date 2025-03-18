package dothttp

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/callerobertsson/resty/utils"
)

// See: https://learn.microsoft.com/en-us/aspnet/core/test/http-files?view=aspnetcore-9.0

type DotHTTP struct {
	Requests []Request
	Env      map[string]string
}

// Supported Verbs.
// All possible: OPTIONS GET HEAD POST PUT PATCH DELETE TRACE CONNECT.
var Verbs = []string{"GET", "PUT", "POST", "DELETE"}

func New() *DotHTTP {
	return &DotHTTP{}
}

func NewFromFile(f string, env map[string]string) (*DotHTTP, error) {
	dotHTTP := DotHTTP{Env: env}
	err := dotHTTP.LoadHTTPFile(f)
	return &dotHTTP, err
}

func (dotHTTP *DotHTTP) LoadHTTPFile(f string) error {
	// Read file line by line
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Add extra comment to avoid extra checking
	lines = append(lines, "# EOF")

	return dotHTTP.LoadHTTPFileLines(lines)
}

func (dotHTTP *DotHTTP) LoadHTTPFileLines(lines []string) error {
	rs := []Request{}

	vars := map[string]string{}

	// Add env file vars
	for k, v := range dotHTTP.Env {
		vars[k] = v
	}

	currRequest := Request{Headers: map[string]string{}}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		switch {
		case strings.HasPrefix(line, "#"):
			// Comment - ignore (@name comments handled below)
			continue

		case strings.HasPrefix(line, "@"):
			// Variable - add to vars
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				return fmt.Errorf("variable definition on line %d is invalid", i+1)
			}

			vars[strings.TrimSpace(parts[0])] = replaceVars(strings.TrimSpace(parts[1]), vars)

		case hasVerbPrefix(line):
			// Verb - check if current request already created
			if currRequest.Verb != "" {
				// rs = append(rs, currRequest)                        // store previous request
				currRequest = Request{Headers: map[string]string{}} // create new request
			}

			parts := strings.Split(line, " ")
			if len(parts) < 2 {
				// TODO: Handle HTTP version as part 3
				return fmt.Errorf("request on line %d has no url", i+1)
			}

			currRequest.Verb = strings.TrimSpace(parts[0])
			currRequest.URLFormat = strings.TrimSpace(parts[1])
			currRequest.URL = replaceVars(currRequest.URLFormat, vars)
			currRequest.Variables = copyMap(vars) // vars up to current request

			if i > 0 {
				currRequest.Name = parseName(lines[i-1])
			}
			if currRequest.Name == "" {
				currRequest.Name = currRequest.URL
			}

			// Read header values (if any)
			i++
			newI, headers := readHeaders(i, lines)
			currRequest.Headers = headers

			if newI > len(lines)-1 {
				rs = append(rs, currRequest) // store previous request
				break
			}

			i = newI // update index to after header

			// Read body data after the empty line (if any) but
			// only if next line isn't a comment or variable
			if strings.TrimSpace(lines[i]) == "" && !strings.HasPrefix(lines[i+1], "#") && !strings.HasPrefix(lines[i+1], "@") {
				newNewI, body := readBody(i+1, lines)
				i = newNewI + 1 // for-loop i++ adjustment

				currRequest.Body = replaceVars(strings.TrimSpace(body), vars)
			}

			rs = append(rs, currRequest) // store previous request
		}
	}

	dotHTTP.Requests = rs

	return nil
}

func readHeaders(i int, lines []string) (int, map[string]string) {
	hs := map[string]string{}

	for ; i < len(lines); i++ {
		if !hasHeaderValue(lines[i]) {
			return i, hs
		}

		parts := strings.Split(lines[i], ":")
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.ToLower(strings.TrimSpace(parts[1]))

		hs[key] = val
	}

	return i, hs
}

func readBody(i int, lines []string) (int, string) {
	body := ""

	if i > len(lines) {
		return len(lines), ""
	}

	for ; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "#") {
			return i - 1, ""
		}

		if strings.TrimSpace(lines[i]) == "" {
			return i - 1, body
		}
		body += lines[i] + "\n"
	}

	return i - 1, body
}

func replaceVars(s string, vars map[string]string) string {
	// {{var}}
	r := s

	for k, v := range vars {
		k = strings.TrimLeft(k, "@")
		r = strings.ReplaceAll(r, "{{"+k+"}}", v)
	}

	return r
}

func parseName(s string) string {
	// ### @name <NAME_TEXT>

	parts := strings.Split(s, "@name")
	if len(parts) < 2 {
		return s
	}

	return strings.TrimSpace(parts[1])
}

func copyMap(in map[string]string) map[string]string {
	out := map[string]string{}

	for k, v := range in {
		out[k] = v
	}

	return out
}

func hasVerbPrefix(s string) bool {
	verb := strings.Split(s, " ")[0]

	return utils.ContainsString(Verbs, verb)
}

func hasHeaderValue(s string) bool {
	return len(strings.Split(s, ":")) == 2
}
