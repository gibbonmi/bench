#!/usr/bin/env bash
# bench — the operational substrate for the Bench workflow.
# Fuses a warm worktree pool (treehouse-lite) with a gated loop (gnhf-lite),
# where the gate is the external oracle: a shift only commits on green.
#
# The worktree lifecycle, the gated loop, and the gate resolution/record all live in
# the Go core now (internal/worktree, internal/shift, internal/gate); this wrapper is
# routing plus a one-glance run_gate adapter. Gate config resolution, owned by
# internal/gate, is: (1) ./.bench/gate.sh, (2) $BENCH_GATE, (3) auto-detect.
set -euo pipefail

# Exported so the Go core's worktree/shift/gate commands resolve the same pool home the
# shell did (they read BENCH_HOME from the environment).
export BENCH_HOME="${BENCH_HOME:-$HOME/.bench}"

# ---- gate: the oracle -------------------------------------------------------
# run_gate — the one-glance adapter over the Go core's `bench gate-run`. Resolution of
# the gate (.bench/gate.sh → $BENCH_GATE → auto-detect), the run from the repo root, and
# the verdict-cache record all live in internal/gate, so `bench gate` and the Stop hook's
# `<wrapper> gate` path share exactly one resolver — never a second live implementation.
# exec, not run: the gate case is terminal and the binary owns the exit code and the
# record; a missing binary exits 127 via route_binary, which writes no forged verdict.
run_gate() { route_porcelain gate-run; }

gate_command() {
  case "${2:-}" in
    "") run_gate ;;
    pin) route_porcelain gate-pin "${@:3}" ;;
    *) run_gate ;;
  esac
}

# Resolve where the canonical kit lives (parent of this script's bin/), following
# symlinks without relying on GNU-only `readlink -f`.
resolve_script_path() {
  local source="${BASH_SOURCE[0]:-$0}" dir target
  while [[ -L "$source" ]]; do
    dir="$(cd -P "$(dirname "$source")" >/dev/null 2>&1 && pwd)"
    target="$(readlink "$source")"
    [[ "$target" == /* ]] && source="$target" || source="$dir/$target"
  done
  dir="$(cd -P "$(dirname "$source")" >/dev/null 2>&1 && pwd)"
  printf '%s/%s\n' "$dir" "$(basename "$source")"
}

kit_dir() {
  local script
  script="$(resolve_script_path)"
  cd "$(dirname "$script")/.." && pwd
}

# ---- strangler router: send a ported subcommand to the Go binary ------------
# The one seam that grows across the port: later slices add subcommand names to the
# dispatch, never a second resolver. platform_pkg maps this host to its
# @redbench/<os>-<arch> package (npm os/cpu spelling); an off-matrix host returns
# non-zero so the caller can emit the "unsupported platform" error instead of naming
# a package that does not exist.
platform_pkg() {
  local os arch
  case "$(uname -s)" in
    Darwin) os=darwin ;;
    Linux)  os=linux ;;
    *) return 1 ;;
  esac
  case "$(uname -m)" in
    arm64|aarch64) arch=arm64 ;;
    x86_64|amd64)  arch=x64 ;;
    *) return 1 ;;
  esac
  printf '@redbench/%s-%s' "$os" "$arch"
}

platform_suffix() {
  local pkg
  pkg="$(platform_pkg)" || return 1
  printf '%s\n' "${pkg#@redbench/}"
}

package_version() {
  local kit="${1:-$(kit_dir)}" pkg_json line
  pkg_json="$kit/package.json"
  [[ -f "$pkg_json" ]] || return 1
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" =~ \"version\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
      printf '%s\n' "${BASH_REMATCH[1]}"
      return 0
    fi
  done < "$pkg_json"
  return 1
}

# main_tree_kit <kit> — when <kit> sits inside a linked git worktree, echo the same
# kit path re-anchored under the main worktree's root; echo nothing when <kit> is the
# main tree itself, outside any repo, or the mapping is degenerate. Linked worktrees
# carry the tracked tree but not untracked artifacts (dist/, node_modules/), so the
# binary that serves a worktree lives in the main tree. Always returns 0: a failed
# resolution degrades to "no extra candidates", never to a caller-visible error.
main_tree_kit() {
  local kit="$1" common wt_root main
  common="$(git -C "$kit" rev-parse --git-common-dir 2>/dev/null)" || return 0
  [[ -n "$common" ]] || return 0
  [[ "$common" == /* ]] || common="$kit/$common"
  wt_root="$(git -C "$kit" rev-parse --show-toplevel 2>/dev/null)" || return 0
  [[ -n "$wt_root" ]] || return 0
  main="$(dirname "$common")${kit#"$wt_root"}"
  [[ "$main" != "$kit" ]] && printf '%s\n' "$main"
  return 0
}

# bench_binary_path — echo the resolved Go binary path, or fail with a distinct exit
# code the caller maps to a message. Resolution order: (1) repo-local dev build (kit
# checkout), (2) the platform package bundled under the wrapper's node_modules, (3) the
# hoisted sibling npm produces for global installs — tried first against the kit dir
# itself, then re-anchored at the main tree when the kit dir is a linked git worktree
# (worktrees carry the tracked wrapper but not the untracked binary). First executable,
# non-empty match wins; a present-but-empty or non-executable file is treated as missing
# (never named), so a torn build falls through rather than resolving to a non-runnable
# path. Exit 2 = off-matrix platform (no package to name); 127 = no binary present for
# this platform. One source of both the platform→package mapping and the resolution
# order, shared by route_binary (which execs the path) and the status adapters (which
# capture its output).
bench_binary_path() {
  local kit="${1:-$(kit_dir)}" pkg suffix version cache c k main
  platform_pkg >/dev/null || return 2
  pkg="$(platform_pkg)"
  suffix="$(platform_suffix)"
  main="$(main_tree_kit "$kit")"
  for k in "$kit" ${main:+"$main"}; do
    for c in "$k/dist/bench" "$k/node_modules/$pkg/bin/bench" "$k/../$pkg/bin/bench"; do
      [[ -x "$c" && -s "$c" ]] && { printf '%s\n' "$c"; return 0; }
    done
  done
  version="$(package_version "$kit" 2>/dev/null || true)"
  if [[ -n "$version" ]]; then
    cache="$BENCH_HOME/cache/bin/$version/$suffix/bench"
    [[ -x "$cache" && -s "$cache" ]] && { printf '%s\n' "$cache"; return 0; }
  fi
  return 127
}

repair_binary() {
  local kit="$1" wrapper="$2" script version suffix
  script="$(dirname "$wrapper")/bench-repair-binary.mjs"
  if ! command -v node >/dev/null 2>&1; then
    echo "bench: repair skipped because node is not on PATH" >&2
    return 1
  fi
  [[ -f "$script" ]] || return 1
  version="$(package_version "$kit" 2>/dev/null || true)"
  suffix="$(platform_suffix 2>/dev/null || true)"
  [[ -n "$version" && -n "$suffix" ]] || return 1
  node "$script" "$kit" "$(platform_pkg)" "$version" "$suffix"
}

# route_binary <subcommand> [args...] — resolve and exec the Go binary, passing the
# whole argv through. The one seam the strangler grows: later slices add subcommand
# names to the dispatch below, never a second resolver.
route_binary() {
  local allow_repair=0 bin rc kit wrapper
  if [[ "${1:-}" == "--repair" ]]; then
    allow_repair=1
    shift
  fi
  kit="${BENCH_KIT:-$(kit_dir)}"
  wrapper="$(resolve_script_path)"
  bin="$(bench_binary_path "$kit")" && rc=0 || rc=$?
  case "$rc" in
    0) BENCH_KIT="${BENCH_KIT:-$kit}" BENCH_WRAPPER="${BENCH_WRAPPER:-$wrapper}" exec "$bin" "$@" ;;
    2)
      echo "bench: unsupported platform: $(uname -s | tr '[:upper:]' '[:lower:]')/$(uname -m)" >&2
      exit 2
      ;;
    *)
      if [[ "$allow_repair" == 1 ]]; then
        repair_binary "$kit" "$wrapper" || true
        bin="$(bench_binary_path "$kit")" && BENCH_KIT="${BENCH_KIT:-$kit}" BENCH_WRAPPER="${BENCH_WRAPPER:-$wrapper}" exec "$bin" "$@"
      fi
      echo 'bench: no binary for this platform — build it from a clone: bash scripts/go-build.sh "$PWD" dist/bench' >&2
      exit 127
      ;;
  esac
}

route_porcelain() {
  route_binary --repair "$@"
}

adoption_route() {
  local kit
  kit="${BENCH_KIT:-$(kit_dir)}"
  if [[ ! -d "$kit/.agents/commands" || ! -f "$kit/AGENTS.md" ]]; then
    echo "bench: link/init/doctor must run from the real Bench kit; no source asset tree at $kit" >&2
    echo "bench: use the installed 'bench' command or the source kit's bin/bench.sh" >&2
    exit 1
  fi
  route_porcelain "$@"
}

case "${1:-help}" in
  version)  route_porcelain "$@" ;;
  --version|-v) shift; route_porcelain version "$@" ;;
  gate)     gate_command "$@" ;;
  doctor)   adoption_route "$@" ;;
  worktree) route_porcelain "$@" ;;
  shift)    route_porcelain "$@" ;;
  commit)   route_porcelain "$@" ;;
  spec)     route_porcelain "$@" ;;
  link)     adoption_route "$@" ;;
  init)     adoption_route "$@" ;;
  unlink)   route_porcelain "$@" ;;
  models)   route_porcelain "$@" ;;
  outline)  route_porcelain "$@" ;;
  structure) route_porcelain "$@" ;;
  idea)     route_porcelain "$@" ;;
  roadmap)  route_porcelain "$@" ;;
  status)   route_porcelain "$@" ;;
  dashboard) route_porcelain "$@" ;;
  canary)   route_porcelain "$@" ;;
  learnings) route_porcelain "$@" ;;
  maps)     route_porcelain "$@" ;;
  guards)   route_porcelain "$@" ;;
  diff)     route_porcelain "$@" ;;
  coverage) route_porcelain "$@" ;;
  tree-hash) route_binary "$@" ;;
  gate-run) route_binary "$@" ;;
  gate-phases) route_binary "$@" ;;
  guard-git) route_binary "$@" ;;
  resolve-model) route_binary "$@" ;;
  check-agent-line) route_binary "$@" ;;
  stop-verdict) route_binary "$@" ;;
  worktree-pool) route_binary "$@" ;;
  worktree-lease-file) route_binary "$@" ;;
  help|--help|-h) cat <<EOF
bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants.
  bench link [copy|symlink]  safely wire the kit into this repo for every harness
  bench init                 scaffold .bench/gate.sh in the current repo
  bench unlink [--dry-run]   remove the per-repo Bench footprint the manifest records
  bench models               list advisory model-id candidates for the line binding
  bench structure            flag oversized files + crowded dirs (wire into the gate)
  bench idea "<text>"        park an out-of-scope idea in IDEAS.md (commit to nothing)
  bench roadmap              print the roadmap + drain status (IDEAS.md, learnings)
  bench status               ambient dashboard: what needs attention + the next action
  bench dashboard [--stdout] write a self-contained HTML snapshot of the board (--stdout emits it)
  bench canary [root]        run the gate against known-broken fixtures
  bench learnings            open journal entries as a TOON table (date, title)
  bench maps                 unresolved decision-map tickets as TOON (map, ticket, type, state)
  bench guards               every guard's deny surface as TOON (guard, boundary, denies)
  bench diff                 review base + changed files as TOON (--full appends log + diff body)
  bench coverage <spec>      acceptance-coverage state and rows as TOON (--check to validate)
  bench outline [path]       locate candidate seams (file:line) as TOON; does not identify the project's blessed seams
  bench doctor [--fix]       report (and repair) the PATH shim under a node version manager
  bench gate                 run the project gate (the oracle)
  bench gate pin             pin HEAD's .bench tree for pre-push verification
  bench worktree             warm, isolated worktree subshell
  bench worktree clean       remove clean out-of-pool worktrees after confirmation
  bench shift "<objective>"  gated loop in a pooled worktree; commit on green
  bench commit -m <msg> <path>...  gate, then commit named paths on green (--spec flips its status)
  bench spec implemented <slug>    flip a spec's Status: staged line to implemented
  bench spec retire <slug>         delete a merged spec + its review pickup (validated)
  bench spec history <slug>        retire/delete commits for a spec, newest first (TOON)
  bench version              print the installed benchkit version (os/arch)
EOF
  ;;
  # An unrecognized token (a typo, not a help request) is not this shell's job to
  # explain: route it to the Go binary so its own default-case "unknown subcommand"
  # handler (cmd/bench/main.go's run()) renders the message and exit 2 — one source
  # of that message, not a second copy here.
  *) route_binary "$@" ;;
esac
