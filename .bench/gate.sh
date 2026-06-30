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

fail=0
err() { echo "gate: $*" >&2; fail=1; }

# 1. Shell scripts parse.
for f in bin/bench.sh .bench/hooks/*.sh; do
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

# 1d. `bench init` must scaffold the self-learning journal. AGENTS.md tells agents to
#     append to .bench/learnings.md and /resynthesize reads it; if init does not
#     create it, the contract points at a file that never exists. Exercise the real
#     init path in a throwaway repo rather than grepping for the literal.
tmp="$(mktemp -d)"
( cd "$tmp" && git init -q && bash "$root/bin/bench.sh" init >/dev/null 2>&1 )
[ -f "$tmp/.bench/learnings.md" ] || err "bench init does not scaffold .bench/learnings.md (self-learning journal)"
rm -rf "$tmp"

# 1e. `bench link` must safely incorporate Bench into an existing repo. This is the
#     adoption seam: exercise the real CLI against throwaway repos so clobbering,
#     duplicate managed blocks, unsafe conflicts, missing hook adapters, and wrong
#     default link mode all fail at the oracle.
check_link_contract() {
  local repo="$1"
  ( cd "$repo" && bash "$root/bin/bench.sh" link ) >"$repo/link.out" 2>&1
}

count_literal() {
  local needle="$1" file="$2"
  grep -oF "$needle" "$file" 2>/dev/null | wc -l | tr -d ' '
}

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  check_link_contract "$tmp"
  [ -f AGENTS.md ] || { echo "fresh link did not create AGENTS.md"; exit 1; }
  [ "$(count_literal '<!-- bench:start -->' AGENTS.md)" = "1" ] || { echo "fresh link did not create exactly one managed start marker"; exit 1; }
  [ "$(count_literal '<!-- bench:end -->' AGENTS.md)" = "1" ] || { echo "fresh link did not create exactly one managed end marker"; exit 1; }
  [ -f .bench/BENCH.md ] || { echo "fresh link did not install .bench/BENCH.md"; exit 1; }
  [ -f .bench/link-manifest.tsv ] || { echo "fresh link did not write link manifest"; exit 1; }
  [ -f .agents/commands/build.md ] || { echo "fresh link did not install portable commands"; exit 1; }
  [ -f .agents/skills/seams/SKILL.md ] || { echo "fresh link did not install portable skills"; exit 1; }
  [ -f .claude/README.md ] || { echo "fresh link did not install Claude adapter README"; exit 1; }
  grep -qF '.agents/' .claude/README.md || { echo "Claude adapter README does not explain .agents"; exit 1; }
  grep -qF '.bench/hooks/' .claude/README.md || { echo "Claude adapter README does not explain shared hooks"; exit 1; }
  [ -e .claude/commands/build.md ] || { echo "fresh link did not install Claude command adapter"; exit 1; }
  [ -e .claude/skills/seams/SKILL.md ] || { echo "fresh link did not install Claude skill adapter"; exit 1; }
  [ -f .codex/hooks.json ] || { echo "fresh link did not install Codex hook adapter"; exit 1; }
  [ -f .bench/hooks/block-dangerous-git.sh ] || { echo "fresh link did not install shared hook scripts"; exit 1; }
  [ -f .bench/hooks/session-start.sh ] || { echo "fresh link did not install the SessionStart hook"; exit 1; }
  grep -q 'SessionStart' .claude/settings.json || { echo "fresh link .claude/settings.json has no SessionStart wiring"; exit 1; }
  [ -x .git/hooks/pre-push ] || { echo "fresh link did not install git pre-push hook"; exit 1; }
  [ ! -L .agents/commands/build.md ] || { echo "default link mode symlinked portable commands"; exit 1; }
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  [ "$(count_literal '<!-- bench:start -->' AGENTS.md)" = "1" ] || { echo "relink duplicated managed Bench block"; exit 1; }
) || err "bench link safe fresh/relink contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  printf 'PROJECT RULES\n' > AGENTS.md
  check_link_contract "$tmp"
  grep -qF 'PROJECT RULES' AGENTS.md || { echo "existing AGENTS.md content was clobbered"; exit 1; }
  [ "$(count_literal '<!-- bench:start -->' AGENTS.md)" = "1" ] || { echo "existing AGENTS.md did not get exactly one managed block"; exit 1; }
) || err "bench link existing AGENTS.md contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  mkdir -p .agents/commands
  printf 'project command\n' > .agents/commands/build.md
  if bash "$root/bin/bench.sh" link >link.out 2>&1; then
    echo "link succeeded despite a project-owned command conflict"; exit 1
  fi
  grep -qi 'conflict' link.out || { echo "conflict output did not explain the conflict"; exit 1; }
  grep -qF 'project command' .agents/commands/build.md || { echo "conflicting project command was overwritten"; exit 1; }
  [ ! -f .bench/link-manifest.tsv ] || { echo "conflicting link wrote a manifest despite failing"; exit 1; }
) || err "bench link conflict contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  check_link_contract "$tmp"
  printf '\nlocal edit\n' >> .agents/commands/build.md
  if bash "$root/bin/bench.sh" link >relink.out 2>&1; then
    echo "relink overwrote a locally modified managed file"; exit 1
  fi
  grep -qi 'modified' relink.out || { echo "modified-managed output did not explain the local edit"; exit 1; }
) || err "bench link modified-managed contract failed ($(cat "$tmp/relink.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

# 1f. `bench idea` / `bench roadmap` — the capture-and-forget roadmap sink. Exercise the
#     real CLI in a throwaway repo: roadmap reports empty when absent, idea creates
#     ROADMAP.md and appends a dated line, a no-arg idea errors without appending a blank
#     entry, and roadmap prints parked ideas. A regression here means the capture path
#     silently stopped working.
tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  bash "$root/bin/bench.sh" roadmap | grep -qi 'empty' || { echo "roadmap on absent file did not report empty"; exit 1; }
  bash "$root/bin/bench.sh" idea "ship dark mode" >/dev/null 2>&1
  [ -f ROADMAP.md ] || { echo "idea did not create ROADMAP.md"; exit 1; }
  grep -qE '^- [0-9]{4}-[0-9]{2}-[0-9]{2}  ship dark mode$' ROADMAP.md || { echo "idea entry not dated '- YYYY-MM-DD  <text>'"; exit 1; }
  before="$(wc -l < ROADMAP.md)"
  if bash "$root/bin/bench.sh" idea >/dev/null 2>&1; then echo "no-arg idea succeeded; should error"; exit 1; fi
  after="$(wc -l < ROADMAP.md)"
  [ "$before" = "$after" ] || { echo "no-arg idea appended a blank entry"; exit 1; }
  bash "$root/bin/bench.sh" roadmap | grep -qF 'ship dark mode' || { echo "roadmap did not print the parked idea"; exit 1; }
  # The capture guarantee is "all words after the subcommand are the idea" — the join
  # is $* not $1, so exercise the unquoted multi-word form a single-arg test can't tell apart.
  bash "$root/bin/bench.sh" idea capture all the words >/dev/null 2>&1
  grep -qE '^- [0-9]{4}-[0-9]{2}-[0-9]{2}  capture all the words$' ROADMAP.md || { echo "idea did not join unquoted multi-word args (\$* not \$1)"; exit 1; }
  # Empty is the absence of parked ideas whether the file is missing or present-but-blank.
  : > ROADMAP.md
  bash "$root/bin/bench.sh" roadmap | grep -qi 'empty' || { echo "roadmap on present-but-empty file did not report empty"; exit 1; }
) || err "bench idea/roadmap contract failed"
rm -rf "$tmp"

# 1g. `bench status` — the ambient-feedback surface renderer. Construct repo state in
#     throwaway repos and assert the rendered output: the all-clear line, the gate
#     signal resolved from the cache (red / stale / silent), each signal's action
#     string, the zero-severity roadmap footer, and the five-row budget with a `+k more`
#     tail. Deterministic plain shell, so it is fully contract-testable here.
gci() { git -c user.email=bench@local -c user.name=bench "$@"; }
# A — clean repo → all-clear, no footer.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qiF 'clean — nothing pending' <<<"$out" || { echo "clean repo did not report all-clear"; exit 1; }
) || err "bench status clean contract failed"
rm -rf "$tmp"
# B — clean + committed ROADMAP.md → footer present, never the lead.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  printf -- '- 2026-06-30  an idea\n' > ROADMAP.md; gci add -A; gci commit -q -m s
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qiF 'clean — nothing pending' <<<"$out" || { echo "footer repo lost the all-clear line"; exit 1; }
  grep -qF 'parked — bench roadmap' <<<"$out" || { echo "roadmap footer missing"; exit 1; }
  if grep -qE '^▶.*bench roadmap' <<<"$out"; then echo "roadmap footer became the lead"; exit 1; fi
) || err "bench status footer contract failed"
rm -rf "$tmp"
# C — gate cache present but sha != HEAD → stale row, and NOT a clean bill.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  printf 'green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n' > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qF 're-run the gate' <<<"$out" || { echo "stale gate cache did not surface re-run"; exit 1; }
  if grep -qiF 'clean — nothing pending' <<<"$out"; then echo "stale green read as a clean bill"; exit 1; fi
) || err "bench status stale-gate contract failed"
rm -rf "$tmp"
# D — gate cache green AND fresh (sha == HEAD) → gate silent → all-clear.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  printf 'green %s 2026-06-30T00:00:00Z\n' "$(gci rev-parse HEAD)" > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qiF 'clean — nothing pending' <<<"$out" || { echo "fresh-green gate was not silent"; exit 1; }
) || err "bench status fresh-green contract failed"
rm -rf "$tmp"
# E — decision-map marker alone → the /grill → /spec action string.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir decisions; printf '### Answer\n— (deferred)\n' > decisions/x.md; gci add -A; gci commit -q -m s
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qF '/grill → /spec' <<<"$out" || { echo "unresolved decision map did not surface /grill"; exit 1; }
) || err "bench status decisions contract failed"
rm -rf "$tmp"
# F — six signals firing → gate red leads; budget caps at five rows + `+1 more`; the
#     lowest-priority signal (the decision map) is dropped under the tail.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf -- '- a [open]\n' > .bench/learnings.md
  seq 401 | sed 's/^/x = /' > big.py
  mkdir decisions; printf '### Answer\n— (deferred)\n' > decisions/x.md
  gci add -A; gci commit -q -m s
  printf 'red %s 2026-06-30T00:00:00Z\n' "$(gci rev-parse HEAD)" > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  echo dirty > dirty.txt
  gci worktree add -q --detach "$tmp/wt2" HEAD 2>/dev/null
  out="$(bash "$root/bin/bench.sh" status)"
  head -1 <<<"$out" | grep -qF 'fix before commit' || { echo "red gate did not lead the budget case"; exit 1; }
  grep -qF '+1 more' <<<"$out" || { echo "six signals did not trigger the +k more tail"; exit 1; }
  grep -qF '/resynthesize' <<<"$out" || { echo "learnings dropped from the top five"; exit 1; }
  grep -qF 'split (seams)' <<<"$out" || { echo "structure dropped from the top five"; exit 1; }
  grep -qF 'commit on green / push' <<<"$out" || { echo "git signal action string missing"; exit 1; }
  grep -qF 'resume or clean up' <<<"$out" || { echo "worktree signal action string missing"; exit 1; }
  if grep -qF '/grill → /spec' <<<"$out"; then echo "lowest-priority signal not dropped under the budget"; exit 1; fi
  rows="$(grep -cE '^  [a-z]' <<<"$out")"
  [ "$rows" -le 5 ] || { echo "budget exceeded five rows ($rows)"; exit 1; }
) || err "bench status budget contract failed"
rm -rf "$tmp"
# G — the Stop hook records the gate verdict to the git-dir cache, in the format
#     `bench status` reads back. Run it armed (BENCH_SHIFT=1) in a throwaway repo.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  BENCH_SHIFT=1 BENCH_STOP_CHECKED=0 bash "$root/.bench/hooks/stop.sh" >/dev/null 2>&1 || true
  cache="$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  [ -f "$cache" ] || { echo "Stop hook did not write the gate cache"; exit 1; }
  grep -qE '^(green|red) [0-9a-f]+ [0-9T:Z-]+$' "$cache" || { echo "gate cache not <status> <sha> <iso8601>"; exit 1; }
) || err "bench status gate-cache write contract failed"
rm -rf "$tmp"
# H — BENCH_LEARNINGS_FLOOR moves the open-learnings threshold (a single open learning is
#     silent at floor 2, surfaced at floor 1).
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf -- '- a [open]\n' > .bench/learnings.md; gci add -A; gci commit -q -m s
  hi="$(BENCH_LEARNINGS_FLOOR=2 bash "$root/bin/bench.sh" status)"
  if grep -qF '/resynthesize' <<<"$hi"; then echo "floor=2 still surfaced a single open learning"; exit 1; fi
  lo="$(BENCH_LEARNINGS_FLOOR=1 bash "$root/bin/bench.sh" status)"
  grep -qF '/resynthesize' <<<"$lo" || { echo "floor=1 did not surface the open learning"; exit 1; }
) || err "bench status learnings-floor contract failed"
rm -rf "$tmp"

# 2. JSON is valid.
for f in package.json .claude/settings.json .codex/hooks.json; do
  node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"))' "$f" \
    || err "invalid JSON in $f"
done

# 3. Every skill carries YAML frontmatter (first line is the --- fence).
for f in .agents/skills/*/SKILL.md; do
  [ "$(head -1 "$f")" = "---" ] || err "$f missing frontmatter"
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
    ".agents/commands/build.md",
    ".agents/skills/seams/SKILL.md",
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
#    d) the roadmap promotion seam — /start-ideation must name ROADMAP.md and the
#       auto-remove-on-map-creation behavior, or the only path that drains a parked idea
#       silently rots. The capture sink (bench idea) is useless without this graduation.
si=".agents/commands/start-ideation.md"
grep -qF 'ROADMAP.md' "$si" || err "/start-ideation does not reference ROADMAP.md (roadmap promotion seam)"
grep -qiE 'remove|delete' "$si" || err "/start-ideation does not describe removing a promoted roadmap entry"

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
