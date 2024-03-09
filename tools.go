//go:build tools

package tools

import (
	// generate protobuf code for Go
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	// generate type-safe html
	_ "github.com/a-h/templ/cmd/templ"
)
