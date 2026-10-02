package compatibility

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestCompatibilityFingerprint(t *testing.T) {
	baseline := Context{
		Interface:         CodexCLI,
		Repository:        Fact{Value: "/workspace/bench", Source: "git"},
		Environment:       Fact{Value: "linux/wsl2", Source: "runtime"},
		ConfigurationHome: Fact{Value: "/home/user/.codex", Source: "HOME"},
		ActiveRuntime:     Fact{Value: "runtime-a", Source: "process"},
		PolicyProvenance:  Fact{Value: "hooks-v1", Source: "repository"},
	}
	want := Fingerprint(baseline)
	if want == "" {
		t.Fatal("baseline fingerprint is empty")
	}
	for name, mutate := range map[string]func(*Context){
		"runtime":   func(c *Context) { c.ActiveRuntime.Value = "runtime-b" },
		"config":    func(c *Context) { c.ConfigurationHome.Value = "/mnt/c/Users/user/.codex" },
		"workspace": func(c *Context) { c.Repository.Value = "/workspace/other" },
		"policy":    func(c *Context) { c.PolicyProvenance.Value = "hooks-v2" },
	} {
		t.Run(name, func(t *testing.T) {
			changed := baseline
			mutate(&changed)
			if got := Fingerprint(changed); got == want {
				t.Fatalf("fingerprint after %s change = %q, want a different value", name, got)
			}
		})
	}
}

func TestCompatibilityRuntimeProvenance(t *testing.T) {
	input := compatibilityInput(t)
	input.Context.ActiveRuntime = Fact{Source: "active process"}
	input.Context.LauncherVersion = Fact{Value: "codex-cli 1.2.3", Source: "PATH launcher"}
	report := Inspect(input)
	if got := contextValue(report, "active-runtime"); got != "unknown" {
		t.Fatalf("active runtime = %q, want unknown", got)
	}
	if got := contextValue(report, "launcher-version"); got != "codex-cli 1.2.3" {
		t.Fatalf("launcher version = %q", got)
	}
}

func TestCompatibilityHomes(t *testing.T) {
	cli := compatibilityInput(t)
	cli.Context.Interface = CodexCLI
	cli.Context.ConfigurationHome = Fact{Value: "/home/user/.codex", Source: "HOME default"}
	desktop := cli
	desktop.Context.Interface = CodexDesktop
	desktop.Context.ConfigurationHome = Fact{Value: "/mnt/c/Users/user/.codex", Source: "CODEX_HOME"}
	if got := contextValue(Inspect(cli), "configuration-home"); got != "/home/user/.codex" {
		t.Fatalf("CLI configuration home = %q", got)
	}
	if got := contextValue(Inspect(desktop), "configuration-home"); got != "/mnt/c/Users/user/.codex" {
		t.Fatalf("desktop configuration home = %q", got)
	}
}

func TestCompatibilityUnknownConfig(t *testing.T) {
	input := compatibilityInput(t)
	input.Configuration = FileFact{Path: "/config.toml", Source: "desktop config", State: bounds.StateUnreadable, Reason: "permission denied"}
	row := Inspect(input).Checks[0]
	if row.State != StateUnknown || !strings.Contains(row.Action, "desktop config") {
		t.Fatalf("configuration row = %#v, want unknown with source", row)
	}
}

func TestCompatibilityPresentConfigIsNotEffective(t *testing.T) {
	input := compatibilityInput(t)
	row := Inspect(input).Checks[0]
	if row.State != StateUnknown || !strings.Contains(row.Action, "effective") {
		t.Fatalf("readable declaration produced %#v; effective settings remain unknown", row)
	}
}

func TestCompatibilityDeclaredHook(t *testing.T) {
	input := compatibilityInput(t)
	input.HookDeclared = true
	report := Inspect(input)
	if len(report.Live) == 0 || report.Live[0].Capability != "hook-behavior" {
		t.Fatalf("live rows = %#v, want hook behavior", report.Live)
	}
	for _, row := range report.Checks {
		if row.Check == "hook-behavior" && row.State == StateOK {
			t.Fatalf("declared hook became a local pass: %#v", row)
		}
	}
}

func TestCompatibilityFileKinds(t *testing.T) {
	root := t.TempDir()
	fifo := filepath.Join(root, "settings.fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	dangling := filepath.Join(root, "dangling.toml")
	if err := os.Symlink("missing.toml", dangling); err != nil {
		t.Fatal(err)
	}
	for name, path := range map[string]string{"fifo": fifo, "dangling": dangling} {
		t.Run(name, func(t *testing.T) {
			fact := ReadFile(path, name+" source")
			if fact.State != bounds.StateWrongType && fact.State != bounds.StateUnreadable {
				t.Fatalf("file state = %q, want bounded refusal", fact.State)
			}
		})
	}
}

func TestCompatibilityAbsentEmpty(t *testing.T) {
	root := t.TempDir()
	absent := ReadFile(filepath.Join(root, "absent.toml"), "config")
	emptyPath := filepath.Join(root, "empty.toml")
	if err := os.WriteFile(emptyPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	empty := ReadFile(emptyPath, "config")
	if absent.State != bounds.StateAbsent || empty.State != bounds.StateEmpty || absent.State == empty.State {
		t.Fatalf("absent/empty states = %q/%q", absent.State, empty.State)
	}
	if a, e := configurationCheck(absent).Action, configurationCheck(empty).Action; a == e {
		t.Fatalf("absent and empty actions both = %q", a)
	}
}

func TestCompatibilityPersonalSettings(t *testing.T) {
	root := t.TempDir()
	paths := []string{filepath.Join(root, "one.toml"), filepath.Join(root, "two.toml")}
	for i, body := range []string{"model = \"alpha\"\n", "model = \"beta\"\n"} {
		if err := os.WriteFile(paths[i], []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	left := compatibilityInput(t)
	right := left
	left.Configuration = ReadFile(paths[0], "config")
	right.Configuration = ReadFile(paths[1], "config")
	if Fingerprint(left.Context) != Fingerprint(right.Context) {
		t.Fatal("personal model preference changed compatibility fingerprint")
	}
	for _, row := range Inspect(left).Checks {
		if row.State == StateFailed {
			t.Fatalf("personal settings produced a failed check: %#v", row)
		}
	}
	if !reflect.DeepEqual(Inspect(left).Checks, Inspect(right).Checks) {
		t.Fatal("personal model preference created a compatibility conflict")
	}
}

func TestCompatibilityOutput(t *testing.T) {
	input := compatibilityInput(t)
	report := Inspect(input)
	report.Context[0].Value = "unsafe\x1bcell"
	if output, err := report.Render(); err == nil || output != "" {
		t.Fatalf("control-bearing render = (%q, %v), want bounded refusal", output, err)
	}
	for _, value := range []string{"01abc", "tab\tcell", "line\ncell", "return\rcell"} {
		report := Inspect(input)
		report.Context[0].Value = value
		if _, err := report.Render(); err != nil {
			t.Fatalf("Render(%q): %v", value, err)
		}
	}
}

func compatibilityInput(t *testing.T) Input {
	t.Helper()
	config := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(config, []byte("hooks = true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return Input{
		Context: Context{
			Interface:         CodexCLI,
			Repository:        Fact{Value: "/workspace/bench", Source: "git"},
			Environment:       Fact{Value: "linux/wsl2", Source: "runtime"},
			ConfigurationHome: Fact{Value: filepath.Dir(config), Source: "HOME default"},
			ActiveRuntime:     Fact{Source: "active process"},
			LauncherVersion:   Fact{Value: "codex-cli 1.2.3", Source: "PATH launcher"},
			PolicyProvenance:  Fact{Value: "hooks-v1", Source: "repository"},
		},
		Configuration: ReadFile(config, "config.toml"),
		GlobalBench:   true,
		HookAction:    "observe the hook in the selected interface",
	}
}

func contextValue(report Report, field string) string {
	for _, row := range report.Context {
		if row.Field == field {
			return row.Value
		}
	}
	return ""
}
