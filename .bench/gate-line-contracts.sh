# Line-routing contracts for the benchkit gate: the tier binding, the Agent-tool
# hook, and the adapter guards that enforce invariant #2's declared line. Checks
# are guarded on the presence of the surface they test, keeping canary fixtures
# attributable; messages are distinct per failure mode for the same reason.

# a) The binding. benchkit is a routed repo: .bench/lines.env must exist, and a
#    file that exists must carry three model-id-shaped tier values — a binding
#    that drifts to empty silently disarms both enforcement surfaces.
if [ ! -f .bench/lines.env ]; then
  err "lines.env missing: .bench/lines.env is the tier binding enforcement reads"
else
  for _lv in BENCH_TIER_TOP BENCH_TIER_MID BENCH_TIER_CHEAP; do
    _val="$(grep -E "^${_lv}=" .bench/lines.env | tail -1 | cut -d= -f2- | tr -d '\r' | sed -E 's/[[:space:]]+$//')"
    if [ -z "$_val" ]; then
      err "lines.env tier unset: $_lv has no value in .bench/lines.env"
    elif ! printf '%s' "$_val" | grep -qE '^claude-[a-z0-9][a-z0-9.-]*$'; then
      err "lines.env tier malformed: $_lv='$_val' is not a model id"
    fi
  done
  # Alias declarations are optional, but a present one must be a bare alias — a
  # trailing inline comment or stray token silently re-bricks in-session
  # delegation (the hook would compare against 'opus  # mid' and deny 'opus').
  for _la in BENCH_ALIAS_TOP BENCH_ALIAS_MID BENCH_ALIAS_CHEAP; do
    grep -qE "^${_la}=" .bench/lines.env || continue
    _val="$(grep -E "^${_la}=" .bench/lines.env | tail -1 | cut -d= -f2- | tr -d '\r' | sed -E 's/[[:space:]]+$//')"
    printf '%s' "$_val" | grep -qE '^[a-z0-9-]+$' \
      || err "lines.env alias malformed: $_la='$_val' is not a bare alias"
  done
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

# c) Hook behavior, exercised with fixture stdin against a controlled temp repo
#    so the check is deterministic regardless of this repo's own binding.
if [ -f .bench/hooks/check-agent-line.sh ]; then
  _hook="$root/.bench/hooks/check-agent-line.sh"
  _hd="$(mktemp -d)"
  ( cd "$_hd" && git init -q )
  mkdir -p "$_hd/.bench"
  printf 'BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\nBENCH_ALIAS_MID=opus\n' >"$_hd/.bench/lines.env"

  printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-opus-4-8"}}' \
    | ( cd "$_hd" && bash "$_hook" ) >/dev/null 2>&1 \
    || err "check-agent-line.sh denies a bound model (allow case broken)"

  printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus"}}' \
    | ( cd "$_hd" && bash "$_hook" ) >/dev/null 2>&1 \
    || err "check-agent-line.sh denies a declared alias (Agent tool speaks aliases)"

  if printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","model":"sonnet"}}' \
    | ( cd "$_hd" && bash "$_hook" ) >/dev/null 2>&1; then
    err "check-agent-line.sh does not deny an undeclared alias"
  fi

  if printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}' \
    | ( cd "$_hd" && bash "$_hook" ) >/dev/null 2>&1; then
    err "check-agent-line.sh does not deny an unbound model"
  fi

  printf '%s' 'not json at all' \
    | ( cd "$_hd" && bash "$_hook" ) >/dev/null 2>&1 \
    || err "check-agent-line.sh does not fail open on malformed stdin"

  printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x"}}' \
    | ( cd "$_hd" && bash "$_hook" ) >/dev/null 2>&1 \
    || err "check-agent-line.sh does not fail open on a missing model field"

  _hd2="$(mktemp -d)"
  ( cd "$_hd2" && git init -q )
  _werr="$( printf '%s' '{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}' \
    | ( cd "$_hd2" && bash "$_hook" ) 2>&1 >/dev/null )" \
    || err "check-agent-line.sh does not fail open without lines.env"
  printf '%s' "$_werr" | grep -qF 'no .bench/lines.env' \
    || err "check-agent-line.sh does not warn on stderr when lines.env is absent"
  rm -rf "$_hd" "$_hd2"
fi

# d) Adapter guards, exercised against a stub harness on PATH — a routed repo
#    must refuse an undeclared or unbound BENCH_MODEL before the harness runs;
#    an unrouted repo must stay a plain pass-through.
if [ -d .bench/adapters ]; then
  _sd="$(mktemp -d)"
  for _stub in claude codex opencode; do
    printf '#!/bin/sh\nprintf '"'"'%%s\\n'"'"' "$@"\n' >"$_sd/$_stub"
    chmod +x "$_sd/$_stub"
  done
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

  # A present guard-source that resolves to a missing file must fail closed —
  # an unguarded passthrough in a routed repo is silent de-enforcement.
  if [ -f .bench/adapters/claude ]; then
    _tmpad="$(mktemp -d)"
    cp .bench/adapters/claude "$_tmpad/claude-adapter"
    if ( cd "$_routed" && BENCH_MODEL=claude-opus-4-8 PATH="$_sd:$PATH" bash "$_tmpad/claude-adapter" "line probe prompt" ) >/dev/null 2>&1; then
      err "adapter claude does not fail closed when _line-guard.sh is missing"
    fi
    rm -rf "$_tmpad"
  fi
  rm -rf "$_sd" "$_routed" "$_unrouted"
fi

# e) Prose anchors: the command that declares the line must route through the
#    rubric, and the shared contract must document the adapter's line surface.
if [ -f .agents/commands/bench-implement-spec.md ]; then
  grep -qF 'craft-line' .agents/commands/bench-implement-spec.md \
    || err "bench-implement-spec.md does not reference craft-line"
fi
if [ -f .bench/BENCH.md ]; then
  grep -qF 'BENCH_MODEL' .bench/BENCH.md \
    || err "BENCH.md adapter contract does not document BENCH_MODEL"
fi
