// Package dothttp unit tests for loading .http-file data.
package dothttp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLoadValidHTTPLines_OneRequest(t *testing.T) {

	var httpFileContent = `### Fake .http File

@var1=var1value

### @name Request 1
GET https://host.com/{{var1}}

`

	expectedDotHTTP := DotHTTP{
		Requests: []Request{
			Request{
				Name:      "Request 1",
				Verb:      "GET",
				URL:       "https://host.com/var1value",
				URLFormat: "https://host.com/{{var1}}",
				Variables: map[string]string{
					"@var1": "var1value",
				},
				Headers: map[string]string{},
				Body:    "",
			},
		},
	}

	bs, _ := json.MarshalIndent(expectedDotHTTP, "", "  ")
	expectedJson := string(bs)

	vut := &DotHTTP{}

	vut.loadHTTPFileLines(strings.Split(httpFileContent, "\n"))

	bs, _ = json.MarshalIndent(vut, "", "  ")
	json := string(bs)

	if expectedJson != string(json) {
		t.Errorf("Expected DotHTTP to be:\n%s\nbut in was:\n%s\n", expectedJson, json)
	}
}

func TestLoadValidHTTPLines_MultipleRequests(t *testing.T) {

	var httpFileContent = `### Fake .http File

@var1=var1value

### @name Request 1
GET https://host.com/{{var1}}

@var2=var2value

### @name Request 2
PUT https://host.com/{{var2}}
accept: application/json

{
  "data": "foo"
}

# this will be ignored
### @name Request 3
#DELETE https://host.com/stuff/123
#accept: text/html

`

	expectedDotHTTP := DotHTTP{
		Requests: []Request{
			Request{
				Name:      "Request 1",
				Verb:      "GET",
				URL:       "https://host.com/var1value",
				URLFormat: "https://host.com/{{var1}}",
				Variables: map[string]string{
					"@var1": "var1value",
					// "@var2": "var2value",
				},
				Headers: map[string]string{},
				Body:    "",
			},
			Request{
				Name:      "Request 2",
				Verb:      "PUT",
				URL:       "https://host.com/var2value",
				URLFormat: "https://host.com/{{var2}}",
				Variables: map[string]string{
					"@var1": "var1value",
					"@var2": "var2value",
				},
				Headers: map[string]string{
					"accept": "application/json",
				},
				Body: "{\n  \"data\": \"foo\"\n}",
			},
		},
	}

	bs, _ := json.MarshalIndent(expectedDotHTTP, "", "  ")
	expectedJson := string(bs)

	vut := &DotHTTP{}

	vut.loadHTTPFileLines(strings.Split(httpFileContent, "\n"))

	bs, _ = json.MarshalIndent(vut, "", "  ")
	json := string(bs)

	if string(expectedJson) != json {
		t.Errorf("Expected DotHTTP to be:\n%s\nbut in was:\n%s\n", expectedJson, json)
	}
}
