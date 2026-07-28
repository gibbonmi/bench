#!/usr/bin/env bash
# Canary fixture: a bench.sh with the `status` ambient surface REMOVED. The gate's
# `bench status` contract (check 1g) must go red against this — proving the contract
# still bites if the renderer is deleted or broken. Targeted EXPECT substring is the
# first failing assertion: a clean repo no longer reports the all-clear line.
set -euo pipefail
case "${1:-help}" in
  *) echo "bench (status feature removed in this fixture)" ;;
esac
