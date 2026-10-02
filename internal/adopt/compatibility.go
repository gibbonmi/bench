package adopt

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"

	"github.com/gibbonmi/bench/internal/compatibility"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/usage"
)

const compatibilityDoctorUsage = "usage: bench doctor --compat <codex-cli|codex-desktop>"

var (
	doctorGrammar = usage.Grammar{
		Cmd:   "bench doctor",
		Help:  "usage: bench doctor [--fix]",
		Flags: []usage.Flag{{Name: "--fix"}},
	}
	compatibilityGrammar = usage.Grammar{
		Cmd:     "bench doctor --compat",
		Help:    compatibilityDoctorUsage,
		MinArgs: 1,
		MaxArgs: 1,
	}
)

func Doctor(args []string, stdout, stderr io.Writer, version string) int {
	if len(args) > 0 && args[0] == "--compat" {
		return compatibilityDoctor(args[1:], stdout, stderr, collectCompatibility)
	}
	return legacyDoctor(args, stdout, stderr, version)
}

func legacyDoctor(args []string, stdout, stderr io.Writer, version string) int {
	parsed, line, code := usage.Parse(doctorGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
		} else {
			fmt.Fprintln(stderr, line)
		}
		return code
	}
	if _, fix := parsed.Flags["--fix"]; fix {
		return doctorFix(stdout, stderr, version)
	}
	return doctorReport(stdout, version)
}

type compatibilityCollector func(compatibility.Interface) compatibility.Input

func compatibilityDoctor(args []string, stdout, stderr io.Writer, collect compatibilityCollector) int {
	parsed, line, code := usage.Parse(compatibilityGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
		} else {
			fmt.Fprintln(stderr, line)
		}
		return code
	}
	selected, ok := compatibility.ParseInterface(parsed.Positionals[0])
	if !ok {
		fmt.Fprintln(stderr, compatibilityDoctorUsage)
		return 2
	}
	report := compatibility.Inspect(collect(selected))
	output, err := report.Render()
	if err != nil {
		fmt.Fprintln(stderr, "compatibility output refused")
		return 1
	}
	fmt.Fprint(stdout, output)
	return report.ExitCode()
}

func collectCompatibility(selected compatibility.Interface) compatibility.Input {
	root, rootErr := git.Root()
	context := compatibility.Context{
		Interface:         selected,
		Repository:        compatibility.Fact{Value: root, Source: "git root"},
		Environment:       compatibility.Fact{Value: runtime.GOOS + "/" + runtime.GOARCH, Source: "runtime"},
		ConfigurationHome: configurationHome(selected),
		ActiveRuntime:     compatibility.Fact{Source: "not observable from doctor subprocess"},
		LauncherVersion:   compatibility.Fact{Source: "not observed"},
	}
	configPath := ""
	if context.ConfigurationHome.Value != "" {
		configPath = filepath.Join(context.ConfigurationHome.Value, "config.toml")
	}
	input := compatibility.Input{
		Context:       context,
		Configuration: compatibility.ReadFile(configPath, "selected interface configuration"),
		GlobalBench:   globalBenchAvailable(),
		HookAction:    "observe the declared hook in the selected interface",
		Live: []compatibility.LiveRow{
			{Capability: "normal-shell", Action: "run a normal-permission shell command in the selected interface"},
			{Capability: "repository-wrapper", Action: "invoke the repository Bench wrapper in the selected interface"},
		},
	}
	if rootErr != nil {
		input.Context.Repository = compatibility.Fact{Source: "git root unavailable"}
		input.ByPathAction = "run the repository .bench/bin/bench.sh by path"
		return input
	}
	hook := compatibility.ReadFile(filepath.Join(root, ".codex", "hooks.json"), "repository hook declaration")
	input.HookDeclared = hook.State == "parsed"
	input.Context.PolicyProvenance = compatibility.Fact{Value: string(hook.State), Source: hook.Path}
	wrapper, _ := linkDestination("bin/bench.sh")
	if gate.KitSourceCheckout(root) {
		wrapper = "bin/bench.sh"
	}
	input.ByPathAction = fmt.Sprintf("run %s doctor --compat %s", sanitize.ShellQuote(filepath.Join(root, filepath.FromSlash(wrapper))), selected)
	input.Assets = compatibilityAssets(root)
	return input
}

func configurationHome(selected compatibility.Interface) compatibility.Fact {
	if home := os.Getenv("CODEX_HOME"); home != "" {
		return compatibility.Fact{Value: home, Source: "CODEX_HOME"}
	}
	if selected == compatibility.CodexDesktop {
		return compatibility.Fact{Source: "CODEX_HOME unavailable for selected desktop interface"}
	}
	home, _ := os.UserHomeDir()
	return compatibility.Fact{Value: filepath.Join(home, ".codex"), Source: "HOME default"}
}

func compatibilityAssets(root string) []compatibility.Asset {
	kit := gate.KitDir()
	source := gate.KitSourceCheckout(root)
	restore := "bench link"
	if source {
		restore = "bench doctor --fix"
	}
	plan, err := buildLinkPlan(kit)
	if err != nil {
		return []compatibility.Asset{{Name: "consumer-payload", File: compatibility.FileFact{State: bounds.StateUnreadable, Reason: err.Error()}, RestoreAction: "inspect the installed Bench payload"}}
	}
	assets := make([]compatibility.Asset, 0, len(plan))
	for _, entry := range plan {
		rel := entry.rel
		if source {
			if entry.kind != "file" {
				continue
			}
			var err error
			rel, err = filepath.Rel(kit, entry.src)
			if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				continue
			}
		}
		path := filepath.Join(root, filepath.FromSlash(rel))
		fact := compatibility.FileFact{Path: path, Source: "canonical payload", State: bounds.StateParsed}
		info, err := os.Stat(path)
		switch {
		case os.IsNotExist(err):
			fact.State = bounds.StateAbsent
		case err != nil:
			fact.State, fact.Reason = bounds.StateUnreadable, err.Error()
		case !info.Mode().IsRegular():
			fact.State, fact.Reason = bounds.StateWrongType, "required asset is not a regular file"
		}
		assets = append(assets, compatibility.Asset{Name: filepath.ToSlash(rel), File: fact, RestoreAction: restore})
	}
	return assets
}

func globalBenchAvailable() bool {
	_, err := exec.LookPath("bench")
	return err == nil
}
