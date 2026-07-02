# Runtime CLI contracts for the benchkit gate. Without a CLI in the tree there is
# nothing to contract-test; the skip is a distinct red so canary fixtures that
# plant a broken CLI stay attributable to their targeted assertion.
[ -f "$root/bin/bench.sh" ] || { err "bench CLI missing (runtime contracts skipped)"; return 0 2>/dev/null || exit 0; }

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
  printf 'green %s 2026-06-30T00:00:00Z\n' "$(gci rev-parse 'HEAD^{tree}')" > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
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
  echo dirty > dirty.txt
  gci worktree add -q --detach "$tmp/wt2" HEAD 2>/dev/null
  # The red verdict must be keyed to the tree as status will see it (dirty file and
  # nested worktree included), or the stale row would outrank the red lead.
  # The throwaway index must live outside the repo, or it would join the tree it hashes.
  tree="$(export GIT_INDEX_FILE="${TMPDIR:-/tmp}/bench-budget-idx.$$"; git read-tree HEAD; git add -A 2>/dev/null; git write-tree; rm -f "$GIT_INDEX_FILE")"
  printf 'red %s 2026-06-30T00:00:00Z\n' "$tree" > "$(gci rev-parse --absolute-git-dir)/bench-last-gate"
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
  # Warm pooled worktrees (released, no lease) are expected state, never a signal;
  # a leased pool entry is in-flight work and must fire.
  export BENCH_HOME="$tmp/.bh"
  rr="$(git rev-parse --show-toplevel)"
  pool="$BENCH_HOME/worktrees/$(basename "$rr")-$(echo "$rr" | cksum | cut -d' ' -f1)"
  mkdir -p "$pool"
  gci worktree add -q --detach "$pool/warm" HEAD 2>/dev/null
  out="$(bash "$root/bin/bench.sh" status)"
  if grep -qF 'resume or clean up' <<<"$out"; then echo "warm pooled worktree surfaced as a signal"; exit 1; fi
  : > "$(git -C "$pool/warm" rev-parse --git-path bench-lease)"
  out="$(bash "$root/bin/bench.sh" status)"
  grep -qF '1 active worktree' <<<"$out" || { echo "leased pooled worktree did not surface"; exit 1; }
) || err "bench status warm-pool contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench bin
  cp "$root"/bin/bench.sh "$root"/bin/bench-link.sh "$root"/bin/bench-status.sh bin/
  chmod +x bin/bench.sh
  printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  gci add -A; gci commit -q -m init
  printf '{}\n' | BENCH_SHIFT=1 bash "$root/.bench/hooks/stop.sh" >/dev/null 2>&1 || true
  cache="$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  [ -f "$cache" ] || { echo "Stop hook did not write the gate cache"; exit 1; }
  grep -qE '^(green|red) [0-9a-f]+ [0-9T:Z-]+$' "$cache" || { echo "gate cache not <status> <tree> <iso8601>"; exit 1; }
) || err "bench status gate-cache write contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench
  printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  gci add -A; gci commit -q -m init
  cache="$(gci rev-parse --absolute-git-dir)/bench-last-gate"
  printf 'green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n' > "$cache"
  bash "$root/bin/bench.sh" gate >/dev/null 2>&1 || { echo "manual green gate run exited nonzero"; exit 1; }
  grep -qF "green $(gci rev-parse 'HEAD^{tree}')" "$cache" || { echo "manual gate run did not record green keyed to the tested tree"; exit 1; }
  out="$(bash "$root/bin/bench.sh" status)"
  if grep -qF 're-run the gate' <<<"$out"; then echo "status still reads stale after a manual gate run"; exit 1; fi
  # Commit-on-green survival: a commit that does not change the tree must not
  # stale the verdict that authorized it.
  gci commit -q --allow-empty -m same-tree
  out="$(bash "$root/bin/bench.sh" status)"
  if grep -qF 're-run the gate' <<<"$out"; then echo "same-tree commit staled a fresh green verdict"; exit 1; fi
  printf '#!/usr/bin/env bash\nexit 1\n' > .bench/gate.sh
  if bash "$root/bin/bench.sh" gate >/dev/null 2>&1; then echo "red gate run exited zero"; exit 1; fi
  # The throwaway index must live outside the repo, or it would join the tree it hashes.
  tree="$(export GIT_INDEX_FILE="${TMPDIR:-/tmp}/bench-vr-idx.$$"; git read-tree HEAD; git add -A 2>/dev/null; git write-tree; rm -f "$GIT_INDEX_FILE")"
  grep -qF "red $tree" "$cache" || { echo "manual red gate run did not record red keyed to the dirty tree"; exit 1; }
) || err "bench gate verdict-record contract failed"
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
  seq 401 | sed 's/^/x=/g' > big.sh
  gci add big.sh
  out="$(bash "$root/bin/bench.sh" structure 2>&1)" && { printf '%s\n' "$out"; echo "shell source over the line budget did not fail structure"; exit 1; }
  grep -qF 'FILE TOO LONG' <<<"$out" || { echo "structure did not report shell file length"; exit 1; }
) || err "bench structure shell-file contract failed"
rm -rf "$tmp"

# Per-path budgets: .bench/structure.budgets grants (or tightens) named paths;
# the value replaces the global cap. Malformed lines warn and fall through to
# the global cap; the file's last line parses without a trailing newline.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench sub
  seq 401 | sed 's/^/x=/' > big.sh
  seq 200 | sed 's/^/y=/' > mid.sh
  for i in $(seq 13); do printf 'z=1\n' > "sub/f$i.sh"; done
  printf '# reviewer grants\nbig.sh 500\nsub/ 20\nweird abc\nmid.sh 100' > .bench/structure.budgets
  gci add -A; gci commit -q -m s
  if out="$(bash "$root/bin/bench.sh" structure 2>&1)"; then
    echo "override below the global cap did not fail structure"; exit 1
  fi
  grep -qF 'ignoring malformed line' <<<"$out" || { echo "malformed budgets line was not warned about"; exit 1; }
  if grep -qF 'big.sh' <<<"$out"; then echo "granted file budget was not applied"; exit 1; fi
  if grep -qF 'DIR CROWDED' <<<"$out"; then echo "granted dir budget was not applied"; exit 1; fi
  grep -qF '200 lines (max 100)   mid.sh' <<<"$out" || { echo "tightening override (on an unterminated last line) was not applied"; exit 1; }
) || err "bench structure budgets contract failed"
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

# Lease hardening: claims are atomic and owned. A dead-pid or aged-out lease is
# reclaimed; a live foreign lease is respected and survives a foreign release;
# a fresh empty lease (writer mid-claim) is not stolen.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  printf '#!/usr/bin/env bash\n: "${BENCH_WT_RECORD:?}"\npwd >> "$BENCH_WT_RECORD"\n' > rec-shell
  chmod +x rec-shell
  record="$tmp/paths"
  run_wt() { BENCH_HOME="$tmp/.bh" BENCH_WT_RECORD="$record" SHELL="$tmp/rec-shell" bash "$root/bin/bench.sh" worktree >/dev/null 2>&1; }
  run_wt
  P="$(head -1 "$record")"
  lease="$(git -C "$P" rev-parse --git-path bench-lease)"
  printf '4194300 2020-01-01T00:00:00Z\n' > "$lease"
  run_wt   # dead pid → reclaim P
  printf '%s 2026-01-01T00:00:00Z\n' "$$" > "$lease"
  run_wt   # live foreign pid → leave P alone
  [ -f "$lease" ] || { echo "live foreign lease was removed by a foreign release"; exit 1; }
  : > "$lease"; touch -t 202001010000 "$lease"
  run_wt   # aged-out empty (legacy/crash) lease → reclaim P
  : > "$lease"
  run_wt   # fresh empty lease (writer mid-claim) → leave P alone
  rm -f "$lease"
  mapfile -t paths < "$record"
  [ "${#paths[@]}" = "5" ] || { echo "expected five worktree runs, got ${#paths[@]}"; exit 1; }
  [ "${paths[1]}" = "$P" ] || { echo "dead-pid lease was not reclaimed"; exit 1; }
  [ "${paths[2]}" != "$P" ] || { echo "live foreign lease was stolen"; exit 1; }
  [ "${paths[3]}" = "$P" ] || { echo "aged-out empty lease was not reclaimed"; exit 1; }
  [ "${paths[4]}" != "$P" ] || { echo "fresh empty lease was stolen"; exit 1; }
) || err "bench worktree lease hardening contract failed"
rm -rf "$tmp"

# Concurrent acquires never share a worktree. Rendezvous shell: each holds its
# lease until both invocations have recorded, so the two acquires overlap.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  cat > rv-shell <<'EOF'
#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}"
pwd >> "$BENCH_WT_RECORD"
for _ in $(seq 100); do
  [ "$(grep -c . "$BENCH_WT_RECORD" 2>/dev/null)" -ge 2 ] && exit 0
  sleep 0.1
done
exit 0
EOF
  chmod +x rv-shell
  record="$tmp/paths"
  BENCH_HOME="$tmp/.bh" BENCH_WT_RECORD="$record" SHELL="$tmp/rv-shell" bash "$root/bin/bench.sh" worktree >/dev/null 2>&1 &
  BENCH_HOME="$tmp/.bh" BENCH_WT_RECORD="$record" SHELL="$tmp/rv-shell" bash "$root/bin/bench.sh" worktree >/dev/null 2>&1 &
  wait
  mapfile -t paths < "$record"
  [ "${#paths[@]}" = "2" ] || { echo "concurrent worktree runs did not both complete"; exit 1; }
  [ "${paths[0]}" != "${paths[1]}" ] || { echo "concurrent acquires shared a worktree"; exit 1; }
) || err "bench worktree concurrent-acquire contract failed"
rm -rf "$tmp"

# shellcheck source=/dev/null
. "$gate_dir/gate-runtime-git-contracts.sh"

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

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench bin
  cp "$root"/bin/bench.sh "$root"/bin/bench-link.sh "$root"/bin/bench-status.sh bin/
  chmod +x bin/bench.sh
  printf '#!/usr/bin/env bash\nexit 1\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  gci add -A; gci commit -q -m init
  printf '{"stop_hook_active":true}\n' | BENCH_SHIFT=1 bash "$root/.bench/hooks/stop.sh" >/dev/null 2>&1; rc=$?
  [ "$rc" = "0" ] || { echo "stop hook ignored stop_hook_active (exit $rc)"; exit 1; }
  printf '{}\n' | BENCH_SHIFT=1 bash "$root/.bench/hooks/stop.sh" >/dev/null 2>&1; rc=$?
  [ "$rc" = "2" ] || { echo "armed red-gate stop not blocked without the flag (exit $rc)"; exit 1; }
) || err "stop hook stop_hook_active contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q; gci commit -q --allow-empty -m init
  out="$(printf '{}\n' | BENCH_SHIFT=1 PATH=/usr/bin:/bin bash "$root/.bench/hooks/stop.sh" 2>&1)"; rc=$?
  [ "$rc" = "0" ] || { echo "missing bench trapped the stop (exit $rc)"; exit 1; }
  grep -qi 'bench' <<<"$out" || { echo "missing-bench stop gave no warning"; exit 1; }
  [ ! -f "$(gci rev-parse --absolute-git-dir)/bench-last-gate" ] || { echo "missing bench forged a gate cache"; exit 1; }
) || err "stop hook missing-bench fail-open contract failed"
rm -rf "$tmp"
