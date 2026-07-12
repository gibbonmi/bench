#!/usr/bin/env bash
# Claude WorktreeCreate/WorktreeRemove adapter. This shim resolves the launcher
# through the shared resolver and passes the official stdin event to the Go adapter;
# it owns no JSON parsing, request identity, ownership, lock, or cleanup policy.
set -uo pipefail

if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: worktree-lifecycle\n'
  printf 'boundary: WorktreeCreate/WorktreeRemove\n'
  printf 'denies: nothing (informational)\n'
  printf 'why: routes Claude worktree events through the deterministic Bench lifecycle\n'
  exit 0
fi

case "${1:-}" in
  create | remove) action="$1" ;;
  *) printf 'usage: worktree-lifecycle.sh create|remove\n' >&2; exit 2 ;;
esac

fail_closed() {
  printf 'bench worktree hook: %s — refusing to run unguarded\n' "$1" >&2
  exit 1
}

# Lifecycle creation/removal is safety-sensitive: a missing resolver or core must
# refuse the event rather than bypass ownership, locking, or recovery policy.
hook_dir="$(cd -- "${BASH_SOURCE[0]%/*}" && pwd -P)" || fail_closed "cannot resolve hook directory"
lib="$hook_dir/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || fail_closed "wrapper resolver missing at $lib"
# shellcheck source=../lib/resolve-bench.sh
. "$lib" || fail_closed "wrapper resolver could not be loaded"
cmd="$(bench_resolve_wrapper)" || fail_closed "bench core not found"

exec "$cmd" worktree-hook "$action"
