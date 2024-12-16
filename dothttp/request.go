package dothttp

type Request struct {
	Name      string            // ### @name The Name
	Verb      string            // The Verb, GET, PUT, POST, ...
	URL       string            // Parsed URL
	URLFormat string            // Url with possible {{tags}}
	Variables map[string]string // Variables available for Request
	Headers   map[string]string // Headers for Request
	Body      string            // Request body data
}

func (r Request) BuildCurlArgs(insecureSSL bool) []string {
	// TODO: Silent? - args := []string{"-s"} // silent
	args := []string{}

	// TODO: add curl config args

	// Insecure
	if insecureSSL {
		args = append(args, "-k") // Insecure
	}

	// Verb (-X) and URL
	args = append(args, []string{"-X", r.Verb, r.URL}...)

	// Headers (-H)
	for k, v := range r.Headers {
		args = append(args, "-H '"+k+": "+v+"'")
	}

	// Body (-d)
	if r.Body != "" {
		args = append(args, "-d '"+r.Body+"'")
	}

	// // URL
	// args = append(args, r.URL)

	return args
}
