#!/usr/bin/env bash
# shellcheck shell=bash
# Compatibility shim for gates that still source .bench/lib/canary-run.sh.
#
# Contract: the sourcing gate defines `root`, `err`, and `fail`, and does not use
# `set -e`. New scaffolded gates call `bench canary "$root"` directly; this file
# remains shipped so already-linked gates keep working after relink.

_canary_lib_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "$_canary_lib_dir/resolve-bench.sh" ]; then
  # shellcheck source=/dev/null
  . "$_canary_lib_dir/resolve-bench.sh"
fi

_canary_bench=
if type bench_resolve_wrapper >/dev/null 2>&1; then
  _canary_bench="$(bench_resolve_wrapper)" || _canary_bench=
fi
if [ -z "$_canary_bench" ]; then
  _canary_bench="$(command -v bench 2>/dev/null)" || _canary_bench=
fi

# `root` is provided by the sourcing gate; see the contract above.
# shellcheck disable=SC2154
if [ -z "$_canary_bench" ]; then
  err "canary sweep failed"
elif ! "$_canary_bench" canary "$root"; then
  err "canary sweep failed"
fi

unset _canary_lib_dir _canary_bench
