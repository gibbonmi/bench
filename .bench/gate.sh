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
for f in bin/*.sh .bench/gate-*.sh .bench/hooks/*.sh; do
  bash -n "$f" || err "bash syntax error in $f"
done

# 1b. Scripts the harness/CLI exec by path are executable in git. A fresh clone or
#     npm install gets the git index mode; if it is 100644 the hooks fail at runtime
#     with "permission denied" (exit 126) on every Stop and every Bash tool call.
for f in bin/bench.sh .bench/hooks/*.sh; do
  mode="$(git ls-files -s "$f" | awk '{print $1}')"
  [ -z "$mode" ] && continue   # untracked — not part of what ships
  [ "$mode" = "100755" ] || err "$f is not executable in git (mode $mode); the harness runs it as a command path"
done

# 1c. The CLI must name the gate/done files that actually exist (.bench/gate.sh,
#     .bench/done.sh). An extensionless reference routes `bench gate` to auto-detect
#     instead of the repo's oracle, so it is a hard error.
bad_refs="$(grep -oE '\.bench/(gate|done)(\.sh)?' bin/bench.sh | grep -vE '\.sh$' | sort -u || true)"
[ -z "$bad_refs" ] || err "bin/bench.sh has extensionless gate/done refs ($(echo "$bad_refs" | tr '\n' ' ')); the contract is .sh"

# shellcheck source=/dev/null
. "$gate_dir/gate-link-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-runtime-contracts.sh"

# 2. JSON is valid.
for f in package.json .claude/settings.json .codex/hooks.json; do
  node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"))' "$f" \
    || err "invalid JSON in $f"
done

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

# 4. Every path in package.json files[] exists (npm-pack integrity).
node -e '
  const fs = require("fs"), p = require("./package.json");
  let bad = 0;
  for (const f of p.files) if (!fs.existsSync(f)) { console.error("gate: package.json files[] missing " + f); bad = 1; }
  process.exit(bad);
' || fail=1

pack_json="$(npm_config_cache="${TMPDIR:-/tmp}/bench-npm-cache" npm pack --dry-run --json 2>/dev/null)" || {
  err "npm pack --dry-run failed"
  pack_json="[]"
}
printf '%s' "$pack_json" | node -e '
  const fs = require("fs");
  const packs = JSON.parse(fs.readFileSync(0, "utf8"));
  const files = new Set((packs[0]?.files ?? []).map(f => f.path));
  let bad = 0;
	  for (const required of [
	    "bin/bench-link.sh",
	    "bin/bench-status.sh",
	    ".agents/commands/bench-implement-spec.md",
    ".agents/skills/bench-craft-seams/SKILL.md",
    ".agents/skills/bench-implement-spec/SKILL.md",
    ".agents/skills/bench-implement-spec/agents/openai.yaml",
    ".bench/BENCH.md",
    ".bench/hooks/stop.sh",
    ".claude/README.md",
    ".codex/hooks.json",
  ]) {
    if (!files.has(required)) {
      console.error("gate: npm package missing " + required);
      bad = 1;
    }
  }
  for (const forbidden of [".claude/settings.local.json"]) {
    if (files.has(forbidden)) {
      console.error("gate: npm package includes local-only file " + forbidden);
      bad = 1;
    }
  }
  process.exit(bad);
' || fail=1

# 5. Kit conformance — AGENTS.md index stays in sync with disk, both directions.
#    a) every skill dir is referenced in AGENTS.md
for d in .agents/skills/*/; do
  name="$(basename "$d")"
  [ -f ".agents/commands/$name.md" ] && continue # command adapters are documented in .bench/BENCH.md
  grep -q "$name" AGENTS.md || err "skill '$name' on disk but not referenced in AGENTS.md"
done
#    b) every skill the index names exists on disk
for name in $(grep -oE 'skills/[a-z0-9-]+/SKILL\.md' AGENTS.md | sed -E 's#skills/([a-z0-9-]+)/SKILL\.md#\1#' | sort -u); do
  [ -d ".agents/skills/$name" ] || err "AGENTS.md indexes skill '$name' with no .agents/skills/$name on disk"
done
#    c) every command file is referenced as /name in AGENTS.md
for f in .agents/commands/*.md; do
  name="$(basename "$f" .md)"
  grep -q "/$name" AGENTS.md || err "command '/$name' on disk but not referenced in AGENTS.md"
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
)
for m in "${ss_markers[@]}"; do
  grep -qF "$m" .bench/BENCH.md || err "shared rule missing from canonical .bench/BENCH.md: \"$m\""
  ! grep -qF "$m" AGENTS.md || err "shared rule duplicated in AGENTS.md (it must live only in .bench/BENCH.md): \"$m\""
done
grep -qF 'canonical in `.bench/BENCH.md`' AGENTS.md || err "AGENTS.md lost its pointer to the canonical .bench/BENCH.md shared rules"

#    g) living docs must name commands that exist now. Historical specs/maps may
#       mention old command names only when explicitly marked on a line by itself, but
#       the cold pickup surface, live maps, and command/skill bodies must not point
#       agents at dead slash commands or Codex $bench-* adapters.
node <<'NODE' || fail=1
const fs = require("fs");
const path = require("path");

let bad = 0;
const err = msg => { console.error("gate: " + msg); bad = 1; };

const commandsDir = ".agents/commands";
const validSlash = new Set();
if (fs.existsSync(commandsDir)) {
  for (const f of fs.readdirSync(commandsDir)) {
    if (f.endsWith(".md")) validSlash.add("/" + path.basename(f, ".md"));
  }
}
for (const external of ["/model"]) validSlash.add(external);

const validCodex = new Set(
  [...validSlash]
    .filter(token => token.startsWith("/bench-"))
    .map(token => "$" + token.slice(1))
);

const files = [];
const addFile = file => {
  if (fs.existsSync(file)) files.push(file);
};
const walk = dir => {
  if (!fs.existsSync(dir)) return;
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, ent.name);
    if (ent.isDirectory()) {
      walk(file);
    } else if (
      ent.name === "SKILL.md" ||
      /\.(md|ya?ml|json|sh)$/.test(ent.name)
    ) {
      files.push(file);
    }
  }
};

for (const f of [
  "README.md",
  "AGENTS.md",
  ".bench/BENCH.md",
  ".bench/learnings.md",
  "CONTEXT.md",
  "HANDOFF.md",
  "CHANGELOG.md",
]) addFile(f);
walk("specs");
walk("decisions");
walk(".agents");

const knownStale = new Set([
  "/resynthesize",
  "/spec",
  "/grill",
  "/start-ideation",
  "/setup",
  "/build",
  "/prep-shift",
  "/fix-bug",
  "/verify-gate",
  "/map",
  "/diagnose",
  "/review",
  "/verify",
  "/shift",
]);
const historicalMarker = /^<!-- command-currency: historical -->$/m;
const slashRef = /(^|[\s([`"'])\/([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])/g;
const codexRef = /(^|[\s([`"'])\$([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])/g;

for (const file of [...new Set(files)].sort()) {
  let text = fs.readFileSync(file, "utf8");
  if (historicalMarker.test(text)) continue;
  if (file === ".bench/learnings.md") {
    text = text.split("<!-- entries below -->")[0];
  }
  if (file === "CHANGELOG.md") {
    text = text.split(/\n## /)[0];
  }
  const lines = text.split(/\n/);
  for (let i = 0; i < lines.length; i++) {
    let m;
    slashRef.lastIndex = 0;
    while ((m = slashRef.exec(lines[i])) !== null) {
      const token = "/" + m[2];
      if (!validSlash.has(token) && (token.startsWith("/bench-") || knownStale.has(token))) {
        err(`stale command reference ${token} in ${file}:${i + 1}`);
      }
    }
    codexRef.lastIndex = 0;
    while ((m = codexRef.exec(lines[i])) !== null) {
      const token = "$" + m[2];
      if (token.startsWith("$bench-") && !validCodex.has(token)) {
        err(`stale Codex adapter reference ${token} in ${file}:${i + 1}`);
      }
    }
  }
}
process.exit(bad);
NODE

#    h) shipped cold-pickup docs that list CLI commands must list the real subcommands
#       from bin/bench.sh. HANDOFF.md ships in the npm package, and .bench/BENCH.md is
#       the operating guide installed into consumer repos.
node <<'NODE' || fail=1
const fs = require("fs");

let bad = 0;
const err = msg => { console.error("gate: " + msg); bad = 1; };
const bench = fs.readFileSync("bin/bench.sh", "utf8");
const commands = [...bench.matchAll(/^  ([a-z][a-z-]*)\)\s/gm)].map(m => m[1]).sort();
for (const file of ["HANDOFF.md", ".bench/BENCH.md"]) {
  if (!fs.existsSync(file)) continue;
  const text = fs.readFileSync(file, "utf8");
  for (const cmd of commands) {
    if (!text.includes(`bench ${cmd}`)) err(`${file} does not list CLI command 'bench ${cmd}'`);
  }
}
process.exit(bad);
NODE

#    i) acceptance coverage maps are now part of the feature-build workflow. These
#       anchors are intentionally structural: they prove the command/skill surfaces
#       still carry the contract, while semantic completeness stays a review/dogfood
#       responsibility.
require_anchor() {
  file="$1"
  needle="$2"
  grep -qF "$needle" "$file" || err "$file missing acceptance coverage anchor: $needle"
}
require_anchor ".agents/commands/bench-write-spec.md" "acceptance coverage map"
require_anchor ".agents/commands/bench-write-spec.md" "why it catches the failure"
require_anchor ".agents/commands/bench-write-spec.md" "red signal"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "acceptance row"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "not TDD-able"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "call count"
require_anchor ".agents/commands/bench-implement-spec.md" "coverage table"
require_anchor ".agents/commands/bench-implement-spec.md" "already covered"
require_anchor ".agents/commands/bench-implement-spec.md" "turning red-to-green"
require_anchor ".agents/commands/bench-review-implementation.md" "acceptance coverage map"
require_anchor ".agents/commands/bench-review-implementation.md" "mapped behavior"

grep -qF 'session-start.sh' README.md || err "README layout omits .bench/hooks/session-start.sh"
grep -qF 'bench.sh' README.md || err "README layout omits the real bin/bench.sh filename"
grep -qF 'benchkit.md' README.md || err "README layout omits projects/benchkit.md"
! grep -qF '│   └── bench                 #' README.md || err "README layout still names bin/bench instead of bin/bench.sh"

# 6. shellcheck — stronger shell lint, best-effort (runs only when installed).
if command -v shellcheck >/dev/null 2>&1; then
  shellcheck -S warning bin/bench.sh .bench/hooks/*.sh || err "shellcheck reported issues"
fi

# 7. Canary — prove the gate's own checks still bite. For each known-broken fixture in
#    tests/canary/, run THIS gate against it in a throwaway repo and assert it goes red
#    WITH the fixture's targeted error substring. A fixture that stops biting means a
#    check rotted into an always-pass. Attribution is by substring, not isolation: a
#    minimal fixture over-fails on unrelated checks, and that is fine — we only assert
#    the targeted message is present. BENCH_CANARY_INNER marks the inner run so this
#    check skips itself and never recurses.
if [ "${BENCH_CANARY_INNER:-0}" != "1" ] && [ -d tests/canary ]; then
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
