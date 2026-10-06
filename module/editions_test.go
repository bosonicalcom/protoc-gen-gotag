package module_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

// TestEditions runs gotag against the untagged protoc-gen-go output in testdata
// (regenerate it with `make testdata`) and compares the result with the tagged
// files in example/.
func TestEditions(t *testing.T) {
	files := []string{
		"example/editions2023/editions2023.proto",
		"example/editions2024/editions2024.proto",
	}

	raw, err := os.ReadFile("testdata/editions.binpb")
	if err != nil {
		t.Fatal(err)
	}

	set := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(raw, set); err != nil {
		t.Fatal(err)
	}

	req, err := proto.Marshal(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: files,
		Parameter:      proto.String("paths=source_relative"),
		ProtoFile:      set.GetFile(),
	})
	if err != nil {
		t.Fatal(err)
	}

	// gotag reads the generated go files relative to the working directory.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir("testdata"); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	out := &bytes.Buffer{}
	pgs.Init(
		pgs.ProtocInput(bytes.NewReader(req)),
		pgs.ProtocOutput(out),
		pgs.FileSystem(afero.NewMemMapFs()),
	).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()

	res, err := module.SetEditionSupport(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	resp := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(res, resp); err != nil {
		t.Fatal(err)
	}

	if resp.Error != nil {
		t.Fatalf("plugin error: %s", resp.GetError())
	}
	if got := resp.GetMinimumEdition(); got != int32(module.MinimumEdition) {
		t.Errorf("minimum edition = %v, want %v", got, module.MinimumEdition)
	}
	if got := resp.GetMaximumEdition(); got != int32(module.MaximumEdition) {
		t.Errorf("maximum edition = %v, want %v", got, module.MaximumEdition)
	}

	if len(resp.GetFile()) != len(files) {
		t.Fatalf("expected %d generated files, got %d", len(files), len(resp.GetFile()))
	}

	for _, f := range resp.GetFile() {
		want, err := os.ReadFile(filepath.Join(wd, "..", f.GetName()))
		if err != nil {
			t.Fatal(err)
		}

		if f.GetContent() != string(want) {
			t.Errorf("%s: output does not match tagged example", f.GetName())
		}
	}
}
