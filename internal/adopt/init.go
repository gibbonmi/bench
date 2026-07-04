package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/git"
)

func Init(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintln(stderr, "usage: bench init")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "not in a git repo")
		return 1
	}
	kit := kitDir()
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
	canaryDir := filepath.Join(root, "tests", "canary", "example")
	if _, err := os.Lstat(canaryDir); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(canaryDir, "files"), 0o755); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := os.WriteFile(filepath.Join(canaryDir, "files", "DO-NOT-SHIP"), []byte("Presence of this file at the repo root makes the seed example check fail.\n"), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := os.WriteFile(filepath.Join(canaryDir, "EXPECT"), []byte("example check: DO-NOT-SHIP marker file present\n"), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "scaffolded tests/canary/example - the seed canary; copy it for each real check")
	}
	learnings := filepath.Join(root, ".bench", "learnings.md")
	if _, err := os.Lstat(learnings); os.IsNotExist(err) {
		if err := os.WriteFile(learnings, []byte(scaffoldLearnings()), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "scaffolded .bench/learnings.md - the self-learning journal")
	}
	fmt.Fprintln(stdout, "see projects/<name>.md in the Bench kit for the profile template")
	return 0
}

func scaffoldGate() string {
	return `#!/usr/bin/env bash
# The external oracle for this repo - correctness only. Exit 0 = done is allowed.
# No ` + "`set -e`" + `: the canary runner reads $? after a subshell that is meant to fail.
set -uo pipefail
root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "gate: not in a git repo" >&2; exit 3; }
cd "$root"
gate_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fail=0
err() { echo "gate: $*" >&2; fail=1; }

# Stack checks - replace the sentinel below with what fits your project, e.g.:
#   mypy . && pytest -q && ruff check .        || err "python checks failed"
#   pnpm -s typecheck && pnpm -s test && pnpm -s lint  || err "js checks failed"

# Sentinel - keeps a fresh gate RED until you configure it, so you cannot commit real
# work against an empty gate. Delete this one line once a real check exists above.
err "configure .bench/gate.sh - replace this sentinel with real checks"  # BENCH_SENTINEL

# Example check + its seed canary (tests/canary/example) are the pattern to copy: run a
# check, err on failure, and add a canary that proves it bites. This one fails if a
# forbidden marker file is present; the seed fixture plants that file so the canary can
# prove the check still bites. Replace it with your real checks - and keep a canary for
# each, because the runner turns the gate red if tests/canary/ ever goes empty.
[ -e DO-NOT-SHIP ] && err "example check: DO-NOT-SHIP marker file present"

# Structural debt is NOT checked here. ` + "`bench shift`" + ` runs ` + "`bench structure`" + ` after the
# loop and refactors only when a file or dir is over budget. Uncomment to hard-block
# structure at every commit:
#   bench structure || err "structure over budget"

# Canary - prove the checks above still bite, and that the harness itself is present.
# Guard the source: a missing runner cannot fire its own absent-check, so deleting it
# would otherwise pass silently green - deleting the runner is as red as deleting the
# fixtures. Keep .bench/lib/canary-run.sh installed (bench link and bench init both do).
if [ -f "$gate_dir/lib/canary-run.sh" ]; then
  # shellcheck source=/dev/null
  . "$gate_dir/lib/canary-run.sh"
else
  err "canary runner missing (.bench/lib/canary-run.sh) - run 'bench link' or 'bench init' to reinstall it"
fi

if [ "$fail" -eq 0 ]; then echo "gate: green"; else echo "gate: red" >&2; fi
exit "$fail"
`
}

func scaffoldLearnings() string {
	return `# Learnings - usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate - name the ` + "`bench`" + `
subcommand it wants to be). You capture; the reviewer
decides. ` + "`/bench-integrate-learnings`" + ` reviews the open entries, promotes the
generalizable ones into the kit with sign-off, and prunes them: a resolved entry
leaves this file, and its verdict (promoted or dismissed, one line of why) is
recorded in the integration commit and CHANGELOG. The journal holds open entries
only; history lives in git. Never rewrite a kit rule yourself - that is the whole
point of capturing here instead.

Format per entry:

## <date> - <short title>  [open]
- **What happened:** ...
- **Right behavior:** ...
- **Proposed rule change:** ... (or "none")

An entry leaves this file only via /bench-integrate-learnings.

<!-- entries below -->
`
}
