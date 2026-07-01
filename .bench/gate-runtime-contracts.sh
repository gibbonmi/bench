# Runtime CLI contracts for the benchkit gate.

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
  bash "$root/bin/bench.sh" idea capture all the words >/dev/null 2>&1
  grep -qE '^- [0-9]{4}-[0-9]{2}-[0-9]{2}  capture all the words$' ROADMAP.md || { echo "idea did not join unquoted multi-word args (\$* not \$1)"; exit 1; }
  # A hand-edited ROADMAP.md whose last line lacks a trailing newline must not swallow
  # the next entry onto the same physical line (that also undercounts the '^- ' tally
  # roadmap/status read). Regression for the newline-normalization in idea().
  printf -- '- 2026-06-01  hand added' > ROADMAP.md
  bash "$root/bin/bench.sh" idea "after handedit" >/dev/null 2>&1
  [ "$(grep -c '^- ' ROADMAP.md)" = "2" ] || { echo "idea merged onto a newline-less last line ('^- ' count != 2)"; exit 1; }
  : > ROADMAP.md
  bash "$root/bin/bench.sh" roadmap | grep -qi 'empty' || { echo "roadmap on present-but-empty file did not report empty"; exit 1; }
) || err "bench idea/roadmap contract failed"
rm -rf "$tmp"

gci() { git -c user.email=bench@local -c user.name=bench "$@"; }

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  mkdir -p .bench sub
  printf '{"ok":true}\n' > package.json
  printf '#!/usr/bin/env bash\n[ -f package.json ]\n' > .bench/gate.sh
  chmod +x .bench/gate.sh
  ( cd sub && bash "$root/bin/bench.sh" gate ) || { echo ".bench/gate.sh did not run from repo root"; exit 1; }
) || err "bench gate repo-root cwd contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  cd "$tmp"
  git init -q
  mkdir sub
  printf '{"ok":true}\n' > package.json
  printf '#!/usr/bin/env bash\n[ -f package.json ]\n' > gate-root.sh
  chmod +x gate-root.sh
  ( cd sub && BENCH_GATE=./gate-root.sh bash "$root/bin/bench.sh" gate ) || { echo "BENCH_GATE did not run from repo root"; exit 1; }
) || err "bench gate BENCH_GATE cwd contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qiF 'clean — nothing pending' <<<"$out" || { echo "clean repo did not report all-clear"; exit 1; }
) || err "bench status clean contract failed"
rm -rf "$tmp"

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

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  printf 'green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n' > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qF 're-run the gate' <<<"$out" || { echo "stale gate cache did not surface re-run"; exit 1; }
  if grep -qiF 'clean — nothing pending' <<<"$out"; then echo "stale green read as a clean bill"; exit 1; fi
) || err "bench status stale-gate contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  printf 'green %s 2026-06-30T00:00:00Z\n' "$(gci rev-parse HEAD)" > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qiF 'clean — nothing pending' <<<"$out" || { echo "fresh-green gate was not silent"; exit 1; }
) || err "bench status fresh-green contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir decisions; printf '### Answer\n— (deferred)\n' > decisions/x.md; gci add -A; gci commit -q -m s
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qF 'craft-grill → /bench-write-spec' <<<"$out" || { echo "unresolved decision map did not surface craft-grill"; exit 1; }
) || err "bench status decisions contract failed"
rm -rf "$tmp"

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
  grep -qF '/bench-integrate-learnings' <<<"$out" || { echo "learnings dropped from the top five"; exit 1; }
  grep -qF 'split (craft-seams)' <<<"$out" || { echo "structure dropped from the top five"; exit 1; }
  grep -qF 'commit on green / push' <<<"$out" || { echo "git signal action string missing"; exit 1; }
  grep -qF 'resume or clean up' <<<"$out" || { echo "worktree signal action string missing"; exit 1; }
  if grep -qF 'craft-grill → /bench-write-spec' <<<"$out"; then echo "lowest-priority signal not dropped under the budget"; exit 1; fi
  rows="$(grep -cE '^  [a-z]' <<<"$out")"
  [ "$rows" -le 5 ] || { echo "budget exceeded five rows ($rows)"; exit 1; }
) || err "bench status budget contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  BENCH_SHIFT=1 BENCH_STOP_CHECKED=0 bash "$root/.bench/hooks/stop.sh" >/dev/null 2>&1 || true
  cache="$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  [ -f "$cache" ] || { echo "Stop hook did not write the gate cache"; exit 1; }
  grep -qE '^(green|red) [0-9a-f]+ [0-9T:Z-]+$' "$cache" || { echo "gate cache not <status> <sha> <iso8601>"; exit 1; }
) || err "bench status gate-cache write contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench; printf -- '- a [open]\n' > .bench/learnings.md; gci add -A; gci commit -q -m s
  hi="$(BENCH_LEARNINGS_FLOOR=2 bash "$root/bin/bench.sh" status)"
  if grep -qF '/bench-integrate-learnings' <<<"$hi"; then echo "floor=2 still surfaced a single open learning"; exit 1; fi
  lo="$(BENCH_LEARNINGS_FLOOR=1 bash "$root/bin/bench.sh" status)"
  grep -qF '/bench-integrate-learnings' <<<"$lo" || { echo "floor=1 did not surface the open learning"; exit 1; }
) || err "bench status learnings-floor contract failed"
rm -rf "$tmp"

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
  [ "$(gci rev-list --count "$before_head..$branch")" = "1" ] || { echo "shift branch has the wrong commit count"; exit 1; }
  if gci worktree list --porcelain | grep -qF "branch refs/heads/$branch"; then echo "released worktree still holds the shift branch"; exit 1; fi
  [ -z "$(find "$home" -name bench-lease -print 2>/dev/null)" ] || { echo "shift worktree lease was not released"; exit 1; }
  rm -rf "$home"
) || err "bench shift worktree-isolation contract failed"
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
  seq 401 | sed 's/^/x=/g' > big.sh
  gci add big.sh
  out="$(bash "$root/bin/bench.sh" structure 2>&1)" && { printf '%s\n' "$out"; echo "shell source over the line budget did not fail structure"; exit 1; }
  grep -qF 'FILE TOO LONG' <<<"$out" || { echo "structure did not report shell file length"; exit 1; }
) || err "bench structure shell-file contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p "space dir"
  for i in $(seq 1 13); do printf 'x=%s\n' "$i" > "space dir/file$i.sh"; done
  gci add "space dir"
  out="$(BENCH_MAX_DIR_FILES=12 bash "$root/bin/bench.sh" structure 2>&1)" && { printf '%s\n' "$out"; echo "crowded path-with-spaces directory did not fail structure"; exit 1; }
  grep -qF 'space dir/' <<<"$out" || { printf '%s\n' "$out"; echo "structure did not preserve the directory with spaces"; exit 1; }
  if grep -qF '   ./' <<<"$out" || grep -qF '   dir/' <<<"$out"; then
    printf '%s\n' "$out"; echo "structure split a directory with spaces"; exit 1
  fi
) || err "bench structure path-with-spaces contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  printf 'ignored/\n' > .gitignore; gci add .gitignore; gci commit -q -m ignore
  cat > wt-shell <<'EOF'
#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}"
pwd >> "$BENCH_WT_RECORD"
lease="$(git rev-parse --git-path bench-lease)"
[ -f "$lease" ] || { echo "lease missing"; exit 7; }
[ ! -e dirty.txt ] || { echo "dirty file carried into reused worktree"; exit 8; }
[ ! -e ignored/leak.txt ] || { echo "ignored artifact carried into reused worktree"; exit 9; }
echo dirty > dirty.txt
mkdir -p ignored; echo leak > ignored/leak.txt
EOF
  chmod +x wt-shell
  record="$tmp/paths"
  BENCH_HOME="$tmp/.bh" BENCH_WT_RECORD="$record" SHELL="$tmp/wt-shell" bash "$root/bin/bench.sh" worktree >wt1.out 2>&1
  BENCH_HOME="$tmp/.bh" BENCH_WT_RECORD="$record" SHELL="$tmp/wt-shell" bash "$root/bin/bench.sh" worktree >wt2.out 2>&1
  mapfile -t paths < "$record"
  [ "${#paths[@]}" = "2" ] || { echo "worktree shell did not run twice"; exit 1; }
  [ "${paths[0]}" = "${paths[1]}" ] || { echo "worktree pool did not reuse a clean released path"; exit 1; }
  [ ! -f "$(git -C "${paths[1]}" rev-parse --git-path bench-lease)" ] || { echo "worktree lease was not removed on release"; exit 1; }
  [ ! -f "${paths[1]}/dirty.txt" ] || { echo "worktree release did not clean dirty files"; exit 1; }
  [ ! -f "${paths[1]}/ignored/leak.txt" ] || { echo "worktree release did not clean ignored artifacts"; exit 1; }
) || err "bench worktree lease/reuse contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  guard="$root/.bench/hooks/block-dangerous-git.sh"
  run_guard() {
    local command="$1" out rc
    out="$(printf '{"tool_input":{"command":"%s"}}\n' "$command" | bash "$guard" 2>&1)" && rc=0 || rc=$?
    [ "$rc" = "2" ] || { printf '%s\n' "$out"; echo "guard allowed $command (exit $rc)"; exit 1; }
    grep -qF 'BLOCKED:' <<<"$out" || { echo "guard block output was not actionable"; exit 1; }
  }
  run_guard 'git push'
  run_guard 'git -C . push'
  run_guard 'git -C /tmp reset --hard'
  run_guard 'git -C . clean -fd'
  run_guard 'git branch -D old-work'
  run_guard 'git rebase main'
  run_guard 'git checkout -- README.md'
  run_guard 'git checkout --pathspec-from-file=paths.txt'
  run_guard 'git checkout --pathspec-from-file paths.txt'
  run_guard 'git restore README.md'
  run_guard 'git restore --pathspec-from-file=paths.txt'
  run_guard 'git restore --pathspec-from-file paths.txt'
  allowed="$(printf '{"tool_input":{"command":"git -C . status --short"}}\n' | bash "$guard" 2>&1)" || { printf '%s\n' "$allowed"; echo "guard blocked harmless git status"; exit 1; }
) || err "block-dangerous-git global-option contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u
  repo="$tmp/repo"; bin_dir="$tmp/bin"; shim="$tmp/shim"
  mkdir -p "$repo" "$bin_dir" "$shim"
  ln -s "$root/bin/bench.sh" "$bin_dir/bench"
  cat > "$shim/readlink" <<'EOF'
#!/usr/bin/env bash
if [ "${1:-}" = "-f" ]; then exit 1; fi
/usr/bin/readlink "$@"
EOF
  chmod +x "$shim/readlink"
  cd "$repo"
  git init -q
  PATH="$shim:/usr/bin:/bin" "$bin_dir/bench" link >"$tmp/link.out" 2>&1
  [ -f .bench/BENCH.md ] || { echo "symlinked bench did not resolve kit dir without readlink -f"; exit 1; }
) || err "bench symlinked kit-dir portability contract failed"
rm -rf "$tmp"
