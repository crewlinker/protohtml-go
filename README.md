# protohtml-go 
[![Check](https://github.com/crewlinker/protohtml-go/actions/workflows/checks.yml/badge.svg)](https://github.com/crewlinker/protohtml-go/actions/workflows/checks.yml)

Protobuf code generation for type-safe routing, links generation and HTML rendering

## Features
- Specification-first creation of routes
- Generates route patterns for use with the new go 1.22 *ServerMux
- Generates http.Handler with type-safe access to route parameters, headers and query parameters
- Generates type-safe generation of links
- Renders HTML using the Templ templates
- Validation using https://github.com/bufbuild/protovalidate-go
- Can use any standard middleware