package module_test

import (
	"bytes"
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

func TestExtract(t *testing.T) {
	raw, err := os.ReadFile("../debug/code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

	// The request was captured with protoc-gen-debug's own parameters, so swap in
	// the ones gotag needs to locate tagger/tagger.pb.go from the repository root.
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(raw, req); err != nil {
		t.Fatal(err)
	}
	req.Parameter = proto.String("paths=source_relative")

	raw, err = proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	out := &bytes.Buffer{}
	pgs.Init(
		pgs.ProtocInput(bytes.NewReader(raw)),
		pgs.ProtocOutput(out),
		pgs.FileSystem(afero.NewMemMapFs()),
	).RegisterModule(module.New()).Render()

	resp := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(out.Bytes(), resp); err != nil {
		t.Fatal(err)
	}

	if resp.Error != nil {
		t.Fatalf("plugin error: %s", resp.GetError())
	}

	if len(resp.GetFile()) != 1 || resp.GetFile()[0].GetName() != "tagger/tagger.pb.go" {
		t.Fatalf("expected tagger/tagger.pb.go to be generated, got %v", resp.GetFile())
	}
}
