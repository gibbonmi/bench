package adopt

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setup_test.go pins FT227's seeded gate input manifest: `bench setup` leaves every
// adopted repository with the environment declaration its installed wrapper needs.

// wantSeededGateInputs is authored independently of scaffoldGateInputs on purpose.
// Dropping a name from the seeded environment list reds nothing else in the tree. The
// gate that would notice lives in the adopted repository, not here. This is the one place
// outside the function where the bytes are spelled out.
const wantSeededGateInputs = `{
  "schema": 1,
  "closure": "local",
  "environment": ["BENCH_HOME", "BENCH_KIT", "BENCH_RUN_BINARY", "HOME"],
  "paths": [],
  "tools": ["bash", "basename", "dirname", "git", "readlink", "uname"]
}
`

// runSetupYes runs a non-interactive converge and fails on anything but the two success
// shapes. Those shapes are 0, or 3 for a zero-signal repository whose fail-closed gate
// stub leaves the doctor legitimately red. Neither says anything about the seed itself;
// the callers assert that on the written bytes.
func runSetupYes(t *testing.T) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	notTTY := func(io.Reader) bool { return false }
	code := setup([]string{"--yes"}, strings.NewReader(""), &stdout, &stderr, "1.0.0", notTTY)
	if code != 0 && code != 3 {
		t.Fatalf("setup --yes exit = %d, stderr:\n%s", code, stderr.String())
	}
}

func readGateInputs(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".bench", "gate-inputs.json"))
	if err != nil {
		t.Fatalf("read .bench/gate-inputs.json: %v", err)
	}
	return string(data)
}

// TestSetupSeedsGateInputs pins SD1.
func TestSetupSeedsGateInputs(t *testing.T) {
	root := setupPromptTestRepo(t)
	runSetupYes(t)
	if got := readGateInputs(t, root); got != wantSeededGateInputs {
		t.Fatalf("seeded gate-inputs.json =\n%s\nwant:\n%s", got, wantSeededGateInputs)
	}
}

// TestSetupPreservesOperatorGateInputs pins SD3: an operator-authored manifest
// survives byte-identical and is never recorded as a managed row.
func TestSetupPreservesOperatorGateInputs(t *testing.T) {
	root := setupPromptTestRepo(t)
	operator := "{\n  \"schema\": 1,\n  \"closure\": \"local\",\n  \"environment\": [\"HOME\"],\n  \"paths\": [],\n  \"tools\": [\"bash\"]\n}\n"
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(operator), 0o644); err != nil {
		t.Fatal(err)
	}
	runSetupYes(t)
	if got := readGateInputs(t, root); got != operator {
		t.Fatalf("operator-authored gate-inputs.json was rewritten:\n%s", got)
	}
	manifest, err := os.ReadFile(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		t.Fatalf("read link-manifest.tsv: %v", err)
	}
	if strings.Contains(string(manifest), "gate-inputs.json") {
		t.Fatalf("gate-inputs.json recorded in link-manifest.tsv:\n%s", manifest)
	}
}

// TestSetupTwiceLeavesGateInputsIdentical pins SD4.
func TestSetupTwiceLeavesGateInputsIdentical(t *testing.T) {
	root := setupPromptTestRepo(t)
	runSetupYes(t)
	first := readGateInputs(t, root)
	runSetupYes(t)
	if second := readGateInputs(t, root); second != first {
		t.Fatalf("second setup rewrote gate-inputs.json:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// TestSetupSeedsGateInputsInZeroSignalRepo pins SD5: the seed does not depend on a
// detected ecosystem — the fail-closed-stub path gets it too.
func TestSetupSeedsGateInputsInZeroSignalRepo(t *testing.T) {
	root := setupBareTestRepo(t)
	for _, name := range []string{"go.mod", "Makefile", "package.json", "Cargo.toml"} {
		if fileExists(filepath.Join(root, name)) {
			t.Fatalf("zero-signal fixture unexpectedly carries %s", name)
		}
	}
	runSetupYes(t)
	if got := readGateInputs(t, root); got != wantSeededGateInputs {
		t.Fatalf("zero-signal seeded gate-inputs.json =\n%s\nwant:\n%s", got, wantSeededGateInputs)
	}
}

func TestSetupGateRejectsIgnoredDeclaredInputs(t *testing.T) {
	for _, test := range []struct {
		name, declared, want string
		red                  bool
	}{
		{name: "space", declared: `"ignored\u0020file"`, want: "gate input path ignored file is gitignored", red: true},
		{name: "glob", declared: `"literal[*].txt"`, want: "gate input path literal[*].txt is gitignored", red: true},
		{name: "undeclared", declared: "", want: "gate: green"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := setupBareTestRepo(t)
			runSetupYes(t)

			gate := filepath.Join(root, ".bench", "gate.sh")
			data, err := os.ReadFile(gate)
			if err != nil {
				t.Fatalf("read gate: %v", err)
			}
			if err := os.WriteFile(gate, []byte(strings.Replace(string(data), "err \"configure .bench/gate.sh - replace this sentinel with real checks\"  # "+SentinelMarker+"\n", "", 1)), 0o755); err != nil {
				t.Fatalf("retire gate sentinel: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored file\nliteral[[][*][]].txt\nordinary-ignored\n"), 0o644); err != nil {
				t.Fatalf("write .gitignore: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, "package.json"), []byte("{}\n"), 0o644); err != nil {
				t.Fatalf("write package.json: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, "ordinary-ignored"), []byte("allowed\n"), 0o644); err != nil {
				t.Fatalf("write ordinary ignored file: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, "fixture_test.go"), []byte("package fixture\n"), 0o644); err != nil {
				t.Fatalf("write fixture package: %v", err)
			}
			manifest := `{"schema":1,"closure":"local","environment":["BENCH_HOME","HOME"],"paths":[` + test.declared + `],"tools":["bash","basename","dirname","git","readlink","uname"]}` + "\n"
			if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(manifest), 0o644); err != nil {
				t.Fatalf("write gate inputs: %v", err)
			}

			cmd := exec.Command("bash", gate)
			cmd.Dir = root
			out, err := cmd.CombinedOutput()
			if (err != nil) != test.red || !strings.Contains(string(out), test.want) {
				t.Fatalf("gate = (%v, %q), want red=%t and %q", err, out, test.red, test.want)
			}
		})
	}
}

func TestSetupGateDeclaresConsumerOnlyValidityScope(t *testing.T) {
	root := setupBareTestRepo(t)
	runSetupYes(t)
	gate, err := os.ReadFile(filepath.Join(root, ".bench", "gate.sh"))
	if err != nil {
		t.Fatal(err)
	}
	want := `BENCH_CONFORMANCE_CONSUMER_ONLY=1 "$bench" test --check load-validity-metadata`
	if !strings.Contains(string(gate), want) {
		t.Fatalf("scaffolded gate missing consumer-only validity intent %q:\n%s", want, gate)
	}
}
