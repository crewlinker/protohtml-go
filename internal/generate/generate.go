// Package generate implements the code generator.
package generate

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sort"

	. "github.com/dave/jennifer/jen" //nolint:revive,stylecheck

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	parsedPatternsVarName = "parsedPatterns"
)

const (
	httpatternpkg = "github.com/crewlinker/protohtml-go/internal/httppattern"
	phttppkg      = "github.com/crewlinker/protohtml-go/phtml"
)

// parsedPatternsKey standardizes on the key name of the pre-parsed routes.
func parsedPatternsKey(svc *Service, route *Route) string {
	return fmt.Sprintf("%s.%s", svc.GoName, route.GoName)
}

// foreachKeySortedErr iterates over the items of a map with string keys, sorted by the key.
func foreachKeySortedErr[T any](m map[string]T, f func(string, T) error) error {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if err := f(k, m[k]); err != nil {
			return err
		}
	}

	return nil
}

// foreachKeySorted iterates over the items of a map with string keys, sorted by the key.
func foreachKeySorted[T any](m map[string]T, f func(string, T)) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		f(k, m[k])
	}
}

// generateInit will perform logic when the package is initializing.
func generateInit(file *File, pkg *Package) error {
	file.Commentf("%s hold all parsed route patterns, done once when the package initializes", parsedPatternsVarName)
	file.Var().Id(parsedPatternsVarName).Op("=").Map(String()).Op("*").Qual(httpatternpkg, "Pattern").Values()

	block := []Code{}

	foreachKeySorted(pkg.Services, func(_ string, svc *Service) {
		foreachKeySorted(svc.Routes, func(_ string, route *Route) {
			block = append(block,
				Id("parsedPatterns").Index(Lit(parsedPatternsKey(svc, route))).
					Op("=").Qual(httpatternpkg, "MustParsePattern").Call(Lit(route.StrPattern)))
		})
	})

	file.Func().Id("init").Params().Block(block...)

	return nil
}

// standardizes on var naming from param names.
func paramNameToVarIdent(s string) string {
	return strcase.ToLowerCamel(s)
}

// generateURLGeneration generates the code for generation URLs for a route.
func generateURLGeneration(file *File, pkg *Package, svc *Service) error {
	return foreachKeySortedErr(svc.Routes, func(_ string, route *Route) error {
		req, ok := pkg.Requests[route.RequestIdent]
		if !ok {
			panic("encountered unknown request ident: " + route.RequestIdent.String())
		}

		// setup code that requires iteration of the params
		paramCode, assignCode := []Code{}, []Code{}
		for _, pathParamName := range req.PathParamsAsInPattern {
			param := req.Params[pathParamName]
			paramCode = append(paramCode, Id(paramNameToVarIdent(pathParamName)).Id(param.GoBasicType))

			assignCode = append(assignCode, Id("x").Dot(param.GoName).Op("=").Id(paramNameToVarIdent(pathParamName)))
		}

		// we only generate the variadic opt argument if the message has any optional parameters.
		var initOrPanicArg *Statement
		if len(req.Params)-len(req.PathParamsAsInPattern) > 0 {
			paramCode = append(paramCode, Id("opt").
				Op("...").Op("*").Qual(string(req.GoIdent.GoImportPath), req.GoIdent.GoName))
			initOrPanicArg = Id("opt").Op("...")
		}

		// generate the actual method
		file.Commentf(route.GoName + "URL generates a url given the parameterse.")
		file.Func().Params(Id("h").Op("*").Id(svc.GoName+"Handlers")).Id(route.GoName+"URL").
			Params(paramCode...).
			Params(String(), Error()).
			Block(
				// generate initialization of the request message struct
				Id("x").Op(":=").Qual(phttppkg, "FirstInitOrPanic").Index(
					Qual(string(req.GoIdent.GoImportPath), req.GoIdent.GoName),
				).Call(initOrPanicArg),
				// assing path params from method arguments
				Block(assignCode...),
				// call the shared GenerateURL method
				Return(Id("h").Dot("phtml").Dot("GenerateURL").Call(
					Id("x"),
					Id(parsedPatternsVarName).Index(Lit(parsedPatternsKey(svc, route))),
				)),
			)

		return nil
	})
}

// generateHandlerSets generates the handler sets.
func generateHandlerSets(file *File, pkg *Package) error {
	return foreachKeySortedErr(pkg.Services, func(_ string, svc *Service) error {
		// implementation interface.
		file.Commentf(svc.GoName + " describes the route handler implementation.")
		file.Type().Id(svc.GoName).Interface()

		// handlers struct
		file.Commentf(svc.GoName + "Handlers provides methods for serving our routes.")
		file.Type().Id(svc.GoName+"Handlers").Struct(
			Id("impl").Id(svc.GoName),
			Id("phtml").Op("*").Qual(phttppkg, "PHTML"),
		)

		// handlers constructor
		file.Commentf("New" + svc.GoName + "Handlers constructs the handler set.")
		file.Func().Id("New" + svc.GoName + "Handlers").
			Params(Id("impl").Id(svc.GoName)).
			Params(Op("*").Id(svc.GoName + "Handlers")).
			Block(Return(Op("&").Id(svc.GoName + "Handlers").Values(Dict{
				Id("impl"):  Id("impl"),
				Id("phtml"): Qual(phttppkg, "New").Call(),
			})))

		// generate the handlers method that returns the pattern
		foreachKeySorted(svc.Routes, func(_ string, route *Route) {
			file.Commentf(route.GoName + "Pattern returns the pattern for the Go 1.22 mux.")
			file.Func().Params(Id("h").Op("*").Id(svc.GoName + "Handlers")).Id(route.GoName + "Pattern").
				Params().
				Params(String()).
				Block(Return(Lit(route.StrPattern)))
		})

		// generate the method that generate urls.
		if err := generateURLGeneration(file, pkg, svc); err != nil {
			return fmt.Errorf("[%s] failed to generate url method: %w", svc.GoName, err)
		}

		return nil
	})
}

// generatePackage generates for a single package.
func generatePackage(w io.Writer, pkg *Package) error {
	file := NewFile(string(pkg.GoPackageName))
	file.HeaderComment("Code generated by protocgenpgxm. DO NOT EDIT.")

	if err := generateInit(file, pkg); err != nil {
		return fmt.Errorf("failed to generate init: %w", err)
	}

	if err := generateHandlerSets(file, pkg); err != nil {
		return fmt.Errorf("failed to generate handler set: %w", err)
	}

	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}

	return nil
}

// Generate generates protohtml code.
func Generate(plugin *protogen.Plugin) (map[string]*Package, error) {
	blueprint, err := preGenPlugin(plugin)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-generate: %w", err)
	}

	pkgs := map[string]*Package{}
	for _, pkg := range blueprint.Packages {
		if pkg.IsEmpty() {
			continue // skip if empty.
		}

		fname := filepath.Join(pkg.Dir, string(pkg.GoPackageName), packageNameSuffix+".gen.go")
		fdata := bytes.NewBuffer(nil)

		if err := generatePackage(fdata, pkg); err != nil {
			return nil, fmt.Errorf("[%s] failed to generate: %w", pkg.GoPackageName, err)
		}

		pkg.Result = fdata
		pkgs[fname] = pkg
	}

	return pkgs, nil
}
