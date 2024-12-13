package dothttp

type DotHTTP struct {
	// TODO: implement

	Requests []Request
}

type Request struct {
	Name      string // ### @name The Name
	Verb      string // TODO: should this be an "enum"?
	URLFormat string // Url with possible {{tags}}
	ParsedURL string // Url with possible {{tags}}
	Variables map[string]string
	Headers   map[string]string
}

func New() *DotHTTP {
	return &DotHTTP{}
}

func NewFromFile(f string) (*DotHTTP, error) {
	dh := DotHTTP{}
	err := dh.LoadHTTPFile(f)
	return &dh, err
}

func (dh *DotHTTP) LoadHTTPFile(f string) error {
	// TODO: Build Request List

	// TODO: Remove dummy requests
	dh.Requests = []Request{
		{"Sunet", "GET", "https://www.sunet.se", "tbd", map[string]string{}, map[string]string{}},
		{"Funet", "DELETE", "https://www.funet.se/{{ep}}", "tbd", map[string]string{}, map[string]string{}},
	}

	return nil
}
