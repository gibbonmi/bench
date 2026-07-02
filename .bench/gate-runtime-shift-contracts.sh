# Runtime contracts for the `bench shift` gated loop and its harness adapters.
# Split out of gate-runtime-contracts.sh to keep each fragment under budget; the
# gate sources it in its own shell (shares $root, $gate_dir, err(), fail). Without
# a CLI in the tree there is nothing to contract-test; the skip is a distinct red
# so canary fixtures that plant a broken CLI stay attributable to their assertion.
[ -f "$root/bin/bench.sh" ] || { err "bench CLI missing (runtime contracts skipped)"; return 0 2>/dev/null || exit 0; }

gci() { git -c user.email=bench@local -c user.name=bench "$@"; }

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  before_branch="$(gci branch --show-current)"
  before_status="$(gci status --porcelain)"
  out="$(BENCH_AGENT=true BENCH_MAX_ITERS=1 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift noop 2>&1)" || true
  grep -qF 'shift done' <<<"$out" || { echo "bench shift loop did not complete (run_gate exec'd the gate?)"; exit 1; }
  branch="$(sed -n 's/^■ shift done: \([^,]*\),.*/\1/p' <<<"$out")"
  [ -n "$branch" ] || { echo "shift summary did not name the branch"; exit 1; }
  [ "$(gci branch --show-current)" = "$before_branch" ] || { echo "bench shift changed the main checkout branch"; exit 1; }
  [ "$(gci status --porcelain)" = "$before_status" ] || { echo "bench shift dirtied the main checkout"; exit 1; }
  gci rev-parse --verify "$branch" >/dev/null || { echo "shift branch was not preserved"; exit 1; }
  if gci worktree list --porcelain | grep -qF "branch refs/heads/$branch"; then echo "released worktree still holds the shift branch"; exit 1; fi
  [ -z "$(find "$home" -name bench-lease -print 2>/dev/null)" ] || { echo "shift worktree lease was not released"; exit 1; }
  rm -rf "$home"
) || err "bench shift gated-loop contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
printf 'shifted\n' > shifted.txt
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  before_branch="$(gci branch --show-current)"
  before_head="$(gci rev-parse HEAD)"
  before_status="$(gci status --porcelain)"
  out="$(BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=1 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift isolate 2>&1)" || true
  grep -qF '1 committed iteration(s)' <<<"$out" || { echo "shift did not commit the green iteration"; exit 1; }
  branch="$(sed -n 's/^■ shift done: \([^,]*\),.*/\1/p' <<<"$out")"
  [ -n "$branch" ] || { echo "shift summary did not name the branch"; exit 1; }
  [ "$(gci branch --show-current)" = "$before_branch" ] || { echo "shift changed the main checkout branch"; exit 1; }
  [ "$(gci rev-parse HEAD)" = "$before_head" ] || { echo "shift moved main checkout HEAD"; exit 1; }
  [ "$(gci status --porcelain)" = "$before_status" ] || { echo "shift dirtied the main checkout"; exit 1; }
  gci cat-file -e "$branch:shifted.txt" || { echo "shift branch does not contain the committed work"; exit 1; }
  [ "$(gci config "branch.$branch.benchBase")" = "$before_head" ] \
    || { echo "shift did not record the pre-shift HEAD in branch.<name>.benchBase"; exit 1; }
  [ "$(gci rev-list --count "$before_head..$branch")" = "1" ] || { echo "shift branch has the wrong commit count"; exit 1; }
  if gci worktree list --porcelain | grep -qF "branch refs/heads/$branch"; then echo "released worktree still holds the shift branch"; exit 1; fi
  [ -z "$(find "$home" -name bench-lease -print 2>/dev/null)" ] || { echo "shift worktree lease was not released"; exit 1; }
  rm -rf "$home"
) || err "bench shift worktree-isolation contract failed"
rm -rf "$tmp"

# Staging contract: an iteration commit carries exactly what the agent touched.
# The gate here drops an unignored byproduct — it must stay out of the commit in
# the same iteration (snapshot precedes the gate run) and in the next one
# (pre-agent dirt is subtracted). Touched paths keep spaces and glob characters.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench
  printf '#!/usr/bin/env bash\nprintf cache > gate-artifact.txt\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n+1)); printf '%s\n' "$n" > count
printf 'work\n' > "step $n [a].txt"
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  out="$(BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=2 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift stage-touched 2>&1)" || true
  grep -qF '2 committed iteration(s)' <<<"$out" || { echo "stage-touched shift did not commit both iterations"; exit 1; }
  branch="$(sed -n 's/^■ shift done: \([^,]*\),.*/\1/p' <<<"$out")"
  gci cat-file -e "$branch:step 1 [a].txt" || { echo "touched path with space+glob chars was not staged"; exit 1; }
  gci cat-file -e "$branch:step 2 [a].txt" || { echo "second iteration's touched path was not staged"; exit 1; }
  if gci cat-file -e "$branch:gate-artifact.txt" 2>/dev/null; then echo "gate byproduct rode into an iteration commit"; exit 1; fi
  rm -rf "$home"
) || err "bench shift stage-touched contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 1\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
printf 'red\n' > red.txt
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  before_branch="$(gci branch --show-current)"
  before_head="$(gci rev-parse HEAD)"
  before_status="$(gci status --porcelain)"
  out="$(BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=1 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift red-rollback 2>&1)" || true
  grep -qF 'red gate' <<<"$out" || { echo "red shift did not report rollback"; exit 1; }
  branch="$(sed -n 's/^■ shift done: \([^,]*\),.*/\1/p' <<<"$out")"
  [ -n "$branch" ] || { echo "red shift summary did not name the branch"; exit 1; }
  [ "$(gci branch --show-current)" = "$before_branch" ] || { echo "red shift changed the main checkout branch"; exit 1; }
  [ "$(gci rev-parse HEAD)" = "$before_head" ] || { echo "red shift moved main checkout HEAD"; exit 1; }
  [ "$(gci status --porcelain)" = "$before_status" ] || { echo "red shift dirtied the main checkout"; exit 1; }
  if gci cat-file -e "$branch:red.txt" 2>/dev/null; then echo "red shift preserved rolled-back work"; exit 1; fi
  [ "$(gci rev-list --count "$before_head..$branch")" = "0" ] || { echo "red shift branch gained a commit"; exit 1; }
  [ -z "$(find "$home" -name bench-lease -print 2>/dev/null)" ] || { echo "red shift worktree lease was not released"; exit 1; }
  rm -rf "$home"
) || err "bench shift red-rollback isolation contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  seq 401 | sed 's/^/x = /' > preexisting.py
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  out="$(BENCH_AGENT=true BENCH_MAX_ITERS=1 BENCH_REFACTOR_ITERS=1 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift noop 2>&1)" || true
  grep -qF 'shift done' <<<"$out" || { echo "shift with unrelated structural debt did not complete"; exit 1; }
  if grep -qF 'refactor phase' <<<"$out"; then echo "pre-existing structural debt triggered refactor phase"; exit 1; fi
  rm -rf "$home"
) || err "bench shift touched-scope structure contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
if [ ! -f made-big ]; then
  seq 401 | sed 's/^/x = /' > touched.py
  : > made-big
fi
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  out="$(BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=1 BENCH_REFACTOR_ITERS=3 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift make-big 2>&1)" || true
  grep -qF 'refactor phase' <<<"$out" || { echo "touched over-budget file did not trigger refactor phase"; exit 1; }
  grep -qF 'refactor 1 made no staged change' <<<"$out" || { echo "no-op refactor pass did not report no staged change"; exit 1; }
  if grep -qF 'refactor 2/' <<<"$out"; then echo "no-op refactor pass did not exit early"; exit 1; fi
  if grep -qF 'refactor 1 committed' <<<"$out"; then echo "no-op refactor pass reported a phantom commit"; exit 1; fi
  if grep -qF '/improve-codebase-architecture' <<<"$out"; then echo "shift fallback suggests an unbundled command"; exit 1; fi
  branch="$(sed -n 's/^■ shift done: \([^,]*\),.*/\1/p' <<<"$out")"
  [ "$(gci rev-list --count HEAD.."$branch")" = "1" ] || { echo "no-op refactor created an unexpected commit"; exit 1; }
  rm -rf "$home"
) || err "bench shift refactor no-op contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > interrupt-agent <<'EOF'
#!/usr/bin/env bash
kill -INT "$PPID"
exit 130
EOF
  chmod +x interrupt-agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  interrupt_log="$(mktemp)"
  if BENCH_AGENT="$tmp/interrupt-agent" BENCH_MAX_ITERS=2 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift interrupt >"$interrupt_log" 2>&1; then
    echo "interrupted shift exited successfully"; exit 1
  fi
  rm -f "$interrupt_log"
  [ ! -e .bench-objective ] || { echo "interrupted shift left .bench-objective"; exit 1; }
  [ ! -e .bench-notes.md ] || { echo "interrupted shift left .bench-notes.md"; exit 1; }
  [ -z "$(find "$home" -name bench-lease -print 2>/dev/null)" ] || { echo "interrupted shift left a leased worktree"; exit 1; }
  sleep 1
  out="$(BENCH_AGENT=true BENCH_MAX_ITERS=1 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift after-interrupt 2>&1)" || { printf '%s\n' "$out"; exit 1; }
  grep -qF 'shift done' <<<"$out" || { echo "follow-up shift after interrupt did not complete"; exit 1; }
  rm -rf "$home"
) || err "bench shift interrupt cleanup contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench
  printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  printf '#!/usr/bin/env bash\n[ -f step1.txt ]\n' > .bench/done.sh; chmod +x .bench/done.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n + 1))
printf '%s\n' "$n" > count
printf '%s\n' "$n" > "step$n.txt"
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  out="$(BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=3 BENCH_HOME="$home" bash "$root/bin/bench.sh" shift done-early 2>&1)" || true
  grep -qF 'objective met.' <<<"$out" || { echo ".bench/done.sh did not mark the objective met"; exit 1; }
  if grep -qF 'iteration 2/3' <<<"$out"; then echo ".bench/done.sh did not stop before the second iteration"; exit 1; fi
  grep -qF '1 committed iteration(s)' <<<"$out" || { echo "shift summary did not report the committed iteration count"; exit 1; }
  branch="$(sed -n 's/^■ shift done: \([^,]*\),.*/\1/p' <<<"$out")"
  gci cat-file -e "$branch:step1.txt" || { echo "done.sh shift branch does not contain the completed work"; exit 1; }
  rm -rf "$home"
) || err "bench shift done.sh early-completion contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench
  printf '#!/usr/bin/env bash\n[ ! -f junk.txt ]\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
n=0; [ -f "$BENCH_TEST_STATE/count" ] && n="$(cat "$BENCH_TEST_STATE/count")"
n=$((n+1)); printf '%s\n' "$n" > "$BENCH_TEST_STATE/count"
if [ "$n" = 1 ]; then
  printf 'tried A, broke gate\n' >> .bench-notes.md
  printf 'junk\n' > junk.txt
else
  [ -f .bench-notes.md ] && grep -q 'tried A' .bench-notes.md && printf 'notes-survived\n' >> "$BENCH_TEST_STATE/report"
  [ -f .bench-objective ] && printf 'objective-survived\n' >> "$BENCH_TEST_STATE/report"
  printf 'ok\n' > done.txt
fi
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"; state="$(mktemp -d)"
  BENCH_TEST_STATE="$state" BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=2 BENCH_HOME="$home" \
    bash "$root/bin/bench.sh" shift survive >/dev/null 2>&1 || true
  grep -q 'notes-survived' "$state/report" 2>/dev/null || { echo "red rollback wiped .bench-notes.md"; exit 1; }
  grep -q 'objective-survived' "$state/report" 2>/dev/null || { echo "red rollback wiped .bench-objective"; exit 1; }
  rm -rf "$home" "$state"
) || err "bench shift scratch-survival contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > agent <<'EOF'
#!/usr/bin/env bash
printf '%s\n@@@@\n' "$1" >> "$BENCH_TEST_PROMPTS"
if [ ! -f made-big ]; then seq 401 | sed 's/^/x = /' > touched.py; : > made-big; fi
EOF
  chmod +x agent
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"; prompts="$tmp/prompts.txt"
  BENCH_TEST_PROMPTS="$prompts" BENCH_AGENT="$tmp/agent" BENCH_MAX_ITERS=1 BENCH_REFACTOR_ITERS=1 BENCH_HOME="$home" \
    bash "$root/bin/bench.sh" shift make-big >/dev/null 2>&1 || true
  refactor="$(sed -n '/@@@@/,$p' "$prompts" 2>/dev/null)"
  grep -qF 'touched.py' <<<"$refactor" || { echo "refactor prompt does not name the flagged touched files"; exit 1; }
  ! grep -qF 'Run `bench structure` to see the flagged files' <<<"$refactor" || { echo "refactor prompt still points at repo-wide structure output"; exit 1; }
  rm -rf "$home"
) || err "bench shift refactor-prompt scope contract failed"
rm -rf "$tmp"

# Adapter contract: bench shift drives a configured adapter executable, passing the
# generated prompt as its single positional argument with BENCH_SHIFT=1 armed. There
# is no default harness — misconfiguration fails fast in a preflight, before any
# agent or gate run.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"
  if out="$(env -u BENCH_AGENT BENCH_HOME="$home" bash "$root/bin/bench.sh" shift probe 2>&1)"; then
    echo "shift with no BENCH_AGENT succeeded; should error"; exit 1
  fi
  grep -qF 'BENCH_AGENT' <<<"$out" || { echo "unconfigured-adapter error does not name BENCH_AGENT"; exit 1; }
  grep -qiE 'configure.*adapter|adapter.*configure' <<<"$out" || { echo "unconfigured-adapter error is not a configure-your-adapter message"; exit 1; }
  if grep -qF 'iteration 1/' <<<"$out"; then echo "unconfigured adapter still entered the loop"; exit 1; fi
  if BENCH_AGENT= BENCH_HOME="$home" bash "$root/bin/bench.sh" shift probe >/dev/null 2>&1; then
    echo "shift with empty BENCH_AGENT succeeded; should error"; exit 1
  fi
  if out="$(BENCH_AGENT=/no/such/adapter BENCH_HOME="$home" bash "$root/bin/bench.sh" shift probe 2>&1)"; then
    echo "shift with a missing adapter path succeeded; should error"; exit 1
  fi
  grep -qiF 'not executable' <<<"$out" || { echo "missing-adapter error does not say not executable"; exit 1; }
  if grep -qF 'iteration 1/' <<<"$out"; then echo "missing adapter still entered the loop"; exit 1; fi
  # a shell keyword resolves via command -v but is not an executable file; the
  # preflight must reject it rather than let every iteration exec-fail silently
  if out="$(BENCH_AGENT=if BENCH_HOME="$home" bash "$root/bin/bench.sh" shift probe 2>&1)"; then
    echo "shift with a shell-keyword adapter succeeded; should error"; exit 1
  fi
  grep -qiF 'not executable' <<<"$out" || { echo "shell-keyword adapter error does not say not executable"; exit 1; }
  rm -rf "$home"
) || err "bench shift adapter preflight contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  cat > adapter <<'EOF'
#!/usr/bin/env bash
{
  printf 'argc=%s\n' "$#"
  printf 'shift_env=%s\n' "${BENCH_SHIFT:-unset}"
  printf '%s\n@@@@\n' "$1"
} >> "$BENCH_TEST_RECORD"
EOF
  chmod +x adapter
  gci add -A; gci commit -q -m init
  home="$(mktemp -d)"; record="$tmp/record.txt"
  BENCH_TEST_RECORD="$record" BENCH_AGENT="$tmp/adapter" BENCH_MAX_ITERS=1 BENCH_HOME="$home" \
    bash "$root/bin/bench.sh" shift adapter-arg-probe >/dev/null 2>&1 || true
  [ -f "$record" ] || { echo "adapter was never invoked"; exit 1; }
  grep -qF 'argc=1' "$record" || { echo "prompt was not the adapter's single positional argument"; exit 1; }
  grep -qF 'shift_env=1' "$record" || { echo "BENCH_SHIFT=1 not armed on the adapter call"; exit 1; }
  if grep -qxF -- '-p' "$record"; then echo "loop still passes the Claude-specific -p flag"; exit 1; fi
  grep -qF 'adapter-arg-probe' "$record" || { echo "objective missing from the adapter argument"; exit 1; }
  grep -qF 'You are one iteration of a Bench shift' "$record" || { echo "prompt head missing from \$1"; exit 1; }
  grep -qF 'decides if it counts' "$record" || { echo "multi-line prompt tail missing from \$1 (prompt split or re-tokenized?)"; exit 1; }
  rm -rf "$home"
) || err "bench shift adapter single-argument contract failed"
rm -rf "$tmp"

(
  set -u
  for a in claude codex opencode; do
    f="$root/.bench/adapters/$a"
    [ -f "$f" ] || { echo "reference adapter missing: .bench/adapters/$a"; exit 1; }
    [ -x "$f" ] || { echo "reference adapter not executable: .bench/adapters/$a"; exit 1; }
    bash -n "$f" || { echo "reference adapter has a syntax error: .bench/adapters/$a"; exit 1; }
    grep -qE '^exec ' "$f" || { echo "reference adapter $a does not exec its harness (exit code must pass through)"; exit 1; }
    grep -qF '"$1"' "$f" || { echo "reference adapter $a does not pass the prompt as \$1"; exit 1; }
  done
  grep -qF 'claude -p -- "$1"' "$root/.bench/adapters/claude" || { echo "claude adapter does not map the prompt to claude -p (behind the -- sentinel)"; exit 1; }
  grep -qF 'codex exec -- "$1"' "$root/.bench/adapters/codex" || { echo "codex adapter does not map the prompt to codex exec (behind the -- sentinel)"; exit 1; }
  grep -qF 'opencode run -- "$1"' "$root/.bench/adapters/opencode" || { echo "opencode adapter does not map the prompt to opencode run (behind the -- sentinel)"; exit 1; }
) || err "reference adapter files contract failed"
