package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

// TestExecutedTagCensus pins row CE3: the census a citation grades against comes from
// the resolved phase table, not from a copy of the tag list.
func TestExecutedTagCensus(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		goMod    bool
		want     []TagSet
	}{
		{
			name:  "kit root holds the untagged set and the system set",
			goMod: true,
			want:  []TagSet{{}, {SystemPhaseName}},
		},
		{
			name:     "a manifest-declared custom tag joins the census",
			manifest: `{"phases":[{"name":"test","argv":["go","test","-tags=fixturetag","./..."]}]}`,
			want:     []TagSet{{"fixturetag"}},
		},
		{
			name:     "a root with no test phase yields an empty census",
			manifest: `{"phases":[{"name":"vet","argv":["go","vet","./..."]}]}`,
			want:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(baselinePolicyEnv, "")
			root := censusRoot(t, test.goMod, test.manifest)
			got, err := ExecutedTagCensus(root, root)
			if err != nil {
				t.Fatalf("ExecutedTagCensus: %v", err)
			}
			assertCensus(t, got, test.want)
		})
	}
}

// censusRoot writes one graded root. The root is its own kit, which is the arrangement
// the built-in table calls the kit root.
func censusRoot(t *testing.T, goMod bool, manifest string) string {
	t.Helper()
	root := t.TempDir()
	if goMod {
		writeCensusFile(t, filepath.Join(root, "go.mod"), "module censusfixture\n\ngo 1.24\n")
	}
	if manifest != "" {
		path := filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		writeCensusFile(t, path, manifest)
	}
	return root
}

func writeCensusFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// assertCensus compares the census to the expected sets. normalizeTags returns a
// non-nil slice for every set, so a deep comparison is exact here.
func assertCensus(t *testing.T, got, want []TagSet) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("census = %v, want %v", got, want)
	}
}
