package generate

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/crewlinker/protohtml-go/internal/httppattern"
	phtmlv1 "github.com/crewlinker/protohtml-go/phtml/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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

	Services  map[string]*Service
	Requests  map[protogen.GoIdent]*Request
	Responses map[protogen.GoIdent]*Response
}

// IsEmpty returns true if the package definition shows
// there is no code to generate.
func (p *Package) IsEmpty() bool {
	return len(p.Services) == 0
}

// Param describes a request parameter.
type Param struct {
	GoName   string
	Source   phtmlv1.Source
	DescKind protoreflect.Kind
}

// Request message.
type Request struct {
	GoIdent protogen.GoIdent
	Params  map[protoreflect.Name]*Param
}

// Response message.
type Response struct {
	GoIdent protogen.GoIdent
}

// Service groups a number of routes.
type Service struct {
	GoName string
	Routes map[string]*Route
}

// Route describes a route that is served
// by the generated handlers.
type Route struct {
	GoName      string
	StrPattern  string
	Pattern     *httppattern.Pattern
	InputIdent  protogen.GoIdent
	OutputIdent protogen.GoIdent
}

// AssertPathParamField asserts if the field description is suited for a path parameter.
func AssertPathParamField(card protoreflect.Cardinality, kind protoreflect.Kind) error {
	if card != protoreflect.Optional {
		return fmt.Errorf("path parameter field must have default cardinality, has: %q", card)
	}

	switch kind {
	case protoreflect.BoolKind, protoreflect.Int32Kind, protoreflect.Sint32Kind,
		protoreflect.Uint32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind,
		protoreflect.Uint64Kind, protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
		protoreflect.FloatKind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind,
		protoreflect.DoubleKind, protoreflect.StringKind, protoreflect.BytesKind:
		return nil // all basic
	case protoreflect.MessageKind, protoreflect.EnumKind, protoreflect.GroupKind:
		return fmt.Errorf("path parameter field must be basic kind, has: %q", kind)
	default:
		return fmt.Errorf("unsupported kind for path parameter, got: %q", kind)
	}
}

// preGenRequest pre-generates any request message.
func preGenRequest(_ *Package, inp *protogen.Message) (*Request, error) {
	req := &Request{
		GoIdent: inp.GoIdent,
		Params:  map[protoreflect.Name]*Param{},
	}

	for _, fld := range inp.Fields {
		popts := paramOpts(fld)
		if popts == nil {
			continue // not a field for a parameter.
		}

		// path parameters have constraints we assert in the pre-generation phase.
		if popts.GetSource() == phtmlv1.Source_SOURCE_PATH {
			if err := AssertPathParamField(fld.Desc.Cardinality(), fld.Desc.Kind()); err != nil {
				return nil, fmt.Errorf("[%s] failed to assert as path param: %w", fld.GoName, err)
			}
		}

		req.Params[fld.Desc.Name()] = &Param{
			GoName:   fld.GoName,
			Source:   popts.GetSource(),
			DescKind: fld.Desc.Kind(),
		}
	}

	// @TODO assert that all path parameters that the pattern expect are defined.

	return req, nil
}

// preGenRoute pre-generates any method.
func preGenRoute(pkg *Package, genMethod *protogen.Method, ropts *phtmlv1.RouteOptions) (route *Route, err error) {
	pat, err := httppattern.ParsePattern(ropts.GetPattern())
	if err != nil {
		return nil, fmt.Errorf("failed to parse route pattern '%s': %w",
			ropts.GetPattern(), err)
	}

	req, ok := pkg.Requests[genMethod.Input.GoIdent]
	if !ok {
		req, err = preGenRequest(pkg, genMethod.Input)
		if err != nil {
			return nil, fmt.Errorf("failed to pre-generate request from input: %w", err)
		}

		pkg.Requests[genMethod.Input.GoIdent] = req
	}

	resp, ok := pkg.Responses[genMethod.Output.GoIdent]
	if !ok {
		resp = &Response{GoIdent: genMethod.Output.GoIdent}
		pkg.Responses[genMethod.Output.GoIdent] = resp
	}

	return &Route{
		GoName:      genMethod.GoName,
		Pattern:     pat,
		StrPattern:  ropts.GetPattern(),
		InputIdent:  req.GoIdent,
		OutputIdent: resp.GoIdent,
	}, nil
}

// preGen a service.
func preGenService(pkg *Package, genSvc *protogen.Service) (err error) {
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
				GoName: genSvc.GoName,
				Routes: map[string]*Route{},
			}

			pkg.Services[svc.GoName] = svc
		}

		svc.Routes[genMethod.GoName], err = preGenRoute(pkg, genMethod, ropts)
		if err != nil {
			return fmt.Errorf("[%s] failed to pre-gen method: %w", genMethod.GoName, err)
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
				Requests:      map[protogen.GoIdent]*Request{},
				Responses:     map[protogen.GoIdent]*Response{},
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
