package module_test

import (
	"bytes"
	"flag"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/fatih/structtag"

	"github.com/bosonicalcom/protoc-gen-gotag/module"
)

var replaceOut = flag.Bool("tag-rep", false, "")

func TestRetag(t *testing.T) {
	fs := token.NewFileSet()

	n, err := parser.ParseFile(fs, "./test/input.txt", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	module.Retag(n, map[string]map[string]*structtag.Tags{
		"Simple": {
			"Single":   tagMust(structtag.Parse(`sql:"-,omitempty"`)),
			"Multiple": tagMust(structtag.Parse(`xml:"-,omitempty" sql:"ke,op" bson:"ke,op"`)),
			"None":     tagMust(structtag.Parse(`json:"none,omitempty"`)),
		},
	})

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fs, n); err != nil {
		t.Fatal(err)
	}

	if *replaceOut {
		f, err := os.Create("./test/golden.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		if _, err := io.Copy(f, &buf); err != nil {
			t.Fatal(err)
		}

		return
	}

	out, err := ioutil.ReadFile("./test/golden.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out, buf.Bytes()) {
		t.Error("output does not match golden file")
	}
}

func TestRetagOpaque(t *testing.T) {
	const src = `package main

type Opaque struct {
	xxx_hidden_Tagged   string ` + "`" + `protobuf:"bytes,1,opt,name=tagged"` + "`" + `
	xxx_hidden_Untagged string ` + "`" + `protobuf:"bytes,2,opt,name=untagged"` + "`" + `
}
`

	fs := token.NewFileSet()
	n, err := parser.ParseFile(fs, "opaque.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	// Fields without any tags to apply must not trip the opaque check.
	if err := module.Retag(n, map[string]map[string]*structtag.Tags{
		"Opaque": {"Untagged": tagMust(structtag.Parse(``))},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = module.Retag(n, map[string]map[string]*structtag.Tags{
		"Opaque": {"Tagged": tagMust(structtag.Parse(`json:"tagged"`))},
	})
	if err == nil || !strings.Contains(err.Error(), "API_OPEN") {
		t.Fatalf("expected opaque API error, got: %v", err)
	}
}

func tagMust(t *structtag.Tags, err error) *structtag.Tags {
	if err != nil {
		panic(err)
	}
	return t
}
