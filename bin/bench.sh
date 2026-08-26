#!/usr/bin/env bash
# bench is the operational substrate for the Bench workflow. It fuses a warm
# worktree pool (treehouse-lite) with a gated loop (gnhf-lite), where the gate is the
# external oracle: a shift commits only on green.
#
# The worktree lifecycle, the gated loop, and the gate resolution and record all live
# in the Go core now — internal/worktree, internal/shift, internal/gate. This wrapper
# only routes, plus a one-glance run_gate adapter. internal/gate owns gate config
# resolution: first ./.bench/gate.sh, then $BENCH_GATE, then auto-detect.
set -euo pipefail

# This export lets the Go core's worktree, shift, and gate commands resolve the same
# pool home the shell did; they read BENCH_HOME from the environment. The `:?` form
# names the missing input itself: under `set -u`, a bare $HOME dies on "HOME: unbound
# variable", a variable name with no action for an adopter meeting it on their first
# command. This form is reached only when BENCH_HOME is unset, so the existing
# precedence is unchanged.
export BENCH_HOME="${BENCH_HOME:-${HOME:?the Bench pool home needs BENCH_HOME set, or HOME set to derive it from}/.bench}"

# ---- gate: the oracle -------------------------------------------------------
# run_gate is the one-glance adapter over the Go core's `bench gate-run`. Resolution
# of the gate — .bench/gate.sh, then $BENCH_GATE, then auto-detect — the run from the
# repo root, and the verdict-cache record all live in internal/gate. So `bench gate`
# and the Stop hook's `<wrapper> gate` path share exactly one resolver, never a
# second live implementation.
#
# This function execs rather than runs: the gate case is terminal, and the binary
# owns the exit code and the record. A missing binary exits 127 through
# route_binary, which writes no forged verdict.
run_gate() { route_porcelain gate-run "$@"; }

gate_usage() { printf 'usage: bench gate [--fresh|pin]\n'; }

gate_command() {
  case "$#" in
    1) run_gate ;;
    2)
      case "$2" in
        pin) route_porcelain gate-pin ;;
        --fresh) run_gate --fresh ;;
        --help|-h|help) gate_usage ;;
        *) gate_usage >&2; return 2 ;;
      esac
      ;;
    *)
      if [[ "$2" == pin ]]; then
        route_porcelain gate-pin "${@:3}"
      fi
      gate_usage >&2
      return 2
      ;;
  esac
}

# This resolves where the canonical kit lives — the parent of this script's bin/ —
# by following symlinks without relying on the GNU-only `readlink -f`. It caps the
# walk at about 40 hops, the readlink -f / SYMLOOP_MAX convention, so a symlink cycle
# to the wrapper fails fast with a structured error instead of chasing targets
# forever. Unlike opening the file, `readlink` never ELOOPs on a cyclic target, so an
# uncapped loop here can spin even where the OS itself would refuse to open the
# path.
resolve_script_path() {
  local source="${BASH_SOURCE[0]:-$0}" dir target hops=0
  while [[ -L "$source" ]]; do
    hops=$((hops + 1))
    if (( hops > 40 )); then
      printf 'bench: symlink cycle or too many levels resolving %s (stopped after %d hops)\n' "$source" "$hops" >&2
      exit 1
    fi
    dir="$(cd -P "$(dirname "$source")" >/dev/null 2>&1 && pwd)"
    target="$(readlink "$source")"
    [[ "$target" == /* ]] && source="$target" || source="$dir/$target"
  done
  dir="$(cd -P "$(dirname "$source")" >/dev/null 2>&1 && pwd)"
  printf '%s/%s\n' "$dir" "$(basename "$source")"
}

# git_common_dir_abs <dir> echoes <dir>'s git common directory as an absolute path.
# It fails non-zero when <dir> is outside a repository or git is unavailable.
# `rev-parse --git-common-dir` answers relatively, a bare `.git`, inside a plain
# checkout, so this anchors the result at <dir> before it returns it. Two of these
# results compare equal only for two working trees of the same repository.
git_common_dir_abs() {
  local dir="$1" common base
  command -v git >/dev/null 2>&1 || return 1
  common="$(git -C "$dir" rev-parse --git-common-dir 2>/dev/null)" || return 1
  [[ -n "$common" ]] || return 1
  if [[ "$common" != /* ]]; then
    base="$(cd -P "$dir" >/dev/null 2>&1 && pwd)" || return 1
    common="$base/$common"
  fi
  common="$(cd -P "$(dirname "$common")" >/dev/null 2>&1 && pwd)/$(basename "$common")" || return 1
  printf '%s\n' "$common"
}

# kit_dir echoes the kit this invocation serves. It is normally the parent of this
# script's bin/, the tree the wrapper file lives in. The exception: when the current
# directory sits in a git worktree of that same repository, the kit is that
# worktree's top level. An operator who runs the main checkout's `bench` from PATH
# with a worktree as CWD is working on the worktree, and a BENCH_KIT naming the main
# checkout while the graded root is the worktree makes the gate drop its kit-only
# phases over a clean tree.
#
# Two conditions guard the exception, and a linked project repo fails both. The kit
# must be its own tree's top level: an adopted repo carries the kit at
# `<repo>/.bench`, a subdirectory, so its wrapper keeps naming `<repo>/.bench`.
# Sameness means identity of the git common directory, never a path prefix, so a CWD
# in an unrelated repository is never mistaken for a second working tree of this one.
# Every git failure — no repo, an unrelated repo, no git on PATH — falls back to the
# wrapper's own tree.
kit_dir() {
  local script wrapper_kit wrapper_root cwd_common wrapper_common cwd_root
  script="$(resolve_script_path)"
  wrapper_kit="$(cd -P "$(dirname "$script")/.." && pwd)"
  wrapper_root="$(git -C "$wrapper_kit" rev-parse --show-toplevel 2>/dev/null)" || { printf '%s\n' "$wrapper_kit"; return 0; }
  wrapper_root="$(cd -P "$wrapper_root" >/dev/null 2>&1 && pwd)" || { printf '%s\n' "$wrapper_kit"; return 0; }
  [[ "$wrapper_root" == "$wrapper_kit" ]] || { printf '%s\n' "$wrapper_kit"; return 0; }
  cwd_common="$(git_common_dir_abs .)" || { printf '%s\n' "$wrapper_kit"; return 0; }
  wrapper_common="$(git_common_dir_abs "$wrapper_kit")" || { printf '%s\n' "$wrapper_kit"; return 0; }
  [[ "$cwd_common" == "$wrapper_common" ]] || { printf '%s\n' "$wrapper_kit"; return 0; }
  cwd_root="$(git -C . rev-parse --show-toplevel 2>/dev/null)" || { printf '%s\n' "$wrapper_kit"; return 0; }
  [[ -n "$cwd_root" ]] || { printf '%s\n' "$wrapper_kit"; return 0; }
  cwd_root="$(cd -P "$cwd_root" >/dev/null 2>&1 && pwd)" || { printf '%s\n' "$wrapper_kit"; return 0; }
  printf '%s\n' "$cwd_root"
}

# recover_source_go_path repairs only this wrapper process when a source kit needs
# Go and the harness environment is partial. The bounded clean-login probe supplies
# data only. Every invalid result leaves PATH unchanged for the selected binary.
recover_source_go_path() {
  local kit="$1" timeout_cmd env_cmd bash_cmd raw marker status output executable dir old_path
  [[ -f "$kit/go.mod" && -f "$kit/scripts/go-build.sh" ]] || return 0
  command -v go >/dev/null 2>&1 && return 0
  timeout_cmd="$(command -v timeout 2>/dev/null)" || return 0
  env_cmd="$(command -v env 2>/dev/null)" || return 0
  bash_cmd="$(command -v bash 2>/dev/null)" || return 0
  marker=$'\036'
  raw="$(
    status=0
    "$timeout_cmd" -s KILL 2 "$env_cmd" -u ENVMAN_LOAD "$bash_cmd" -lc 'command -v go' 2>/dev/null || status=$?
    printf '%s%d' "$marker" "$status"
  )"
  status="${raw##*"$marker"}"
  output="${raw%"$marker"*}"
  [[ "$status" == 0 && "$output" == *$'\n' ]] || return 0
  executable="${output%$'\n'}"
  [[ -n "$executable" && "$executable" != *$'\n'* && ! "$executable" =~ [[:cntrl:]] ]] || return 0
  [[ "$executable" == /* && "${executable##*/}" == go && -f "$executable" && -x "$executable" ]] || return 0
  dir="${executable%/*}"
  [[ -n "$dir" ]] || dir=/
  old_path="$PATH"
  PATH="$dir:$PATH"
  export PATH
  if [[ "$(command -v go 2>/dev/null)" != "$executable" ]]; then
    PATH="$old_path"
    export PATH
  fi
}

# ---- strangler router: send a ported subcommand to the Go binary ------------
# This dispatch is the one seam that grows across the port: a later slice adds
# subcommand names here, never a second resolver. platform_pkg maps this host to its
# @redbench/<os>-<arch> package, using npm's os/cpu spelling. An off-matrix host
# returns non-zero, so the caller can emit the "unsupported platform" error instead
# of naming a package that does not exist.
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
  local kit="${1:-$(kit_dir)}" pkg_json manifest line version
  pkg_json="$kit/package.json"
  if [[ -f "$pkg_json" ]]; then
    while IFS= read -r line || [[ -n "$line" ]]; do
      if [[ "$line" =~ \"version\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
        version="${BASH_REMATCH[1]}"
        valid_package_version "$version" || return 1
        printf '%s\n' "$version"
        return 0
      fi
    done < "$pkg_json"
  fi
  manifest="$kit/link-manifest.tsv"
  if [[ -f "$manifest" ]]; then
    while IFS= read -r line || [[ -n "$line" ]]; do
      if [[ "$line" == \#kit$'\t'* ]]; then
        version="${line#*$'\t'}"
        valid_package_version "$version" || return 1
        printf '%s\n' "$version"
        return 0
      fi
    done < "$manifest"
  fi
  return 1
}

valid_package_version() {
  [[ "$1" =~ ^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$ ]]
}

# main_tree_kit <kit> handles the case where <kit> sits inside a linked git worktree
# — kit_dir names one whenever the invocation's CWD is one. It echoes the same kit
# path re-anchored under the main worktree's root. It echoes nothing when <kit> is
# the main tree itself, is outside any repo, or the mapping is degenerate.
#
# A linked worktree carries the tracked tree but not untracked artifacts (dist/,
# node_modules/), so the binary that serves a worktree lives in the main tree. This
# function always returns 0: a failed resolution degrades to "no extra candidates",
# never to a caller-visible error.
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

# bench_binary_path echoes the resolved Go binary path, or fails with a distinct
# exit code the caller maps to a message. The resolution order is: the platform
# package bundled under the wrapper's node_modules; the hoisted sibling npm produces
# for a global install; the exact version/target cache; then a repo-local dev build.
# Each candidate is tried against the kit dir itself, then re-anchored at the main
# tree when the kit dir is a linked git worktree, because a worktree carries the
# tracked wrapper but not the untracked binary.
#
# The first executable, non-empty match wins. A present-but-empty or non-executable
# file counts as missing, and this function never names it, so a torn build falls
# through instead of resolving to a non-runnable path. Exit 2 means an off-matrix
# platform with no package to name; exit 127 means no binary exists for this
# platform. This is the one source of both the platform-to-package mapping and the
# resolution order; route_binary, which execs the path, and the status adapters,
# which capture its output, both share it.
bench_binary_path() {
  local kit="${1:-$(kit_dir)}" pkg suffix version cache c k main
  platform_pkg >/dev/null || return 2
  pkg="$(platform_pkg)"
  suffix="$(platform_suffix)"
  main="$(main_tree_kit "$kit")"
  for k in "$kit" ${main:+"$main"}; do
    for c in "$k/node_modules/$pkg/bin/bench" "$k/../$pkg/bin/bench"; do
      [[ -x "$c" && -s "$c" ]] && { printf '%s\n' "$c"; return 0; }
    done
  done
  version="$(package_version "$kit" 2>/dev/null || true)"
  if [[ -n "$version" ]]; then
    cache="$BENCH_HOME/cache/bin/$version/$suffix/bench"
    [[ -x "$cache" && -s "$cache" ]] && { printf '%s\n' "$cache"; return 0; }
  fi
  for k in "$kit" ${main:+"$main"}; do
    c="$k/dist/bench"
    [[ -x "$c" && -s "$c" ]] && { printf '%s\n' "$c"; return 0; }
  done
  return 127
}

repair_binary() {
  local kit="$1" wrapper="$2" mode="${3:-repair}" script version suffix
  if [[ "$mode" == repair ]]; then
    if [[ "${BENCH_OFFLINE:-}" == 1 ]]; then
      echo "bench: repair suppressed by BENCH_OFFLINE=1" >&2
      return 1
    fi
    if [[ -n "${BENCH_NO_REPAIR:-}" ]]; then
      echo "bench: repair disabled by BENCH_NO_REPAIR" >&2
      return 1
    fi
  fi
  script="$(dirname "$wrapper")/bench-repair-binary.mjs"
  if ! command -v node >/dev/null 2>&1; then
    echo "bench: repair skipped because node is not on PATH" >&2
    return 1
  fi
  [[ -f "$script" ]] || return 1
  version="$(package_version "$kit" 2>/dev/null || true)"
  suffix="$(platform_suffix 2>/dev/null || true)"
  [[ -n "$version" && -n "$suffix" ]] || return 1
  node "$script" "$mode" "$kit" "$(platform_pkg)" "$version" "$suffix"
}

repair_command() {
  local mode=repair kit wrapper
  case "$#:${1:-}" in
    0:) ;;
    1:--prune) mode=prune ;;
    *)
      echo 'usage: bench repair [--prune]' >&2
      exit 2
      ;;
  esac
  kit="${BENCH_KIT:-$(kit_dir)}"
  wrapper="$(resolve_script_path)"
  repair_binary "$kit" "$wrapper" "$mode"
}

# route_binary <subcommand> [args...] resolves and execs the Go binary, passing the
# whole argv through. This is the one seam the strangler grows: a later slice adds
# subcommand names to the dispatch below, never a second resolver.
route_binary() {
  local bin rc kit wrapper repair_rc physical
  kit="${BENCH_KIT:-$(kit_dir)}"
  wrapper="$(resolve_script_path)"
  recover_source_go_path "$kit"
  if [[ -n "${BENCH_RUN_BINARY+x}" ]]; then
    bin="${BENCH_RUN_BINARY:-}"
    case "$bin" in
      /*) ;;
      *) echo 'bench: inherited BENCH_RUN_BINARY is not absolute' >&2; exit 1 ;;
    esac
    if [[ ! -f "$bin" || ! -x "$bin" || -L "$bin" ]]; then
      echo 'bench: inherited BENCH_RUN_BINARY is not a regular executable' >&2
      exit 1
    fi
    physical="$(cd -P "$(dirname "$bin")" 2>/dev/null && pwd)/$(basename "$bin")"
    if [[ "$physical" != "$bin" ]]; then
      echo 'bench: inherited BENCH_RUN_BINARY is not a cleaned physical path' >&2
      exit 1
    fi
    BENCH_KIT="$kit" BENCH_WRAPPER="${BENCH_WRAPPER:-$wrapper}" exec "$bin" "$@"
  fi
  bin="$(bench_binary_path "$kit")" && rc=0 || rc=$?
  case "$rc" in
    0) BENCH_KIT="${BENCH_KIT:-$kit}" BENCH_WRAPPER="${BENCH_WRAPPER:-$wrapper}" exec "$bin" "$@" ;;
    2)
      echo "bench: unsupported platform: $(uname -s | tr '[:upper:]' '[:lower:]')/$(uname -m)" >&2
      exit 2
      ;;
    *)
      if [[ "${BENCH_ALLOW_IMPLICIT_REPAIR:-}" == 1 && "${BENCH_REPAIR:-}" == 1 ]]; then
        repair_binary "$kit" "$wrapper" && repair_rc=0 || repair_rc=$?
        [[ "$repair_rc" == 130 || "$repair_rc" == 143 ]] && exit "$repair_rc"
        bin="$(bench_binary_path "$kit")" && BENCH_KIT="${BENCH_KIT:-$kit}" BENCH_WRAPPER="${BENCH_WRAPPER:-$wrapper}" exec "$bin" "$@"
      fi
      echo 'bench: no pinned binary for this platform — reinstall redbench or run bench repair' >&2
      exit 127
      ;;
  esac
}

route_porcelain() {
  BENCH_ALLOW_IMPLICIT_REPAIR=1 route_binary "$@"
}

# ---- worktree land: the stable promotion owner ------------------------------
# Public landing runs only under the installed promotion broker. The route refuses
# every inherited routing override before any repository read, then authenticates the
# broker through the installation manifest beside this wrapper — path, version,
# platform, and executable digest. `bench doctor --fix` (and so the release install)
# publishes the manifest and broker together. The installer owns the platform fact and
# writes it; the route requires the field but never derives a second copy to compare
# against, because the digest binds the exact executable this host runs. Current-
# directory state and repository executables never join this selection, so repository
# code cannot authorize its own publication.
land_repair_advice() {
  echo "bench: run 'bench doctor --fix' (or reinstall redbench) to republish the promotion broker" >&2
}

file_sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    return 1
  fi
}

land_route() {
  local script bindir install manifest key value broker='' version='' platform='' digest='' installed actual
  if [[ -n "${BENCH_KIT+x}" ]]; then
    echo 'bench: worktree land does not honor inherited BENCH_KIT; unset it and re-run' >&2
    exit 1
  fi
  if [[ -n "${BENCH_RUN_BINARY+x}" ]]; then
    echo 'bench: worktree land does not honor inherited BENCH_RUN_BINARY; unset it and re-run' >&2
    exit 1
  fi
  if [[ -n "${BENCH_WRAPPER+x}" ]]; then
    echo 'bench: worktree land does not honor inherited BENCH_WRAPPER; unset it and re-run' >&2
    exit 1
  fi
  script="$(resolve_script_path)"
  bindir="$(cd -P "$(dirname "$script")" && pwd)"
  install="$(cd -P "$bindir/.." && pwd)"
  manifest="$bindir/bench-broker.manifest"
  if [[ ! -f "$manifest" ]]; then
    echo "bench: no promotion-broker manifest at $manifest" >&2
    land_repair_advice
    exit 127
  fi
  while IFS=$'\t' read -r key value || [[ -n "$key" ]]; do
    case "$key" in
      path) broker="$value" ;;
      version) version="$value" ;;
      platform) platform="$value" ;;
      sha256) digest="$value" ;;
    esac
  done < "$manifest"
  if [[ -z "$broker" || -z "$version" || -z "$platform" || -z "$digest" ]]; then
    echo "bench: promotion-broker manifest at $manifest is incomplete" >&2
    land_repair_advice
    exit 127
  fi
  installed="$(package_version "$install" 2>/dev/null || true)"
  if [[ -n "$installed" && "$installed" != "$version" ]]; then
    echo "bench: promotion broker version $version does not match installed package $installed" >&2
    land_repair_advice
    exit 127
  fi
  [[ "$broker" == /* ]] || broker="$install/$broker"
  if [[ ! -f "$broker" || ! -x "$broker" || ! -s "$broker" || -L "$broker" ]]; then
    echo "bench: promotion broker at $broker is not a regular executable" >&2
    land_repair_advice
    exit 127
  fi
  actual="$(file_sha256 "$broker" 2>/dev/null || true)"
  if [[ -z "$actual" || "$actual" != "$digest" ]]; then
    echo "bench: promotion broker at $broker does not match its manifest digest" >&2
    land_repair_advice
    exit 127
  fi
  # The authenticated broker builds the prospective executable. Give it the same
  # bounded toolchain recovery that ordinary wrapper routes receive.
  recover_source_go_path "$install"
  exec "$broker" "$@"
}

adoption_route() {
  local kit
  kit="${BENCH_KIT:-$(kit_dir)}"
  if [[ ! -d "$kit/.agents/commands" || ( ! -f "$kit/.bench/BENCH.md" && ! -f "$kit/AGENTS.md" ) ]]; then
    echo "bench: link/init/doctor must run from the real Bench kit; no source asset tree at $kit" >&2
    echo "bench: use the installed 'bench' command or the source kit's bin/bench.sh" >&2
    exit 1
  fi
  route_porcelain "$@"
}

case "${1:-}" in
  "") route_porcelain status --route ;;
  version)  route_porcelain "$@" ;;
  --version|-v) shift; route_porcelain version "$@" ;;
  gate)     gate_command "$@" ;;
  doctor)   adoption_route "$@" ;;
  repair)   shift; repair_command "$@" ;;
  worktree)
    if [[ "${2:-}" == land ]]; then
      land_route "$@"
    fi
    route_porcelain "$@"
    ;;
  resume-clean) route_porcelain "$@" ;;
  shift)    route_porcelain "$@" ;;
  commit)   route_porcelain "$@" ;;
  spec)     route_porcelain "$@" ;;
  setup)    adoption_route "$@" ;;
  link)     adoption_route "$@" ;;
  init)     adoption_route "$@" ;;
  unlink)   route_porcelain "$@" ;;
  upgrade)  adoption_route "$@" ;;
  models)   route_porcelain "$@" ;;
  outline)  route_porcelain "$@" ;;
  structure) route_porcelain "$@" ;;
  skills-index) route_porcelain "$@" ;;
  idea)     route_porcelain "$@" ;;
  roadmap)  route_porcelain "$@" ;;
  status)   route_porcelain "$@" ;;
  handoff)  route_porcelain "$@" ;;
  commands) route_binary "$@" ;;
  dashboard) route_porcelain "$@" ;;
  canary)   route_porcelain "$@" ;;
  anchors)  route_porcelain "$@" ;;
  learnings) route_porcelain "$@" ;;
  maps)     route_porcelain "$@" ;;
  guards)   route_porcelain "$@" ;;
  diff)     route_porcelain "$@" ;;
  preflight) route_porcelain "$@" ;;
  coverage) route_porcelain "$@" ;;
  test)     route_porcelain "$@" ;;
  tree-hash) route_binary "$@" ;;
  gate-run) route_binary "$@" ;;
  freshness-check) route_binary "$@" ;;
  gate-phases) route_binary "$@" ;;
  gate-go) route_binary "$@" ;;
  release-preflight) route_binary "$@" ;;
  prep-release) route_binary "$@" ;;
  release)  route_binary "$@" ;;
  guard-git) route_binary "$@" ;;
  guard-bench-follow-on) route_binary "$@" ;;
  session-inspect) route_binary "$@" ;;
  resolve-model) route_binary "$@" ;;
  check-agent-line) route_binary "$@" ;;
  stop-verdict) route_binary "$@" ;;
  worktree-pool) route_binary "$@" ;;
  worktree-lease-file) route_binary "$@" ;;
  worktree-hook) route_porcelain "$@" ;;
  help|--help|-h) shift; route_porcelain help "$@" ;;
  # Explaining an unrecognized token — a typo, not a help request — is not this
  # shell's job. Route it to the Go binary, so its own default-case "unknown
  # subcommand" handler (cmd/bench/main.go's run()) renders the message and exit 2.
  # This keeps one source of that message, not a second copy here.
  *) route_binary "$@" ;;
esac
