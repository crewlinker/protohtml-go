package generate

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/crewlinker/protohtml-go/internal/httppattern"
	"google.golang.org/protobuf/compiler/protogen"
)

// packageNameSuffix determines the sub-package's full name.
const packageNameSuffix = "phtml"

// Blueprint defines the code generation.
type Blueprint struct {
	Packages map[string]*Package
}

// Package represent a single package to generate for.
type Package struct {
	Dir           string
	GoImportPath  protogen.GoImportPath
	GoPackageName protogen.GoPackageName
	Result        *bytes.Buffer

	Services map[string]*Service
}

// IsEmpty returns true if the package definition shows
// there is no code to generate.
func (p *Package) IsEmpty() bool {
	return len(p.Services) == 0
}

// Service of routes.
type Service struct {
	GoName  string
	Methods map[string]*Method
}

// Method describes a route that is served
// by the generated handlers.
type Method struct {
	GoName  string
	Pattern *httppattern.Pattern
}

// preGen a service.
func preGenService(pkg *Package, genSvc *protogen.Service) error {
	for _, genMethod := range genSvc.Methods {
		ropts := routeOpts(genMethod)
		if ropts == nil {
			continue
		}

		// create a service in a ad-hoc fashion if the options
		// show that a method is declared a route.
		svc, ok := pkg.Services[genSvc.GoName]
		if !ok {
			svc = &Service{
				GoName:  genSvc.GoName,
				Methods: map[string]*Method{},
			}

			pkg.Services[svc.GoName] = svc
		}

		// parse pattern
		pat, err := httppattern.ParsePattern(ropts.GetPattern())
		if err != nil {
			return fmt.Errorf("[%s] failed to parse route pattern '%s': %w", genMethod.GoName, ropts.GetPattern(), err)
		}

		// add our representation of a method.
		svc.Methods[genMethod.GoName] = &Method{
			GoName:  genMethod.GoName,
			Pattern: pat,
		}
	}

	return nil
}

// preGenPlugin will create the code generation "blueprint": a domain-specific representation
// that can then be used gemerate the actual code.
func preGenPlugin(plugin *protogen.Plugin) (*Blueprint, error) {
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
				GoPackageName: plugf.GoPackageName + packageNameSuffix,
				GoImportPath:  plugf.GoImportPath,
				Services:      map[string]*Service{},
			}
		}

		for _, service := range plugf.Services {
			if err := preGenService(pkg, service); err != nil {
				return nil, fmt.Errorf("[%s] failed to pre-gen service: %w", service.GoName, err)
			}
		}

		blueprint.Packages[dir] = pkg
	}

	return blueprint, nil
}
