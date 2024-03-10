package phtml

import (
	"context"
	"net/http"

	"github.com/crewlinker/protohtml-go/phtml/httppattern"
	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

// ErrorHandlerFunc types a function to render an error.
type ErrorHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)

// PHTML is embeddded in any generated handlers.
type PHTML struct {
	dec ValuesDecoder
	enc ValuesEncoder

	NewMovrServiceHandlerSetM ErrorHandlerFunc
}

// New initializes the handler set.
func New() *PHTML {
	dec := form.NewDecoder()
	dec.SetTagName("json")
	enc := form.NewEncoder()
	enc.SetTagName("json")

	return &PHTML{
		dec: dec,
		enc: enc,
	}
}

// HandleParseRequestError takes care of an error while parsing the request.
func (bh PHTML) HandleParseRequestError(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// HandleImplementationError takes care of an error from the implementation call.
func (bh PHTML) HandleImplementationError(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// HandleRenderResponseError takes care of writing a rendering error to the response.
func (bh PHTML) HandleRenderResponseError(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// ParseRequest parses a requests's req.PostForm into a protobuf message 'msg'.
func (bh PHTML) ParseRequest(
	msg proto.Message,
	req *http.Request,
	pathParamsNames ...string,
) error {
	// @TODO make the package level function private
	return parseRequest(bh.dec, msg, req, pathParamsNames...)
}

// GenerateURL will turn a request proto message into a url.
func (bh PHTML) GenerateURL(msg any, pat *httppattern.Pattern) (string, error) {
	// @TODO make the package level function private
	return generateURL(bh.enc, msg, pat)
}
