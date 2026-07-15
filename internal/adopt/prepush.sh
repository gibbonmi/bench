#!/usr/bin/env bash
# bench:managed-pre-push
# name: pre-push
# boundary: pre-push
# denies: direct push to the protected branch; .bench drift when pinned
# why: the merge belongs to the reviewer; agents open a PR instead of pushing the protected branch
# Installed by 'bench link'. The merge is the human's; agents don't push __BENCH_DEFAULT_BRANCH__.
# Read the pin state once, at the top: the drift clause is armed only when a non-empty pin exists.
pin_path="$(git rev-parse --git-path bench-gate-pin 2>/dev/null || true)"
pin_tree=""
if [[ -n "$pin_path" && -f "$pin_path" ]]; then
  IFS= read -r pin_tree < "$pin_path" || true
fi
# Resolve the protected branch live so a repo linked before its remote existed baked a
# fabricated default, so query origin/HEAD and fall back to the baked token only when it
# is unresolvable (no remote, or origin/HEAD unset).
protected="__BENCH_DEFAULT_BRANCH__"
live_head="$(git symbolic-ref --short refs/remotes/origin/HEAD 2>/dev/null || true)"
if [[ -n "$live_head" ]]; then
  protected="${live_head#origin/}"
fi
if [[ -z "$pin_tree" ]]; then
  echo "bench: gate unpinned - run 'bench gate pin' to enable .bench drift checks." >&2
fi
ref_lines=()
while IFS= read -r line || [ -n "$line" ]; do
  ref_lines+=("$line")
  read -r _ local_oid _ _ <<< "$line"
  if [[ -n "$pin_tree" && ! "$local_oid" =~ ^0+$ ]]; then
    if ! bench_tree="$(git rev-parse "$local_oid:.bench" 2>/dev/null)"; then
      echo "blocked: pushed commit has no .bench tree. Review the gate change, then run 'bench gate pin'." >&2
      exit 1
    fi
    if [[ "$bench_tree" != "$pin_tree" ]]; then
      echo "blocked: pushed .bench tree does not match bench gate pin. Review the gate change, then run 'bench gate pin'." >&2
      exit 1
    fi
  fi
done
for line in "${ref_lines[@]}"; do
  read -r _ _ remote_ref _ <<< "$line"
  if [[ "$remote_ref" == "refs/heads/$protected" ]]; then
    echo "blocked: direct push to $protected. Open a PR or merge it yourself." >&2
    exit 1
  fi
done
exit 0
