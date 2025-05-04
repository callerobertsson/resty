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

	expectedSlice := []string{
		"-X",
		"VERB",
		"https://base.url/path",
		"-H",
		"hkey: hval",
		"-d",
		"{\"a\": 123}",
	}

	expectedSecure := strings.Join(expectedSlice, ", ")

	got := strings.Join(vut.BuildCurlArgs(false), ", ")

	if expectedSecure != got {
		t.Errorf("BuildCurlArgs() expected\n  %q\nbut got\n  %q", expectedSecure, got)
	}

	expectedUnsecure := strings.Join(append([]string{"-k"}, expectedSlice...), ", ")

	got = strings.Join(vut.BuildCurlArgs(true), ", ")

	if expectedUnsecure != got {
		t.Errorf("BuildCurlArgs() expected\n  %q\nbut got\n  %q", expectedUnsecure, got)
	}
}
