//go:build system

package systemtest

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	toonlib "github.com/toon-format/toon-go"

	"github.com/gibbonmi/bench/internal/compatibility"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	"github.com/gibbonmi/bench/internal/toon"
)

func TestCompatibilityMissingPath(t *testing.T) {
	result, repo, wrapper := runCompatibilityDoctor(t)
	if result.code != 1 {
		t.Fatalf("doctor --compat without global Bench = %d, want 1; stdout=%q stderr=%q", result.code, result.stdout, result.stderr)
	}
	for _, want := range []string{"global-bench,failed", wrapper, filepath.Clean(repo)} {
		if !strings.Contains(result.stdout, want) {
			t.Fatalf("doctor --compat output missing %q: %q", want, result.stdout)
		}
	}
}

func TestCompatibilityPaths(t *testing.T) {
	result, repo, _ := runCompatibilityDoctor(t)
	expected, err := toon.Table("context", []string{"field", "value", "source"}, [][]string{{"repository", filepath.Clean(repo), "git root"}})
	if err != nil {
		t.Fatal(err)
	}
	_, row, found := strings.Cut(expected, "\n")
	if !found || !strings.Contains(result.stdout, row) {
		t.Fatalf("doctor --compat did not reach repository %q: (%d, %q, %q)", repo, result.code, result.stdout, result.stderr)
	}
}

type compatibilitySession struct {
	repo, wrapper, working string
	homes                  map[compatibility.Interface]string
	overrides              []string
}

func runCompatibilityDoctor(t *testing.T) (processResult, string, string) {
	fixture := newCompatibilitySession(t)
	return fixture.run(t, compatibility.CodexCLI), fixture.repo, fixture.wrapper
}

func newCompatibilitySession(t *testing.T) compatibilitySession {
	t.Helper()
	repo := owner.repos[0]
	if result := owner.runSelected(repo, "link", "copy"); result.code != 0 {
		t.Fatalf("link fixture = %d; stdout=%q stderr=%q", result.code, result.stdout, result.stderr)
	}
	working := filepath.Join(repo, "path with spaces [and] glob*")
	if err := os.MkdirAll(working, 0o755); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	fixture := compatibilitySession{
		repo: repo, working: working, wrapper: filepath.Join(repo, ".bench", "bin", "bench.sh"),
		homes: map[compatibility.Interface]string{},
		overrides: []string{
			"BENCH_HOME=" + filepath.Join(home, ".bench"), "BENCH_KIT=" + owner.kit,
			"BENCH_RUN_BINARY=" + owner.selected.path, "HOME=" + home,
			"PATH=" + privateToolPath(t, "git", "bash", "uname", "dirname", "basename", "readlink", "tr"),
		},
	}
	for _, selected := range compatibility.Interfaces() {
		configHome := filepath.Join(home, string(selected)+" home [personal]*")
		if selected == compatibility.CodexCLI {
			configHome = filepath.Join(home, ".codex")
		}
		fixture.homes[selected] = configHome
		if err := os.MkdirAll(configHome, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(configHome, "config.toml"), []byte("model = \"personal\"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, configHome := range fixture.homes {
		nested := filepath.Join(configHome, "personal", "nested")
		if err := os.MkdirAll(nested, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(nested, "sentinel"), []byte("preserve personal state\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return fixture
}

func (fixture compatibilitySession) run(t *testing.T, selected compatibility.Interface) processResult {
	t.Helper()
	return fixture.runWithHome(t, selected, fixture.homes[selected])
}

func (fixture compatibilitySession) runWithHome(t *testing.T, selected compatibility.Interface, home string) processResult {
	t.Helper()
	homeOverride := "CODEX_HOME"
	if home != "" {
		homeOverride += "=" + home
	}
	overrides := append(append([]string{}, fixture.overrides...), homeOverride)
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(fixture.working, overrides, "bash", fixture.wrapper, "doctor", "--compat", string(selected))
}

func compatibilityContext(t *testing.T, result processResult) map[string]string {
	t.Helper()
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("compatibility result = %#v", result)
	}
	output := result.stdout
	if spill, ok := responseboundtest.Find(output); ok {
		output = readSpillFile(t, spill.Path)
	}
	decoded, err := toonlib.DecodeString(output)
	if err != nil {
		t.Fatal(err)
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		t.Fatalf("compatibility document = %T", decoded)
	}
	rows, ok := document["context"].([]any)
	if !ok {
		t.Fatalf("compatibility context = %T", document["context"])
	}
	values := map[string]string{}
	for _, row := range rows {
		cells, ok := row.(map[string]any)
		if !ok {
			t.Fatalf("compatibility row = %T", row)
		}
		field, fieldOK := cells["field"].(string)
		value, valueOK := cells["value"].(string)
		if !fieldOK || !valueOK {
			t.Fatalf("compatibility row = %#v", cells)
		}
		values[field] = value
	}
	return values
}

func TestCompatibilityCollectorInterfaces(t *testing.T) {
	fixture := newCompatibilitySession(t)
	for _, selected := range compatibility.Interfaces() {
		t.Run(string(selected), func(t *testing.T) {
			context := compatibilityContext(t, fixture.run(t, selected))
			if context["interface"] != string(selected) {
				t.Fatalf("selected interface = %q, want %q", context["interface"], selected)
			}
		})
	}
}

func TestCompatibilityCollectorReadOnly(t *testing.T) {
	fixture := newCompatibilitySession(t)
	before := map[string]map[string]compatibilityHomeEntry{}
	for _, home := range fixture.homes {
		before[home] = snapshotCompatibilityHome(t, home)
	}
	type homeCase struct {
		name     string
		selected compatibility.Interface
		home     string
	}
	cases := []homeCase{}
	for _, selected := range compatibility.Interfaces() {
		cases = append(cases, homeCase{string(selected), selected, fixture.homes[selected]})
	}
	cases = append(cases, homeCase{"CLI HOME fallback", compatibility.CodexCLI, ""})
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			context := compatibilityContext(t, fixture.runWithHome(t, test.selected, test.home))
			if context["configuration-home"] != fixture.homes[test.selected] {
				t.Errorf("selected home = %q, want %q", context["configuration-home"], fixture.homes[test.selected])
			}
			for home, want := range before {
				if got := snapshotCompatibilityHome(t, home); !reflect.DeepEqual(got, want) {
					t.Errorf("inspection changed configuration home %s: got %#v, want %#v", home, got, want)
				}
			}
		})
	}
}

type compatibilityHomeEntry struct {
	mode    os.FileMode
	content string
}

func snapshotCompatibilityHome(t *testing.T, home string) map[string]compatibilityHomeEntry {
	t.Helper()
	entries := map[string]compatibilityHomeEntry{}
	err := filepath.WalkDir(home, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(home, path)
		if err != nil {
			return err
		}
		var content string
		switch {
		case info.Mode().IsRegular():
			data, readErr := os.ReadFile(path)
			content, err = string(data), readErr
		case info.Mode()&os.ModeSymlink != 0:
			content, err = os.Readlink(path)
		}
		entries[rel] = compatibilityHomeEntry{info.Mode(), content}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return entries
}

func TestCompatibilityCollectorPolicyFingerprint(t *testing.T) {
	fixture := newCompatibilitySession(t)
	hookPath := filepath.Join(fixture.repo, ".codex", "hooks.json")
	original, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(hookPath, original, 0o644); err != nil {
			t.Error(err)
		}
	})
	previous := compatibilityContext(t, fixture.run(t, compatibility.CodexCLI))["compatibility-fingerprint"]
	for _, policy := range []struct{ path, body string }{
		{hookPath, `{"hooks":{"SessionStart":[]}}`},
		{filepath.Join(fixture.homes[compatibility.CodexCLI], "config.toml"), "approval_policy = \"untrusted\"\n"},
	} {
		if err := os.WriteFile(policy.path, []byte(policy.body), 0o600); err != nil {
			t.Fatal(err)
		}
		current := compatibilityContext(t, fixture.run(t, compatibility.CodexCLI))["compatibility-fingerprint"]
		if current == "" || current == previous {
			t.Errorf("policy change at %s retained fingerprint %q", policy.path, current)
		}
		previous = current
	}
}
