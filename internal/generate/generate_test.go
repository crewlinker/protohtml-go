package generate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/crewlinker/protohtml-go/internal/generate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerate(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/generate")
}

var _ = Describe("generate", func() {
	It("should not generate any files if our options are not used at all", func() {
		Expect(AssertTestByName("example_empty", "v1", BeNil())).To(BeEmpty())
	})

	DescribeTable("table", func(vname string, expErr OmegaMatcher, filef func(files Files)) {
		name := CurrentSpecReport().LeafNodeText
		files := AssertTestByName(name, vname, expErr)
		filef(files)
	},
		Entry("example1", "v1", BeNil(), AssertFilesWithOnDisk),
		Entry("example_err", "v1", MatchError(MatchRegexp(`field must be basic kind`)), AssertFilesEmpty),
		Entry("example_err", "v2", MatchError(MatchRegexp(`path parameter field must have default cardinality`)), AssertFilesEmpty),
		Entry("example_err", "v3", MatchError(MatchRegexp(`failed to parse route pattern`)), AssertFilesEmpty),
		Entry("example_err", "v4", MatchError(MatchRegexp(`don't match the path parameters defined in the request message`)), AssertFilesEmpty),
		Entry("example_err", "v5", MatchError(MatchRegexp(` must define a templ component for rendering`)), AssertFilesEmpty),
	)
})

// AssertFilesWithOnDisk asserts with what is generated by running the generator
// outside of the tests.
func AssertFilesWithOnDisk(files Files) {
	GinkgoHelper()
	for name, pkg := range files {
		data, err := os.ReadFile(filepath.Join("..", "..", strings.ReplaceAll(name, ".go", ".exp")))
		Expect(err).ToNot(HaveOccurred())
		Expect(pkg.Result.String()).To(Equal(string(data)))
	}
}

// AssertFilesEmpty asserts that there are nog generated files.
func AssertFilesEmpty(files Files) {
	Expect(files).To(BeEmpty())
}

// Files provides a more readable type for our tests.
type Files map[string]*generate.Package

// helper for asserting code generation for a test.
func AssertTestByName(name, vdir string, expErr OmegaMatcher) Files {
	opts, req := protogen.Options{}, pluginpb.CodeGeneratorRequest{}

	snapshot, err := os.ReadFile(filepath.Join("..", "..", "examples", name, vdir, "phtml_request.bin"))
	Expect(err).ToNot(HaveOccurred())

	Expect(proto.Unmarshal(snapshot, &req)).To(Succeed())

	plugin, err := opts.New(&req)
	Expect(err).ToNot(HaveOccurred())

	actFiles, actErr := generate.Generate(plugin)
	Expect(actErr).To(expErr)

	return actFiles
}
