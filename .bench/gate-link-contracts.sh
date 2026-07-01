# Link/adoption contracts for the benchkit gate.

tmp="$(mktemp -d)"
( cd "$tmp" && git init -q && bash "$root/bin/bench.sh" init >/dev/null 2>&1 )
[ -f "$tmp/.bench/learnings.md" ] || err "bench init does not scaffold .bench/learnings.md (self-learning journal)"
rm -rf "$tmp"

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
  [ -x .bench/bin/bench.sh ] || { echo "fresh link did not install local hook CLI .bench/bin/bench.sh"; exit 1; }
  [ -f .bench/bin/bench-link.sh ] || { echo "fresh link did not install local hook CLI link helper"; exit 1; }
  [ -f .bench/bin/bench-status.sh ] || { echo "fresh link did not install local hook CLI status helper"; exit 1; }
  [ -f .bench/link-manifest.tsv ] || { echo "fresh link did not write link manifest"; exit 1; }
  [ -f .agents/commands/bench-implement-spec.md ] || { echo "fresh link did not install portable commands"; exit 1; }
  [ -f .agents/skills/bench-craft-seams/SKILL.md ] || { echo "fresh link did not install portable skills"; exit 1; }
  [ -f .agents/skills/bench-implement-spec/SKILL.md ] || { echo "fresh link did not install Codex command adapter skills"; exit 1; }
  [ -f .agents/skills/bench-implement-spec/agents/openai.yaml ] || { echo "fresh link did not install Codex command adapter metadata"; exit 1; }
  [ -f .claude/README.md ] || { echo "fresh link did not install Claude adapter README"; exit 1; }
  grep -qF '.agents/' .claude/README.md || { echo "Claude adapter README does not explain .agents"; exit 1; }
  grep -qF '.bench/hooks/' .claude/README.md || { echo "Claude adapter README does not explain shared hooks"; exit 1; }
  [ -e .claude/commands/bench-implement-spec.md ] || { echo "fresh link did not install Claude command adapter"; exit 1; }
  [ -e .claude/skills/bench-craft-seams/SKILL.md ] || { echo "fresh link did not install Claude skill adapter"; exit 1; }
  [ -f .codex/hooks.json ] || { echo "fresh link did not install Codex hook adapter"; exit 1; }
  [ -f .bench/hooks/block-dangerous-git.sh ] || { echo "fresh link did not install shared hook scripts"; exit 1; }
  [ -f .bench/hooks/session-start.sh ] || { echo "fresh link did not install the SessionStart hook"; exit 1; }
  grep -q 'SessionStart' .claude/settings.json || { echo "fresh link .claude/settings.json has no SessionStart wiring"; exit 1; }
  [ -x .git/hooks/pre-push ] || { echo "fresh link did not install git pre-push hook"; exit 1; }
  [ ! -L .agents/commands/bench-implement-spec.md ] || { echo "default link mode symlinked portable commands"; exit 1; }
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
  printf 'PROJECT BEFORE\n<!-- bench:end -->\nPROJECT MIDDLE\n<!-- bench:start -->\nPROJECT AFTER\n' > AGENTS.md
  if bash "$root/bin/bench.sh" link >link.out 2>&1; then
    echo "link succeeded despite reversed Bench managed block markers"; exit 1
  fi
  grep -qi 'malformed' link.out || { echo "malformed marker output did not explain the conflict"; exit 1; }
  grep -qF 'PROJECT AFTER' AGENTS.md || { echo "malformed marker failure still rewrote project-owned text"; exit 1; }
) || err "bench link malformed marker contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  mkdir -p .agents/commands
  printf 'project command\n' > .agents/commands/bench-implement-spec.md
  if bash "$root/bin/bench.sh" link >link.out 2>&1; then
    echo "link succeeded despite a project-owned command conflict"; exit 1
  fi
  grep -qi 'conflict' link.out || { echo "conflict output did not explain the conflict"; exit 1; }
  grep -qF 'project command' .agents/commands/bench-implement-spec.md || { echo "conflicting project command was overwritten"; exit 1; }
  [ ! -f .bench/link-manifest.tsv ] || { echo "conflicting link wrote a manifest despite failing"; exit 1; }
) || err "bench link conflict contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  check_link_contract "$tmp"
  printf '\nlocal edit\n' >> .agents/commands/bench-implement-spec.md
  if bash "$root/bin/bench.sh" link >relink.out 2>&1; then
    echo "relink overwrote a locally modified managed file"; exit 1
  fi
  grep -qi 'modified' relink.out || { echo "modified-managed output did not explain the local edit"; exit 1; }
) || err "bench link modified-managed contract failed ($(cat "$tmp/relink.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  check_link_contract "$tmp"
  mkdir -p .bench
  printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh
  chmod +x .bench/gate.sh
  git -c user.email=bench@local -c user.name=bench add -A
  git -c user.email=bench@local -c user.name=bench commit -q -m linked
  PATH=/usr/bin:/bin .bench/hooks/session-start.sh >/dev/null
  BENCH_SHIFT=1 BENCH_STOP_CHECKED=0 PATH=/usr/bin:/bin .bench/hooks/stop.sh >/dev/null 2>&1
) || err "bench linked hooks local-CLI contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"
