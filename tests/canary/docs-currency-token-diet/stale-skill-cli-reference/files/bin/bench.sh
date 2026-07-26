#!/usr/bin/env bash
# Canary fixture CLI: implements only `worktree`, while the fixture's skill
# reference file names a `bench worktee` that no route provides.
case "${1:-}" in
  worktree) echo ok ;;
esac
