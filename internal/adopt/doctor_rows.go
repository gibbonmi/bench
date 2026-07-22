package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// doctorRow evaluates one independent per-harness health signal. Each row prints its
// own ok/red line and answers only for itself — never folded into one aggregate green
// — so a broken row (a preserved CLAUDE.md without its imports, a missing profile)
// stays visible even when every other signal, including the shim, is healthy. Adding a
// row is one more entry in doctorRows plus one evaluator function.
type doctorRow struct {
	label string
	eval  func(root string) (ok bool, message string)
}

var doctorRows = []doctorRow{
	{"AGENTS.md", evalAgentsRow},
	{"CLAUDE.md", evalClaudeRow},
	{"gate", evalGateRow},
	{"profile", evalProfileRow},
	{"repo-local bench", evalRepoLocalBenchRow},
	{"setup pointers", evalSetupPointersRow},
}

// reportDoctorRows renders every per-harness row when doctor runs inside a git
// worktree that bench has touched (a `.bench` directory exists — `bench init` or
// `bench link` has run), and returns whether any row is red. An untouched plain git
// repo has nothing bench-adopted to check yet, so the rows stay silent rather than
// reporting every asset absent — the shim/missing row already covers that case.
func reportDoctorRows(stdout io.Writer) bool {
	root, err := git.Root()
	if err != nil {
		return false // not in a git worktree - nothing to check
	}
	if info, err := os.Stat(filepath.Join(root, ".bench")); err != nil || !info.IsDir() {
		return false // bench has not touched this repo yet
	}
	red := false
	for _, row := range doctorRows {
		ok, message := row.eval(root)
		if ok {
			fmt.Fprintf(stdout, "  ok: %s\n", message)
			continue
		}
		fmt.Fprintf(stdout, "  red: %s\n", message)
		red = true
	}
	return red
}

func evalAgentsRow(root string) (bool, string) {
	path := filepath.Join(root, "AGENTS.md")
	// A FIFO/socket/device must never be opened for read - see isSpecialFile's doc
	// comment. Route it to a red row instead of a hang, same posture as the link
	// transaction's own AGENTS.md/CLAUDE.md special-file guard.
	if isSpecialFile(path) {
		return false, "AGENTS.md is a special file, not a regular file (run bench link)"
	}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, "AGENTS.md absent - no Bench managed block (run bench link)"
	}
	if err != nil {
		return false, fmt.Sprintf("AGENTS.md: %v", err)
	}
	content := string(b)
	if verr := validateAgentsContent(content); verr != nil {
		return false, "AGENTS.md: " + verr.Error()
	}
	scan := scanMarkers(content)
	if len(scan.starts) == 0 {
		return false, "AGENTS.md has no Bench managed block (run bench link)"
	}
	return true, "AGENTS.md carries the Bench managed block"
}

func evalClaudeRow(root string) (bool, string) {
	path := filepath.Join(root, "CLAUDE.md")
	if isSpecialFile(path) {
		return false, "CLAUDE.md is a special file, not a regular file (run bench link)"
	}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, "CLAUDE.md absent (run bench link)"
	}
	if err != nil {
		return false, fmt.Sprintf("CLAUDE.md: %v", err)
	}
	if !claudeHasImports(string(b)) {
		return false, "CLAUDE.md is present but lacks the Bench import lines (run bench link)"
	}
	return true, "CLAUDE.md carries the Bench import lines"
}

func evalGateRow(root string) (bool, string) {
	path := filepath.Join(root, ".bench", "gate.sh")
	ok, msg := evalExecutableRow(path, ".bench/gate.sh", "gate present and executable at .bench/gate.sh", "run bench init")
	if !ok {
		return ok, msg
	}
	// The fail-closed scaffold (bench init's scaffoldGate, and bench setup's
	// zero-signal stub - the same content, one source) carries this sentinel until a
	// real check replaces it. A fabricated green oracle must never hide behind a
	// green doctor row either: the row stays red until the sentinel is gone, not just
	// until the gate happens to be present and executable.
	if content, err := os.ReadFile(path); err == nil && strings.Contains(string(content), benchSentinelMarker) {
		return false, ".bench/gate.sh is still the unconfigured fail-closed stub (replace the " + benchSentinelMarker + " sentinel with real checks)"
	}
	return ok, msg
}

// evalSetupPointersRow validates the pointer setup's own next-action print relies on
// (FT76's deferred row-11 sub-row): the /bench-setup-repo command file finishSetup
// tells a converged user to go continue. A missing pointer would make that print a
// dead end, so doctor catches it rather than a cold user discovering it by hand.
func evalSetupPointersRow(root string) (bool, string) {
	path := filepath.Join(root, ".agents", "commands", "bench-setup-repo.md")
	if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
		return true, "setup-emitted pointer present at .agents/commands/bench-setup-repo.md"
	}
	return false, "setup-emitted pointer missing: .agents/commands/bench-setup-repo.md (run bench link)"
}

func evalProfileRow(root string) (bool, string) {
	entries, err := os.ReadDir(filepath.Join(root, "projects"))
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				return true, "profile present at projects/" + e.Name()
			}
		}
	}
	return false, "no projects/<name>.md profile (see projects/<name>.md in the kit for the template)"
}

func evalRepoLocalBenchRow(root string) (bool, string) {
	path := filepath.Join(root, ".bench", "bin", "bench.sh")
	return evalExecutableRow(path, ".bench/bin/bench.sh", "repo-local bench resolvable at .bench/bin/bench.sh", "run bench link")
}

// evalExecutableRow is the shared present-and-executable check behind both the gate
// row and the repo-local bench row — one evaluator shape, two callers.
func evalExecutableRow(path, label, okMessage, remedy string) (bool, string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, fmt.Sprintf("%s absent (%s)", label, remedy)
	}
	if err != nil {
		return false, fmt.Sprintf("%s: %v", label, err)
	}
	if info.IsDir() || info.Mode().Perm()&0o111 == 0 {
		return false, fmt.Sprintf("%s is present but not executable", label)
	}
	return true, okMessage
}
