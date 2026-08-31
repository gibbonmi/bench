package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/mod/modfile"
)

// otelPinnedModules are the OpenTelemetry modules the kit declares directly. The seam
// record pins one version across the three, so a split version is a footprint defect.
var otelPinnedModules = []string{
	"go.opentelemetry.io/otel",
	"go.opentelemetry.io/otel/trace",
	"go.opentelemetry.io/otel/sdk",
}

const otelPinnedVersion = "v1.46.0"

// excludedFootprintModules name the module families the dependency standard keeps out of
// the tree. Every official OTLP exporter drags protobuf in, and its transport drags gRPC
// in, so the seam record carries a hand-written encoder instead.
var excludedFootprintModules = []string{"protobuf", "grpc"}

// checkOtelFootprint grades the kit's own go.mod. The dependency standard is a kit-owned
// policy, so this check rides the kit-compliance subject rather than the graded root. A
// subject without a go.mod is not a Go module, so it carries no footprint to grade.
func checkOtelFootprint(kitRoot string) []string {
	data, err := os.ReadFile(filepath.Join(kitRoot, "go.mod"))
	if err != nil {
		return nil
	}
	return otelFootprintDiags(string(data))
}

func otelFootprintDiags(gomod string) []string {
	parsed, err := modfile.Parse("go.mod", []byte(gomod), nil)
	if err != nil {
		return []string{"go.mod does not parse: " + err.Error()}
	}

	versions := map[string]string{}
	var diags []string
	for _, require := range parsed.Require {
		path := require.Mod.Path
		versions[path] = require.Mod.Version
		for _, excluded := range excludedFootprintModules {
			if moduleNames(path, excluded) {
				diags = append(diags, "go.mod names an excluded module: "+path)
			}
		}
	}

	for _, module := range otelPinnedModules {
		version, present := versions[module]
		switch {
		case !present:
			diags = append(diags, "go.mod does not require "+module)
		case version != otelPinnedVersion:
			diags = append(diags, module+" is pinned at "+version+", not "+otelPinnedVersion)
		}
	}
	return diags
}

// moduleNames reports whether a module path carries the name as a path element or as a
// host label, so `google.golang.org/protobuf` and `google.golang.org/grpc` both match
// while an unrelated path that merely embeds the letters does not.
func moduleNames(path, name string) bool {
	for _, element := range strings.FieldsFunc(path, func(r rune) bool { return r == '/' || r == '.' }) {
		if element == name {
			return true
		}
	}
	return false
}

const otelFootprintFixture = `module example.com/subject

go 1.25.0

require (
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
)
`

func TestOtelFootprintCheckBites(t *testing.T) {
	if diags := otelFootprintDiags(otelFootprintFixture); len(diags) != 0 {
		t.Fatalf("the pinned fixture is not clean: %v", diags)
	}
	for name, mutation := range map[string]string{
		"protobuf": "\nrequire google.golang.org/protobuf v1.36.10 // indirect\n",
		"gRPC":     "\nrequire google.golang.org/grpc v1.76.0 // indirect\n",
	} {
		if len(otelFootprintDiags(otelFootprintFixture+mutation)) == 0 {
			t.Errorf("the footprint check stayed green with a %s module in go.mod", name)
		}
	}
	if len(otelFootprintDiags(strings.ReplaceAll(otelFootprintFixture, otelPinnedVersion, "v1.45.0"))) == 0 {
		t.Error("the footprint check stayed green with an off-pin OpenTelemetry version")
	}
	if len(otelFootprintDiags(strings.ReplaceAll(otelFootprintFixture, "go.opentelemetry.io/otel/sdk v1.46.0\n\t", ""))) == 0 {
		t.Error("the footprint check stayed green with a dropped OpenTelemetry module")
	}
}
