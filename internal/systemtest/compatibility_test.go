//go:build system

package systemtest

import (
	"github.com/gibbonmi/bench/internal/adopt/transaction"
	"github.com/gibbonmi/bench/internal/compatibility/compatibilitytest"
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

func compatibilityOutput(t *testing.T, result processResult) string {
	t.Helper()
	output := result.stdout
	if spill, ok := responseboundtest.Find(output); ok {
		output = readSpillFile(t, spill.Path)
	}
	return output
}

func compatibilityContext(t *testing.T, result processResult) map[string]string {
	t.Helper()
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("compatibility result = %#v", result)
	}
	output := compatibilityOutput(t, result)
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
	before := map[string]map[string]compatibilitytest.HomeEntry{}
	for _, home := range fixture.homes {
		before[home] = compatibilitytest.SnapshotHome(t, home)
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
				if got := compatibilitytest.SnapshotHome(t, home); !reflect.DeepEqual(got, want) {
					t.Errorf("inspection changed configuration home %s: got %#v, want %#v", home, got, want)
				}
			}
		})
	}
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

func TestCompatibilitySetup(t *testing.T) {
	fixture := newCompatibilitySession(t)
	before := map[string]map[string]compatibilitytest.HomeEntry{}
	for _, home := range fixture.homes {
		before[home] = compatibilitytest.SnapshotHome(t, home)
	}
	hook := filepath.Join(fixture.repo, ".codex", "hooks.json")
	data, err := os.ReadFile(hook)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(hook, data, 0o644); err != nil {
			t.Error(err)
		}
	})
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	result := owner.runAt(fixture.repo, fixture.overrides, owner.selected.path, "setup", "--yes")
	for home, want := range before {
		if got := compatibilitytest.SnapshotHome(t, home); !reflect.DeepEqual(got, want) {
			t.Errorf("setup changed personal home %s", home)
		}
	}
	installed, err := os.ReadFile(hook)
	if err != nil || string(installed) != string(data) {
		t.Fatalf("setup failed to install shared integration: %v, %#v", err, result)
	}
	output := compatibilityOutput(t, result)
	if !strings.Contains(output, "interface qualification remains pending") {
		t.Fatalf("setup omitted pending qualification: %#v", result)
	}
}

func TestCompatibilityInterruptedRepair(t *testing.T) {
	fixture := newCompatibilitySession(t)
	paths := []string{".bench/BENCH.md", ".bench/BENCH-reference.md", ".codex/hooks.json"}
	for _, rel := range paths {
		path := filepath.Join(fixture.repo, rel)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Error(err)
			}
		})
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	overrides := append(append([]string{}, fixture.overrides...), "BENCH_LINK_FAULT=interrupt:3")
	result := owner.runAt(fixture.repo, overrides, owner.selected.path, "doctor", "--compat", "codex-cli", "--fix")
	if result.code != 97 {
		t.Fatalf("repair did not interrupt at publication: %#v", result)
	}
	touched := 0
	for _, rel := range paths {
		if _, err := os.Stat(filepath.Join(fixture.repo, rel)); err == nil {
			touched++
		}
	}
	if touched != 2 {
		t.Fatalf("interruption did not leave two published targets: %d", touched)
	}
	home := filepath.Dir(fixture.homes[compatibility.CodexCLI])
	records, err := filepath.Glob(filepath.Join(home, ".bench", "compatibility-repairs", "*", "*", "record.json"))
	if err != nil || len(records) != 1 {
		t.Fatalf("interrupted recovery records: %v, %v", records, err)
	}
	id := filepath.Base(filepath.Dir(records[0]))
	result = owner.runAt(fixture.repo, fixture.overrides, owner.selected.path, "doctor", "--compat", "codex-cli", "--undo", id)
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("fresh-process recovery failed: %#v", result)
	}
	for _, rel := range paths {
		if _, err := os.Lstat(filepath.Join(fixture.repo, rel)); !os.IsNotExist(err) {
			t.Errorf("fresh-process undo did not restore absence: %s: %v", rel, err)
		}
	}
}

func TestCompatibilityConcurrentWriters(t *testing.T) {
	fixture := newCompatibilitySession(t)
	hook := filepath.Join(fixture.repo, ".git", "hooks", "pre-push")
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	fixed := owner.runAt(fixture.repo, fixture.overrides, owner.selected.path, "doctor", "--compat", "codex-cli", "--fix")
	if fixed.code != 1 || fixed.stderr != "" {
		t.Fatalf("repair fixture: %#v", fixed)
	}
	home := filepath.Dir(fixture.homes[compatibility.CodexCLI])
	records, err := filepath.Glob(filepath.Join(home, ".bench", "compatibility-repairs", "*", "*", "record.json"))
	if err != nil || len(records) != 1 {
		t.Fatalf("repair fixture records: %v, %v", records, err)
	}
	id := filepath.Base(filepath.Dir(records[0]))
	alias := filepath.Join(t.TempDir(), "shared-hooks")
	if err := os.Symlink(filepath.Dir(hook), alias); err != nil {
		t.Fatal(err)
	}
	other := owner.repos[1]
	compatibilitytest.RemovePreserving(t, filepath.Join(other, ".bench", "gate.sh"))
	systemGitOutput(t, other, "config", "core.hooksPath", alias)
	t.Cleanup(func() { systemGitOutput(t, other, "config", "--unset", "core.hooksPath") })
	if result := owner.runSelected(other, "link", "copy"); result.code != 0 {
		t.Fatalf("second repo fixture: %#v", result)
	}
	for _, repo := range []string{fixture.repo, other} {
		manifest := filepath.Join(repo, ".bench", "link-manifest.tsv")
		data, err := os.ReadFile(manifest)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.WriteFile(manifest, data, 0o644); err != nil {
				t.Error(err)
			}
		})
		_, body, _ := strings.Cut(string(data), "\n")
		if err := os.WriteFile(manifest, []byte("#kit\t0.0.0\n"+body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lease, err := transaction.Lock([]string{hook})
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	before := compatibilitytest.SnapshotHome(t, filepath.Dir(hook))
	for _, repo := range []string{fixture.repo, other} {
		for _, args := range [][]string{{"setup", "--yes"}, {"link", "copy"}, {"upgrade"}, {"unlink"}, {"doctor", "--fix"}, {"doctor", "--compat", "codex-cli", "--fix"}} {
			overrides := append(append([]string{}, fixture.overrides...), "BENCH_HOME="+t.TempDir(), "BENCH_WRAPPER="+fixture.wrapper)
			result := owner.runAt(repo, overrides, owner.selected.path, args...)
			if result.code == 0 || !strings.Contains(compatibilityOutput(t, result)+result.stderr, "competing writer") {
				t.Errorf("writer bypassed shared exclusion: %s %v: code=%d stderr=%q stdout=%q", repo, args, result.code, result.stderr, result.stdout)
			}
		}
	}
	result := owner.runAt(fixture.repo, fixture.overrides, owner.selected.path, "doctor", "--compat", "codex-cli", "--undo", id)
	if result.code != 1 || !strings.Contains(compatibilityOutput(t, result)+result.stderr, "competing writer") {
		t.Errorf("undo bypassed shared exclusion: %#v", result)
	}
	if after := compatibilitytest.SnapshotHome(t, filepath.Dir(hook)); !reflect.DeepEqual(before, after) {
		t.Error("competing writer changed the shared hook")
	}
	if result := owner.runSelected(owner.repos[2], "link", "copy"); result.code != 0 {
		t.Errorf("disjoint writer was blocked: %#v", result)
	}
}
