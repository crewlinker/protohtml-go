package phtml_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	example1v1 "github.com/crewlinker/protohtml-go/examples/example1/v1"
	"github.com/crewlinker/protohtml-go/internal/httppattern"
	"github.com/crewlinker/protohtml-go/phtml"
	"github.com/go-playground/form/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPhtml(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "phtml")
}

var _ = Describe("phtml", func() {
	var dec *form.Decoder
	var enc *form.Encoder

	BeforeEach(func() {
		dec = form.NewDecoder()
		dec.SetTagName("json")
		enc = form.NewEncoder()
		enc.SetTagName("json")
	})

	It("should allow parsing requests into protobuf message", func() {
		req := httptest.NewRequest(http.MethodGet, "/users/1?show_address=true", nil)
		req.SetPathValue("user_id", "1223-1111")
		msg := &example1v1.ShowOneUserRequest{}

		Expect(phtml.ParseRequest(dec, msg, req, "user_id")).To(Succeed())

		Expect(msg.GetShowAddress()).To(BeTrue())
		Expect(msg.GetUserId()).To(Equal("1223-1111"))
	})

	It("should allow generating url from protobuf message", func() {
		pat, err := httppattern.ParsePattern(`/users/{user_id}`)
		Expect(err).ToNot(HaveOccurred())

		msg := &example1v1.ShowOneUserRequest{UserId: "1223-1111", ShowAddress: true}

		uri, err := phtml.GenerateURL(enc, msg, pat)
		Expect(err).ToNot(HaveOccurred())
		Expect(uri).To(Equal(`/users/1223-1111?show_address=true`))
	})

	It("should use FirstInitOrPanic corrrectly", func() {
		v1 := phtml.FirstInitOrPanic[example1v1.ShowOneUserRequest](nil)
		Expect(v1.GetUserId()).To(Equal(""))

		v2 := phtml.FirstInitOrPanic(&example1v1.ShowOneUserRequest{UserId: "foo"})
		Expect(v2.GetUserId()).To(Equal("foo"))

		Expect(func() {
			phtml.FirstInitOrPanic(
				&example1v1.ShowOneUserRequest{},
				&example1v1.ShowOneUserRequest{})
		}).To(PanicWith(MatchRegexp(`too many elements`)))
	})
})
