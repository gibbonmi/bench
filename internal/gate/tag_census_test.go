package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

// rootToken names the graded root inside an expected census entry. The root is a
// temporary directory, so a table entry writes this token and the runner resolves it
// against the root that case actually got.
const rootToken = "<root>"

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
				{Name: "test", Tags: TagSet{}, Packages: []string{"./..."}, Dir: rootToken},
				{
					Name:     SystemPhaseName,
					Tags:     TagSet{SystemPhaseName},
					Packages: []string{"./internal/systemtest"},
					Dir:      rootToken,
					Env:      []string{SystemRootEnv + "=" + rootToken},
				},
			},
		},
		{
			name:     "a manifest-declared phase keeps its package operands",
			manifest: `{"phases":[{"name":"test","argv":["go","test","-tags=fixturetag","./..."]}]}`,
			want: []TestExecution{
				{Name: "test", Tags: TagSet{"fixturetag"}, Packages: []string{"./..."}, Dir: rootToken},
			},
		},
		{
			name: "equal tags retain separate package operands",
			manifest: `{"phases":[
				{"name":"first","argv":["go","test","-tags=fixturetag","./first"]},
				{"name":"second","argv":["go","test","-tags=fixturetag","./second"]}
			]}`,
			want: []TestExecution{
				{Name: "first", Tags: TagSet{"fixturetag"}, Packages: []string{"./first"}, Dir: rootToken},
				{Name: "second", Tags: TagSet{"fixturetag"}, Packages: []string{"./second"}, Dir: rootToken},
			},
		},
		{
			// Go resolves a relative package operand after it changes to the -C directory,
			// so the effective directory is the phase's own directory joined with it.
			name:     "a -C setting joins onto the phase directory",
			manifest: `{"phases":[{"name":"test","argv":["go","-C","sub","test","./..."],"dir":"pkg"}]}`,
			want: []TestExecution{{
				Name:     "test",
				Tags:     TagSet{},
				Packages: []string{"./..."},
				Dir:      rootToken + "/pkg/sub",
				GoC:      rootToken + "/pkg/sub",
			}},
		},
		{
			// A phase's own environment decides which files the toolchain compiles, so it
			// travels with the phase rather than being replaced by the host's.
			name:     "a phase environment travels with its census entry",
			manifest: `{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"GOARCH":"386"}}]}`,
			want: []TestExecution{{
				Name:     "test",
				Tags:     TagSet{},
				Packages: []string{"./..."},
				Dir:      rootToken,
				Env:      []string{"GOARCH=386"},
			}},
		},
		{
			// Every one of these flags carries its value in the separated form, which is
			// the form that turns the value into a package operand when the flag table
			// omits the flag. `go list` then rejects the whole phase, and every citation
			// the phase could have supplied false-reds.
			name: "separated flag values are not package operands",
			manifest: `{"phases":[{"name":"test","argv":[
				"go","test","-skip","TestSlow","-o","/dev/null","-exec","echo",
				"-ldflags","-s -w","-outputdir","/tmp","-timeout","30s","./..."
			]}]}`,
			want: []TestExecution{
				{Name: "test", Tags: TagSet{}, Packages: []string{"./..."}, Dir: rootToken},
			},
		},
		{
			// Go's flag package accepts one or two leading dashes for the same flag, and a
			// project writes its manifest argv by hand. The two-dash form must consume its
			// own value rather than donate it to the package list.
			name: "a two-dash flag consumes its separated value",
			manifest: `{"phases":[{"name":"test","argv":[
				"go","test","--skip","TestSlow","--timeout","30s","./..."
			]}]}`,
			want: []TestExecution{
				{Name: "test", Tags: TagSet{}, Packages: []string{"./..."}, Dir: rootToken},
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
			assertCensus(t, got, resolveCensus(root, test.want))
		})
	}
}

// TestExecutedCensusReadsTheRaceDriverArgv checks the value-taking flag table against the
// one production argv that depends on it. raceDriverArgv ends with `-run <pattern>` after
// its package operands, so a table that omits -run reads the filter pattern as a package
// and makes `go list` reject the whole phase.
func TestExecutedCensusReadsTheRaceDriverArgv(t *testing.T) {
	t.Setenv(baselinePolicyEnv, "")
	root := censusRoot(t, true, "")
	// The race phase materializes only where the root declares one of its sentinel tests.
	writeCensusFile(t, censusPath(t, root, "internal", "worktree", "race_test.go"),
		"package worktree\n\nfunc "+raceTests[0].Name+"() {}\n")

	got, err := ExecutedTestCensus(root, root)
	if err != nil {
		t.Fatalf("ExecutedTestCensus: %v", err)
	}
	race, found := entryNamed(got, canary.PhaseRace)
	if !found {
		t.Fatalf("census = %v, want a %s entry", got, canary.PhaseRace)
	}
	// The expectation reads the race registry that raceDriverArgv builds the phase from,
	// so a registry edit moves the argv and this expectation together. It repeats the
	// driver's own deduplication, which keeps the first occurrence of each path. An
	// operand list that leaked the -run pattern or a build flag's value reds here.
	var want []string
	for _, test := range raceTests {
		if !contains(want, test.PackagePath) {
			want = append(want, test.PackagePath)
		}
	}
	if !reflect.DeepEqual(race.Packages, want) {
		t.Fatalf("race packages = %q, want %q", race.Packages, want)
	}
}

func entryNamed(census []TestExecution, name string) (TestExecution, bool) {
	for _, entry := range census {
		if entry.Name == name {
			return entry, true
		}
	}
	return TestExecution{}, false
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
		writeCensusFile(t, censusPath(t, root, filepath.FromSlash(canary.PhaseManifestPath)), manifest)
	}
	return root
}

// censusPath joins one root-relative path and creates the directories that lead to it.
func censusPath(t *testing.T, root string, elements ...string) string {
	t.Helper()
	path := filepath.Join(append([]string{root}, elements...)...)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeCensusFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// resolveCensus renders an expected census against the root its case ran in. A directory
// reads as rootToken plus a slash path; an environment entry carries the token inside its
// value.
func resolveCensus(root string, want []TestExecution) []TestExecution {
	var resolved []TestExecution
	for _, entry := range want {
		entry.Dir = resolveRoot(root, entry.Dir)
		entry.GoC = resolveRoot(root, entry.GoC)
		var env []string
		for _, value := range entry.Env {
			env = append(env, resolveRoot(root, value))
		}
		entry.Env = env
		resolved = append(resolved, entry)
	}
	return resolved
}

func resolveRoot(root, value string) string {
	if rest, found := strings.CutPrefix(value, rootToken); found {
		return filepath.Join(root, filepath.FromSlash(rest))
	}
	return strings.ReplaceAll(value, rootToken, root)
}

// assertCensus compares the census to the expected sets. normalizeTags returns a
// non-nil slice for every set, so a deep comparison is exact here.
func assertCensus(t *testing.T, got, want []TestExecution) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("census = %v, want %v", got, want)
	}
}
