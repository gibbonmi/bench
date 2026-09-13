package testrepo_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/testrepo"
)

// HP1: one declaration produces both requested scripts and canonical inputs.
func TestGateFixtureDerivesInputs(t *testing.T) {
	t.Parallel()
	root := filepath.Join(t.TempDir(), "repo [x] 'quoted'")
	f := testrepo.NewGateFixture(t.TempDir())
	f.Environment = []string{"FIXTURE_VALUE"}
	f.Paths = []string{"declared.txt"}
	body := "printf '%s' \"$FIXTURE_VALUE\"\n"
	if err := f.Write(root, body, body); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"gate.sh", "gate-prospective.sh"} {
		cmd := exec.Command(filepath.Join(root, ".bench", name))
		cmd.Env = append(os.Environ(), "FIXTURE_VALUE=shared body")
		if output, err := cmd.CombinedOutput(); err != nil || string(output) != "shared body" {
			t.Fatalf("%s = %q, %v, want shared body", name, output, err)
		}
	}
	want := "{\"schema\":1,\"closure\":\"local\",\"environment\":[\"FIXTURE_VALUE\"],\"paths\":[\"declared.txt\"],\"tools\":[]}\n"
	if got, err := os.ReadFile(filepath.Join(root, ".bench", "gate-inputs.json")); err != nil || string(got) != want {
		t.Fatalf("manifest = %q, %v, want %q", got, err, want)
	}
}

// HP3: a rich host cannot satisfy a command absent from the declaration.
func TestGateFixtureRestrictsAmbientPath(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	f := testrepo.NewGateFixture(filepath.Join(t.TempDir(), "tools [x] 'quoted'"))
	body := "set -eu\nprintf 'declared\\n' | " + f.Command("cat") + "\n"
	if err := f.Write(root, body, body+"sed --version\n"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		wantErr bool
	}{
		{"gate.sh", false}, {"gate-prospective.sh", true},
	} {
		output, err := exec.Command(filepath.Join(root, ".bench", tc.name)).CombinedOutput()
		if (err != nil) != tc.wantErr || !strings.HasPrefix(string(output), "declared\n") {
			t.Fatalf("%s = %q, %v, want failure %v", tc.name, output, err, tc.wantErr)
		}
	}
}

// HP2: duplicate requests still produce one sorted row of ambient commands.
func TestGateFixtureDerivesToolClosure(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	f := testrepo.NewGateFixture(t.TempDir())
	body := f.Command("sed") + " -n '1p' input\n" +
		f.Command("cat") + " input\n" + f.Command("sed") + " -n '1p' input\n"
	if err := f.Write(root, body, ""); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".bench", "gate-inputs.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inputs struct {
		Tools []string `json:"tools"`
	}
	if err := json.Unmarshal(data, &inputs); err != nil {
		t.Fatal(err)
	}
	if want := []string{"cat", "sed"}; !reflect.DeepEqual(inputs.Tools, want) {
		t.Fatalf("manifest tools = %q, want %q", inputs.Tools, want)
	}
}

// HP3: a reused directory cannot supply an undeclared command.
func TestGateFixtureRefusesContaminatedPath(t *testing.T) {
	t.Parallel()
	path := t.TempDir()
	if err := os.Symlink("/bin/sh", filepath.Join(path, "undeclared")); err != nil {
		t.Fatal(err)
	}
	f := testrepo.NewGateFixture(path)
	if err := f.Write(t.TempDir(), "undeclared -c 'exit 0'\n", ""); err == nil || !strings.Contains(err.Error(), "undeclared") {
		t.Fatalf("contaminated path error = %v, want undeclared-command refusal", err)
	}
}

func TestGateFixtureWritesOnlyRequestedPathsAndCanRepeat(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, ordinary, prospective string
	}{
		{"ordinary", "ordinary", ""}, {"prospective", "", "prospective"},
		{"distinct", "ordinary", "prospective"}, {"empty", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			f := testrepo.NewGateFixture(t.TempDir())
			body := func(value string) string {
				if value == "" {
					return ""
				}
				return "printf " + value + " | " + f.Command("cat") + "\n"
			}
			ordinary, prospective := body(tc.ordinary), body(tc.prospective)
			for range 2 {
				f.MustWrite(t, root, ordinary, prospective)
			}
			for name, want := range map[string]string{
				"gate.sh": tc.ordinary, "gate-prospective.sh": tc.prospective,
			} {
				path := filepath.Join(root, ".bench", name)
				if want == "" {
					if _, err := os.Stat(path); !os.IsNotExist(err) {
						t.Fatalf("unrequested %s: %v", name, err)
					}
					continue
				}
				if output, err := exec.Command(path).CombinedOutput(); err != nil || string(output) != want {
					t.Fatalf("%s = %q, %v, want %q", name, output, err, want)
				}
			}
		})
	}
}

func TestGateFixtureRefusesMissingHostCommand(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	root := t.TempDir()
	f := testrepo.NewGateFixture(t.TempDir())
	err := f.Write(root, f.Command("sed")+" --version\n", "")
	if err == nil || !strings.Contains(err.Error(), `ambient command "sed"`) {
		t.Fatalf("missing host resolution = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".bench", "gate.sh")); !os.IsNotExist(err) {
		t.Fatalf("missing tool published a script: %v", err)
	}
}

func TestGateFixtureRefusesNonAmbientNames(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", ".", "..", "/bin/sh", "../sh", "dir/sh", "bad\nname"} {
		f := testrepo.NewGateFixture(t.TempDir())
		err := f.Write(t.TempDir(), f.Command(name)+"\n", "")
		if err == nil || !strings.Contains(err.Error(), "invalid ambient command") {
			t.Fatalf("command %q error = %v, want ambient-name refusal", name, err)
		}
	}
}
