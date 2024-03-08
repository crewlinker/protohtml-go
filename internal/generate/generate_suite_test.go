package generate_test

import (
	"os"
	"path/filepath"
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

// use a large set of examples to test our code generator.
var _ = DescribeTable("generate", func(expErr OmegaMatcher) {
	name := CurrentSpecReport().LeafNodeText
	AssertTestByName(name, expErr)
},
	Entry("example1", BeNil()),
)

// helper for asserting code generation for a test.
func AssertTestByName(name string, expErr OmegaMatcher) {
	opts, req := protogen.Options{}, pluginpb.CodeGeneratorRequest{}

	snapshot, err := os.ReadFile(filepath.Join("..", "..", "examples", name, "v1", "phtml_request.bin"))
	Expect(err).ToNot(HaveOccurred())

	Expect(proto.Unmarshal(snapshot, &req)).To(Succeed())

	plugin, err := opts.New(&req)
	Expect(err).ToNot(HaveOccurred())

	actFiles, actErr := generate.Generate(plugin)
	Expect(actErr).To(expErr)

	_ = actFiles
}
