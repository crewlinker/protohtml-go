package example1v1phtml_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	example1v1 "github.com/crewlinker/protohtml-go/examples/example1/v1"
	"github.com/crewlinker/protohtml-go/examples/example1/v1/example1v1phtml"
	"github.com/go-playground/form/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPhtml(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "examples/example1/v1")
}

var _ = Describe("handling", func() {
	var mux *http.ServeMux
	var hs *example1v1phtml.MovrServiceHandlerSetM
	var enc *form.Encoder
	var dec *form.Decoder

	BeforeEach(func() {
		enc = form.NewEncoder()
		enc.SetTagName("json")
		dec = form.NewDecoder()
		dec.SetTagName("json")
		hs = example1v1phtml.NewMovrServiceHandlerSetM(&MovrImpl1{}, dec, enc)

		mux = http.NewServeMux()
		mux.Handle(hs.ShowOneUserPattern(), hs.ShowOneUserHandler())
	})

	It("should serve the generated handler", func() {
		resp, req := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/user/111-34", nil)
		mux.ServeHTTP(resp, req)

		Expect(resp.Body.String()).To(Equal(`<div>hello, John Doe</div>`))
	})
})

type MovrImpl1 struct{}

func (MovrImpl1) ShowOneUser(context.Context, *example1v1.ShowOneUserRequest) (*example1v1.ShowOneUserResponse, error) {
	return &example1v1.ShowOneUserResponse{
		FirstName: "John",
		LastName:  "Doe",
	}, nil
}
