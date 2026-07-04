# Line-routing contracts for the benchkit gate: the tier binding, the Agent-tool
# hook, and the adapter guards that enforce invariant #2's declared line. Checks
# are guarded on the presence of the surface they test, keeping canary fixtures
# attributable; messages are distinct per failure mode for the same reason.
#
# The tier-binding parse now lives in the Go core (internal/lines), reached through
# `bench check-agent-line`/`bench resolve-model`; the shell parser (lib/lines-env.sh)
# and its hostile-input pin (`a0`) retired into internal/lines table tests. The gate
# keeps its own shape judgment (the regexes below); it reads the three tier values
# through the binary's `--describe-binding` (the enforcement's own view, so a Go
# parser drift surfaces here) and reads the optional aliases directly from the
# known-good lines.env.

# The tier binding is read through the REAL kit's freshly-built core (the Go layer
# rebuilds dist/bench before these fragments run), not $root's — so a canary fixture
# that plants a binding/hook regression without shipping a CLI is still graded through
# the enforcement's own view. The binary reads the CWD repo's lines.env, which is the
# fixture under grade during the canary sweep and benchkit itself in the normal run.
_realbench="$(cd "$gate_dir/.." 2>/dev/null && pwd)/bin/bench.sh"

# a) The binding. benchkit is a routed repo: .bench/lines.env must exist, and a
#    file that exists must carry three model-id-shaped tier values — a binding
#    that drifts to empty silently disarms both enforcement surfaces. The tier
#    values are read through `bench check-agent-line --describe-binding`; the shape
#    checks and prose-sync stay the gate's own judgment.
if [ ! -f .bench/lines.env ]; then
  err "lines.env missing: .bench/lines.env is the tier binding enforcement reads"
elif [ ! -f "$_realbench" ]; then
  : # no core to read the binding through (minimal fixture) — skip shape checks
else
  _db="$(bash "$_realbench" check-agent-line --describe-binding 2>/dev/null)"
  # describe-binding (routed) is: "Agent delegation off the bound line (top=X mid=Y cheap=Z)"
  # with `-` for an empty tier; extract each value from the enforcement's own view.
  _lv_db() { printf '%s' "$_db" | sed -n "s/.*$1=\\([^ )]*\\).*/\\1/p"; }
  # Aliases are optional and not part of the denies clause; read them straight from
  # the known-good binding file (a bare `grep|sed`, not the retired hostile-input parser).
  _lv_file() { grep -E "^[[:space:]]*$1=" .bench/lines.env 2>/dev/null | tail -n1 | sed 's/^[[:space:]]*[^=]*=//'; }
  for _lv in top mid cheap; do
    _val="$(_lv_db "$_lv")"
    if [ -z "$_val" ] || [ "$_val" = "-" ]; then
      err "lines.env tier unset: BENCH_TIER_$(printf '%s' "$_lv" | tr '[:lower:]' '[:upper:]') has no value in .bench/lines.env"
    elif ! printf '%s' "$_val" | grep -qE '^claude-[a-z0-9][a-z0-9.-]*$'; then
      err "lines.env tier malformed: $_lv='$_val' is not a model id"
    fi
  done
  # Alias declarations are optional, but a present one must be a bare alias — a
  # trailing inline comment or stray token silently re-bricks in-session
  # delegation (the hook would compare against 'opus  # mid' and deny 'opus').
  for _la in BENCH_ALIAS_TOP BENCH_ALIAS_MID BENCH_ALIAS_CHEAP; do
    grep -qE "^[[:space:]]*${_la}=" .bench/lines.env || continue
    _val="$(_lv_file "$_la")"
    printf '%s' "$_val" | grep -qE '^[a-z0-9-]+$' \
      || err "lines.env alias malformed: $_la='$_val' is not a bare alias"
  done
  # The profile's Lines prose restates this binding for cold sessions and asks to
  # be kept in sync; this is the check behind that promise — a stale prose copy
  # makes a session declare the line from one binding while the hooks enforce
  # another. Guarded on the profile existing so minimal fixtures skip it.
  if [ -f projects/benchkit.md ]; then
    for _lv in top mid cheap; do
      _val="$(_lv_db "$_lv")"
      { [ -z "$_val" ] || [ "$_val" = "-" ]; } && continue   # unset/malformed already reported above
      grep -qF "$_val" projects/benchkit.md \
        || err "profile Lines prose stale: projects/benchkit.md does not name bound model id '$_val' (BENCH_TIER_$(printf '%s' "$_lv" | tr '[:lower:]' '[:upper:]') in lines.env)"
    done
    for _la in BENCH_ALIAS_TOP BENCH_ALIAS_MID BENCH_ALIAS_CHEAP; do
      grep -qE "^[[:space:]]*${_la}=" .bench/lines.env || continue
      _val="$(_lv_file "$_la")"
      [ -n "$_val" ] || continue
      grep -qF "${_la}=${_val}" projects/benchkit.md \
        || err "profile Lines prose stale: projects/benchkit.md does not carry alias declaration ${_la}=${_val}"
    done
  fi
fi

# b) Hook wiring. A hook script without its settings.json matcher is decorative.
if [ -f .claude/settings.json ]; then
  node -e '
const fs = require("fs");
let cfg;
try { cfg = JSON.parse(fs.readFileSync(".claude/settings.json", "utf8")); } catch (e) { process.exit(0); } // invalid JSON is check 2 to report
const pre = (cfg.hooks && cfg.hooks.PreToolUse) || [];
const entry = pre.find(e => e.matcher === "Agent");
if (!entry || !(entry.hooks || []).some(h => (h.command || "").includes("check-agent-line.sh"))) process.exit(1);
' || err "claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh"
fi

# c) Hook behavior, exercised with fixture stdin against a controlled temp repo so
#    the check is deterministic regardless of this repo's own binding. The shim
#    resolves the wrapper and pipes the envelope to `bench check-agent-line`, which
#    owns the verdict; a `bench` on PATH (pointing at the real wrapper) makes the
#    core reachable from the fixture repos, whose git root the binary reads for the
#    binding. Assertions are the same allow/deny/degraded ones as before the port.
if [ -f .bench/hooks/check-agent-line.sh ] && [ -f "$_realbench" ]; then
  _hook="$root/.bench/hooks/check-agent-line.sh"
  _bindir="$(mktemp -d)"
  printf '#!/bin/sh\nexec %s "$@"\n' "$_realbench" >"$_bindir/bench"
  chmod +x "$_bindir/bench"
  _run_hook() { PATH="$_bindir:$PATH" bash "$_hook"; }   # cwd is set by each caller subshell

  _hd="$(mktemp -d)"
  ( cd "$_hd" && git init -q )
  mkdir -p "$_hd/.bench"
  printf 'BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\nBENCH_ALIAS_MID=opus\n' >"$_hd/.bench/lines.env"

  printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-opus-4-8"}}' \
    | ( cd "$_hd" && _run_hook ) >/dev/null 2>&1 \
    || err "check-agent-line.sh denies a bound model (allow case broken)"

  printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus"}}' \
    | ( cd "$_hd" && _run_hook ) >/dev/null 2>&1 \
    || err "check-agent-line.sh denies a declared alias (Agent tool speaks aliases)"

  if printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","model":"sonnet"}}' \
    | ( cd "$_hd" && _run_hook ) >/dev/null 2>&1; then
    err "check-agent-line.sh does not deny an undeclared alias"
  fi

  if printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}' \
    | ( cd "$_hd" && _run_hook ) >/dev/null 2>&1; then
    err "check-agent-line.sh does not deny an unbound model"
  fi

  printf '%s' 'not json at all' \
    | ( cd "$_hd" && _run_hook ) >/dev/null 2>&1 \
    || err "check-agent-line.sh does not fail open on malformed stdin"

  printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x"}}' \
    | ( cd "$_hd" && _run_hook ) >/dev/null 2>&1 \
    || err "check-agent-line.sh does not fail open on a missing model field"

  _hd2="$(mktemp -d)"
  ( cd "$_hd2" && git init -q )
  _werr="$( printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}' \
    | ( cd "$_hd2" && _run_hook ) 2>&1 >/dev/null )" \
    || err "check-agent-line.sh does not fail open without lines.env"
  printf '%s' "$_werr" | grep -qF 'no .bench/lines.env' \
    || err "check-agent-line.sh does not warn on stderr when lines.env is absent"

  # A hook copied without the shared wrapper resolver (.bench/lib/resolve-bench.sh)
  # must fail open like every other broken-hook case — a broken hook never bricks
  # delegation — and say why on stderr. (The re-pointed slice-4 failure mode: a
  # missing lib must not let the shim error before its fail-open rim.)
  _hd3="$(mktemp -d)"
  cp "$_hook" "$_hd3/check-agent-line.sh"
  _werr="$( printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}' \
    | ( cd "$_hd" && PATH="$_bindir:$PATH" bash "$_hd3/check-agent-line.sh" ) 2>&1 >/dev/null )" \
    || err "check-agent-line.sh does not fail open when the wrapper resolver lib is missing"
  printf '%s' "$_werr" | grep -qF 'wrapper resolver missing' \
    || err "check-agent-line.sh does not warn on stderr when the wrapper resolver lib is missing"

  # An incomplete binding (a tier key present but empty) is a partial oracle:
  # the hook must fail open with the incomplete-binding warning, never deny
  # against half a binding.
  _hd4="$(mktemp -d)"
  ( cd "$_hd4" && git init -q )
  mkdir -p "$_hd4/.bench"
  printf 'BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n' >"$_hd4/.bench/lines.env"
  _werr="$( printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}' \
    | ( cd "$_hd4" && _run_hook ) 2>&1 >/dev/null )" \
    || err "check-agent-line.sh does not fail open on an incomplete binding"
  printf '%s' "$_werr" | grep -qF 'unset or empty' \
    || err "check-agent-line.sh does not warn on stderr about an incomplete binding"
  rm -rf "$_hd" "$_hd2" "$_hd3" "$_hd4" "$_bindir"
fi

# d) Adapter guards, exercised against a stub harness on PATH — a routed repo must
#    refuse an undeclared or unbound BENCH_MODEL before the harness runs; an
#    unrouted repo must stay a plain pass-through. The adapters exec `bench
#    resolve-model` (the Go verdict), so a `bench` on PATH (real wrapper) makes the
#    core reachable; the resolution reads the cwd repo's binding.
if [ -d .bench/adapters ]; then
  _sd="$(mktemp -d)"
  for _stub in claude codex opencode; do
    printf '#!/bin/sh\nprintf '"'"'%%s\\n'"'"' "$@"\n' >"$_sd/$_stub"
    chmod +x "$_sd/$_stub"
  done
  # A `bench` on PATH that routes to the real kit's core, so the adapter's `bench
  # resolve-model` reaches the Go verdict. A broken fixture adapter that ignores the
  # binding never calls it.
  if [ -f "$_realbench" ]; then
    printf '#!/bin/sh\nexec %s "$@"\n' "$_realbench" >"$_sd/bench"
    chmod +x "$_sd/bench"
  fi
  _routed="$(mktemp -d)"; ( cd "$_routed" && git init -q )
  mkdir -p "$_routed/.bench"
  printf 'BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n' >"$_routed/.bench/lines.env"
  _unrouted="$(mktemp -d)"; ( cd "$_unrouted" && git init -q )

  for _a in claude codex opencode; do
    [ -f ".bench/adapters/$_a" ] || { err "adapter missing from .bench/adapters: $_a"; continue; }
    _ap="$root/.bench/adapters/$_a"

    # Dash-leading probe: the prompt must survive as a positional even when it
    # looks like a flag (`--` sentinel), alongside the bound model.
    _out="$( cd "$_routed" && BENCH_MODEL=claude-opus-4-8 PATH="$_sd:$PATH" bash "$_ap" "--line probe prompt" 2>/dev/null )"
    { printf '%s' "$_out" | grep -qF 'claude-opus-4-8' && printf '%s' "$_out" | grep -qF -- '--line probe prompt'; } \
      || err "adapter $_a does not pass BENCH_MODEL and a dash-leading prompt to the harness in a routed repo"

    if ( cd "$_routed" && env -u BENCH_MODEL PATH="$_sd:$PATH" bash "$_ap" "line probe prompt" ) >/dev/null 2>&1; then
      err "adapter $_a does not refuse undeclared BENCH_MODEL in a routed repo"
    fi

    if ( cd "$_routed" && BENCH_MODEL=claude-nonexistent-9 PATH="$_sd:$PATH" bash "$_ap" "line probe prompt" ) >/dev/null 2>&1; then
      err "adapter $_a does not refuse an unbound BENCH_MODEL in a routed repo"
    fi

    _out="$( cd "$_unrouted" && env -u BENCH_MODEL PATH="$_sd:$PATH" bash "$_ap" "line probe prompt" 2>/dev/null )"
    printf '%s' "$_out" | grep -qF 'line probe prompt' \
      || err "adapter $_a does not pass through in an unrouted repo"

    _out="$( cd "$_unrouted" && BENCH_MODEL=claude-anything-7 PATH="$_sd:$PATH" bash "$_ap" "line probe prompt" 2>/dev/null )"
    { printf '%s' "$_out" | grep -qF 'claude-anything-7' && printf '%s' "$_out" | grep -qF 'line probe prompt'; } \
      || err "adapter $_a does not pass an explicit BENCH_MODEL through in an unrouted repo"
  done

  # An adapter copied WITHOUT the shared wrapper resolver (../lib/resolve-bench.sh)
  # must fail closed — an unguarded passthrough in a routed repo is silent
  # de-enforcement.
  if [ -f .bench/adapters/claude ]; then
    _tmpad="$(mktemp -d)"
    cp .bench/adapters/claude "$_tmpad/claude-adapter"
    if ( cd "$_routed" && BENCH_MODEL=claude-opus-4-8 PATH="$_sd:$PATH" bash "$_tmpad/claude-adapter" "line probe prompt" ) >/dev/null 2>&1; then
      err "adapter claude does not fail closed when resolve-bench.sh is missing"
    fi
    rm -rf "$_tmpad"
  fi

  # An adapter WITH its resolver lib but no reachable bench core must also refuse in
  # a routed repo — it must never launch when it cannot resolve the binding. Copied
  # into a ../lib layout and run with a PATH holding no bench, resolve-bench.sh finds
  # no wrapper, so the adapter fails closed.
  if [ -f .bench/adapters/claude ] && [ -f .bench/lib/resolve-bench.sh ]; then
    _tmpad="$(mktemp -d)"
    mkdir -p "$_tmpad/adapters" "$_tmpad/lib"
    cp .bench/adapters/claude "$_tmpad/adapters/claude"
    cp .bench/lib/resolve-bench.sh "$_tmpad/lib/resolve-bench.sh"
    if ( cd "$_routed" && BENCH_MODEL=claude-opus-4-8 PATH=/usr/bin:/bin bash "$_tmpad/adapters/claude" "line probe prompt" ) >/dev/null 2>&1; then
      err "adapter claude does not fail closed when the bench core is unreachable"
    fi
    rm -rf "$_tmpad"
  fi

  # An incomplete binding (a tier key present but empty) is a partial oracle:
  # resolve-model must warn and fall back to the unrouted passthrough — it
  # neither refuses nor enforces half a binding.
  if [ -f .bench/adapters/claude ] && [ -f "$_realbench" ]; then
    _partial="$(mktemp -d)"; ( cd "$_partial" && git init -q )
    mkdir -p "$_partial/.bench"
    printf 'BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n' >"$_partial/.bench/lines.env"
    _out="$( cd "$_partial" && BENCH_MODEL=claude-anything-7 PATH="$_sd:$PATH" bash "$root/.bench/adapters/claude" "line probe prompt" 2>/dev/null )"
    { printf '%s' "$_out" | grep -qF 'claude-anything-7' && printf '%s' "$_out" | grep -qF 'line probe prompt'; } \
      || err "adapter claude does not fall back to passthrough on an incomplete binding"
    _werr="$( cd "$_partial" && BENCH_MODEL=claude-anything-7 PATH="$_sd:$PATH" bash "$root/.bench/adapters/claude" "line probe prompt" 2>&1 >/dev/null )"
    printf '%s' "$_werr" | grep -qF 'unset or empty' \
      || err "adapter claude does not warn about an incomplete binding"
    rm -rf "$_partial"
  fi
  rm -rf "$_sd" "$_routed" "$_unrouted"
fi

# e) Prose anchors: the command that declares the line must route through the
#    rubric, and the shared contract must document the adapter's line surface.
if [ -f .agents/commands/bench-implement-spec.md ]; then
  grep -qF 'craft-line' .agents/commands/bench-implement-spec.md \
    || err "bench-implement-spec.md does not reference craft-line"
fi
if [ -f .agents/commands/bench-write-spec.md ]; then
  grep -qF 'craft-line' .agents/commands/bench-write-spec.md \
    || err "bench-write-spec.md does not reference craft-line"
  grep -qF 'model and effort' .agents/commands/bench-write-spec.md \
    || err "bench-write-spec.md does not mandate per-story model and effort"
fi
# The shift adapter contract (with BENCH_MODEL) moved to the on-demand
# .bench/BENCH-reference.md in the token diet; that is where it must be documented.
if [ -f .bench/BENCH-reference.md ]; then
  grep -qF 'BENCH_MODEL' .bench/BENCH-reference.md \
    || err "BENCH-reference.md adapter contract does not document BENCH_MODEL"
fi
