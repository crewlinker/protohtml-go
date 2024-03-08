package generate

import (
	"io"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
)

// Blueprint defines the code generation.
type Blueprint struct {
	Packages map[string]*Package
}

// Package represent a single package to generate for.
type Package struct {
	Dir           string
	GoImportPath  protogen.GoImportPath
	GoPackageName protogen.GoPackageName
	Result        io.Reader

	Routes map[string]Route
}

// IsEmpty returns true if the package definition shows
// there is no code to generate.
func (p *Package) IsEmpty() bool {
	return len(p.Routes) == 0
}

// Route describes a route that is served
// by the generated handlers.
type Route struct{}

// preGenerate will create the code generation "blueprint": a domain-specific representation
// that can then be used gemerate the actual code.
func preGenerate(plugin *protogen.Plugin) (*Blueprint, error) {
	blueprint := &Blueprint{Packages: map[string]*Package{}}

	// iterate over all files, this can span multiple directories (packages)
	for _, name := range plugin.Request.GetFileToGenerate() {
		dir := filepath.Dir(name)
		plugf := plugin.FilesByPath[name]

		// if we already know about the package, add to it.
		pkg, ok := blueprint.Packages[dir]
		if !ok {
			pkg = &Package{
				Dir:           dir,
				GoPackageName: plugf.GoPackageName,
				GoImportPath:  plugf.GoImportPath,
			}
		}

		blueprint.Packages[dir] = pkg
	}

	return blueprint, nil
}
