// Package phtml implements shared functionality for generated code.
package phtml

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewlinker/protohtml-go/internal/httppattern"
	"google.golang.org/protobuf/proto"
)

// ValuesDecoder should allow for decoding values.
type ValuesDecoder interface {
	Decode(v any, values url.Values) (err error)
}

// ValuesEncoder should allow for encoding values.
type ValuesEncoder interface {
	Encode(v any) (values url.Values, err error)
}

// PathReader allows for reading parameter from the
// url Path. This is impelmented by *http.Request.
type PathReader interface {
	PathValue(name string) string
}

// ParseRequest parses a requests's req.PostForm into a protobuf message 'msg'.
func ParseRequest(
	dec ValuesDecoder,
	msg proto.Message,
	req *http.Request,
	pathParamsNames ...string,
) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	// @TODO get "pathParams" from method on "msg"

	for _, name := range pathParamsNames {
		req.Form.Set(name, req.PathValue(name))
	}

	if err := dec.Decode(msg, req.Form); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	return nil
}

// GenerateURL will turn a request proto message into a url.
func GenerateURL(
	enc ValuesEncoder,
	msg proto.Message,
	pat *httppattern.Pattern,
) (string, error) {
	vals, err := enc.Encode(msg)
	if err != nil {
		return "", fmt.Errorf("failed to encode request message: %w", err)
	}

	// @TODO get "pathParams" from method on 'msg'?
	pathParamNames := httppattern.ParamNames(pat)

	// after encoding, extract the path params the query
	pathVals := make([]string, 0, len(pathParamNames))
	for _, name := range pathParamNames {
		pathVals = append(pathVals, vals.Get(name))

		vals.Del(name) // so the path params don't end up in the query
	}

	// then, build the uri.
	uri, err := httppattern.Build(pat, pathVals...)
	if err != nil {
		return uri, fmt.Errorf("failed to build url from pattern: %w", err)
	}

	// if there are any query params, add them.
	if qry := vals.Encode(); qry != "" {
		return uri + "?" + vals.Encode(), nil
	}

	return uri, nil
}
