// Package main implements our protobuf code generator.
package main

import (
	"flag"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// programs entrypoint.
func main() {
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(run)
}

// run the plugin logic.
func run(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	log.Println("hello, world")

	return nil
}
