#!/usr/bin/env bash
# .bench/gate.sh — the oracle for benchkit. Exits 0 when the kit is shippable.
#
# benchkit has no test/type/lint stack; it is shell + markdown + JSON consumed by
# harnesses. "Shippable" therefore means kit conformance: scripts parse, JSON is
# valid, skills carry frontmatter, the published file list resolves, and the
# AGENTS.md index stays in sync with the skills/commands on disk.
#
# This file is NOT in package.json files[], so it never ships to consumers.
set -uo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "gate: not in a git repo" >&2; exit 3; }
cd "$root"
gate_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

fail=0
err() { echo "gate: $*" >&2; fail=1; }

# 1. Shell scripts parse.
for f in bin/*.sh .bench/gate-*.sh .bench/skills-index.sh .bench/hooks/*.sh .bench/lib/*.sh; do
  bash -n "$f" || err "bash syntax error in $f"
done

# 1b. Scripts the harness/CLI exec by path are executable in git. A fresh clone or
#     npm install gets the git index mode; if it is 100644 the hooks fail at runtime
#     with "permission denied" (exit 126) on every Stop and every Bash tool call.
for f in bin/bench.sh .bench/hooks/*.sh .bench/adapters/*; do
  mode="$(git ls-files -s "$f" | awk '{print $1}')"
  [ -z "$mode" ] && continue   # untracked — not part of what ships
  [ "$mode" = "100755" ] || err "$f is not executable in git (mode $mode); the harness runs it as a command path"
done

# 1c. Python programs in the hooks tree parse (the git-guard analyzer). Same
#     best-effort posture as shellcheck: runs only when python3 is present, so a
#     syntax error is caught at gate time, not on the next Bash tool call.
#     Bytecode goes to a scratch cfile — the gate never writes __pycache__ into
#     the tree.
if command -v python3 >/dev/null 2>&1; then
  for f in .bench/hooks/*.py; do
    [ -e "$f" ] || continue
    python3 - "$f" <<'PYCHECK' || err "py_compile failed for $f"
import os, py_compile, sys, tempfile
try:
    with tempfile.TemporaryDirectory() as d:
        py_compile.compile(sys.argv[1], cfile=os.path.join(d, "c.pyc"), doraise=True)
except (py_compile.PyCompileError, OSError) as e:
    print(e, file=sys.stderr)
    sys.exit(1)
PYCHECK
  done
fi

# 1d. The CLI must name the gate/done files that actually exist (.bench/gate.sh,
#     .bench/done.sh). An extensionless reference routes `bench gate` to auto-detect
#     instead of the repo's oracle, so it is a hard error.
bad_refs="$(grep -oE '\.bench/[A-Za-z0-9_.-]+' bin/bench.sh | grep -xE '\.bench/(gate|done)' | sort -u || true)"
[ -z "$bad_refs" ] || err "bin/bench.sh has extensionless gate/done refs ($(echo "$bad_refs" | tr '\n' ' ')); the contract is .sh"

# The shared fixture harness the contract fragments call (`contract` — one
# source for provision/report/cleanup); must precede every fragment.
# shellcheck source=/dev/null
. "$gate_dir/gate-contract-runner.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-link-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-runtime-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-runtime-shift-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-line-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-axi-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-axi-guards-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-axi-wave2-contracts.sh"

# 2. JSON is valid. A missing file errs distinctly; "invalid JSON" fires only on a
#    real parse failure of a file that exists (so the canary attributes it precisely).
for f in package.json .claude/settings.json .codex/hooks.json; do
  if [ ! -f "$f" ]; then
    err "JSON file missing: $f"
  elif ! node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"))' "$f" 2>/dev/null; then
    err "invalid JSON in $f"
  fi
done

# 2b. The Codex hook adapter must carry the two enforcement layers Codex actually
#     runs (verified against codex-cli: it loads .codex/hooks.json, fires PreToolUse
#     with tool_name "Bash" and Stop with stop_hook_active, and shells out the
#     command so the git-rev-parse substitution resolves). A drift here silently
#     drops the git guard or the stop gate under Codex.
if [ -f .codex/hooks.json ]; then
  node <<'NODE' || fail=1
const fs = require("fs");
let bad = 0;
const err = m => { console.error("gate: " + m); bad = 1; };
let cfg;
try { cfg = JSON.parse(fs.readFileSync(".codex/hooks.json", "utf8")); }
catch { process.exit(0); } // invalid JSON already reported by check 2
const hooks = (cfg && cfg.hooks) || {};
const cmds = event => (hooks[event] || [])
  .flatMap(g => (g.hooks || []).map(h => h.command || ""));
const stop = cmds("Stop");
if (!stop.some(c => c.includes(".bench/hooks/stop.sh")))
  err("codex hooks.json Stop event does not run .bench/hooks/stop.sh");
const pre = (hooks.PreToolUse || []);
if (!pre.some(g => g.matcher === "Bash"))
  err("codex hooks.json PreToolUse Bash matcher missing");
if (!cmds("PreToolUse").some(c => c.includes(".bench/hooks/block-dangerous-git.sh")))
  err("codex hooks.json PreToolUse does not run .bench/hooks/block-dangerous-git.sh");
process.exit(bad);
NODE
fi

# 3. Every skill carries YAML frontmatter (first line is the --- fence).
for f in .agents/skills/*/SKILL.md; do
  [ "$(head -1 "$f")" = "---" ] || err "$f missing frontmatter"
done

# 3b. Craft guidance is model-invoked, not a human-run Bench phase. Its installed
#     directories stay bench-craft-* for stable paths, but the visible skill names
#     must be craft-* so `$bench` menus show only human-run phase adapters.
skill_name() {
  awk -F': *' '$1 == "name" { print $2; exit }' "$1"
}
for d in .agents/skills/bench-craft-*/; do
  base="$(basename "$d")"
  expected="craft-${base#bench-craft-}"
  actual="$(skill_name "$d/SKILL.md")"
  [ "$actual" = "$expected" ] || err "craft skill '$base' visible name is '$actual'; expected '$expected'"
done
for f in .agents/skills/*/SKILL.md; do
  dir="$(basename "$(dirname "$f")")"
  actual="$(skill_name "$f")"
  [ -f ".agents/commands/$dir.md" ] && continue # command adapters intentionally expose bench-* names
  case "$actual" in
    bench-*) err "non-command skill '$dir' uses bench-* visible name '$actual'" ;;
  esac
done

# 4. Package and installable-surface contracts.
# shellcheck source=/dev/null
. "$gate_dir/gate-package-contracts.sh"

# 5. Kit conformance — the .bench/BENCH.md index (the shipped operating guide)
#    stays in sync with disk, both directions.
#    a+b) the skills index is generated from skill frontmatter; verify the
#       committed block equals the generated one (presence, wording, and order),
#       with drift attributed per skill. Replaces the old two-way name checks.
bash "$gate_dir/skills-index.sh" --check || fail=1
#    c) every command file is referenced as /name in .bench/BENCH.md
for f in .agents/commands/*.md; do
  name="$(basename "$f" .md)"
  grep -q "/$name" .bench/BENCH.md || err "command '/$name' on disk but not referenced in .bench/BENCH.md"
done
#    d) every command has an explicit Codex skill adapter. Codex does not scan
#       .agents/commands as an invocation surface, so each command phase needs a
#       thin $bench-* skill that delegates to the canonical command file.
for f in .agents/commands/*.md; do
  name="$(basename "$f" .md)"
  adapter=".agents/skills/$name/SKILL.md"
  metadata=".agents/skills/$name/agents/openai.yaml"
  [ -f "$adapter" ] || { err "command '$name' has no Codex adapter skill at $adapter"; continue; }
  grep -qE "^name:[[:space:]]*$name$" "$adapter" || err "Codex adapter '$name' frontmatter name does not match command"
  grep -qF ".agents/commands/$name.md" "$adapter" || err "Codex adapter '$name' does not reference .agents/commands/$name.md"
  [ -f "$metadata" ] || { err "Codex adapter '$name' missing agents/openai.yaml explicit-invocation metadata"; continue; }
  grep -qF "allow_implicit_invocation: false" "$metadata" || err "Codex adapter '$name' does not disable implicit invocation"
  grep -qF "\$$name" .bench/BENCH.md || err "Codex adapter '$name' is not documented in .bench/BENCH.md"
done
#    e) the roadmap promotion seam — /bench-shape-idea must name ROADMAP.md and the
#       auto-remove-on-map-creation behavior, or the only path that drains a parked idea
#       silently rots. The capture sink (bench idea) is useless without this graduation.
si=".agents/commands/bench-shape-idea.md"
grep -qF 'ROADMAP.md' "$si" || err "/bench-shape-idea does not reference ROADMAP.md (roadmap promotion seam)"
grep -qiE 'remove|delete' "$si" || err "/bench-shape-idea does not describe removing a promoted roadmap entry"
#    f) shared platform rules are single-sourced. The four invariants and the
#       communication rules are canonical in .bench/BENCH.md and referenced from
#       AGENTS.md — never copied back into AGENTS.md. Each marker must live in BENCH.md
#       and be absent from AGENTS.md (the drift this guards), and AGENTS.md must keep its
#       pointer to the canonical file.
ss_markers=(
  "you never grade your own work"
  "Declare the line before a long run"
  "Document for the teammate who just walked in"
  "One small change at a time, repo stays green"
  "Clear beats dense"
  "You are the worker; I am the reviewer"
  "Right-size the process"
  "never silently rewrite your own rules"
)
for m in "${ss_markers[@]}"; do
  grep -qF "$m" .bench/BENCH.md || err "shared rule missing from canonical .bench/BENCH.md: \"$m\""
  ! grep -qF "$m" AGENTS.md || err "shared rule duplicated in AGENTS.md (it must live only in .bench/BENCH.md): \"$m\""
done
grep -qF 'canonical in `.bench/BENCH.md`' AGENTS.md || err "AGENTS.md lost its pointer to the canonical .bench/BENCH.md shared rules"

# shellcheck source=/dev/null
. "$gate_dir/gate-docs-contracts.sh"

# 6. shellcheck — stronger shell lint, best-effort (runs only when installed).
if command -v shellcheck >/dev/null 2>&1; then
  shellcheck -S warning bin/bench.sh .bench/hooks/*.sh .bench/lib/*.sh || err "shellcheck reported issues"
fi

# 7. Canary — prove the gate's own checks still bite. For each known-broken fixture in
#    tests/canary/, run THIS gate against it in a throwaway repo and assert it goes red
#    WITH the fixture's targeted error substring. A fixture that stops biting means a
#    check rotted into an always-pass. Attribution is by substring, not isolation: a
#    minimal fixture over-fails on unrelated checks, and that is fine — we only assert
#    the targeted message is present. BENCH_CANARY_INNER marks the inner run so this
#    check skips itself and never recurses.
if [ "${BENCH_CANARY_INNER:-0}" != "1" ] && [ -d tests/canary ]; then
  # Attribution baseline: an EXPECT that also matches a completely empty fixture
  # proves nothing about its planted regression — the canary is vacuous and the
  # check it guards can rot into an always-pass unnoticed.
  d0="$(mktemp -d)"
  ( cd "$d0" && git init -q && BENCH_CANARY_INNER=1 bash "$root/.bench/gate.sh" ) >"$d0/out" 2>&1 || true
  for fx in tests/canary/*/; do
    name="$(basename "$fx")"
    [ -f "$fx/EXPECT" ] || continue
    if grep -qF "$(cat "$fx/EXPECT")" "$d0/out"; then
      err "canary '$name' EXPECT is vacuous (also matches an empty fixture)"
    fi
  done
  rm -rf "$d0"
  for fx in tests/canary/*/; do
    name="$(basename "$fx")"
    [ -f "$fx/EXPECT" ] || { err "canary fixture '$name' has no EXPECT file"; continue; }
    [ -d "$fx/files" ]  || { err "canary fixture '$name' has no files/ tree"; continue; }
    exp="$(cat "$fx/EXPECT")"
    d="$(mktemp -d)"
    cp -r "$fx/files/." "$d/"
    # Fixtures store dot-dirs as dot-<name> so the harness doesn't load a fixture's
    # .claude/skills as real skills; restore them to .<name> for the inner gate.
    for dd in "$d"/dot-*; do [ -e "$dd" ] && mv "$dd" "$d/.${dd##*/dot-}"; done
    ( cd "$d" && git init -q && BENCH_CANARY_INNER=1 bash "$root/.bench/gate.sh" ) >"$d/out" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ] || ! grep -qF "$exp" "$d/out"; then
      err "canary '$name' did not bite (want red + \"$exp\"; got exit $rc)"
    fi
    rm -rf "$d"
  done
fi

if [ "$fail" -eq 0 ]; then echo "gate: green"; else echo "gate: red" >&2; fi
exit "$fail"
