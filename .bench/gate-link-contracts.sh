# Link/adoption contracts for the benchkit gate. Skipped as a distinct red when
# the tree carries no CLI, keeping canary fixtures attributable.
[ -f "$root/bin/bench.sh" ] || { err "bench CLI missing (link contracts skipped)"; return 0 2>/dev/null || exit 0; }

tmp="$(mktemp -d)"
( cd "$tmp" && git init -q && bash "$root/bin/bench.sh" init >/dev/null 2>&1 )
[ -f "$tmp/.bench/learnings.md" ] || err "bench init does not scaffold .bench/learnings.md (self-learning journal)"

# bench init scaffolds a working, self-defending gate — the tripwire (shared canary
# runner + seed fixture + a red-until-configured sentinel) that gives a fresh consumer
# repo an automated defense against a self-weakened gate. These run only in the outer
# gate: they exercise scaffolded gates (each spawns its own inner canary run), too
# heavy to repeat during benchkit's own inner fixture runs, and they are behavioral
# contracts, not canary-attributed checks.
if [ "${BENCH_CANARY_INNER:-0}" != "1" ]; then
  # Story 5 — the live harness ships: shared runner + seed fixture.
  [ -f "$tmp/.bench/lib/canary-run.sh" ]    || err "bench init does not install the canary runner (.bench/lib/canary-run.sh)"
  [ -f "$tmp/tests/canary/example/EXPECT" ] || err "bench init does not scaffold the seed canary fixture (tests/canary/example/EXPECT)"
  [ -d "$tmp/tests/canary/example/files" ]  || err "bench init does not scaffold the seed canary fixture files/ tree"

  # Story 4 — a fresh scaffolded gate is red until configured (the sentinel).
  ( cd "$tmp" && bash .bench/gate.sh ) >"$tmp/g.sentinel" 2>&1; rc=$?
  { [ "$rc" -ne 0 ] && grep -qF 'configure .bench/gate.sh' "$tmp/g.sentinel"; } \
    || err "bench init scaffolds a gate that is not red-until-configured (want red + sentinel)"

  # Story 6 — sentinel removed → green: the example check passes clean, the seed canary
  # bites its fixture, and its EXPECT is not vacuous. A single exit-0 proves all three.
  ( cd "$tmp" && grep -v BENCH_SENTINEL .bench/gate.sh > .bench/gate.next \
      && mv .bench/gate.next .bench/gate.sh && chmod +x .bench/gate.sh \
      && bash .bench/gate.sh ) >"$tmp/g.green" 2>&1; rc=$?
  [ "$rc" -eq 0 ] || err "scaffolded harness not green after sentinel removed (seed canary vacuous or not biting)"

  # Story 2 — deleting tests/canary/ turns the configured gate red: the lazy escape is loud.
  ( cd "$tmp" && rm -rf tests/canary && bash .bench/gate.sh ) >"$tmp/g.absent" 2>&1; rc=$?
  { [ "$rc" -ne 0 ] && grep -qF 'canary harness absent' "$tmp/g.absent"; } \
    || err "err-if-absent: deleting tests/canary/ did not turn the gate red"

  # Story 3 — an emptied harness is caught the same as a deleted one.
  ( cd "$tmp" && mkdir -p tests/canary && bash .bench/gate.sh ) >"$tmp/g.empty" 2>&1; rc=$?
  { [ "$rc" -ne 0 ] && grep -qF 'canary harness absent' "$tmp/g.empty"; } \
    || err "err-if-absent: an empty tests/canary/ did not turn the gate red"

  # Story 8 — the inner-run guard holds: BENCH_CANARY_INNER=1 skips the whole block, so
  # err-if-absent never fires recursively (which would break every inner fixture run).
  ( cd "$tmp" && BENCH_CANARY_INNER=1 bash .bench/gate.sh ) >"$tmp/g.inner" 2>&1
  grep -qF 'canary harness absent' "$tmp/g.inner" \
    && err "err-if-absent fired during an inner run (BENCH_CANARY_INNER guard leaks)"
fi
rm -rf "$tmp"

# Story 7 — a second `bench init` never clobbers a configured gate. Fresh tmp so the
# check is independent of the mutations above; guarded for the same inner-run reason.
if [ "${BENCH_CANARY_INNER:-0}" != "1" ]; then
  tmp="$(mktemp -d)"
  ( cd "$tmp" && git init -q && bash "$root/bin/bench.sh" init >/dev/null 2>&1 \
      && printf '# configured by hand\n' >> .bench/gate.sh \
      && bash "$root/bin/bench.sh" init >/dev/null 2>&1 )
  grep -qF '# configured by hand' "$tmp/.bench/gate.sh" \
    || err "a second bench init clobbered an existing .bench/gate.sh"
  rm -rf "$tmp"
fi

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
  [ -f .bench/BENCH-reference.md ] || { echo "fresh link did not install .bench/BENCH-reference.md (the operating guide points to it)"; exit 1; }
  ! grep -qF '@.bench/BENCH-reference.md' CLAUDE.md || { echo "fresh link @-imported the reference file; it must stay on-demand"; exit 1; }
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
  [ ! -e .claude/skills/bench-implement-spec ] || { echo "fresh link installed a Codex-only phase adapter into .claude/skills (duplicates the /bench-implement-spec menu entry)"; exit 1; }
  [ -f .codex/hooks.json ] || { echo "fresh link did not install Codex hook adapter"; exit 1; }
  [ -f .bench/hooks/block-dangerous-git.sh ] || { echo "fresh link did not install shared hook scripts"; exit 1; }
  [ -x .bench/adapters/claude ] || { echo "fresh link did not install executable reference adapters"; exit 1; }
  [ -f .bench/lib/lines-env.sh ] || { echo "fresh link did not install the shared tier parser lib"; exit 1; }
  [ -f .bench/hooks/session-start.sh ] || { echo "fresh link did not install the SessionStart hook"; exit 1; }
  grep -q 'SessionStart' .claude/settings.json || { echo "fresh link .claude/settings.json has no SessionStart wiring"; exit 1; }
  [ -x .git/hooks/pre-push ] || { echo "fresh link did not install git pre-push hook"; exit 1; }
  grep -qF '@.bench/BENCH.md' CLAUDE.md || { echo "fresh link CLAUDE.md does not import .bench/BENCH.md"; exit 1; }
  [ ! -L .agents/commands/bench-implement-spec.md ] || { echo "default link mode symlinked portable commands"; exit 1; }
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  [ "$(count_literal '<!-- bench:start -->' AGENTS.md)" = "1" ] || { echo "relink duplicated managed Bench block"; exit 1; }
  printf '# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n' > CLAUDE.md
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  grep -qF '@.bench/BENCH.md' CLAUDE.md || { echo "relink did not retrofit the legacy bench-generated CLAUDE.md"; exit 1; }
  printf '# Custom\n\nproject-owned claude config\n' > CLAUDE.md
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  grep -qF 'project-owned claude config' CLAUDE.md || { echo "relink rewrote a project-owned CLAUDE.md"; exit 1; }
  ! grep -qF '@.bench/BENCH.md' CLAUDE.md || { echo "relink injected an import into a project-owned CLAUDE.md"; exit 1; }
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
  printf '{}\n' | BENCH_SHIFT=1 PATH=/usr/bin:/bin .bench/hooks/stop.sh >/dev/null 2>&1
) || err "bench linked hooks local-CLI contract failed ($(cat "$tmp/link.out" 2>/dev/null | tail -n 1))"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  kitcopy="$tmp/kit[1]"
  mkdir -p "$kitcopy/.bench"
  cp -R "$root/bin" "$root/.agents" "$root/.claude" "$root/.codex" "$kitcopy/"
  cp "$root/.bench/BENCH.md" "$kitcopy/.bench/BENCH.md"
  cp "$root/.bench/BENCH-reference.md" "$kitcopy/.bench/BENCH-reference.md"
  cp -R "$root/.bench/hooks" "$kitcopy/.bench/hooks"
  mkdir "$tmp/repo"; cd "$tmp/repo"; git init -q
  BENCH_KIT="$kitcopy" bash "$root/bin/bench.sh" link >link.out 2>&1 \
    || { tail -1 link.out; echo "link from a metachar kit path failed"; exit 1; }
  [ -x .bench/bin/bench.sh ] || { echo "metachar kit path scattered installed files"; exit 1; }
  [ -f .agents/commands/bench-implement-spec.md ] || { echo "metachar kit path lost portable commands"; exit 1; }
) || err "bench link metachar kit-path contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q main-repo
  cd main-repo
  git -c user.email=bench@local -c user.name=bench commit -q --allow-empty -m init
  git worktree add -q "$tmp/wt" -b side HEAD
  cd "$tmp/wt"
  bash "$root/bin/bench.sh" link >link.out 2>&1 \
    || { tail -1 link.out; echo "link inside a linked worktree failed"; exit 1; }
  hooks="$(git rev-parse --git-path hooks)"
  [ -x "$hooks/pre-push" ] || { echo "worktree link did not install pre-push in the effective hooks dir"; exit 1; }
) || err "bench link worktree contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  git config core.hooksPath .husky
  bash "$root/bin/bench.sh" link >link.out 2>&1 \
    || { tail -1 link.out; echo "link with core.hooksPath failed"; exit 1; }
  [ -x .husky/pre-push ] || { echo "pre-push not installed into configured hooksPath"; exit 1; }
  grep -q 'bench:managed-pre-push' .husky/pre-push || { echo "hooksPath pre-push is not bench-managed"; exit 1; }
) || err "bench link hooksPath contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  git init -q --bare -b master remote.git
  git init -q -b master repo && cd repo
  git -c user.email=bench@local -c user.name=bench commit -q --allow-empty -m init
  git remote add origin "$tmp/remote.git"
  git push -q origin master
  git fetch -q origin
  git symbolic-ref -d refs/remotes/origin/HEAD 2>/dev/null || true
  bash "$root/bin/bench.sh" link >link.out 2>&1 \
    || { tail -1 link.out; echo "link with unset origin/HEAD failed"; exit 1; }
  grep -q 'refs/heads/master' "$(git rev-parse --git-path hooks)/pre-push" \
    || { echo "pre-push guards the wrong branch when origin/HEAD is unset"; exit 1; }
) || err "bench link default-branch resolution contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  cat > AGENTS.md <<'EOF'
# Project rules

How Bench marks its block:

```
<!-- bench:start -->
managed content example
<!-- bench:end -->
```

KEEP-ME project text.
EOF
  bash "$root/bin/bench.sh" link >link.out 2>&1 \
    || { tail -1 link.out; echo "link failed on fenced marker docs"; exit 1; }
  grep -qF 'KEEP-ME project text.' AGENTS.md || { echo "fenced-marker link lost project text"; exit 1; }
  grep -qF 'managed content example' AGENTS.md || { echo "fenced example content was rewritten"; exit 1; }
  bash "$root/bin/bench.sh" link >relink.out 2>&1 \
    || { tail -1 relink.out; echo "relink failed on fenced marker docs"; exit 1; }
  grep -qF 'managed content example' AGENTS.md || { echo "relink consumed the fenced example"; exit 1; }
  [ "$(grep -cF '## Bench' AGENTS.md)" = "1" ] || { echo "fenced markers caused duplicate managed blocks"; exit 1; }
) || err "bench link fenced-marker contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  git config core.hooksPath .husky
  mkdir .husky
  printf '#!/bin/sh\nexit 0\n' > .husky/pre-push
  if bash "$root/bin/bench.sh" link >link.out 2>&1; then
    echo "link succeeded over a non-managed pre-push in hooksPath"; exit 1
  fi
  grep -qi 'conflict' link.out || { echo "hooksPath conflict output did not explain the conflict"; exit 1; }
  grep -qF 'exit 0' .husky/pre-push || { echo "hooksPath conflict overwrote the project hook"; exit 1; }
) || err "bench link hooksPath conflict contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  cat > AGENTS.md <<'EOF'
# Project rules

Broken docs with an unclosed fence:

```
<!-- bench:start -->
<!-- bench:end -->

KEEP-ME text after the unclosed fence.
EOF
  if bash "$root/bin/bench.sh" link >link.out 2>&1; then
    echo "link succeeded despite an unclosed fence around Bench markers"; exit 1
  fi
  grep -qi 'fence' link.out || { echo "unclosed-fence output did not explain the conflict"; exit 1; }
  grep -qF 'KEEP-ME text after the unclosed fence.' AGENTS.md || { echo "unclosed-fence failure rewrote project text"; exit 1; }
  [ "$(grep -cF '## Bench' AGENTS.md)" = "0" ] || { echo "unclosed-fence link still installed a managed block"; exit 1; }
) || err "bench link unclosed-fence contract failed"
rm -rf "$tmp"
