package adopt

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/compatibility"
)

func TestCompatibilityInterface(t *testing.T) {
	for _, selected := range []compatibility.Interface{compatibility.CodexCLI, compatibility.CodexDesktop} {
		t.Run(string(selected), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := compatibilityDoctor([]string{string(selected)}, &stdout, &stderr, func(got compatibility.Interface) compatibility.Input {
				if got != selected {
					t.Fatalf("selected interface = %q, want %q", got, selected)
				}
				return compatibilityFixture(selected)
			})
			if code != 1 || stderr.String() != "" || !strings.Contains(stdout.String(), "interface,"+string(selected)+",argument") {
				t.Fatalf("compatibility report = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestCompatibilityReadOnly(t *testing.T) {
	roots := []string{t.TempDir(), t.TempDir()}
	for i, root := range roots {
		if err := os.WriteFile(filepath.Join(root, "config.toml"), []byte("model = \"personal-"+string(rune('a'+i))+"\"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	before := []string{readCompatibilityFile(t, roots[0]), readCompatibilityFile(t, roots[1])}
	for i, selected := range []compatibility.Interface{compatibility.CodexCLI, compatibility.CodexDesktop} {
		var stdout, stderr bytes.Buffer
		input := compatibilityFixture(selected)
		input.Context.ConfigurationHome = compatibility.Fact{Value: roots[i], Source: "selected home"}
		input.Configuration = compatibility.ReadFile(filepath.Join(roots[i], "config.toml"), "selected config")
		if code := compatibilityDoctor([]string{string(selected)}, &stdout, &stderr, func(compatibility.Interface) compatibility.Input { return input }); code != 1 {
			t.Fatalf("doctor --compat %s = %d; stderr=%q", selected, code, stderr.String())
		}
	}
	for i, root := range roots {
		if got := readCompatibilityFile(t, root); got != before[i] {
			t.Fatalf("configuration home %d changed: %q", i, got)
		}
	}
}

func TestCompatibilityMissingAsset(t *testing.T) {
	input := compatibilityFixture(compatibility.CodexCLI)
	input.Assets = []compatibility.Asset{{
		Name:          ".codex/hooks.json",
		File:          compatibility.FileFact{State: bounds.StateAbsent},
		RestoreAction: "bench link",
	}}
	var stdout, stderr bytes.Buffer
	if code := compatibilityDoctor([]string{"codex-cli"}, &stdout, &stderr, func(compatibility.Interface) compatibility.Input { return input }); code != 1 {
		t.Fatalf("missing asset exit = %d, want 1", code)
	}
	if !strings.Contains(stdout.String(), ".codex/hooks.json") || !strings.Contains(stdout.String(), "failed,bench link") {
		t.Fatalf("missing asset output = %q", stdout.String())
	}
}

func TestCompatibilityGrammar(t *testing.T) {
	called := false
	for _, args := range [][]string{{"unknown"}, nil, {"codex-cli", "extra"}, {"--fix"}} {
		var stdout, stderr bytes.Buffer
		code := compatibilityDoctor(args, &stdout, &stderr, func(compatibility.Interface) compatibility.Input {
			called = true
			return compatibility.Input{}
		})
		if code != 2 {
			t.Fatalf("compatibilityDoctor(%q) = %d, want 2", args, code)
		}
	}
	if called {
		t.Fatal("invalid compatibility grammar read configuration")
	}
}

func TestCompatibilityInspectionBoundary(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := compatibilityFixture(compatibility.CodexCLI)
	input.HookDeclared = true
	if code := compatibilityDoctor([]string{"codex-cli"}, &stdout, &stderr, func(compatibility.Interface) compatibility.Input { return input }); code != 1 {
		t.Fatalf("unknown effective configuration exit = %d; stderr=%q", code, stderr.String())
	}
	if strings.Contains(strings.ToLower(stdout.String()), "qualified") || !strings.Contains(stdout.String(), "live[") {
		t.Fatalf("inspection boundary output = %q", stdout.String())
	}
}

func TestCompatibilityLegacyDoctor(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"--bogus"}, {"--fix"}} {
		legacy := runLegacyFixture(t, args, legacyDoctor)
		current := runLegacyFixture(t, args, Doctor)
		if legacy != current {
			t.Fatalf("Doctor(%q) = %#v, pre-compatibility baseline = %#v", args, current, legacy)
		}
	}
}

type doctorResult struct {
	Code           int
	Stdout, Stderr string
}

func runLegacyFixture(t *testing.T, args []string, invoke func([]string, io.Writer, io.Writer, string) int) doctorResult {
	t.Helper()
	root := t.TempDir()
	wrapper := filepath.Join(root, "bench-wrapper")
	if err := os.WriteFile(wrapper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	t.Setenv("HOME", root)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench"))
	t.Setenv("BENCH_WRAPPER", wrapper)
	t.Setenv("PATH", "/usr/bin:/bin")
	t.Setenv("NVM_DIR", "")
	t.Setenv("ASDF_DATA_DIR", "")
	t.Setenv("FNM_DIR", "")
	t.Setenv("VOLTA_HOME", "")
	t.Setenv("HOMEBREW_PREFIX", "")
	var stdout, stderr bytes.Buffer
	code := invoke(args, &stdout, &stderr, "1.0.0")
	normalize := func(value string) string { return strings.ReplaceAll(value, root, "<root>") }
	return doctorResult{Code: code, Stdout: normalize(stdout.String()), Stderr: normalize(stderr.String())}
}

func compatibilityFixture(selected compatibility.Interface) compatibility.Input {
	return compatibility.Input{
		Context: compatibility.Context{
			Interface:         selected,
			Repository:        compatibility.Fact{Value: "/workspace/bench", Source: "git root"},
			Environment:       compatibility.Fact{Value: "linux/amd64", Source: "runtime"},
			ConfigurationHome: compatibility.Fact{Value: "/home/user/.codex", Source: "selected home"},
			ActiveRuntime:     compatibility.Fact{Source: "not observable"},
			LauncherVersion:   compatibility.Fact{Value: "codex-cli 1.2.3", Source: "PATH launcher"},
			PolicyProvenance:  compatibility.Fact{Value: "parsed", Source: ".codex/hooks.json"},
		},
		Configuration: compatibility.FileFact{Path: "/home/user/.codex/config.toml", Source: "selected config", State: bounds.StateParsed},
		GlobalBench:   true,
		HookAction:    "observe the declared hook in the selected interface",
		Live:          []compatibility.LiveRow{{Capability: "normal-shell", Action: "run a normal-permission command"}},
	}
}

func readCompatibilityFile(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestCompatibilityCanonicalAssetsWithoutManifest(t *testing.T) {
	kit, root := t.TempDir(), t.TempDir()
	t.Setenv("BENCH_KIT", kit)
	for _, dir := range []string{filepath.Join(kit, ".codex"), filepath.Join(root, ".bench", "dist")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(kit, ".codex", "hooks.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "dist", "bench"), []byte{0xff, 0xfe, 0}, 0o700); err != nil {
		t.Fatal(err)
	}
	foundHook, foundBinary := false, false
	for _, asset := range compatibilityAssets(root) {
		switch asset.Name {
		case ".codex/hooks.json":
			foundHook = true
			if asset.File.State != bounds.StateAbsent {
				t.Fatalf("missing hook = %#v", asset)
			}
		case ".bench/dist/bench":
			foundBinary = true
			if asset.File.State != bounds.StateParsed {
				t.Fatalf("binary asset = %#v", asset)
			}
		}
	}
	if !foundHook || !foundBinary {
		t.Fatalf("canonical asset census missing hook=%v binary=%v", foundHook, foundBinary)
	}
}

func TestCompatibilityMissingConfigurationHome(t *testing.T) {
	t.Setenv("CODEX_HOME", "")
	t.Setenv("HOME", "")
	for _, selected := range []compatibility.Interface{compatibility.CodexCLI, compatibility.CodexDesktop} {
		got := configurationHome(selected)
		if got.Value != "" || got.Source == "" {
			t.Errorf("%s missing home = %#v, want unknown with provenance", selected, got)
		}
	}
}
