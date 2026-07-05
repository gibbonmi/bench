#!/usr/bin/env bash
# Canary fixture: a bench.sh with the roadmap capture sink REMOVED. The gate's
# `bench idea`/`bench roadmap` contract (check 1f) must go red against this — proving
# the contract still bites if someone deletes or breaks the feature. Targeted EXPECT
# substring is the first failing assertion: roadmap no longer reports "empty".
set -euo pipefail
case "${1:-help}" in
  *) echo "bench (roadmap feature removed in this fixture)" ;;
esac
