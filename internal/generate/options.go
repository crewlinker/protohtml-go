package generate

import (
	phtmlv1 "github.com/crewlinker/protohtml-go/phtml/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// routeOpts returns our plugin specific route options.
func routeOpts(m *protogen.Method) *phtmlv1.RouteOptions {
	opts, hasOpts := m.Desc.Options().(*descriptorpb.MethodOptions)
	if !hasOpts {
		return nil
	}

	ext, hasOpts := proto.GetExtension(opts, phtmlv1.E_Route).(*phtmlv1.RouteOptions)
	if !hasOpts {
		return nil
	}

	return ext
}
