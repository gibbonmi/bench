package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/releaseevidence"
)

// doctorRow evaluates one independent per-harness health signal. Each row prints its own
// ok/red line and answers only for itself, never folded into one aggregate green. A
// broken row stays visible even when every other signal, including the shim, is healthy.
// Examples are a preserved CLAUDE.md without its imports, or a missing profile. An empty
// message is a row with nothing to say about this repo, and it prints no line at all.
// This is how a row scoped to an optional asset stays quiet where the asset is absent.
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
	{"binary seal", evalBinarySealRow},
	{"promotion broker", evalBrokerManifestRow},
}

// reportDoctorRows renders every per-harness row when doctor runs inside a git worktree
// that bench has touched. A touched worktree is one where a `.bench` directory exists,
// because `bench init` or `bench link` has run. It returns whether any row is red. An
// untouched plain git repo has nothing bench-adopted to check yet, so the rows stay
// silent rather than reporting every asset absent. The shim/missing row already covers
// that case.
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
	if KitSourceCheckout(root) {
		return true, "kit source checkout - AGENTS.md is the source agreement; no managed block applies"
	}
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
	// The fail-closed scaffold carries this sentinel until a real check replaces it. Bench
	// init's scaffoldGate and bench setup's zero-signal stub are the same content, one
	// source. A fabricated green oracle must never hide behind a green doctor row. The row
	// stays red until the sentinel is gone, not just until the gate is present and
	// executable.
	if content, err := os.ReadFile(path); err == nil && strings.Contains(string(content), SentinelMarker) {
		return false, ".bench/gate.sh is still the unconfigured fail-closed stub (replace the " + SentinelMarker + " sentinel with real checks)"
	}
	return ok, msg
}

// evalSetupPointersRow validates the pointer setup's own next-action print relies on.
// finishSetup tells a converged user to continue at the /bench-setup-repo command file. A
// missing pointer would make that print a dead end. Doctor catches it rather than a cold
// user discovering it by hand.
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
	if KitSourceCheckout(root) {
		return true, "kit source checkout - the launcher is bin/bench.sh; no .bench/bin copy applies"
	}
	path := filepath.Join(root, ".bench", "bin", "bench.sh")
	return evalExecutableRow(path, ".bench/bin/bench.sh", "repo-local bench resolvable at .bench/bin/bench.sh", "run bench link")
}

// evalBindingMigrationRow reports the hand rewrite that moves a retired
// BENCH_TIER_*/BENCH_ALIAS_* binding onto the BENCH_<HARNESS>_<TIER> matrix. Those keys
// bind nothing, so a file carrying only them is a repo with no line binding at all, not
// one on an older schema. This is why the row is red and why it names every retired key
// it finds. An unnamed key is a silent survivor of the migration. Reporting is the whole
// remedy. Doctor never writes .bench/lines.env, because the binding is reviewer-owned,
// and a tool that edited it would decide the line on the reviewer's behalf.
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
// keys that would bind it. An unadopted harness is not a defect. Its adapter refuses to
// launch, the safe posture for a family nobody chose, so the row is green. What it
// prevents is an operator meeting that refusal with no route forward.
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

// evalBinarySealRow grades the published dist/bench against the sources it was built
// from. Every other row can be green beside a binary that answers from last week's code,
// because nothing in the ordinary loop executes dist/bench. The verdict and its rebuild
// sentence both come from freshness, so this row derives no digest of its own and never
// reads an mtime. The row is scoped to a repository that declares build inputs, and it
// stays plain where no binary is published, because an absent artifact is not a defect.
func evalBinarySealRow(root string) (bool, string) {
	if !freshness.DeclaresBuildInputs(root) {
		return true, ""
	}
	executable := freshness.PublishedExecutable(root)
	if _, err := os.Lstat(executable); err != nil {
		return true, "dist/bench not published in this checkout - no seal to verify"
	}
	if err := freshness.Verify(root, executable); err != nil {
		return false, err.Error()
	}
	return true, "dist/bench seal matches the current build inputs"
}

// evalBrokerManifestRow predicts the landing's exit 127 before an operator meets it. The
// wrapper's land route is shell that runs before any binary is trusted, so it cannot call
// this code; the row therefore applies the same five predicates in the same order, with
// the same wording. One conformance expectation lists those five reasons and holds both
// derivations to them. The manifest is read beside the resolved wrapper, which is where
// the install and repair owner publishes it.
//
// An absent manifest is named, not red, and it names no path. The resolved wrapper
// belongs to the environment rather than to the repository under doctor, so a checkout
// opened beside a wrapper nobody installed a broker for has nothing to repair. A red
// there would fail every such repository, and setup, on a state the repository does not
// own. Every manifest that does exist is graded, and a defect in it stays red.
func evalBrokerManifestRow(string) (bool, string) {
	const remedy = " - run bench doctor --fix to republish the promotion broker"
	bindir := filepath.Dir(resolvedWrapper())
	path := filepath.Join(bindir, brokermanifest.Name)
	if _, err := os.Stat(path); err != nil {
		return true, "no promotion-broker manifest at the resolved wrapper; a landing there would refuse" + remedy
	}
	fields, err := brokermanifest.Read(path)
	if err != nil {
		return false, fmt.Sprintf("promotion-broker manifest at %s is incomplete%s", path, remedy)
	}
	install := filepath.Dir(bindir)
	// The route compares only when it can read the installed package version, so an
	// install root without a package.json leaves this predicate silent rather than red.
	if installed, err := releaseevidence.ReadPackageVersion(install); err == nil && installed != fields["version"] {
		return false, fmt.Sprintf("promotion broker version %s does not match installed package %s%s", fields["version"], installed, remedy)
	}
	broker := fields["path"]
	if !filepath.IsAbs(broker) {
		broker = filepath.Join(install, broker)
	}
	info, err := os.Lstat(broker)
	if err != nil || !info.Mode().IsRegular() || info.Size() == 0 || info.Mode().Perm()&0o111 == 0 {
		return false, fmt.Sprintf("promotion broker at %s is not a regular executable%s", broker, remedy)
	}
	digest, err := brokermanifest.Digest(broker)
	if err != nil || digest != fields["sha256"] {
		return false, fmt.Sprintf("promotion broker at %s does not match its manifest digest%s", broker, remedy)
	}
	return true, "promotion broker authenticated at " + broker
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
