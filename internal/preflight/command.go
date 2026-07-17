package preflight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

type commandFailure struct{ failure Failure }

func (e commandFailure) Error() string { return e.failure.Message }

type runner struct {
	root          string
	mode          Mode
	binaryVersion string
	stderr        io.Writer
	identity      Identity
}

func Command(args []string, binaryVersion string, stderr io.Writer) int {
	mode, focused, profile, usage := parseArgs(args)
	if usage != nil {
		emitFailure(stderr, *usage)
		return 2
	}
	if _, err := exec.LookPath("git"); err != nil {
		emitFailure(stderr, Failure{Kind: "tool", Message: "required tool is missing or not executable: git"})
		return 1
	}
	root, err := repoRoot()
	if err != nil {
		emitFailure(stderr, Failure{Kind: "input", Message: "run release-preflight inside a Git repository"})
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := &runner{root: root, mode: mode, binaryVersion: binaryVersion, stderr: stderr}
	if err := r.populateBaseIdentity(); err != nil {
		failure := failureFrom(err, "identity")
		emitFailure(stderr, failure)
		return 1
	}
	results := r.run(ctx, focused)
	status := terminalStatus(results)
	scope := scopeFor(focused)
	run := RunEvidence{Mode: mode, Scope: scope, Identity: r.identity, Profile: profile, Phases: results}
	finalizeErr := FinalizeEvidence(ctx, root, run)
	if err := finalizeErr; err != nil {
		var intentErr *releaseIntentError
		if errors.As(err, &intentErr) {
			emitFailure(stderr, Failure{Kind: "requirement", Message: intentErr.Error()})
			return 1
		}
		emitFailure(stderr, Failure{Kind: "evidence", Message: "could not promote complete preflight evidence: " + err.Error()})
		for _, result := range results {
			if result.Failure != nil {
				emitFailure(stderr, *result.Failure)
			}
		}
		return 1
	}
	for _, result := range results {
		fmt.Fprintf(stderr, "release-preflight: %s %s\n", result.Name, result.Status)
	}
	if status != StatusGreen {
		for _, result := range results {
			if result.Failure != nil {
				emitFailure(stderr, *result.Failure)
				break
			}
		}
		return 1
	}
	return 0
}

func scopeFor(focused string) Scope {
	if focused != "" {
		return ScopeFocused
	}
	return ScopePreflight
}

func intPointer(value int) *int { return &value }

func parseArgs(args []string) (Mode, string, Profile, *Failure) {
	mode := Mode("")
	focused := ""
	profile := Profile("")
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--mode":
			i++
			if i >= len(args) {
				return "", "", "", usageFailure()
			}
			mode = Mode(args[i])
		case "--phase":
			i++
			if i >= len(args) {
				return "", "", "", usageFailure()
			}
			focused = args[i]
		case "--profile":
			i++
			if i >= len(args) {
				return "", "", "", usageFailure()
			}
			profile = Profile(args[i])
		default:
			return "", "", "", usageFailure()
		}
	}
	if mode != ModeVerify && mode != ModePublish {
		return "", "", "", usageFailure()
	}
	if focused != "" && !contains(PhaseNames(mode), focused) {
		return "", "", "", usageFailure()
	}
	if profile != "" && profile != ProfilePublic && profile != ProfileBank {
		return "", "", "", usageFailure()
	}
	if mode == ModePublish && profile == "" {
		return "", "", "", usageFailure()
	}
	return mode, focused, profile, nil
}
func usageFailure() *Failure {
	return &Failure{Kind: "usage", Message: "usage: bench release-preflight --mode verify|publish [--profile public|bank] [--phase <name>]"}
}
func emitFailure(w io.Writer, failure Failure) {
	data, _ := json.Marshal(failure)
	fmt.Fprintln(w, string(data))
}

func repoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (r *runner) run(ctx context.Context, focused string) []Result {
	names := PhaseNames(r.mode)
	if focused != "" {
		names = []string{focused}
	}
	results := make([]Result, 0, len(names))
	statuses := make(map[string]Status, len(names))
	for _, name := range names {
		definition, _ := phaseDefinition(name)
		blocked := false
		for _, required := range definition.Requires {
			if statuses[required] != StatusGreen {
				blocked = true
			}
		}
		if focused == "" && blocked {
			results = append(results, Result{Name: name, Status: StatusNotRun})
			statuses[name] = StatusNotRun
			continue
		}
		result := r.runPhase(ctx, name)
		results = append(results, result)
		statuses[name] = result.Status
		if result.Status == StatusInterrupted {
			for _, rest := range names[len(results):] {
				results = append(results, Result{Name: rest, Status: StatusNotRun})
			}
			break
		}
	}
	return results
}

func (r *runner) runPhase(ctx context.Context, name string) Result {
	var err error
	exitCode := 0
	definition, ok := phaseDefinition(name)
	if !ok {
		err = commandFailure{Failure{Kind: "phase", Message: "unknown preflight phase " + name}}
	}
	switch definition.Handler {
	case "identity":
		err = r.checkIdentity(ctx)
	case "ancestry":
		err = r.checkAncestry(ctx)
	case "changelog":
		err = r.checkChangelog()
	case "scanner":
		err = r.runVulnerability(ctx)
	case "external":
		exitCode, err = r.runExternal(ctx, name)
	case "":
	default:
		err = commandFailure{Failure{Kind: "phase", Message: "unknown preflight handler for " + name}}
	}
	if err == nil {
		return Result{Name: name, Status: StatusGreen, ExitCode: &exitCode}
	}
	if errors.Is(err, context.Canceled) {
		failure := Failure{Kind: "interrupted", Message: "preflight interrupted; child process group cancelled"}
		return Result{Name: name, Status: StatusInterrupted, Failure: &failure}
	}
	failure := failureFrom(err, name)
	if exitCode == 0 {
		exitCode = 1
	}
	return Result{Name: name, Status: StatusRed, ExitCode: &exitCode, Failure: &failure}
}

func (r *runner) runExternal(ctx context.Context, name string) (int, error) {
	if err := r.validatePhaseInputs(ctx, name); err != nil {
		return 1, err
	}
	definition, ok := phaseDefinition(name)
	if !ok {
		return 1, commandFailure{Failure{Kind: "phase", Message: "unknown preflight phase " + name}}
	}
	argv := phaseArgv(r.root, definition)
	if override := os.Getenv("BENCH_PREFLIGHT_" + strings.ToUpper(name)); override != "" {
		argv = []string{override}
	}
	if len(argv) == 0 {
		return 1, commandFailure{Failure{Kind: "phase", Message: "unknown preflight phase " + name}}
	}
	if _, err := exec.LookPath(argv[0]); err != nil {
		return 1, commandFailure{Failure{Kind: "tool", Message: "required tool is missing or not executable: " + argv[0]}}
	}
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = r.root
	cmd.Stdout = r.stderr
	cmd.Stderr = r.stderr
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	if err := cmd.Start(); err != nil {
		return 1, commandFailure{Failure{Kind: "tool", Message: "could not start required phase tool: " + argv[0]}}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			return 0, nil
		}
		if ee := new(exec.ExitError); errors.As(err, &ee) {
			return ee.ExitCode(), commandFailure{Failure{Kind: "phase", Message: name + " phase failed"}}
		}
		return 1, err
	case <-ctx.Done():
		if cmd.Process != nil {
			if runtime.GOOS != "windows" {
				_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			} else {
				_ = cmd.Process.Kill()
			}
		}
		<-done
		return 0, ctx.Err()
	}
}

func (r *runner) populateBaseIdentity() error {
	commit, err := gitOutput(r.root, "rev-parse", "HEAD")
	if err != nil {
		return commandFailure{Failure{Kind: "input", Message: "could not resolve source HEAD"}}
	}
	toolchain, err := readToolchain(r.root)
	if err != nil {
		return commandFailure{Failure{Kind: "input", Message: err.Error()}}
	}
	r.identity.SourceCommit = &commit
	packageVersion, err := readPackageVersion(r.root)
	if err != nil {
		return commandFailure{failure: Failure{Kind: "input", Message: "package.json version is unreadable"}}
	}
	r.identity.PackageVersion = &packageVersion
	r.identity.Toolchain = &toolchain
	return nil
}

func (r *runner) validatePhaseInputs(ctx context.Context, name string) error {
	requireTool := func(tool string) error {
		if _, err := exec.LookPath(tool); err != nil {
			return commandFailure{Failure{Kind: "tool", Message: "required tool is missing or not executable: " + tool}}
		}
		return nil
	}
	definition, ok := phaseDefinition(name)
	if !ok {
		return commandFailure{Failure{Kind: "phase", Message: "unknown preflight phase " + name}}
	}
	for _, tool := range definition.Tools {
		if err := requireTool(tool); err != nil {
			return err
		}
	}
	if definition.ExactToolchain {
		actual, err := exec.CommandContext(ctx, "go", "env", "GOVERSION").Output()
		if err != nil || r.identity.Toolchain == nil || strings.TrimSpace(string(actual)) != "go"+*r.identity.Toolchain {
			return commandFailure{Failure{Kind: "tool", Message: "actual Go version must equal go.mod toolchain patch"}}
		}
	}
	for _, rel := range definition.Inputs {
		info, err := os.Lstat(filepath.Join(r.root, filepath.FromSlash(rel)))
		if err != nil || !info.Mode().IsRegular() {
			return commandFailure{Failure{Kind: "input", Message: "required repository input is missing or not a regular file: " + rel}}
		}
	}
	return nil
}

func phaseArgv(root string, definition PhaseDefinition) []string {
	argv := append([]string{}, definition.Argv...)
	for i := range argv {
		argv[i] = strings.ReplaceAll(argv[i], "{root}", root)
	}
	return argv
}

func failureFrom(err error, phase string) Failure {
	var cf commandFailure
	if errors.As(err, &cf) {
		return cf.failure
	}
	return Failure{Kind: "phase", Message: phase + " phase failed: " + err.Error()}
}
