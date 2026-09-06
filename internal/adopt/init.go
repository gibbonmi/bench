package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/toon"
)

func Init(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintln(stderr, "usage: bench init")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	kit := gate.KitDir()
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	runnerSrc := filepath.Join(kit, ".bench", "lib", "canary-run.sh")
	runnerDest := filepath.Join(root, ".bench", "lib", "canary-run.sh")
	if _, err := os.Lstat(runnerDest); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(runnerDest), 0o755); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := copyFile(runnerSrc, runnerDest); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	gatePath := filepath.Join(root, ".bench", "gate.sh")
	if _, err := os.Lstat(gatePath); os.IsNotExist(err) {
		if err := os.WriteFile(gatePath, []byte(scaffoldGate()), 0o755); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "scaffolded .bench/gate.sh - red until you replace the sentinel with real checks")
	}
	journal := filepath.Join(root, filepath.FromSlash(learnings.JournalPath))
	if _, err := os.Lstat(journal); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(journal), 0o755); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := os.WriteFile(journal, []byte(scaffoldLearnings()), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "scaffolded capture/learnings.md - the self-learning journal")
	}
	if added, err := ignoreLocalCapture(root); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	} else if added > 0 {
		fmt.Fprintf(stdout, "added %d local capture entries to .gitignore - these files stay out of history\n", added)
	}
	fmt.Fprintln(stdout, "see projects/<name>.md in the Bench kit for the profile template")
	return 0
}

// localCaptureIgnoreLines are the project-local capture files a linked repo must not track.
// Their content enters history only through a reviewed drain, so bench init ignores them.
// landing.LocalCapturePath is the matching predicate; the adopt tests hold the two lists in
// agreement, so a new capture file cannot be ignored here and stay unrecognized there.
var localCaptureIgnoreLines = []string{
	handoffdoc.DocumentPath,
	roadmap.IdeasFile,
	learnings.JournalPath,
}

// ignoreLocalCapture appends the missing localCaptureIgnoreLines entries to the root
// .gitignore and reports how many it added. It creates the file when absent, and it
// terminates a final line that carries no newline before it appends. A line already present
// is left alone, so a re-run of bench init leaves the file byte-identical.
func ignoreLocalCapture(root string) (int, error) {
	path := filepath.Join(root, ".gitignore")
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	present := map[string]bool{}
	for _, line := range strings.Split(string(existing), "\n") {
		present[strings.TrimSpace(line)] = true
	}
	var appended strings.Builder
	for _, line := range localCaptureIgnoreLines {
		if !present[line] {
			appended.WriteString(line)
			appended.WriteString("\n")
			present[line] = true
		}
	}
	if appended.Len() == 0 {
		return 0, nil
	}
	updated := existing
	if len(updated) > 0 && updated[len(updated)-1] != '\n' {
		updated = append(updated, '\n')
	}
	updated = append(updated, appended.String()...)
	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return 0, err
	}
	return strings.Count(appended.String(), "\n"), nil
}

// SentinelMarker is the fail-closed-stub trailer scaffoldGate embeds as a shell comment.
// Doctor's gate row, setup's zero-signal message, and out-of-package readers that retire
// the sentinel line all key off this same literal. This keeps the gate script, the doctor
// row, and the printed remedy from drifting apart.
const SentinelMarker = "BENCH_SENTINEL"

// gateScriptPreamble is the one shebang/set/git-root-guard preamble every generated
// gate.sh carries. scaffoldGate's fail-closed stub and setup's detected-ecosystem script
// both compose it here rather than re-authoring it. This keeps the opening lines of a
// written gate.sh from drifting between the two writers; see setupGateScript in setup.go
// for the second caller. comment is the writer-identifying line(s) that follow the
// shebang, each already newline-terminated.
func gateScriptPreamble(comment string) string {
	return `#!/usr/bin/env bash
` + comment + `set -uo pipefail
root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "error: not in a git repository - run inside a Bench-linked repo" >&2; exit 3; }
cd "$root"
`
}

func scaffoldGate() string {
	return gateScriptPreamble(`# This is the external oracle for this repo, and it checks correctness only. Exit 0
# means the work is done.
# This gate skips `+"`set -e`"+`: the canary command may fail, and the gate still keeps
# collecting errors.
`) + `fail=0
err() { echo "gate: $*" >&2; fail=1; }

# Add stack checks that fit your project. Replace the sentinel below, for example:
#   mypy . && pytest -q && ruff check .        || err "python checks failed"
#   pnpm -s typecheck && pnpm -s test && pnpm -s lint  || err "js checks failed"

# This sentinel keeps a fresh gate red until you configure it, so you cannot commit
# real work against an empty gate. Delete this line once a real check exists above.
err "configure .bench/gate.sh - replace this sentinel with real checks"  # ` + SentinelMarker + `

# This gate does not check structural debt. ` + "`bench shift`" + ` runs ` + "`bench structure`" + ` after the
# loop, and it refactors only when a file or directory is over budget. Uncomment this
# line to hard-block structure at every commit:
#   bench structure || err "structure over budget"

# This validates the configured fixture inventory, but only once the repo has one. A
# project with no tests/canary directory skips validation entirely. A directory that
# exists, even empty, falls through to bench canary and gets validated. Resolve the
# repo-local CLI before a global bench on PATH, so a machine with no global bench
# still reaches the inventory command.
bench="$(dirname "$0")/bin/bench.sh"; [ -x "$bench" ] || bench=bench
if [ -d "$root/tests/canary" ]; then
  "$bench" canary "$root" || err "canary inventory validation failed"
fi

if [ "$fail" -eq 0 ]; then echo "gate: green"; else echo "gate: red" >&2; fi
exit "$fail"
`
}

func scaffoldLearnings() string {
	return learnings.JournalSchemaHeading + `

Append one entry when you deviate from the workflow, or when you make a
judgment call you're unsure about. Append one too when you catch a
should-have-asked in hindsight, or when you catch yourself assembling the
same ad-hoc check twice. That is a codification candidate; name the ` + "`bench`" + `
subcommand it wants to be. You capture; the reviewer decides.

` + "`/bench-drain`" + ` verdicts every open entry in its reviewed batch diff.
Work-shaped and rule-shaped entries become roadmap items; rule-shaped ones
build later under the synthesis discipline. The rest are dismissed with one
line of why. A resolved entry leaves this file, and its verdict is recorded
in the drain commit. The journal holds open entries only; history lives in
git. Never rewrite a kit rule yourself; that is the whole point of capturing
here instead.

Format per entry (the bench learning verb writes it; never type the heading by hand):

` + learnings.FormatEntry("<date>", "<short title>", "...", "...", "...") + `
An entry leaves this file only via /bench-drain.

` + learnings.JournalEntriesMarker + "\n"
}
