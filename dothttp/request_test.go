// Package dothttp_test Request unit tests.
package dothttp_test

import (
	"strings"
	"testing"

	"github.com/callerobertsson/resty/dothttp"
)

func TestBuildCurlArgs(t *testing.T) {
	vut := dothttp.Request{
		Name:      "Fake Request",
		Verb:      "VERB",
		URL:       "https://base.url/path",
		URLFormat: "https://{{baseurl}}/path",
		Variables: map[string]string{
			"@baseurl": "base.url",
			"@notused": "notused",
		},
		Headers: map[string]string{
			"hkey": "hval",
		},
		Body: "{\"a\": 123}",
	}

	expected := strings.Join([]string{
		"-k",
		"-X",
		"VERB",
		"https://base.url/path",
		"-H 'hkey: hval'",
		"-d '{\"a\": 123}'",
	}, ", ")

	got := strings.Join(vut.BuildCurlArgs(), ", ")

	if expected != got {
		t.Errorf("BuildCurlArgs() expected\n  %q\nbut got\n  %q", expected, got)
	}
}
