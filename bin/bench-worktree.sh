# Worktree pool support for bench.sh: warm, isolated, reusable worktrees under
# $BENCH_HOME, guarded by atomic owned leases. Sourced by bin/bench.sh; uses its
# repo_root/default_branch helpers and $BENCH_HOME.

# One-glance adapters over the Go core (internal/worktree), which owns the pool-path
# (cksum key) and lease-file (git-path) conventions. No logic here — the facts live in
# one place across the slice-3→7 window; slice 7 folds these call sites into Go.
worktree_pool() {
  "$(bench_binary_path)" worktree-pool "$1"
}

worktree_lease_file() {
  "$(bench_binary_path)" worktree-lease-file "$1"
}

# Atomically claim a worktree's lease. The claim is an O_EXCL (noclobber)
# create recording "<pid> <utc-time>"; on a race the loser fails and scans on.
# An existing lease is reclaimable only when its owner is provably gone: a
# recorded pid that no longer runs, or unreadable/legacy content older than a
# minute by mtime (fresh-and-empty is a writer mid-claim, not a zombie). The
# reclaim mv is what makes two reclaimers safe — only one mover wins.
worktree_lease_claim() {
  local lease pid
  lease="$(worktree_lease_file "$1")"
  if ( set -C; printf '%s %s\n' "$$" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$lease" ) 2>/dev/null; then
    return 0
  fi
  pid="$(awk '{print $1; exit}' "$lease" 2>/dev/null || true)"
  if [[ "$pid" =~ ^[0-9]+$ ]]; then
    kill -0 "$pid" 2>/dev/null && return 1
  else
    [[ -n "$(find "$lease" -mmin +1 2>/dev/null)" ]] || return 1
  fi
  mv "$lease" "$lease.stale.$$" 2>/dev/null || return 1
  rm -f "$lease.stale.$$"
  ( set -C; printf '%s %s\n' "$$" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$lease" ) 2>/dev/null
}

worktree_acquire() {
  local root reset_ref reset_mode pool wt
  root="$1"
  reset_ref="${2:-}"
  reset_mode="${3:-hard}"
  pool="$(worktree_pool "$root")"; mkdir -p "$pool"
  git -C "$root" fetch -q origin 2>/dev/null || true
  for d in "$pool"/*/; do
    [[ -d "$d/.git" || -f "$d/.git" ]] || continue
    [[ -z "$(git -C "$d" status --porcelain 2>/dev/null)" ]] || continue
    worktree_lease_claim "$d" || continue
    wt="$d"; break
  done
  # A fresh dir is claimable by other scanners the moment `worktree add`
  # finishes, so a lost claim here just means someone else now owns it — mint
  # another rather than sharing or failing.
  local tries=0 cand
  while [[ -z "${wt:-}" && "$tries" -lt 3 ]]; do
    tries=$((tries+1))
    cand="$pool/$(date +%s)-$$-$tries"
    git -C "$root" worktree add -q --detach "$cand" "origin/$(default_branch)" 2>/dev/null \
      || git -C "$root" worktree add -q --detach "$cand" || break
    worktree_lease_claim "$cand" && wt="$cand"
  done
  [[ -n "${wt:-}" ]] || { echo "could not lease a pool worktree" >&2; return 1; }
  git -C "$wt" switch -q --detach >/dev/null 2>&1 || true
  if [[ -n "$reset_ref" ]]; then
    if ! git -C "$wt" reset -q --hard "$reset_ref"; then
      [[ "$reset_mode" == soft ]] || return 1
      git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
    fi
  else
    git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
  fi
  git -C "$wt" clean -qfdx
  printf '%s\n' "${wt%/}"
}

worktree_release() {
  local wt="$1" lease pid
  [[ -n "$wt" ]] || return 0
  lease="$(worktree_lease_file "$wt")"
  # Only the owner cleans and unleases: a stale-reclaimed worktree belongs to
  # its new owner now, and the old owner's deferred trap must leave it alone.
  pid="$(awk '{print $1; exit}' "$lease" 2>/dev/null || true)"
  if [[ "$pid" =~ ^[0-9]+$ && "$pid" != "$$" ]] && kill -0 "$pid" 2>/dev/null; then
    return 0
  fi
  rm -f "$lease"
  git -C "$wt" switch -q --detach >/dev/null 2>&1 || true
  git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
  git -C "$wt" clean -qfdx >/dev/null 2>&1 || true
  rm -f "$lease"
}

worktree() {
  local root target wt
  root="$(repo_root)"
  target="origin/$(default_branch)"
  wt="$(worktree_acquire "$root" "$target" soft)"
  echo "🪵 worktree: $wt  (exit to release)" >&2
  ( cd "$wt" && "${SHELL:-bash}" ) || true
  worktree_release "$wt"
  echo "🪵 released" >&2
}
