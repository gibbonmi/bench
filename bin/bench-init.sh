#!/usr/bin/env bash
# `bench init` — scaffold a repo's .bench/ oracle. The setup half of adoption
# (bench-link.sh installs the shared kit; this scaffolds the project-owned gate,
# its canary harness, and the learnings journal). Sourced by bench.sh; relies on
# repo_root() and kit_dir() defined there.

init() {
  local root kit; root="$(repo_root)"; mkdir -p "$root/.bench"
  kit="$(kit_dir)"
  # The canary runner is shared, not scaffolded inline (one source for benchkit's gate
  # and every consumer's). Install it so the scaffolded gate below can source it even
  # when `bench init` is run without a prior `bench link`.
  if [[ ! -e "$root/.bench/lib/canary-run.sh" ]]; then
    mkdir -p "$root/.bench/lib"
    cp "$kit/.bench/lib/canary-run.sh" "$root/.bench/lib/canary-run.sh"
  fi
  if [[ ! -e "$root/.bench/gate.sh" ]]; then
    cat > "$root/.bench/gate.sh" <<'EOF'
#!/usr/bin/env bash
# The external oracle for this repo — correctness only. Exit 0 = done is allowed.
# No `set -e`: the canary runner reads $? after a subshell that is meant to fail.
set -uo pipefail
root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "gate: not in a git repo" >&2; exit 3; }
cd "$root"
gate_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fail=0
err() { echo "gate: $*" >&2; fail=1; }

# Stack checks — replace the sentinel below with what fits your project, e.g.:
#   mypy . && pytest -q && ruff check .        || err "python checks failed"
#   pnpm -s typecheck && pnpm -s test && pnpm -s lint  || err "js checks failed"

# Sentinel — keeps a fresh gate RED until you configure it, so you cannot commit real
# work against an empty gate. Delete this one line once a real check exists above.
err "configure .bench/gate.sh — replace this sentinel with real checks"  # BENCH_SENTINEL

# Example check + its seed canary (tests/canary/example) are the pattern to copy: run a
# check, err on failure, and add a canary that proves it bites. This one fails if a
# forbidden marker file is present; the seed fixture plants that file so the canary can
# prove the check still bites. Replace it with your real checks — and keep a canary for
# each, because the runner turns the gate red if tests/canary/ ever goes empty.
[ -e DO-NOT-SHIP ] && err "example check: DO-NOT-SHIP marker file present"

# Structural debt is NOT checked here. `bench shift` runs `bench structure` after the
# loop and refactors only when a file or dir is over budget. Uncomment to hard-block
# structure at every commit:
#   bench structure || err "structure over budget"

# Canary — prove the checks above still bite, and that the harness itself is present.
# Guard the source: a missing runner cannot fire its own absent-check, so deleting it
# would otherwise pass silently green — deleting the runner is as red as deleting the
# fixtures. Keep .bench/lib/canary-run.sh installed (bench link and bench init both do).
if [ -f "$gate_dir/lib/canary-run.sh" ]; then
  # shellcheck source=/dev/null
  . "$gate_dir/lib/canary-run.sh"
else
  err "canary runner missing (.bench/lib/canary-run.sh) — run 'bench link' or 'bench init' to reinstall it"
fi

if [ "$fail" -eq 0 ]; then echo "gate: green"; else echo "gate: red" >&2; fi
exit "$fail"
EOF
    chmod +x "$root/.bench/gate.sh"
    echo "scaffolded .bench/gate.sh — red until you replace the sentinel with real checks"
  fi
  # Seed canary — a worked example proving the harness bites. Keep at least one fixture:
  # the runner turns the gate red when tests/canary/ is absent or empty.
  if [[ ! -e "$root/tests/canary/example" ]]; then
    mkdir -p "$root/tests/canary/example/files"
    printf 'Presence of this file at the repo root makes the seed example check fail.\n' \
      > "$root/tests/canary/example/files/DO-NOT-SHIP"
    printf 'example check: DO-NOT-SHIP marker file present\n' \
      > "$root/tests/canary/example/EXPECT"
    echo "scaffolded tests/canary/example — the seed canary; copy it for each real check"
  fi
  if [[ ! -e "$root/.bench/learnings.md" ]]; then
    cat > "$root/.bench/learnings.md" <<'EOF'
# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the
generalizable ones into the kit with sign-off, and prunes them: a resolved entry
leaves this file, and its verdict (promoted or dismissed, one line of why) is
recorded in the integration commit and CHANGELOG. The journal holds open entries
only; history lives in git. Never rewrite a kit rule yourself — that is the whole
point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-integrate-learnings.

<!-- entries below -->
EOF
    echo "scaffolded .bench/learnings.md — the self-learning journal"
  fi
  echo "see projects/<name>.md in the Bench kit for the profile template"
}
