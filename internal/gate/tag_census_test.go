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
		want     []TestExecution
	}{
		{
			name:  "kit root holds the untagged set and the system set",
			goMod: true,
			want: []TestExecution{
				{Name: "test", Tags: TagSet{}, Packages: []string{"./..."}},
				{Name: SystemPhaseName, Tags: TagSet{SystemPhaseName}, Packages: []string{"./internal/systemtest"}},
			},
		},
		{
			name:     "a manifest-declared phase keeps its package operands",
			manifest: `{"phases":[{"name":"test","argv":["go","test","-tags=fixturetag","./..."]}]}`,
			want:     []TestExecution{{Name: "test", Tags: TagSet{"fixturetag"}, Packages: []string{"./..."}}},
		},
		{
			name: "equal tags retain separate package operands",
			manifest: `{"phases":[
				{"name":"first","argv":["go","test","-tags=fixturetag","./first"]},
				{"name":"second","argv":["go","test","-tags=fixturetag","./second"]}
			]}`,
			want: []TestExecution{
				{Name: "first", Tags: TagSet{"fixturetag"}, Packages: []string{"./first"}},
				{Name: "second", Tags: TagSet{"fixturetag"}, Packages: []string{"./second"}},
			},
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
			got, err := ExecutedTestCensus(root, root)
			if err != nil {
				t.Fatalf("ExecutedTestCensus: %v", err)
			}
			for i := range got {
				got[i].Dir = ""
				got[i].GoC = ""
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
func assertCensus(t *testing.T, got, want []TestExecution) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("census = %v, want %v", got, want)
	}
}
