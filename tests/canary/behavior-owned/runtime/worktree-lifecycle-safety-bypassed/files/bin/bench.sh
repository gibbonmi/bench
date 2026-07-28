#!/usr/bin/env bash
# Canary: bypass ownership, preservation-before-unlock, and fail-closed locking.
set -euo pipefail

case "${1:-}" in
  resume-clean)
    primary="$(git rev-parse --show-toplevel)"
    while IFS= read -r path; do
      [[ -n "$path" && "$path" != "$primary" ]] || continue
      git worktree unlock "$path" 2>/dev/null || true
      git worktree remove --force "$path"
    done < <(git worktree list --porcelain | awk '$1 == "worktree" { print substr($0, 10) }')
    ;;
  worktree)
    case "${2:-}" in
      create)
        path="$BENCH_HOME/canary-owned"
        mkdir -p "${path%/*}"
        git worktree add -q --detach --lock --reason 'bench canary unsafe' "$path" HEAD
        printf 'worktree_create[1]{path,assignment,state}:\n  "%s",canary,active\n' "$path"
        ;;
      release)
        path="${@: -1}"
        git worktree unlock "$path"
        git worktree remove --force "$path"
        ;;
      *) exit 2 ;;
    esac
    ;;
  *) exit 0 ;;
esac
