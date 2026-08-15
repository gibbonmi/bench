package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/lines"
)

// doctorRow evaluates one independent per-harness health signal. Each row prints its
// own ok/red line and answers only for itself — never folded into one aggregate green
// — so a broken row (a preserved CLAUDE.md without its imports, a missing profile)
// stays visible even when every other signal, including the shim, is healthy. An empty
// message is a row with nothing to say about this repo and prints no line at all, which
// is how a row scoped to an optional asset stays quiet where the asset is absent.
// Adding a row is one more entry in doctorRows plus one evaluator function.
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
	{"binding migration", evalBindingMigrationRow},
	{"binding columns", evalBindingColumnsRow},
	{"worktree admin", evalWorktreeAdminRow},
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
		if message == "" {
			continue
		}
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

// evalBindingMigrationRow reports the hand rewrite that moves a retired
// BENCH_TIER_*/BENCH_ALIAS_* binding onto the BENCH_<HARNESS>_<TIER> matrix. Those keys
// bind nothing, so a file carrying only them is a repo with no line binding at all rather
// than one on an older schema — which is why the row is red and why it names every retired
// key it finds: an unnamed one is a silent survivor of the migration. Reporting is the
// whole remedy. Doctor never writes .bench/lines.env, because the binding is reviewer-owned
// and a tool that edited it would be deciding the line on the reviewer's behalf.
func evalBindingMigrationRow(root string) (bool, string) {
	content, err := os.ReadFile(bindingPath(root))
	if err != nil {
		return true, ""
	}
	rewrites := lines.RetiredKeyRewrites(content)
	if len(rewrites) == 0 {
		return true, ""
	}
	var b strings.Builder
	b.WriteString(".bench/lines.env carries retired binding keys, which bind nothing under the " +
		"BENCH_<HARNESS>_<TIER> matrix; rewrite each by hand (bench doctor never writes this file):")
	for _, r := range rewrites {
		fmt.Fprintf(&b, "\n          %s=%s  ->  %s=%s", r.Retired, r.Value, r.Replacement, r.Value)
	}
	return false, b.String()
}

// evalBindingColumnsRow names every known harness whose column binds no model, with the
// keys that would bind it. An unadopted harness is not a defect — its adapter refuses to
// launch, which is the safe posture for a family nobody chose — so the row is green; what
// it prevents is an operator meeting that refusal with no route forward.
func evalBindingColumnsRow(root string) (bool, string) {
	content, err := os.ReadFile(bindingPath(root))
	if err != nil {
		return true, ""
	}
	binding := lines.ParseBinding(content)
	var unbound []string
	for _, harness := range lines.Harnesses {
		if binding.Complete(harness) {
			continue
		}
		unbound = append(unbound, harness+" (bind "+strings.Join(binding.UnboundKeys(harness), ", ")+")")
	}
	if len(unbound) == 0 {
		return true, "every known harness column is bound in .bench/lines.env"
	}
	return true, "unbound columns in .bench/lines.env: " + strings.Join(unbound, "; ") +
		" — each refuses to launch until it is bound"
}

func evalWorktreeAdminRow(root string) (bool, string) {
	common, err := git.CommonDir(root)
	if err != nil {
		return false, err.Error()
	}
	if err := git.ScanWorktreeAdmin(common); err != nil {
		return false, err.Error()
	}
	return true, ""
}

func bindingPath(root string) string {
	return filepath.Join(root, ".bench", "lines.env")
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
