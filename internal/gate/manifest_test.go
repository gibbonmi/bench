package gate

// The manifest loader at the PhasesCommand seam: which table ran, and what a
// refused manifest prints. Tests assert observable effects — marker files, exit
// codes, stderr substrings — never loader internals.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// stubBuiltinTable installs a built-in table whose one phase writes marker, so a test
// can tell the fallback table apart from a manifest-declared one by that file alone.
func stubBuiltinTable(t *testing.T, marker string) {
	t.Helper()
	original := benchkitPhasesForCommand
	t.Cleanup(func() { benchkitPhasesForCommand = original })
	benchkitPhasesForCommand = func(root, kit string) []Phase {
		return []Phase{{Name: "builtin", Argv: []string{"touch", marker}}}
	}
}

// runManifest runs the command against a fresh root carrying manifest, and reports the
// exit code, stderr, and whether the built-in fallback table ran.
func runManifest(t *testing.T, manifest string) (code int, stderr string, fellBack bool) {
	t.Helper()
	root := t.TempDir()
	marker := filepath.Join(root, "builtin.marker")
	stubBuiltinTable(t, marker)
	writeFile(t, manifestPath(root), manifest)

	var outBuf, errBuf bytes.Buffer
	code = PhasesCommand([]string{root}, &outBuf, &errBuf)
	_, err := os.Stat(marker)
	return code, errBuf.String(), err == nil
}

func TestPhasesCommandLoadsManifest(t *testing.T) {
	root := t.TempDir()
	builtin := filepath.Join(root, "builtin.marker")
	declared := filepath.Join(root, "declared.marker")
	stubBuiltinTable(t, builtin)
	writeFile(t, manifestPath(root), fmt.Sprintf(`{"phases":[{"name":"declared","argv":["touch",%q]}]}`, declared))

	var stdout, stderr bytes.Buffer
	if code := PhasesCommand([]string{root}, &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(declared); err != nil {
		t.Fatalf("the manifest's phase did not run: %v; stdout=%q", err, stdout.String())
	}
	if _, err := os.Stat(builtin); err == nil {
		t.Fatalf("the built-in table ran alongside a present manifest; stdout=%q", stdout.String())
	}
}

func TestPhasesCommandAbsentManifestFallsBack(t *testing.T) {
	root := t.TempDir()
	builtin := filepath.Join(root, "builtin.marker")
	stubBuiltinTable(t, builtin)

	var stdout, stderr bytes.Buffer
	if code := PhasesCommand([]string{root}, &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(builtin); err != nil {
		t.Fatalf("the built-in table did not run without a manifest: %v; stderr=%q", err, stderr.String())
	}
}

func TestManifestEmptyIsRed(t *testing.T) {
	for _, tc := range []struct {
		name     string
		manifest string
	}{
		{name: "blank", manifest: ""},
		{name: "whitespace", manifest: " \n\t\n"},
		{name: "zero-phases", manifest: `{"phases":[]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, stderr, fellBack := runManifest(t, tc.manifest)
			if code != 1 {
				t.Fatalf("PhasesCommand = %d, want 1; stderr=%q", code, stderr)
			}
			if !strings.Contains(stderr, "empty manifest") {
				t.Fatalf("stderr = %q, want the empty-manifest class named", stderr)
			}
			if fellBack {
				t.Fatalf("an empty manifest fell back to the built-in table; stderr=%q", stderr)
			}
		})
	}
}

func TestManifestMalformed(t *testing.T) {
	type variant struct {
		manifest string
		element  string
	}
	for _, tc := range []struct {
		name     string
		class    string
		variants []variant
	}{
		{
			name:  "parse-error",
			class: "parse error",
			variants: []variant{
				{manifest: `{"phases":[{"name":5,"argv":["true"]}]}`, element: "name"},
				{manifest: `{"phases":[{"name":"a","argv":["true"]},]}`, element: "invalid character"},
				{manifest: `{"phases":[{"name":"a","argv":["true"]}]}{"phases":[]}`, element: "trailing content"},
			},
		},
		{
			// A repeated JSON key is valid JSON that silently shadows whatever the
			// earlier one said, so the phase a reader sees declared is not the phase
			// that runs.
			name:  "duplicate-json-key",
			class: "parse error",
			variants: []variant{
				{manifest: `{"phases":[{"name":"a","argv":["true"]}],"phases":[{"name":"b","argv":["true"]}]}`, element: "duplicate name"},
				{manifest: `{"phases":[{"name":"real","argv":["true"],"name":"shadow","argv":["false"]}]}`, element: "duplicate name"},
			},
		},
		{
			name:  "duplicate-name",
			class: "duplicate phase name",
			variants: []variant{
				{manifest: `{"phases":[{"name":"dup","argv":["true"]},{"name":"dup","argv":["true"]}]}`, element: "dup"},
			},
		},
		{
			name:  "invalid-name",
			class: "invalid phase name",
			variants: []variant{
				{manifest: `{"phases":[{"name":"","argv":["true"]}]}`, element: `""`},
				{manifest: `{"phases":[{"name":"two words","argv":["true"]}]}`, element: "two words"},
				{manifest: `{"phases":[{"name":"bel\u0007","argv":["true"]}]}`, element: `\a`},
			},
		},
		{
			name:  "dangling-needs",
			class: "dangling needs edge",
			variants: []variant{
				{manifest: `{"phases":[{"name":"a","argv":["true"],"needs":["ghost"]}]}`, element: "ghost"},
			},
		},
		{
			name:  "cyclic-needs",
			class: "cyclic needs edge",
			variants: []variant{
				{manifest: `{"phases":[{"name":"self","argv":["true"],"needs":["self"]}]}`, element: "self"},
				{manifest: `{"phases":[{"name":"a","argv":["true"],"needs":["b"]},{"name":"b","argv":["true"],"needs":["a"]}]}`, element: "a"},
			},
		},
		{
			name:  "empty-argv",
			class: "empty argv",
			variants: []variant{
				{manifest: `{"phases":[{"name":"noargv","argv":[]}]}`, element: "noargv"},
				{manifest: `{"phases":[{"name":"blankargv","argv":[""]}]}`, element: "blankargv"},
			},
		},
		{
			name:  "escaping-dir",
			class: "escaping dir",
			variants: []variant{
				{manifest: `{"phases":[{"name":"a","argv":["true"],"dir":"../outside"}]}`, element: "../outside"},
				{manifest: `{"phases":[{"name":"a","argv":["true"],"dir":"sub/../.."}]}`, element: "sub/../.."},
				{manifest: `{"phases":[{"name":"a","argv":["true"],"dir":"/etc"}]}`, element: "/etc"},
			},
		},
		{
			name:  "unknown-field-top",
			class: "unknown field",
			variants: []variant{
				{manifest: `{"phases":[{"name":"a","argv":["true"]}],"parallel":2}`, element: "parallel"},
			},
		},
		{
			name:  "unknown-field-phase",
			class: "unknown field",
			variants: []variant{
				{manifest: `{"phases":[{"name":"a","argv":["true"],"need":["b"]}]}`, element: "need"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, v := range tc.variants {
				code, stderr, fellBack := runManifest(t, v.manifest)
				if code != 1 {
					t.Fatalf("PhasesCommand = %d, want 1 for %s; stderr=%q", code, v.manifest, stderr)
				}
				if fellBack {
					t.Fatalf("a malformed manifest fell back to the built-in table for %s", v.manifest)
				}
				if !strings.Contains(stderr, tc.class) {
					t.Fatalf("stderr = %q for %s, want the class %q named", stderr, v.manifest, tc.class)
				}
				if !strings.Contains(stderr, v.element) {
					t.Fatalf("stderr = %q for %s, want the offending element %q named", stderr, v.manifest, v.element)
				}
			}
		})
	}
}

func TestManifestDanglingSymlinkIsRed(t *testing.T) {
	root := t.TempDir()
	builtin := filepath.Join(root, "builtin.marker")
	stubBuiltinTable(t, builtin)
	path := manifestPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.Symlink(filepath.Join(root, "gone.json"), path); err != nil {
		t.Fatalf("symlink %s: %v", path, err)
	}

	var stdout, stderr bytes.Buffer
	if code := PhasesCommand([]string{root}, &stdout, &stderr); code != 1 {
		t.Fatalf("PhasesCommand = %d, want 1; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "dangling symlink") || !strings.Contains(stderr.String(), "gone.json") {
		t.Fatalf("stderr = %q, want the dangling-symlink class and its target named", stderr.String())
	}
	if _, err := os.Stat(builtin); err == nil {
		t.Fatalf("a broken manifest link fell back to the built-in table")
	}
}

// TestManifestSpecialFileIsRed keeps its own deadline: a loader that opens the path
// before it stats blocks on the FIFO forever, and a hang reported as a hang names the
// defect far better than the package timeout killing the whole suite.
func TestManifestSpecialFileIsRed(t *testing.T) {
	root := t.TempDir()
	builtin := filepath.Join(root, "builtin.marker")
	stubBuiltinTable(t, builtin)
	path := manifestPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatalf("mkfifo %s: %v", path, err)
	}

	type outcome struct {
		code   int
		stderr string
	}
	done := make(chan outcome, 1)
	go func() {
		var stdout, stderr bytes.Buffer
		code := PhasesCommand([]string{root}, &stdout, &stderr)
		done <- outcome{code: code, stderr: stderr.String()}
	}()
	select {
	case got := <-done:
		if got.code != 1 {
			t.Fatalf("PhasesCommand = %d, want 1; stderr=%q", got.code, got.stderr)
		}
		if !strings.Contains(got.stderr, "not a regular file") || !strings.Contains(got.stderr, filepath.Join(".bench", "phases.json")) {
			t.Fatalf("stderr = %q, want the special-file class and the manifest path named", got.stderr)
		}
		if _, err := os.Stat(builtin); err == nil {
			t.Fatalf("a FIFO at the manifest path fell back to the built-in table")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("PhasesCommand blocked on a FIFO at the manifest path; the loader must stat before it reads")
	}
}

// TestPhasesCommandManifestFieldsEndToEnd drives all six fields through one manifest,
// each with an effect visible in the run: a loader that decodes only name and argv
// passes every direct-Phase runner test and fails here.
func TestPhasesCommandManifestFieldsEndToEnd(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", sub, err)
	}
	stubBuiltinTable(t, filepath.Join(root, "builtin.marker"))
	alpha := filepath.Join(root, "alpha.marker")
	// beta is declared first so its position cannot supply the ordering, needs alpha
	// twice so the dedupe stays silent, and proves alpha ran by testing for its marker.
	writeFile(t, manifestPath(root), fmt.Sprintf(`{"phases":[
		{"name":"beta","argv":["bash","-c","test -f %s && printf 'probe=%%s\\n' \"$PHASE_PROBE\" && pwd -P"],"needs":["alpha","alpha"],"env":{"PHASE_PROBE":"probe-value"},"dir":"sub"},
		{"name":"opt","argv":["bench-absent-binary-fixture"],"optional":true},
		{"name":"alpha","argv":["touch",%q]}
	]}`, alpha, alpha))

	var stdout, stderr bytes.Buffer
	if code := PhasesCommand([]string{root}, &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(alpha); err != nil {
		t.Fatalf("alpha did not run: %v; stdout=%q", err, stdout.String())
	}
	for _, want := range []string{
		"[beta] probe=probe-value\n",
		"[beta] " + resolved(t, sub) + "\n",
		"phase alpha: green\n",
		"phase beta: green\n",
		"phase opt: skipped (not installed)\n",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q, want it to carry %q", stdout.String(), want)
		}
	}
}

// TestManifestDirResolvesAgainstGradedRoot separates the two trees a linked repo keeps
// apart: the kit checkout the gate runs from and the root it grades. A manifest dir
// names a directory in the graded tree, and an undeclared dir means that tree's root, so
// the kit checkout must appear in neither phase's working directory.
func TestManifestDirResolvesAgainstGradedRoot(t *testing.T) {
	root := t.TempDir()
	kit := t.TempDir()
	t.Setenv("BENCH_KIT", kit)
	sub := filepath.Join(root, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", sub, err)
	}
	stubBuiltinTable(t, filepath.Join(root, "builtin.marker"))
	writeFile(t, manifestPath(root), `{"phases":[
		{"name":"declared","argv":["bash","-c","pwd -P"],"dir":"sub"},
		{"name":"default","argv":["bash","-c","pwd -P"]}
	]}`)

	var stdout, stderr bytes.Buffer
	if code := PhasesCommand([]string{root}, &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"[declared] " + resolved(t, sub) + "\n",
		"[default] " + resolved(t, root) + "\n",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q, want it to carry %q", stdout.String(), want)
		}
	}
	if strings.Contains(stdout.String(), resolved(t, kit)) {
		t.Fatalf("a manifest phase ran inside the kit checkout %s; stdout=%q", kit, stdout.String())
	}
}

func resolved(t *testing.T, path string) string {
	t.Helper()
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("resolve %s: %v", path, err)
	}
	return real
}

// TestPhasesCommandInnerModeManifestDagOrder pins the inner path to the same table
// source: the manifest is loaded, phasesForMode still drops the canary phase, edges
// still order the run, and the byte shape stays summary-free and prefix-free.
func TestPhasesCommandInnerModeManifestDagOrder(t *testing.T) {
	t.Setenv("BENCH_CANARY_INNER", "1")
	t.Setenv("BENCH_CANARY_PHASE", "")
	root := t.TempDir()
	builtin := filepath.Join(root, "builtin.marker")
	stubBuiltinTable(t, builtin)
	writeFile(t, manifestPath(root), `{"phases":[
		{"name":"beta","argv":["bash","-c","printf 'beta\\n'"],"needs":["alpha"]},
		{"name":"canary","argv":["bash","-c","printf 'canary\\n'"]},
		{"name":"alpha","argv":["bash","-c","printf 'alpha\\n'"]}
	]}`)

	var stdout, stderr bytes.Buffer
	if code := PhasesCommand([]string{root}, &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if got, want := stdout.String(), "alpha\nbeta\ngate: green\n"; got != want {
		t.Fatalf("inner stdout = %q, want %q", got, want)
	}
	if _, err := os.Stat(builtin); err == nil {
		t.Fatalf("inner mode ran the built-in table despite a present manifest")
	}
}
