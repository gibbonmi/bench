package preflight

import (
	"bytes"
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
	"time"
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
	mode, focused, usage := parseArgs(args)
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
		emitFailure(stderr, failureFrom(err, "identity"))
		return 1
	}
	results := r.run(ctx, focused)
	status := terminalStatus(results)
	scope := ScopePreflight
	if focused != "" {
		scope = ScopeFocused
	}
	manifest := Manifest{SchemaVersion: 1, Mode: mode, Scope: scope, Status: status, Identity: r.identity, Phases: phaseSummaries(results)}
	if err := PromoteEvidence(root, mode, results, manifest); err != nil {
		emitFailure(stderr, Failure{Kind: "evidence", Message: "could not promote complete preflight evidence: " + err.Error()})
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

func parseArgs(args []string) (Mode, string, *Failure) {
	mode := Mode("")
	focused := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--mode":
			i++
			if i >= len(args) {
				return "", "", usageFailure()
			}
			mode = Mode(args[i])
		case "--phase":
			i++
			if i >= len(args) {
				return "", "", usageFailure()
			}
			focused = args[i]
		default:
			return "", "", usageFailure()
		}
	}
	if mode != ModeVerify && mode != ModePublish {
		return "", "", usageFailure()
	}
	if focused != "" && !contains(PhaseNames(mode), focused) {
		return "", "", usageFailure()
	}
	return mode, focused, nil
}
func usageFailure() *Failure {
	return &Failure{Kind: "usage", Message: "usage: bench release-preflight --mode verify|publish [--phase <name>]"}
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
	artifactGreen := true
	for _, name := range names {
		if name == "smoke" && focused == "" && !artifactGreen {
			results = append(results, Result{Name: name, Status: StatusNotRun})
			continue
		}
		result := r.runPhase(ctx, name)
		results = append(results, result)
		if name == "artifacts" && result.Status != StatusGreen {
			artifactGreen = false
		}
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
	switch name {
	case "identity":
		err = r.checkIdentity(ctx)
	case "ancestry":
		err = r.checkAncestry(ctx)
	case "changelog":
		err = r.checkChangelog()
	case "vulnerability":
		err = r.runVulnerability(ctx)
	default:
		exitCode, err = r.runExternal(ctx, name)
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
	argv := phaseArgv(r.root, name)
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
	if name == "race" || name == "vet" {
		if err := requireTool("go"); err != nil {
			return err
		}
		actual, err := exec.CommandContext(ctx, "go", "env", "GOVERSION").Output()
		if err != nil || r.identity.Toolchain == nil || strings.TrimSpace(string(actual)) != "go"+*r.identity.Toolchain {
			return commandFailure{Failure{Kind: "tool", Message: "actual Go version must equal go.mod toolchain patch"}}
		}
	}
	if name == "artifacts" || name == "smoke" {
		for _, tool := range []string{"bash", "node", "npm"} {
			if err := requireTool(tool); err != nil {
				return err
			}
		}
	}
	required := map[string][]string{
		"gate":      {"bin/bench.sh", ".bench/gate.sh"},
		"artifacts": {"scripts/build-artifacts.sh", "scripts/platforms.json", "scripts/wrapper-assets.json", "package.json"},
		"smoke":     {"scripts/smoke-artifacts.sh", "scripts/platforms.json", "package.json"},
	}
	for _, rel := range required[name] {
		info, err := os.Lstat(filepath.Join(r.root, filepath.FromSlash(rel)))
		if err != nil || !info.Mode().IsRegular() {
			return commandFailure{Failure{Kind: "input", Message: "required repository input is missing or not a regular file: " + rel}}
		}
	}
	return nil
}

func phaseArgv(root, name string) []string {
	switch name {
	case "gate":
		return []string{filepath.Join(root, "bin", "bench.sh"), "gate"}
	case "race":
		return []string{"go", "test", "-race", "-count=1", "./..."}
	case "vet":
		return []string{"go", "vet", "./..."}
	case "artifacts":
		return []string{"bash", filepath.Join(root, "scripts", "build-artifacts.sh"), root, filepath.Join(root, "dist", "artifacts")}
	case "smoke":
		return []string{"bash", filepath.Join(root, "scripts", "smoke-artifacts.sh"), filepath.Join(root, "dist", "artifacts")}
	}
	return nil
}

func failureFrom(err error, phase string) Failure {
	var cf commandFailure
	if errors.As(err, &cf) {
		return cf.failure
	}
	return Failure{Kind: "phase", Message: phase + " phase failed: " + err.Error()}
}
func terminalStatus(results []Result) Status {
	status := StatusGreen
	for _, r := range results {
		if r.Status == StatusInterrupted {
			return StatusInterrupted
		}
		if r.Status != StatusGreen {
			status = StatusRed
		}
	}
	return status
}
func contains(items []string, want string) bool {
	for _, v := range items {
		if v == want {
			return true
		}
	}
	return false
}

func (r *runner) runVulnerability(ctx context.Context) error {
	tool := "govulncheck"
	if override := os.Getenv("BENCH_PREFLIGHT_VULNERABILITY"); override != "" {
		tool = override
	}
	if _, err := exec.LookPath(tool); err != nil {
		return commandFailure{Failure{Kind: "tool", Message: "required tool is missing or not executable: " + tool}}
	}
	cmd := exec.CommandContext(ctx, tool, "-json", "./...")
	cmd.Dir = r.root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.MultiWriter(r.stderr, &stderr)
	runErr := cmd.Run()
	ids, err := findingIDs(stdout.Bytes())
	if err != nil {
		return commandFailure{Failure{Kind: "phase", Message: err.Error()}}
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if !errors.As(runErr, &exitErr) || exitErr.ExitCode() != 3 {
			return commandFailure{Failure{Kind: "phase", Message: "govulncheck scanner failed"}}
		}
	}
	path := filepath.Join(r.root, "scripts", "vuln-exceptions.json")
	data, err := readRegular(path)
	present := err == nil
	if err != nil && !os.IsNotExist(err) {
		return commandFailure{Failure{Kind: "input", Message: "vulnerability exception policy is unreadable"}}
	}
	today := os.Getenv("BENCH_PREFLIGHT_DATE")
	if today == "" {
		today = time.Now().UTC().Format("2006-01-02")
	}
	if err := ValidateVulnerabilityPolicy(ids, data, present, today); err != nil {
		return commandFailure{Failure{Kind: "phase", Message: err.Error()}}
	}
	return nil
}
