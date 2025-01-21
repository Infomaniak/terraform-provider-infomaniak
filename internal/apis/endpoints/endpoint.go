package endpoints

import (
	"fmt"
	"net/url"
	"regexp"
)

var (
	pathParamRegexp = regexp.MustCompile("{([^}]*)}")
)

// Endpoint represents an Infomaniak's API Endpoint
type Endpoint struct {
	Method string
	Route  string
}

// NewEndpoint is a helper to define a new API Endpoint
func NewEndpoint(method string, route string) *Endpoint {
	return &Endpoint{
		Method: method,
		Route:  route,
	}
}

// CompiledEndpoint represents an Infomaniak's API Endpoint with params & query values.
type CompiledEndpoint struct {
	Endpoint *Endpoint
	URL      string
}

// QueryValues holds key value pairs of query values
type QueryValues map[string]any

// Encode encodes the QueryValues into a string to append to the url
func (q QueryValues) Encode() string {
	values := url.Values{}
	for k, v := range q {
		values.Set(k, fmt.Sprint(v))
	}
	return values.Encode()
}

// Compile compiles an Endpoint to a CompiledEndpoint with the given url params & query values
func (e *Endpoint) Compile(values QueryValues, params ...any) (*CompiledEndpoint, error) {
	path := e.Route

	matches := pathParamRegexp.FindAllStringSubmatchIndex(e.Route, -1)
	if len(matches) != len(params) {
		return nil, fmt.Errorf("should specify the same amount of param than the route needs")
	}

	var offset int
	for i, match := range matches {
		matchStart := match[0] - offset
		matchEnd := match[1] - offset
		paramValue := fmt.Sprint(params[i])

		offset += (matchEnd - matchStart) - len(paramValue)

		path = path[:matchStart] + paramValue + path[matchEnd:]
	}

	query := values.Encode()
	if query != "" {
		query = "?" + query
	}

	return &CompiledEndpoint{
		Endpoint: e,
		URL:      path + query,
	}, nil
}
