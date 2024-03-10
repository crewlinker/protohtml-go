// Package main implements our protobuf code generator.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/crewlinker/protohtml-go/internal/generate"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// snapshot flag.
var snapshots = flag.Bool("snapshots", false, "enable snapshots to be saved for testing")

// programs entrypoint.
func main() {
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(run)
}

// run the plugin logic.
func run(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	// write snapshots if enabled (for testing)
	if *snapshots {
		if dir, err := snapshot(plugin.Request); err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		} else if strings.HasPrefix(dir, "examples/example_err") {
			log.Printf("skipping actual code generation for error example (%s)", dir)

			return nil
		}
	}

	files, err := generate.Generate(plugin)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	for name, pkg := range files {
		genf := plugin.NewGeneratedFile(name, pkg.GoImportPath)

		if _, err := io.Copy(genf, pkg.Result); err != nil {
			return fmt.Errorf("failed to copy code generation results into file: %w", err)
		}
	}

	return nil
}

// snapshot requests in the directories the generate for.
func snapshot(req *pluginpb.CodeGeneratorRequest) (string, error) {
	const perms = 0o600

	dir := filepath.Dir(req.GetFileToGenerate()[0])
	if dir == "phtml/v1" {
		return "", nil // skip for our own plugin options
	}

	fname := filepath.Join(dir, "phtml_request.bin")
	log.Printf("writing request snapshot to: %s", fname)

	bin, err := proto.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := os.WriteFile(fname, bin, perms); err != nil {
		return "", fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return dir, nil
}
