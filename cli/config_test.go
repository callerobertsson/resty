package cli_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/callerobertsson/resty/cli"
)

func TestConfigFromJSONReader(t *testing.T) {
	expectedConfig := cli.Config{
		CurlCommand: "curlcommand",
		Editor:      "editor",
		ColorMode:   true,
		InsecureSSL: false,
	}

	bs, _ := json.MarshalIndent(expectedConfig, "", "  ")
	expectedJSON := string(bs)

	jsonData := `
{
  "CurlCommand": "curlcommand",
  "Editor": "editor",
  "ColorMode": true,
  "InsecureSSL": false
}
`
	r := strings.NewReader(jsonData)

	gotConfig, _ := cli.ConfigFromJSONReader(r)

	bs, _ = json.MarshalIndent(gotConfig, "", "  ")
	gotJSON := string(bs)

	if expectedJSON != gotJSON {
		t.Errorf("Expected Config to be:\n%s\nbut in was:\n%s\n", expectedJSON, gotJSON)
	}
}
