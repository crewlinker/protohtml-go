// Package example1v1phtml provides strict html handling.
package example1v1phtml

import (
	"context"
	"fmt"
	"net/http"

	example1v1 "github.com/crewlinker/protohtml-go/examples/example1/v1"
	"github.com/crewlinker/protohtml-go/phtml"
)

// MovrServiceM describes the interface that needs to be implemented
type MovrServiceM interface {
	ShowOneUser(ctx context.Context, req *example1v1.ShowOneUserRequest) (*example1v1.ShowOneUserResponse, error)
}

// MovrServiceHandlerSetM surfaces the logic for resilient HTML serving.
type MovrServiceHandlerSetM struct {
	base *phtml.PHTML
	impl MovrServiceM
}

// NewMovrServiceHandlerSetM inits the handler set for our service.
func NewMovrServiceHandlerSetM(impl MovrServiceM) *MovrServiceHandlerSetM {
	return &MovrServiceHandlerSetM{
		impl: impl,
		base: phtml.New(),
	}
}

func (hs *MovrServiceHandlerSetM) ShowOneUserPattern() string {
	return `/user/{user_id}`
}

func (hs *MovrServiceHandlerSetM) ShowOneUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req example1v1.ShowOneUserRequest

		if err := hs.base.ParseRequest(&req, r, "user_id"); err != nil {
			hs.base.HandleParseRequestError(ctx, w, r, err)
			return
		}

		resp, err := hs.impl.ShowOneUser(ctx, &req)
		if err != nil {
			hs.base.HandleImplementationError(ctx, w, r, err)
			return
		}

		if err := example1v1.ShowOneUser(resp).Render(ctx, w); err != nil {
			hs.base.HandleRenderResponseError(ctx, w, r, err)
			return
		}
	})
}

// @TODO at least partially generic, move to a shared package. @TODO, how to pass-in/handle optional parameters?
func (hs *MovrServiceHandlerSetM) ShowOneUserURL(userId string, xs ...*example1v1.ShowOneUserRequest) (string, error) {
	x := phtml.FirstInitOrPanic(xs...)

	x.UserId = userId

	uri, err := hs.base.GenerateURL(x, parsedPatterns["MovrService.ShowOneUser"])
	if err != nil {
		return "", fmt.Errorf("failed to generate URL: %w", err)
	}

	return uri, nil
}
