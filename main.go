package main

import (
	"bytes"
	"log"
	"os"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

func main() {
	opt := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
		pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)

	// protoc-gen-star does not know about editions, so the response is
	// buffered and the supported edition range is added before handing it to protoc.
	var out bytes.Buffer

	pgs.Init(
		pgs.DebugEnv("GOTAG_DEBUG"),
		pgs.SupportedFeatures(&opt),
		pgs.ProtocOutput(&out),
	).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()

	res, err := module.SetEditionSupport(out.Bytes())
	if err != nil {
		log.Fatalf("[gotag]: %v", err)
	}

	if _, err := os.Stdout.Write(res); err != nil {
		log.Fatalf("[gotag]: %v", err)
	}
}
