#!/usr/bin/env bash
# bench:managed-pre-push
# name: pre-push
# boundary: pre-push
# denies: direct push to the protected branch
# why: the merge belongs to the reviewer; agents open a PR instead of pushing the protected branch
# 'bench link' installs this hook. The merge belongs to the human; an agent does not
# push __BENCH_DEFAULT_BRANCH__.
# Resolve the protected branch live, because a repo linked before its remote existed
# baked in a fabricated default. Query origin/HEAD, and fall back to the baked token
# only when it is unresolvable — no remote, or origin/HEAD unset.
protected="__BENCH_DEFAULT_BRANCH__"
live_head="$(git symbolic-ref --short refs/remotes/origin/HEAD 2>/dev/null || true)"
if [[ -n "$live_head" ]]; then
  protected="${live_head#origin/}"
fi
# The reviewer who owns the merge can lift the clause for one repository with
# 'git config bench.allowProtectedPush true'.
allow_protected="$(git config --type=bool --get bench.allowProtectedPush 2>/dev/null || true)"
# The second read arm passes a last line with no trailing newline to the clause.
while IFS= read -r line || [ -n "$line" ]; do
  read -r _ _ remote_ref _ <<< "$line"
  if [[ "$allow_protected" != "true" && "$remote_ref" == "refs/heads/$protected" ]]; then
    echo "blocked: direct push to $protected. Open a PR or merge it yourself." >&2
    exit 1
  fi
done
exit 0
