// Code generated by protocgenpgxm. DO NOT EDIT.

package example1v1phtml

import (
	"context"
	v1 "github.com/crewlinker/protohtml-go/examples/example1/v1"
	phtml "github.com/crewlinker/protohtml-go/phtml"
	httppattern "github.com/crewlinker/protohtml-go/phtml/httppattern"
	"net/http"
)

// parsedPatterns hold all parsed route patterns, done once when the package initializes
var parsedPatterns = map[string]*httppattern.Pattern{}

func init() {
	parsedPatterns["AnotherService.ShowOneAddress"] = httppattern.MustParsePattern("/addr/{addr_id}")
	parsedPatterns["MovrService.ShowOneUser"] = httppattern.MustParsePattern("/user/{user_id}")
	parsedPatterns["MovrService.ShowUserAddress"] = httppattern.MustParsePattern("/user/{user_id}/address/{addr_id}")
}

// AnotherService describes the route handler implementation.
type AnotherService interface {
	ShowOneAddress(ctx context.Context, req *v1.ShowOneAddressRequest) (*v1.ShowOneAddressResponse, error)
}

// AnotherServiceHandlers provides methods for serving our routes.
type AnotherServiceHandlers struct {
	impl  AnotherService
	phtml *phtml.PHTML
}

// NewAnotherServiceHandlers constructs the handler set.
func NewAnotherServiceHandlers(impl AnotherService) *AnotherServiceHandlers {
	return &AnotherServiceHandlers{
		impl:  impl,
		phtml: phtml.New(),
	}
}

// ShowOneAddressPattern returns the pattern for the Go 1.22 mux.
func (h *AnotherServiceHandlers) ShowOneAddressPattern() string {
	return "/addr/{addr_id}"
}

// ShowOneAddressURL generates a url given the parameterse.
func (h *AnotherServiceHandlers) ShowOneAddressURL(addrId string) (string, error) {
	x := phtml.FirstInitOrPanic[v1.ShowOneAddressRequest]()
	{
		x.AddrId = addrId
	}
	return h.phtml.GenerateURL(x, parsedPatterns["AnotherService.ShowOneAddress"])
}

// ShowOneAddressHandler returns the http handler for the route.
func (h *AnotherServiceHandlers) ShowOneAddressHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req v1.ShowOneAddressRequest
		if err := h.phtml.ParseRequest(&req, r, "addr_id"); err != nil {
			h.phtml.HandleParseRequestError(ctx, w, r, err)
			return
		}
		resp, err := h.impl.ShowOneAddress(ctx, &req)
		if err != nil {
			h.phtml.HandleImplementationError(ctx, w, r, err)
		}
		if err := v1.ShowOneAddress(resp).Render(ctx, w); err != nil {
			h.phtml.HandleParseRequestError(ctx, w, r, err)
			return
		}
	})
}

// MovrService describes the route handler implementation.
type MovrService interface {
	ShowOneUser(ctx context.Context, req *v1.ShowOneUserRequest) (*v1.ShowOneUserResponse, error)
	ShowUserAddress(ctx context.Context, req *v1.ShowUserAddressRequest) (*v1.ShowUserAddressResponse, error)
}

// MovrServiceHandlers provides methods for serving our routes.
type MovrServiceHandlers struct {
	impl  MovrService
	phtml *phtml.PHTML
}

// NewMovrServiceHandlers constructs the handler set.
func NewMovrServiceHandlers(impl MovrService) *MovrServiceHandlers {
	return &MovrServiceHandlers{
		impl:  impl,
		phtml: phtml.New(),
	}
}

// ShowOneUserPattern returns the pattern for the Go 1.22 mux.
func (h *MovrServiceHandlers) ShowOneUserPattern() string {
	return "/user/{user_id}"
}

// ShowUserAddressPattern returns the pattern for the Go 1.22 mux.
func (h *MovrServiceHandlers) ShowUserAddressPattern() string {
	return "/user/{user_id}/address/{addr_id}"
}

// ShowOneUserURL generates a url given the parameterse.
func (h *MovrServiceHandlers) ShowOneUserURL(userId string, opt ...*v1.ShowOneUserRequest) (string, error) {
	x := phtml.FirstInitOrPanic[v1.ShowOneUserRequest](opt...)
	{
		x.UserId = userId
	}
	return h.phtml.GenerateURL(x, parsedPatterns["MovrService.ShowOneUser"])
}

// ShowUserAddressURL generates a url given the parameterse.
func (h *MovrServiceHandlers) ShowUserAddressURL(userId string, addrId string) (string, error) {
	x := phtml.FirstInitOrPanic[v1.ShowUserAddressRequest]()
	{
		x.UserId = userId
		x.AddrId = addrId
	}
	return h.phtml.GenerateURL(x, parsedPatterns["MovrService.ShowUserAddress"])
}

// ShowOneUserHandler returns the http handler for the route.
func (h *MovrServiceHandlers) ShowOneUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req v1.ShowOneUserRequest
		if err := h.phtml.ParseRequest(&req, r, "user_id"); err != nil {
			h.phtml.HandleParseRequestError(ctx, w, r, err)
			return
		}
		resp, err := h.impl.ShowOneUser(ctx, &req)
		if err != nil {
			h.phtml.HandleImplementationError(ctx, w, r, err)
		}
		if err := v1.ShowOneUser(resp).Render(ctx, w); err != nil {
			h.phtml.HandleParseRequestError(ctx, w, r, err)
			return
		}
	})
}

// ShowUserAddressHandler returns the http handler for the route.
func (h *MovrServiceHandlers) ShowUserAddressHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req v1.ShowUserAddressRequest
		if err := h.phtml.ParseRequest(&req, r, "user_id", "addr_id"); err != nil {
			h.phtml.HandleParseRequestError(ctx, w, r, err)
			return
		}
		resp, err := h.impl.ShowUserAddress(ctx, &req)
		if err != nil {
			h.phtml.HandleImplementationError(ctx, w, r, err)
		}
		if err := v1.ShowUserAddress(resp).Render(ctx, w); err != nil {
			h.phtml.HandleParseRequestError(ctx, w, r, err)
			return
		}
	})
}
