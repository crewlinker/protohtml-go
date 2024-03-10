package example1v1phtml_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	example1v1 "github.com/crewlinker/protohtml-go/examples/example1/v1"
	"github.com/crewlinker/protohtml-go/examples/example1/v1/example1v1phtml"
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
	var hss *example1v1phtml.MovrServiceHandlers

	BeforeEach(func() {
		hss = example1v1phtml.NewMovrServiceHandlers(&MovrImpl1{})

		mux = http.NewServeMux()
		mux.Handle(hss.ShowOneUserPattern(), hss.ShowOneUserHandler())
	})

	It("should serve the generated handler", func() {
		resp, req := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/user/111-34", nil)
		mux.ServeHTTP(resp, req)

		Expect(resp.Body.String()).To(Equal(`<div>hello, John Doe</div>`))
	})

	It("should return the pattern", func() {
		Expect(hss.ShowUserAddressPattern()).To(Equal(`/user/{user_id}/address/{addr_id}`))
	})

	It("should generate url", func() {
		uri1, err := hss.ShowOneUserURL("1111092-9", &example1v1.ShowOneUserRequest{
			ShowAddress: true, UserId: "arg takes presedence",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(uri1).To(Equal(`/user/1111092-9?show_address=true`))

		uri2, err := hss.ShowUserAddressURL("1111092-9", "addr1")
		Expect(err).ToNot(HaveOccurred())
		Expect(uri2).To(Equal(`/user/1111092-9/address/addr1`))
	})
})

type MovrImpl1 struct{}

func (MovrImpl1) ShowOneUser(context.Context, *example1v1.ShowOneUserRequest) (*example1v1.ShowOneUserResponse, error) {
	return &example1v1.ShowOneUserResponse{
		FirstName: "John",
		LastName:  "Doe",
	}, nil
}

func (MovrImpl1) ShowUserAddress(ctx context.Context, req *example1v1.ShowUserAddressRequest) (*example1v1.ShowUserAddressResponse, error) {
	return nil, nil
}
