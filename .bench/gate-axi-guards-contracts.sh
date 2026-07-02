# AXI-conformance contracts for the guard-surface aggregation: `bench guards`
# (and `guards --brief`) plus the session-start hook that injects the brief.
# Split from gate-axi-contracts.sh along its section banners to respect the
# structure budget; guard --describe manifest conformance and the query/parser
# contracts stay there. What this file pins:
#   - `bench guards` emits one TOON row per deny-capable guard, excludes the
#     informational session-start hook, bounds each guard's --describe so a
#     hanging hook cannot stall aggregation (`no manifest (timed out)`), never
#     executes an unmanaged (marker-less) pre-push (`unmanaged (no manifest)`),
#     and block-dangerous-git degrades to `manifest unavailable (python3
#     missing)` when python3 is absent;
#   - session-start injects the guard brief in a linked repo and never blocks
#     (prints nothing, exit 0) outside one.
# Sourced by gate.sh so it shares $root, err(), and $fail with the other
# fragments; fixture provisioning and cleanup are the contract runner's
# (gate-contract-runner.sh). Run the real CLI/hook in a fixture and assert
# stdout shape + exit code — never internals.

[ -f "$root/bin/bench.sh" ] || { err "bench CLI missing (AXI guards contracts skipped)"; return 0 2>/dev/null || exit 0; }

# A linked fixture carries all five deny-capable guards; `bench guards` emits one
# TOON row each, excludes the informational session-start hook, reports a stub that
# ignores --describe as `no manifest`, and an absent pre-push as `not installed`.
contract "AXI guards aggregation contract" <<'BODY'
  bash "$root/bin/bench.sh" link >/dev/null 2>&1 || { echo "link failed setting up guards fixture"; exit 1; }
  out="$(bash "$root/bin/bench.sh" guards)"
  head -1 <<<"$out" | grep -qxF 'guards[5]{guard,boundary,denies}:' || { echo "guards header not count-5 TOON: $(head -1 <<<"$out")"; exit 1; }
  for g in block-dangerous-git check-agent-line stop _line-guard pre-push; do
    grep -qE "^  $g," <<<"$out" || { echo "guards missing deny-capable row: $g"; exit 1; }
  done
  if grep -qE '^  session-start,' <<<"$out"; then echo "informational session-start leaked into guard rows"; exit 1; fi
  # The generated pre-push hook answers --describe with the four keys and
  # exits 0, without reading stdin — so the git-layer guard is part of the surface.
  pp="$(git rev-parse --git-path hooks)/pre-push"
  ppm="$(bash "$pp" --describe </dev/null)"; pprc=$?
  [ "$pprc" = "0" ] || { echo "generated pre-push --describe did not exit 0 (exit $pprc)"; exit 1; }
  for k in name boundary denies why; do
    grep -qE "^$k: " <<<"$ppm" || { echo "generated pre-push --describe manifest missing $k key"; exit 1; }
  done
  # An executable hook that ignores --describe is never skipped silently. Capture
  # first, then grep: `bench guards | grep -q` would SIGPIPE the CLI on an early
  # match and fail under pipefail even though the row is present.
  printf '#!/usr/bin/env bash\ncat >/dev/null\nexit 0\n' > .bench/hooks/extra.sh; chmod +x .bench/hooks/extra.sh
  out="$(bash "$root/bin/bench.sh" guards)"
  grep -qxF '  extra,,no manifest' <<<"$out" || { echo "stub hook not reported as no manifest"; exit 1; }
  rm -f .bench/hooks/extra.sh
  # Absent pre-push is a definitive not-installed row, not an omission.
  rm -f "$(git rev-parse --git-path hooks)/pre-push"
  out="$(bash "$root/bin/bench.sh" guards)"
  grep -qxF '  pre-push,,not installed' <<<"$out" || { echo "absent pre-push not reported as not installed"; exit 1; }
BODY

# guards --brief: one plain line per deny-capable guard plus exactly one footer
# pointing at `bench guards` — the surface session-start injects.
contract "AXI guards --brief contract" <<'BODY'
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  out="$(bash "$root/bin/bench.sh" guards --brief)"
  [ "$(grep -cF 'full manifests: bench guards' <<<"$out")" = "1" ] || { echo "brief footer not exactly one line"; exit 1; }
  [ "$(grep -c . <<<"$out")" = "6" ] || { echo "brief not five guard lines plus one footer (got $(grep -c . <<<"$out"))"; exit 1; }
  grep -qF 'block-dangerous-git: destructive git' <<<"$out" || { echo "brief lost the git guard deny clause"; exit 1; }
BODY

# An unknown argument is a usage error (exit 2); the command resolves the repo
# root from a subdirectory.
contract "AXI guards usage/subdirectory contract" <<'BODY'
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  out="$(bash "$root/bin/bench.sh" guards bogusarg)"; rc=$?
  [ "$rc" = "2" ] || { echo "guards unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "guards unknown arg did not print usage on stdout: $out"; exit 1; }
  mkdir -p sub/deeper
  sout="$(cd sub/deeper && bash "$root/bin/bench.sh" guards)"
  head -1 <<<"$sout" | grep -qF 'guards[' || { echo "guards from a subdirectory lost root resolution: $sout"; exit 1; }
BODY

# A linked fixture under a space-containing path still aggregates.
contract "AXI guards path-with-spaces contract" --space-path <<'BODY'
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  out="$(bash "$root/bin/bench.sh" guards)"
  head -1 <<<"$out" | grep -qF 'guards[' || { echo "space-path guards failed: $out"; exit 1; }
BODY

# A guard whose --describe blocks (here, sleeps well past the bound) must not
# hang aggregation: `bench guards` bounds each --describe and reports the offender
# as `no manifest (timed out)`, returning promptly.
contract "AXI guards --describe timeout-bound contract" <<'BODY'
  mkdir -p .bench/hooks
  printf '#!/usr/bin/env bash\nif [ "${1:-}" = "--describe" ]; then sleep 30; fi\nexit 0\n' > .bench/hooks/slow.sh
  chmod +x .bench/hooks/slow.sh
  start=$(date +%s)
  out="$(bash "$root/bin/bench.sh" guards)"
  elapsed=$(( $(date +%s) - start ))
  [ "$elapsed" -lt 10 ] || { echo "guards did not bound a slow --describe (took ${elapsed}s)"; exit 1; }
  grep -qxF '  slow,,no manifest (timed out)' <<<"$out" || { echo "slow guard not reported as timed out: $out"; exit 1; }
BODY

# An unmanaged (foreign) pre-push must NOT be executed by `bench guards`:
# a marker-less pre-push is reported `unmanaged (no manifest)` and never run
# (proven by a sentinel it would touch on any invocation).
contract "AXI guards unmanaged-pre-push safety contract" <<'BODY'
  sentinel="$tmp/ran-foreign-prepush"
  hooks="$(git rev-parse --git-path hooks)"
  mkdir -p "$hooks"
  printf '#!/usr/bin/env bash\ntouch %q\nexit 1\n' "$sentinel" > "$hooks/pre-push"
  chmod +x "$hooks/pre-push"
  out="$(bash "$root/bin/bench.sh" guards)"
  grep -qxF '  pre-push,,unmanaged (no manifest)' <<<"$out" || { echo "foreign pre-push not reported unmanaged: $out"; exit 1; }
  [ ! -e "$sentinel" ] || { echo "bench guards executed a foreign pre-push"; exit 1; }
BODY

# block-dangerous-git.sh --describe degrades honestly when python3 is
# unreachable: with PATH stripped to a scratch bin holding only bash + coreutils,
# it prints the `manifest unavailable (python3 missing)` denies line and exits 0.
if [ -f "$root/.bench/hooks/block-dangerous-git.sh" ]; then
  contract "AXI block-dangerous-git python3-missing manifest contract" <<'BODY'
    mkdir sbin
    for t in bash cat printf env sed grep head; do
      p="$(command -v "$t" 2>/dev/null || true)"; [ -n "$p" ] && ln -sf "$p" "sbin/$t"
    done
    out="$(PATH="$tmp/sbin" bash "$root/.bench/hooks/block-dangerous-git.sh" --describe </dev/null 2>/dev/null)"; rc=$?
    [ "$rc" = "0" ] || { echo "python3-missing --describe did not exit 0 (exit $rc)"; exit 1; }
    grep -qF 'manifest unavailable (python3 missing)' <<<"$out" || { echo "python3-missing denies line absent: $out"; exit 1; }
BODY
fi

# In a linked repo session-start runs `bench guards --brief` after the dashboard
# and never blocks; PATH excludes any global bench so bench_cmd resolves only the
# repo-local CLI.
contract "session-start guard-brief injection contract" <<'BODY'
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  out="$(PATH=/usr/bin:/bin bash .bench/hooks/session-start.sh 2>&1)"; rc=$?
  [ "$rc" = "0" ] || { echo "session-start blocked in a linked repo (exit $rc)"; exit 1; }
  grep -qF 'full manifests: bench guards' <<<"$out" || { echo "session-start did not inject the guard brief"; exit 1; }
  grep -qE '^bench CLI: .*\.bench/bin/bench\.sh \(bench not on PATH; invoke by path\)$' <<<"$out" \
    || { echo "session-start did not advertise the resolved CLI location: $out"; exit 1; }
BODY

# Outside a repo session-start prints nothing and exits 0 — it must never block
# a session over missing context.
contract "session-start never-blocks-outside-repo contract" --no-repo <<'BODY'
  out="$(PATH=/usr/bin:/bin bash "$root/.bench/hooks/session-start.sh" 2>&1)"; rc=$?
  [ "$rc" = "0" ] || { echo "session-start blocked outside a repo (exit $rc)"; exit 1; }
  [ -z "$out" ] || { echo "session-start printed outside a repo: $out"; exit 1; }
BODY
