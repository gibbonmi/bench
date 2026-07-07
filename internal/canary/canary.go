// Package canary runs the gate against known-broken fixture roots and proves each
// fixture still triggers its targeted diagnostic.
package canary

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
)

const absentHarnessMessage = "canary harness absent — tests/canary/ has no fixtures; the gate cannot prove its own checks bite"

// RunCall is one inner gate invocation: Cwd is the materialized repo under grade,
// Gate is the real gate script from the root being checked, and FixtureDir names
// the source fixture or is empty for the vacuity baseline.
type RunCall struct {
	Cwd        string
	Gate       string
	FixtureDir string
	Env        []string
}

// RunResult captures the inner gate verdict and combined output.
type RunResult struct {
	ExitCode int
	Output   string
}

// Runner runs one inner gate invocation.
type Runner func(RunCall) RunResult

// Run is the `bench canary [root]` command.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 1 {
		fmt.Fprintln(stderr, "usage: bench canary [root]")
		return 2
	}
	root := ""
	if len(args) == 1 && args[0] != "" {
		root = args[0]
	} else {
		var err error
		root, err = git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 1
		}
	}
	if err := Sweep(root, defaultRunner); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

// Sweep runs the canary harness for root. It returns all attributable harness
// failures it observes instead of stopping at the first one.
func Sweep(root string, runner Runner) error {
	fixtures, err := fixtures(filepath.Join(root, "tests", "canary"))
	if err != nil {
		return err
	}
	gate := filepath.Join(root, ".bench", "gate.sh")
	env := innerEnv()

	baselineDir, err := os.MkdirTemp("", "bench-canary-empty-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(baselineDir)
	_ = gitInit(baselineDir)
	baseline := runner(RunCall{Cwd: baselineDir, Gate: gate, Env: env})

	errs := runFixtures(fixtures, baseline.Output, gate, env, runner)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func runFixtures(fixtures []string, baselineOutput, gate string, env []string, runner Runner) []string {
	errs := make([]string, len(fixtures))
	jobs := make(chan int)
	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	if workers > len(fixtures) {
		workers = len(fixtures)
	}

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				errs[idx] = runFixture(fixtures[idx], baselineOutput, gate, env, runner)
			}
		}()
	}
	for idx := range fixtures {
		jobs <- idx
	}
	close(jobs)
	wg.Wait()

	out := errs[:0]
	for _, err := range errs {
		if err != "" {
			out = append(out, err)
		}
	}
	return out
}

func runFixture(fx, baselineOutput, gate string, env []string, runner Runner) string {
	name := filepath.Base(fx)
	expectPath := filepath.Join(fx, "EXPECT")
	filesDir := filepath.Join(fx, "files")

	expBytes, err := os.ReadFile(expectPath)
	if err != nil {
		return fmt.Sprintf("canary fixture '%s' has no EXPECT file", name)
	}
	if info, err := os.Stat(filesDir); err != nil || !info.IsDir() {
		return fmt.Sprintf("canary fixture '%s' has no files/ tree", name)
	}
	expect := trimExpectation(expBytes)
	if strings.Contains(baselineOutput, expect) {
		return fmt.Sprintf("canary '%s' EXPECT is vacuous (also matches an empty fixture)", name)
	}

	work, err := os.MkdirTemp("", "bench-canary-"+name+"-*")
	if err != nil {
		return fmt.Sprintf("canary '%s' setup failed: %v", name, err)
	}
	defer os.RemoveAll(work)
	if err := materialize(filesDir, work); err != nil {
		return fmt.Sprintf("canary '%s' setup failed: %v", name, err)
	}
	_ = gitInit(work)
	result := runner(RunCall{Cwd: work, Gate: gate, FixtureDir: fx, Env: env})
	if result.ExitCode == 0 || !strings.Contains(result.Output, expect) {
		return fmt.Sprintf("canary '%s' did not bite (want red + %q; got exit %d)", name, expect, result.ExitCode)
	}
	return ""
}

func fixtures(dir string) ([]string, error) {
	families, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New(absentHarnessMessage)
	}
	var out []string
	seen := map[string]string{}
	addFixture := func(name, fixtureDir string) error {
		if first := seen[name]; first != "" {
			return fmt.Errorf("canary fixture name %q appears in multiple families; base names must be globally unique", name)
		}
		seen[name] = fixtureDir
		out = append(out, fixtureDir)
		return nil
	}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		familyDir := filepath.Join(dir, family.Name())
		if _, err := os.Stat(filepath.Join(familyDir, "EXPECT")); err == nil {
			if err := addFixture(family.Name(), familyDir); err != nil {
				return nil, err
			}
			continue
		}
		entries, err := os.ReadDir(familyDir)
		if err != nil {
			return nil, err
		}
		for _, ent := range entries {
			if !ent.IsDir() {
				continue
			}
			name := ent.Name()
			fixtureDir := filepath.Join(familyDir, name)
			if err := addFixture(name, fixtureDir); err != nil {
				return nil, err
			}
		}
	}
	if len(out) == 0 {
		return nil, errors.New(absentHarnessMessage)
	}
	sort.Strings(out)
	return out, nil
}

func trimExpectation(data []byte) string {
	return strings.TrimRight(string(data), "\n")
}

// MaterializeFixture copies a canary files/ tree into dst and restores dot- path
// segments to dot directories, matching the real canary sweep.
func MaterializeFixture(src, dst string) error {
	return materialize(src, dst)
}

func materialize(src, dst string) error {
	if err := copyTree(src, dst); err != nil {
		return err
	}
	return restoreDotSegments(dst)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

func restoreDotSegments(root string) error {
	var dirs []string
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), "dot-") {
			dirs = append(dirs, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, old := range dirs {
		newPath := filepath.Join(filepath.Dir(old), "."+strings.TrimPrefix(filepath.Base(old), "dot-"))
		if err := os.Rename(old, newPath); err != nil {
			return err
		}
	}
	return nil
}

func innerEnv() []string {
	env := make([]string, 0, len(os.Environ())+1)
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") || strings.HasPrefix(kv, "BENCH_CANARY_INNER=") {
			continue
		}
		env = append(env, kv)
	}
	return append(env, "BENCH_CANARY_INNER=1")
}

func defaultRunner(call RunCall) RunResult {
	cmd := exec.Command("bash", call.Gate)
	cmd.Dir = call.Cwd
	cmd.Env = call.Env
	r := subprocess.CaptureMerged(cmd)
	output := r.Stdout
	// A spawn failure (ProcessState nil) writes nothing, so append the error.
	if r.Err != nil && cmd.ProcessState == nil {
		output += r.Err.Error()
	}
	return RunResult{ExitCode: r.ExitCode, Output: output}
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	return cmd.Run()
}
